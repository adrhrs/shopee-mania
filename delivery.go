package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		Msg:     "Hello World Production, now with shop crawl",
		Data:    dir,
		Latency: time.Since(t).String(),
	}

	log.Println("log of this endpoint", dir)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func getResult(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	result := GetResult()

	data := BasicResp{
		Msg:     "Triggered",
		Latency: time.Since(t).String(),
		Data:    result,
	}

	log.Println("fetch directory")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func aggResult(w http.ResponseWriter, req *http.Request) {
	t := time.Now()

	result := AggResultV2()

	data := BasicResp{
		Msg:     "Aggregated async",
		Latency: time.Since(t).String(),
		Data:    result,
	}

	log.Println("agg result v2")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func fetchResult(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	msg := "fetch result"
	code := http.StatusOK
	id := req.FormValue("catid")
	result, err := FetchResult(id)
	if err != nil {
		msg = err.Error()
		code = http.StatusInternalServerError
	}

	data := BasicResp{
		Msg:     msg,
		Latency: time.Since(t).String(),
		Data:    result,
	}

	log.Println(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func prepCat(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	msg := "fetch cat"
	code := http.StatusOK

	result := PrepareCategoryClient()

	data := BasicResp{
		Msg:     msg,
		Latency: time.Since(t).String(),
		Data:    result,
	}

	log.Println(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := r.FormValue("path")
	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", f.Name())
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDetail(w http.ResponseWriter, req *http.Request) {
	t := time.Now()

	code := http.StatusOK
	itemID := req.FormValue("itemid")
	shopID := req.FormValue("shopid")
	msg := "fetch detail " + itemID + " " + shopID

	result, err := getDetail(itemID, shopID)
	if err != nil {
		msg = err.Error()
		code = http.StatusInternalServerError
	}

	data := BasicResp{
		Msg:     msg,
		Latency: time.Since(t).String(),
		Data:    result,
	}

	log.Println(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func handleEvaluateBuyer(w http.ResponseWriter, req *http.Request) {
	t := time.Now()

	code := http.StatusOK
	itemID := req.FormValue("itemid")
	shopID := req.FormValue("shopid")
	msg := "fetch buyer " + itemID + " " + shopID

	go EvaluateBuyer(itemID, shopID)

	data := BasicResp{
		Msg:     msg,
		Latency: time.Since(t).String(),
		Data:    "async",
	}

	log.Println(msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func checkClientIP(w http.ResponseWriter, httpServer *http.Request) {
	t := time.Now()
	code := http.StatusOK

	var userIP string
	if len(httpServer.Header.Get("CF-Connecting-IP")) > 1 {
		userIP = httpServer.Header.Get("CF-Connecting-IP")
	} else if len(httpServer.Header.Get("X-Forwarded-For")) > 1 {
		userIP = httpServer.Header.Get("X-Forwarded-For")
	} else if len(httpServer.Header.Get("X-Real-IP")) > 1 {
		userIP = httpServer.Header.Get("X-Real-IP")
	} else {
		userIP = httpServer.RemoteAddr
	}
	data := BasicResp{
		Msg:     "check IP",
		Latency: time.Since(t).String(),
		Data:    userIP,
	}

	log.Println("got requet", userIP)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
