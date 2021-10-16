package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	colUser = [][]string{{
		"user_id", "name", "username", "follower", "portrait", "phone", "phone_verified", "address", "district", "city",
		"item_id", "shop_id", "order_id", "rating", "comment", "create_time",
	}}
)

func EvaluateBuyer(itemID, shopID string) (totalReviewer int, err error) {

	allRatings := []RatingDetail{}
	userReviews := [][]string{}
	for st := 1; st <= 5; st++ {
		for i := 0; i < 5; i++ {
			ratings, err := hitRating(itemID, shopID, strconv.Itoa(st), strconv.Itoa(i*6))
			if err != nil {
				log.Println(err)
			}

			for _, r := range ratings.Data.Ratings {
				allRatings = append(allRatings, r)
				totalReviewer++
			}

			if len(ratings.Data.Ratings) < 6 {
				break
			}
		}
		fmt.Println(st, len(allRatings))
	}

	numJobs := len(allRatings)
	jobs := make(chan RatingDetail, numJobs)
	results := make(chan GetUserWorkerReturn, numJobs)

	for w := 1; w <= 10; w++ {
		go workerGetUser(w, jobs, results)
	}

	for _, j := range allRatings {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		r := <-results
		addr := r.ShopData.Address
		acc := r.ShopData.Account
		ctime := time.Unix(r.RatingDetail.Ctime, 0).Format("2006-01-02 15:04:05")

		row := []string{
			fmt.Sprintf("%v", r.ShopData.Shopid),
			addr.Name,
			acc.Username,
			fmt.Sprintf("%v", r.ShopData.FollowerCount),
			fmt.Sprintf("https://cf.shopee.co.id/file/%v", acc.Portrait),
			addr.Phone,
			fmt.Sprintf("%v", acc.PhoneVerified),
			addr.Address,
			addr.District,
			addr.City,
			fmt.Sprintf("%v", r.RatingDetail.Itemid),
			fmt.Sprintf("%v", r.RatingDetail.Shopid),
			fmt.Sprintf("%v", r.RatingDetail.Orderid),
			fmt.Sprintf("%v", r.RatingDetail.RatingStar),
			r.RatingDetail.Comment,
			ctime,
		}
		userReviews = append(userReviews, row)

		fmt.Println(acc.Username, r.ShopData.FollowerCount)
	}

	writeCSV("user", userReviews, colUser)

	return
}

type GetUserWorkerReturn struct {
	RatingDetail RatingDetail
	ShopData     ShopData
}

func workerGetUser(id int, jobs <-chan RatingDetail, results chan<- GetUserWorkerReturn) {
	for j := range jobs {
		userData, errUserData := hitShopInfo(strconv.Itoa(j.AuthorShopid), j.AuthorUsername, j.Anonymous)
		if errUserData != nil {
			log.Println(errUserData)
		}
		results <- GetUserWorkerReturn{
			RatingDetail: j,
			ShopData:     userData.Data,
		}
	}
}

func EvaluateProductReviewer() (err error) {

	t := time.Now()
	f, err := os.Open(".")
	if err != nil {
		log.Println(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}

	files := []string{}
	for _, file := range fileInfo {
		fName := file.Name()
		if strings.Contains(fName, "aggregated") {
			files = append(files, fName)
		}
	}

	uniqueProduct := make(map[string]string) //product_id:shop_id

	for _, filename := range files {
		file, errOpenFile := os.Open(filename)
		if errOpenFile != nil {
			return
		}
		r := csv.NewReader(file)
		for i := 0; ; i++ {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			uniqueProduct[record[0]] = record[1]
		}
		file.Close()
	}

	log.Println("populated unique product", len(uniqueProduct), time.Since(t).String())

	return
}
