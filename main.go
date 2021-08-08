package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type BasicResp struct {
	Msg  string
	Data interface{}
}

func hello(w http.ResponseWriter, req *http.Request) {

	data := BasicResp{
		Msg:  "Hello World",
		Data: []int{1, 2, 3},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func main() {

	http.HandleFunc("/hello", hello)

	fmt.Printf("Starting server at port 6001\n")
	if err := http.ListenAndServe("127.0.0.1:6001", nil); err != nil {
		log.Fatal(err)
	}

}
