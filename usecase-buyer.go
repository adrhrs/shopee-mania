package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultRatingType = "0"
	layout            = "2006-01-02 15:04:05"
)

func TrackProduct(itemID, shopID string) (reviewCount, userCount int, data []Buyer, err error) {

	filename := generateFilename(itemID, shopID)
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		reviewCount, userCount, err = doTrackProduct(itemID, shopID)
		if err != nil {
			return
		}
		data, err = fetchBuyerInfo(filename)
		if err != nil {
			return
		}
	} else if err != nil {
		log.Println(err)
	} else {
		data, err = fetchBuyerInfo(filename)
		if err != nil {
			return
		}
	}

	return
}

func fetchBuyerInfo(filename string) (data []Buyer, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}

	var (
		id int
		r  = csv.NewReader(file)
	)

	for {
		id++
		record, errRecord := r.Read()
		if errRecord == io.EOF {
			break
		}
		if id == 1 || len(record) < 8 {
			continue
		}
		unixBuyTime, errParse := strconv.ParseInt(record[1], 10, 64)
		if errParse != nil {
			err = errParse
			return
		}
		buyTimeFormat := time.Unix(unixBuyTime, 0).Format(layout)
		unixCreateTime, errParse := strconv.ParseInt(record[7], 10, 64)
		if errParse != nil {
			err = errParse
			return
		}
		createTimeFormat := time.Unix(unixCreateTime, 0).Format(layout)
		userShopID, errParse := strconv.ParseInt(record[0], 10, 64)
		if errParse != nil {
			err = errParse
			return
		}
		followingCount, errParse := strconv.Atoi(record[5])
		if errParse != nil {
			err = errParse
			return
		}
		accAge := int(time.Since(time.Unix(unixCreateTime, 0)).Hours() / 24)

		data = append(data, Buyer{
			UserShopID:              userShopID,
			Username:                record[2],
			Name:                    record[3],
			Location:                record[4],
			FollowingCount:          followingCount,
			Portrait:                record[6],
			BuyingTimeFormatted:     buyTimeFormat,
			UserCreateTimeFormatted: createTimeFormat,
			AccountAge:              accAge,
		})
	}
	file.Close()
	return
}

func doTrackProduct(itemID, shopID string) (reviewCount, userCount int, err error) {

	var (
		t = time.Now()
	)

	log.Printf("start track product %s-%s", itemID, shopID)
	reviewCount, buyerData, err := populateProductReview(itemID, shopID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("finish get review %s, got %v reviews, %v users on %s", itemID, reviewCount, len(buyerData), time.Since(t).String())
	userCount, err = populateProductBuyers(itemID, shopID, buyerData)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("finish get user data %s, got %v user data on %s", itemID, userCount, time.Since(t).String())

	return
}

func populateProductReview(itemID, shopID string) (reviewCount int, buyerData map[int64]int64, err error) {

	var (
		defaultReviewPage      = 100
		defaultReviewCount     = 600
		defaultLimitReview     = 6
		defaultGetReviewWorker = 10
		buyerDataTemp          = make(map[int64]int64)
	)

	//get total review first
	reviewData, err := hitRating(itemID, shopID, defaultRatingType, "0")
	if err != nil {
		return
	}
	availableReview := reviewData.Data.ItemRatingSummary.RatingTotal
	if availableReview < defaultReviewCount {
		defaultReviewPage = availableReview / defaultLimitReview
	}

	var (
		numJobs = defaultReviewPage
		jobs    = make(chan GetReviews, numJobs)
		results = make(chan GetReviews, numJobs)
	)

	for w := 1; w <= defaultGetReviewWorker; w++ {
		go workerGetReviews(w, jobs, results)
	}

	go func() {
		for i := 0; i < defaultReviewPage; i++ {
			jobs <- GetReviews{
				ProductID: itemID,
				ShopID:    shopID,
				Offset:    i,
			}
		}
	}()

	for a := 1; a <= numJobs; a++ {
		r := <-results
		reviewCount += len(r.Reviews)
		if r.Err != nil {
			log.Println("err occured", r.Err)
		} else {
			for _, review := range r.Reviews {
				buyerDataTemp[int64(review.AuthorShopid)] = review.Ctime
			}
		}
	}

	close(jobs)
	close(results)

	buyerData = buyerDataTemp

	return
}

func workerGetReviews(id int, jobs <-chan GetReviews, results chan<- GetReviews) {
	for j := range jobs {
		t := time.Now()
		offsetStr := strconv.Itoa(j.Offset * 6)
		reviews, err := hitRating(j.ProductID, j.ShopID, defaultRatingType, offsetStr)
		j.Reviews = append(j.Reviews, reviews.Data.Ratings...)
		j.Err = err
		j.Latency = time.Since(t).String()
		results <- j
	}
}

func populateProductBuyers(itemID, shopID string, buyerData map[int64]int64) (writtenCount int, err error) {

	var (
		defaultGetUser = 20
		numJobs        = len(buyerData)
		jobs           = make(chan GetUser, numJobs)
		results        = make(chan GetUser, numJobs)
		batch          = 100
	)

	file, err := os.Create(generateFilename(itemID, shopID))
	if err != nil {
		log.Println(err)
		return
	}
	csvW := csv.NewWriter(file)
	csvW.Write([]string{
		"user_shop_id", "buying_time", "user_name",
		"name", "location", "following_count", "portrait",
		"user_create_time",
	})

	for w := 1; w <= defaultGetUser; w++ {
		go workerGetUser(w, jobs, results)
	}

	go func() {
		for userShopID, reviewCtime := range buyerData {
			jobs <- GetUser{
				UserShopID:  strconv.FormatInt(userShopID, 10),
				ReviewCtime: reviewCtime,
			}
		}
	}()

	for a := 1; a <= numJobs; a++ {
		r := <-results
		if r.Err != nil {
			log.Println("err occured", r.Err)
		} else if r.UserData.ShopID > 0 {
			writtenCount++
			csvW.Write([]string{
				fmt.Sprintf("%v", r.UserShopID),
				fmt.Sprintf("%v", r.ReviewCtime),
				fmt.Sprintf("%v", r.UserData.Account.Username),
				fmt.Sprintf("%v", r.UserData.Name),
				fmt.Sprintf("%v", r.UserData.ShopLocation),
				fmt.Sprintf("%v", r.UserData.Account.FollowingCount),
				fmt.Sprintf("%v", r.UserData.Account.Portrait),
				fmt.Sprintf("%v", r.UserData.Ctime),
			})
			if writtenCount%batch == 0 {
				csvW.Flush()
			}
		}
	}

	if writtenCount%batch == 0 {
		csvW.Flush()
	}

	close(jobs)
	close(results)
	return
}

func workerGetUser(id int, jobs <-chan GetUser, results chan<- GetUser) {
	for j := range jobs {
		t := time.Now()
		userData, err := hitShopInfo(j.UserShopID, "", true)
		j.UserData = userData.Data
		j.Err = err
		j.Latency = time.Since(t).String()
		results <- j
	}
}

func generateFilename(itemID, shopID string) string {
	return fmt.Sprintf("buyer-data-%s-%s.csv", itemID, shopID)
}
