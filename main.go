package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func main() {

	//util
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/check", checkClientIP)

	//trigger BE
	http.HandleFunc("/trigger", triggerCrawl)
	http.HandleFunc("/directory", getResult)
	http.HandleFunc("/agg", aggResult)

	//FE related
	http.HandleFunc("/getCategory", prepCat)
	http.HandleFunc("/fetch", fetchResult)
	http.HandleFunc("/fetch-single", fetchResultSingle)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/detail", handleDetail)
	http.HandleFunc("/track", handleTrackProduct)

	//serve static file
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//cron
	cr := cron.New()
	cr.AddFunc("@daily", CrawlWrapper)
	cr.Start()

	//Ready gas gas gas
	fmt.Printf("Starting server at port 6001, up and running\n")
	if err := http.ListenAndServe(":6001", nil); err != nil {
		log.Fatal(err)
	}

}
