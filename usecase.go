package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	col          = [][]string{{"item_id", "shop_id", "cat_id", "price", "sold", "date", "name"}}
	ip           = "dummy"
	defaultLimit = 60
	defaultPage  = 5
	loc, _       = time.LoadLocation("Asia/Jakarta")
	dateFormat   = "2006-01-02 15:04:05"
)

func CrawlByCategory() (result CrawlCronResult) {
	var (
		ip                = strings.Replace(getIP(), ".", "-", -1)
		whitelistedLv1IDs = []int{100017, 100630, 100010}
		catIDs            []int
		start             = time.Now()
		errCount, totProd int
	)

	catIDs = prepareCategory(whitelistedLv1IDs)
	log.Println("prepare crawl for", len(catIDs))

	var numJobs = len(catIDs)
	jobs := make(chan int, numJobs)
	results := make(chan workerDoCrawlReturn, numJobs)

	for w := 1; w <= 10; w++ {
		go workerDoCrawl(w, jobs, results)
	}

	for _, matchID := range catIDs {
		jobs <- matchID
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		r := <-results
		if r.Err != nil {
			errCount++
		}
		totProd += r.TotProd
		writeCSV(fmt.Sprintf("%v-%v-%v-data", ip, time.Now().In(loc).Format("2006-01-02"), r.MatchID), r.Result, col)
	}

	result.CrawlDuration = time.Since(start).String()
	result.TotalCategories = len(catIDs)
	result.TotalProduct = totProd
	result.ErrCount = errCount

	if len(catIDs) > 0 {
		result.AvgProductCount = int(float64(totProd) / float64(len(catIDs)))
	}

	log.Println(result)

	return
}

func workerDoCrawl(id int, jobs <-chan int, results chan<- workerDoCrawlReturn) {
	for j := range jobs {
		var (
			formattedData = [][]string{}
			t             = time.Now()
			errCrawl      error
			totProd       int
		)
		for i := 0; i < defaultPage; i++ {
			resp, _, err := doSearchCrawl(SearchParam{
				Matchid: j,
				Newest:  i * defaultLimit,
			})
			if err != nil {
				errCrawl = err
			}
			for _, r := range resp.Items {
				formattedData = append(formattedData, formatCSV(r.ItemBasic))
				totProd++
			}
		}

		results <- workerDoCrawlReturn{
			Result:  formattedData,
			Dur:     time.Since(t).String(),
			Err:     errCrawl,
			TotProd: totProd,
			MatchID: j,
		}
	}
}

