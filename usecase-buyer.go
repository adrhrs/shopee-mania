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

/// v2

func workerGetReviews(id int, jobs <-chan GetReviews, results chan<- GetReviews) {
	for j := range jobs {
		t := time.Now()
		for i := 0; i < 10; i++ {
			offsetStr := strconv.Itoa(i * 6)
			reviews, err := hitRating(j.ProductID, j.ShopID, "0", offsetStr)
			j.Review = append(j.Review, reviews.Data.Ratings...)
			j.Err = err
			if len(reviews.Data.Ratings) < 6 {
				break
			}
		}
		j.Latency = time.Since(t).String()
		results <- j
	}
}

// 729373 5175509

func EvaluateProductReviewer() (err error) {

	var (
		t             = time.Now()
		uniqueProduct = make(map[string]string)
	)

	uniqueProduct, err = populateUniqueProduct()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("populated unique product", len(uniqueProduct), time.Since(t).String())
	t = time.Now()

	reviewCount, writterCounter, err := populateProductReview(uniqueProduct)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("populated product review", reviewCount, writterCounter, time.Since(t).String())

	return
}

func populateProductReview(uniqueProduct map[string]string) (reviewCount, writterCounter int, err error) {

	numJobs := len(uniqueProduct)
	jobs := make(chan GetReviews, numJobs)
	results := make(chan GetReviews, numJobs)

	file, err := os.Create("reviewer-data.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvW := csv.NewWriter(file)
	csvW.Write([]string{
		"reviewer_id", "item_id", "shop_id",
	})

	const batch = 1000

	for w := 1; w <= 50; w++ {
		go workerGetReviews(w, jobs, results)
	}

	go func() {
		for k, v := range uniqueProduct {
			jobs <- GetReviews{
				ProductID: k,
				ShopID:    v,
			}
		}
	}()

	for a := 1; a <= numJobs; a++ {
		r := <-results
		reviewCount += len(r.Review)
		if r.Err != nil {
			log.Println("err occured", r.Err)
		} else {
			for _, review := range r.Review {
				csvW.Write([]string{
					strconv.Itoa(review.AuthorShopid),
					r.ProductID,
					r.ShopID,
				})
				writterCounter++
				if writterCounter%batch == 0 {
					csvW.Flush()
				}
			}
		}
		if a%batch == 0 {
			log.Println(a, r.ProductID, r.ShopID, len(r.Review), r.Latency)
		}
	}

	if writterCounter%batch != 0 {
		csvW.Flush()
	}
	file.Close()

	close(jobs)
	close(results)

	return
}

func populateUniqueProduct() (uniqueProduct map[string]string, err error) {

	f, err := os.Open(".")
	if err != nil {
		log.Println(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
		return
	}

	files := []string{}
	for _, file := range fileInfo {
		fName := file.Name()
		if strings.Contains(fName, "aggregated") {
			files = append(files, fName)
		}
	}

	uniqueProduct = make(map[string]string) //product_id:shop_id

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

	return
}
