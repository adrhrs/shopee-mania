package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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

func handleEvaluateReviewer(w http.ResponseWriter, req *http.Request) {
	t := time.Now()

	code := http.StatusOK

	go EvaluateProductReviewer()

	data := BasicResp{
		Msg:     "will evaluate all agg files",
		Latency: time.Since(t).String(),
		Data:    "async",
	}

	log.Println("evaluate agg files triggered")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
