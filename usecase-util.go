package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"
)

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

func prepContent(cr, sr CrawlCronResult, ar, sar AggCronResult) (content string) {

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
							so we have average of <b>%v</b> products each category. Shop crawl done for <b>%d</b> shops and got <b>%d</b> products, finish in <b>%s</b>. 
							Category Aggregation for <b>%d</b> days done in <b>%v</b>, including reading files in <b>%v</b> and calculating in <b>%v</b>.
							Shop Aggregation for <b>%d</b> days done in <b>%v</b>, including reading files in <b>%v</b> and calculating in <b>%v</b>.
						</p>
						
	`, dateKey, cr.TotalCategories, cr.CrawlDuration, cr.ErrCount, cr.TotalProduct, cr.AvgProductCount, sr.TotalShopIDs, sr.TotalProduct, sr.CrawlDuration, ar.TotalDays, ar.AggDuration, ar.ReadDuration, ar.CalcDuration, sar.TotalDays, sar.AggDuration, sar.ReadDuration, sar.CalcDuration)

	content = head + body + tail
	return
}
