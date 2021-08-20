package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
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

func CrawlByCategory() {
	var (
		ip                = strings.Replace(getIP(), ".", "-", -1)
		whitelistedLv1IDs = []int{100017, 100630, 100010}
		catIDs            []int
		start             = time.Now()
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
		log.Println(r.MatchID, r.Dur, len(r.Result), r.Err)
		writeCSV(fmt.Sprintf("%v-%v-%v-data", ip, time.Now().In(loc).Format("2006-01-02"), r.MatchID), r.Result, col)
	}

	log.Println(time.Since(start).String())
}

func prepareCategory(whitelistedLv1IDs []int) (catIDs []int) {

	var (
		lv1IDs = make(map[int]int)
	)

	for _, w := range whitelistedLv1IDs {
		lv1IDs[w] = 1
	}

	cats, err := doGetCategories()
	if err != nil {
		log.Println(err)
		return
	}
	for _, lv1 := range cats {
		if _, ok := lv1IDs[lv1.Main.Catid]; ok {
			catIDs = append(catIDs, lv1.Main.Catid)
			for _, lv2 := range lv1.Sub {
				catIDs = append(catIDs, lv2.Catid)
				for _, lv3 := range lv2.SubSub {
					catIDs = append(catIDs, lv3.Catid)
				}
			}
		}
	}

	return
}

func workerDoCrawl(id int, jobs <-chan int, results chan<- workerDoCrawlReturn) {
	for j := range jobs {
		var (
			formattedData = [][]string{}
			t             = time.Now()
			errCrawl      error
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
			}
		}

		results <- workerDoCrawlReturn{
			Result:  formattedData,
			Dur:     time.Since(t).String(),
			Err:     errCrawl,
			MatchID: j,
		}
	}
}

func formatCSV(pr ItemBasicNecessary) (prs []string) {
	pr = fixProduct(pr)
	prs = []string{fmt.Sprintf("%v", pr.Itemid),
		fmt.Sprintf("%v", pr.Shopid),
		fmt.Sprintf("%v", pr.Catid),
		fmt.Sprintf("%v", pr.Price),
		fmt.Sprintf("%v", pr.HistoricalSold),
		fmt.Sprintf("%v", time.Now().In(loc).Format(dateFormat)),
		fmt.Sprintf("%v", pr.Name),
	}
	return
}

func writeCSV(fileName string, data [][]string, column [][]string) (err error) {

	file, err := os.OpenFile(""+fileName+".csv", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)
	strWrite := column
	strWrite = append(strWrite, data...)
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()

	return

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

func AggResult() (files []string) {

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
				if i > 0 {
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

	log.Println(len(uniqueProductCat), len(aggByCat), prepareDur, time.Since(writeStart).String())

	return

}

func formatCSVAGG(pr SimpleProduct, prices, solds []string, extras []float64) (prs []string) {
	prs = []string{fmt.Sprintf("%v", pr.Itemid),
		fmt.Sprintf("%v", pr.Shopid),
		fmt.Sprintf("%v", pr.Catid),
		fmt.Sprintf("%v", pr.Name),
	}
	prs = append(prs, prices...)
	prs = append(prs, solds...)

	for _, e := range extras {
		s := fmt.Sprintf("%.0f", e)
		prs = append(prs, s)
	}
	return
}

func getDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 00, 00, 00, 0, time.UTC)
}
