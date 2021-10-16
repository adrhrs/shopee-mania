package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func triggerCrawl(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	go CrawlByCategory()

	data := BasicResp{
		Msg:     "Triggered",
		Latency: time.Since(t).String(),
	}

	log.Println("crawl triggered")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func aggResult(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	code := http.StatusOK
	msg := "Aggregated async"
	typeCrawl := req.FormValue("type")
	if typeCrawl == "" {
		msg = "fail type"
		code = http.StatusInternalServerError
	}
	go AggResultV2(typeCrawl)

	data := BasicResp{
		Msg:     msg,
		Latency: time.Since(t).String(),
	}

	log.Println("agg result v2")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
