package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func testCron() {
	log.Println("i am a cron")
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/trigger", triggerCrawl)

	cr := cron.New()
	cr.AddFunc("@daily", CrawlByCategory)
	cr.Start()

	fmt.Printf("Starting server at port 6001, up and running\n")
	if err := http.ListenAndServe(":6001", nil); err != nil {
		log.Fatal(err)
	}

}
