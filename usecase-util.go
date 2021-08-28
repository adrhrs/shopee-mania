package main

import (
	"encoding/csv"
	"fmt"
	"log"
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
