package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func handleTrackProduct(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	code := http.StatusOK
	msg := "Success Track Buyer Product"

	itemID := req.FormValue("itemid")
	shopID := req.FormValue("shopid")
	_, _, dataBuyer, err := TrackProduct(itemID, shopID)
	if err != nil {
		msg = err.Error()
		code = http.StatusInternalServerError
	}

	data := BasicResp{
		Msg:     msg,
		Latency: time.Since(t).String(),
		Data:    dataBuyer,
	}

	log.Println("got track product buyer request")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
