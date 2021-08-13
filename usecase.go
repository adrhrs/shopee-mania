package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	col          = [][]string{{"item_id", "shop_id", "cat_id", "price", "sold", "date", "name"}}
	ip           = "dummy"
	defaultLimit = 60
	defaultPage  = 5
	loc, _       = time.LoadLocation("Asia/Jakarta")
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
		fmt.Println(r.MatchID, r.Dur, len(r.Result), r.Err)
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
		fmt.Sprintf("%v", time.Now().In(loc).Format("2006-01-02 15:04:05")),
		fmt.Sprintf("%v", pr.Name),
	}
	return
}

func writeCSV(fileName string, data [][]string, column [][]string) (err error) {

	// dir := os.Getenv("work_dir")
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
