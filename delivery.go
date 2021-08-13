package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func hello(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	data := BasicResp{
		Msg:     "Hello World",
		Data:    dir,
		Latency: time.Since(t).String(),
	}

	log.Println("log of this endpoint", dir)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func triggerCrawl(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	go CrawlByCategory()

	data := BasicResp{
		Msg:     "Triggered",
		Latency: time.Since(t).String(),
	}

	log.Println("crawl triggered")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