func GetResult() (files []string) {

	f, err := os.Open(".")
	if err != nil {
		log.Println(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return

}

func AggResult() (files []string, result AggCronResult) {

	f, err := os.Open(".")
	if err != nil {
		log.Println(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}

	uniqueProductCat := make(map[string]SimpleProduct)
	uniqueDate := make(map[string]time.Time)
	productDate := make(map[string]SimpleProduct)
	aggByCat := make(map[string][][]string)

	prepareStart := time.Now()
	for _, file := range fileInfo {
		if strings.Contains(file.Name(), ".csv") && !(strings.Contains(file.Name(), "aggregated")) {
			files = append(files, file.Name())
			file, _ := os.Open(file.Name())
			r := csv.NewReader(file)
			i := 0
			dCatID := strings.Split(file.Name(), "-")[7]
			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if i > 0 && len(record) > 6 {
					t, _ := time.Parse(dateFormat, record[5])
					dateKey := getDate(t).Format("2006-01-02")
					if dateKey != "0001-01-01" {
						uniqueDate[dateKey] = t
					}

					productKey := fmt.Sprintf("%v-%v", record[0], dCatID)
					uniqueProductCat[productKey] = SimpleProduct{
						Itemid: record[0],
						Shopid: record[1],
						Catid:  dCatID,
						Name:   record[6],
					}

					pdKey := fmt.Sprintf("%v-%v", dateKey, productKey)
					productDate[pdKey] = SimpleProduct{
						Itemid:         record[0],
						Price:          record[3],
						HistoricalSold: record[4],
					}
				}
				i++
			}
			file.Close()
		}
	}
	prepareDur := time.Since(prepareStart).String()

	var dates []time.Time
	for _, v := range uniqueDate {
		dates = append(dates, v)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	i := 0
	writeStart := time.Now()

	for productKey, v := range uniqueProductCat {

		var (
			totPrice, avgPrice, avail, minSold float64
			maxSold, diffSold, estGMV          float64
			prices, solds                      []string
			fSolds, extras                     []float64
		)
		for _, t := range dates {
			dateKey := getDate(t).Format("2006-01-02")
			pdKey := fmt.Sprintf("%v-%v", dateKey, productKey)
			var price, sold = "-", "-"
			if val, ok := productDate[pdKey]; ok {
				price = val.Price
				sold = val.HistoricalSold

				fPrice, _ := strconv.ParseFloat(val.Price, 64)
				fSold, _ := strconv.ParseFloat(val.HistoricalSold, 64)

				fSolds = append(fSolds, fSold)
				totPrice += fPrice
				avail++
			}
			prices = append(prices, price)
			solds = append(solds, sold)
		}

		if len(fSolds) > 0 {
			sort.Slice(fSolds, func(i, j int) bool {
				return fSolds[i] < fSolds[j]
			})
			maxSold = fSolds[len(fSolds)-1]
			minSold = fSolds[0]
			diffSold = maxSold - minSold
		}

		if avail > 0 {
			avgPrice = totPrice / avail
		}
		estGMV = diffSold * avgPrice
		extras = []float64{avail, minSold, maxSold, diffSold, avgPrice, estGMV}

		aggByCat[v.Catid] = append(aggByCat[v.Catid], formatCSVAGG(v, prices, solds, extras))
		i++
	}

	basicCol := []string{"item_id", "shop_id", "cat_id", "name"}
	for _, t := range dates {
		dateKey := getDate(t).Format("2006-01-02")
		basicCol = append(basicCol, "price "+dateKey)
	}
	for _, t := range dates {
		dateKey := getDate(t).Format("2006-01-02")
		basicCol = append(basicCol, "sold "+dateKey)
	}
	basicCol = append(basicCol, "availibility", "min_sold", "max_sold", "delta_sold", "avg_price", "est_gmv")
	col := [][]string{basicCol}

	for k, v := range aggByCat {
		writeCSV("aggregated-"+k, v, col)
	}

	result.AggDuration = time.Since(prepareStart).String()
	result.ReadDuration = prepareDur
	result.CalcDuration = time.Since(writeStart).String()
	result.TotalDays = len(dates)

	log.Println(result)

	return

}

func crawlWrapper() {
	cr := CrawlByCategory()
	// _, ar := AggResult()

	content := prepContent(cr, AggCronResult{})
	err := sendEmail(content)
	if err == nil {
		log.Println("email sent!")
	}
}

func sendEmail(content string) (err error) {

	MIME := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	from := "adrbot8@gmail.com"
	password := "Adrian123."

	to := []string{"adrianafnandika@gmail.com"}
	subject := "Shopee Crawling Report " + getDate(time.Now()).Format("2006-01-02")

	host := "smtp.gmail.com"
	port := "587"
	body := "To: " + to[0] + "\r\nSubject: " + subject + "\r\n" + MIME + "\r\n" + content
	auth := smtp.PlainAuth("", from, password, host)

	err = smtp.SendMail(fmt.Sprintf("%v:%v", host, port), auth, from, to, []byte(body))
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func prepContent(cr CrawlCronResult, ar AggCronResult) (content string) {

	dateKey := getDate(time.Now()).Format("2006-01-02")

	head := `<!DOCTYPE html>
	<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office">
	<head>
	  <meta charset="UTF-8">
	  <meta name="viewport" content="width=device-width,initial-scale=1">
	  <meta name="x-apple-disable-message-reformatting">
	  <title></title>
	  <!--[if mso]>
	  <noscript>
		<xml>
		  <o:OfficeDocumentSettings>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		  </o:OfficeDocumentSettings>
		</xml>
	  </noscript>
	  <![endif]-->
	  <style>
		table, td, div, h1, p {font-family: Arial, sans-serif;}
	  </style>
	</head>
	<body style="margin:0;padding:0;">
	  <table role="presentation" style="width:100%;border-collapse:collapse;border:0;border-spacing:0;background:#ffffff;">
		<tr>
		  <td align="center" style="padding:0;">
			<table role="presentation" style="width:602px;border-collapse:collapse;border:1px solid #cccccc;border-spacing:0;text-align:left;">
			  <tr>
				<td align="center" style="padding:40px 0 30px 0;background:#1c1c1c;">
				  <img src="https://image.flaticon.com/icons/png/512/235/235253.png" alt="" width="300" style="height:auto;display:block;" />
				</td>
			  </tr>
			  <tr>
				<td style="padding:36px 30px 42px 30px;">
				  <table role="presentation" style="width:100%;border-collapse:collapse;border:0;border-spacing:0;">
					<tr>
					  <td style="padding:0 0 36px 0;color:#153643;">`

	tail := `<p style="margin:0;font-size:16px;line-height:24px;font-family:Arial,sans-serif;"><a href="http://188.166.252.251:6001/static/" style="color:#ee4c50;text-decoration:underline;"><b>Visit Site</b></a></p>
					  </td>
					</tr>
				  </table>
				</td>
			  </tr>
			  <tr>
				<td style="padding:20px;background:#eb7734;">
				  <table role="presentation" style="width:100%;border-collapse:collapse;border:0;border-spacing:0;font-size:9px;font-family:Arial,sans-serif;">
					<tr>
					  <td style="padding:0;width:50%;" align="left">
						<p style="margin:0;font-size:14px;line-height:16px;font-family:Arial,sans-serif;color:#ffffff;">
						  &reg; Adrian 2021
						</p>
					  </td>
					</tr>
				  </table>
				</td>
			  </tr>
			</table>
		  </td>
		</tr>
	  </table>
	</body>
	</html>`

	body := fmt.Sprintf(`
	
						<h1 style="font-size:24px;margin:0 0 20px 0;font-family:Arial,sans-serif;">Shopee Crawling Report %s</h1>
						<p style="margin:0 0 12px 0;font-size:16px;line-height:24px;font-family:Arial,sans-serif;text-align: justify;text-justify: inter-word;">
							Crawling done for <b>%d</b> categories, in <b>%s</b>. We encounter <b>%d</b> error during crawling. We got <b>%d</b> total products for today, 
							so we have average of <b>%v</b> products each category. Aggregation for <b>%d</b> days done in <b>%v</b>, including reading files in <b>%v</b> and calculating in <b>%v</b>.
						</p>
						
	`, dateKey, cr.TotalCategories, cr.CrawlDuration, cr.ErrCount, cr.TotalProduct, cr.AvgProductCount, ar.TotalDays, ar.AggDuration, ar.ReadDuration, ar.CalcDuration)

	content = head + body + tail
	return
}
