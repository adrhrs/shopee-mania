package main

import (
	"fmt"
	"log"
	"strconv"
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
