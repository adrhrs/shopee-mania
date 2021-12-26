package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

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

func fetchResultSingle(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	msg := "fetch result single"
	code := http.StatusOK
	catid := req.FormValue("catid")
	itemid := req.FormValue("itemid")
	result, err := FetchResultSingle(catid, itemid)
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
