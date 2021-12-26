package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	attrSize = 10
	histSize = 2
)

func GetResult() (files []string) {

	f, err := os.Open(".")
	if err != nil {
		log.Println(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return

}

func FetchResult(catID string) (results FetchProductResponse, err error) {

	filename := fmt.Sprintf("aggregated-%v.csv", catID)
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	r := csv.NewReader(file)

	var (
		colSize, dateSize int
		cols              []string
		prods             []CastedProduct
	)

	for i := 0; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		cp := CastedProduct{}

		if i > 0 && len(record) == colSize {
			itemID, errConv := strconv.ParseInt(record[0], 10, 64)
			if errConv != nil {
				continue
			}
			shopID, errConv := strconv.ParseInt(record[1], 10, 64)
			if errConv != nil {
				continue
			}
			catID, errConv := strconv.ParseInt(record[2], 10, 64)
			if errConv != nil {
				continue
			}

			cp = CastedProduct{
				Itemid: itemID,
				Shopid: shopID,
				Catid:  catID,
				Name:   record[3],
			}

			for i := 4; i < 4+dateSize; i++ {
				var (
					iPrice, iSold int64
				)
				iPrice, errConv := strconv.ParseInt(record[i], 10, 64)
				if errConv != nil {
					iPrice = -1
				}
				iSold, errConv = strconv.ParseInt(record[i+dateSize], 10, 64)
				if errConv != nil {
					iSold = -1
				}
				cp.Prices = append(cp.Prices, iPrice)
				cp.Solds = append(cp.Solds, iSold)
			}

			avail, errConv := strconv.ParseInt(record[colSize-6], 10, 64)
			if errConv != nil {
				log.Println(errConv)
				avail = -1
			}
			minSold, errConv := strconv.ParseInt(record[colSize-5], 10, 64)
			if errConv != nil {
				minSold = -1
				log.Println(errConv)
			}
			maxSold, errConv := strconv.ParseInt(record[colSize-4], 10, 64)
			if errConv != nil {
				maxSold = -1
				log.Println(errConv)
			}
			deltaSold, errConv := strconv.ParseInt(record[colSize-3], 10, 64)
			if errConv != nil {
				deltaSold = -1
				log.Println(errConv)
			}
			avgPrice, errConv := strconv.ParseInt(record[colSize-2], 10, 64)
			if errConv != nil {
				avgPrice = -1
				log.Println(errConv)
			}
			estGMV, errConv := strconv.ParseInt(record[colSize-1], 10, 64)
			if errConv != nil {
				estGMV = -1
				log.Println(errConv)
			}

			cp.Avail = avail
			cp.MinSold = minSold
			cp.MaxSold = maxSold
			cp.DeltaSold = deltaSold
			cp.AvgPrice = avgPrice
			cp.EstGMV = estGMV

			prods = append(prods, cp)
		} else if i == 0 {
			colSize = len(record)
			dateSize = (colSize - attrSize) / histSize
			cols = record
		}

	}

	results.Product = prods
	results.Total = len(prods)

	dates := cols[4:(4 + (dateSize*histSize)/2)]
	for _, d := range dates {
		results.Dates = append(results.Dates, strings.Replace(d, "price ", "", -1))
	}

	return
}

type Category struct {
	Label string `json:"label"`
	ID    int    `json:"cat_id"`
	Lv    int    `json:"lv"`
}

func PrepareCategoryClient() (resp []Category) {

	var (
		lv1IDs            = make(map[int]int)
		whitelistedLv1IDs = []int{100017, 100630, 100010}
	)

	for _, w := range whitelistedLv1IDs {
		lv1IDs[w] = 1
	}

	cats, err := doGetCategories()
	if err != nil {
		log.Println(err)
		return
	}
	for _, lv1 := range cats {
		if _, ok := lv1IDs[lv1.Main.Catid]; ok {
			resp = append(resp, Category{
				ID:    lv1.Main.Catid,
				Label: lv1.Main.Name,
				Lv:    1,
			})
			for _, lv2 := range lv1.Sub {
				resp = append(resp, Category{
					ID:    lv2.Catid,
					Label: lv1.Main.Name + " - " + lv2.Name,
					Lv:    2,
				})
				for _, lv3 := range lv2.SubSub {
					resp = append(resp, Category{
						ID:    lv3.Catid,
						Label: lv1.Main.Name + " - " + lv2.Name + " - " + lv3.Name,
						Lv:    3,
					})
				}
			}
		}
	}

	return
}

func getDetail(itemID, shopID string) (resp RespDetailV4, err error) {
	resp, err = hitDetailAPI(itemID, shopID)
	return
}

func FetchResultSingle(catID, itemID string) (results FetchProductResponse, err error) {
	catResult, err := FetchResult(catID)
	if err != nil {
		return
	}

	for _, cp := range catResult.Product {
		strItemID := strconv.FormatInt(cp.Itemid, 10)
		if strItemID == itemID {
			results.Product = append(results.Product, cp)
			results.Total++
			break
		}
	}

	results.Dates = catResult.Dates
	return
}
