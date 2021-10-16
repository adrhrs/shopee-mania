package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func main() {

	EvaluateProductReviewer()
	return

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/trigger", triggerCrawl)
	http.HandleFunc("/result", getResult)
	http.HandleFunc("/agg", aggResult)
	http.HandleFunc("/fetch", fetchResult)
	http.HandleFunc("/getCategory", prepCat)
	http.HandleFunc("/download", handleDownload)
	http.HandleFunc("/detail", handleDetail)
	http.HandleFunc("/buyer", handleEvaluateBuyer)
	http.HandleFunc("/reviewer", handleEvaluateReviewer)

	http.HandleFunc("/check", checkClientIP)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	cr := cron.New()
	cr.AddFunc("@daily", CrawlWrapper)
	cr.Start()

	fmt.Printf("Starting server at port 6001, up and running\n")
	if err := http.ListenAndServe(":6001", nil); err != nil {
		log.Fatal(err)
	}

}
