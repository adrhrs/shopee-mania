package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

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

func hello(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	data := BasicResp{
		Msg:     "Hello World Production, new build v4 detail ",
		Data:    dir,
		Latency: time.Since(t).String(),
	}

	log.Println("log of this endpoint", dir)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
