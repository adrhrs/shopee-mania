package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	col = [][]string{{"item_id", "shop_id", "cat_id", "price", "sold", "date", "name"}}
	ip  = "dummy"

	loc, _         = time.LoadLocation("Asia/Jakarta")
	dateFormat     = "2006-01-02 15:04:05"
	dateOnlyFormat = "2006-01-02"

	pageTypeCat        = "search"
	pageTypeShop       = "shop"
	defaultLimitCat    = 60
	defaultLimitShop   = 30
	defaultMaxPageCat  = 5
	defaultMaxPageShop = 50
	aggTypeCat         = "data"
	aggTypeShop        = "shop"

	shopType = 2
	catType  = 1
)

func CrawlWrapper() {
	cr := CrawlByCategory()
	sr := CrawlByShop()
	ar := AggResultV2(aggTypeCat)
	sar := AggResultV2(aggTypeShop)

	content := prepContent(cr, sr, ar, sar)
	err := sendEmail(content)
	if err == nil {
		log.Println("email sent!")
	}
}

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
	jobs := make(chan SearchParam, numJobs)
	results := make(chan workerDoCrawlReturn, numJobs)

	for w := 1; w <= 10; w++ {
		go workerDoCrawl(w, jobs, results)
	}

	for _, matchID := range catIDs {
		jobs <- SearchParam{
			Matchid:   matchID,
			CrawlType: catType,
		}
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

func CrawlByShop() (result CrawlCronResult) {
	var (
		ip                = strings.Replace(getIP(), ".", "-", -1)
		shopIDs           = []int{11487927, 62582411, 39283823}
		start             = time.Now()
		errCount, totProd int
	)

	log.Println("prepare crawl by shop", len(shopIDs))

	var numJobs = len(shopIDs)
	jobs := make(chan SearchParam, numJobs)
	results := make(chan workerDoCrawlReturn, numJobs)

	for w := 1; w <= 10; w++ {
		go workerDoCrawl(w, jobs, results)
	}

	for _, matchID := range shopIDs {
		jobs <- SearchParam{
			Matchid:   matchID,
			CrawlType: shopType,
		}
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		r := <-results
		if r.Err != nil {
			errCount++
		}
		totProd += r.TotProd
		writeCSV(fmt.Sprintf("%v-%v-%v-shop", ip, time.Now().In(loc).Format("2006-01-02"), r.MatchID), r.Result, col)
	}

	result.CrawlDuration = time.Since(start).String()
	result.TotalShopIDs = len(shopIDs)
	result.TotalProduct = totProd
	result.ErrCount = errCount

	if len(shopIDs) > 0 {
		result.AvgProductCount = int(float64(totProd) / float64(len(shopIDs)))
	}

	log.Println(result)

	return
}

func workerDoCrawl(id int, jobs <-chan SearchParam, results chan<- workerDoCrawlReturn) {
	for j := range jobs {
		var (
			formattedData = [][]string{}
			t             = time.Now()
			errCrawl      error
			totProd       int
			maxPage       int
		)

		switch j.CrawlType {
		case catType:
			maxPage = defaultMaxPageCat
			j.Limit = defaultLimitCat
			j.PageType = pageTypeCat
		case shopType:
			maxPage = defaultMaxPageShop
			j.Limit = defaultLimitShop
			j.PageType = pageTypeShop
		}

		for i := 0; i < maxPage; i++ {
			j.Newest = i * j.Limit
			resp, _, err := doSearchCrawl(j)
			if err != nil {
				errCrawl = err
			}
			for _, r := range resp.Items {
				formattedData = append(formattedData, formatCSV(r.ItemBasic))
				totProd++
			}
			if len(resp.Items) == 0 {
				break
			}

		}

		results <- workerDoCrawlReturn{
			Result:  formattedData,
			Dur:     time.Since(t).String(),
			Err:     errCrawl,
			TotProd: totProd,
			MatchID: j.Matchid,
		}
	}
}
