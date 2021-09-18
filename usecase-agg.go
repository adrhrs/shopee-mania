package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/// agg v2
func AggResultV2() (ar AggCronResult) {

	prepareStart := time.Now()
	f, err := os.Open(".")
	if err != nil {
		log.Println(err)
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Println(err)
	}

	//get available date, and map to filename
	var (
		uniqueDate = make(map[string]time.Time)
		catIDs     = make(map[string]int)
		dates      []time.Time
	)

	for _, file := range fileInfo {
		fName := file.Name()
		if strings.Contains(fName, "data") {
			var (
				fNaming         []string
				dateKey, catKey string
			)
			fNaming = strings.Split(fName, "-")
			if len(fNaming) > 8 {

				dateKey = strings.Join(fNaming[4:7], "-")
				t, _ := time.Parse(dateOnlyFormat, dateKey)
				dk := getDate(t).Format("2006-01-02")
				uniqueDate[dk] = t

				catKey = fNaming[7]
				catIDs[catKey]++

			} else {
				err = errors.New("file is corrupted " + fName)
				continue
			}
		}
	}

	for _, v := range uniqueDate {
		dates = append(dates, v)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	basicCol := []string{"item_id", "shop_id", "cat_id", "name"}
	for _, t := range dates {
		dateKey := getDate(t).Format("2006-01-02")
		basicCol = append(basicCol, "price "+dateKey)
	}
	for _, t := range dates {
		dateKey := getDate(t).Format("2006-01-02")
		basicCol = append(basicCol, "sold "+dateKey)
	}
	basicCol = append(basicCol, "availibility", "min_sold", "max_sold", "delta_sold", "avg_price", "est_gmv")
	col := [][]string{basicCol}

	var (
		totProdCatDate    int
		totProdCat        int
		invalidData       int
		readDur, writeDur int64
	)
	for catID := range catIDs {
		var (
			mapProd     = make(map[string]SimpleProduct)
			mapProdDate = make(map[string]SimpleProduct)
			row         [][]string
		)
		readStart := time.Now()
		for _, d := range dates {
			dateKey := getDate(d).Format("2006-01-02")
			filename := fmt.Sprintf("188-166-252-251-%s-%s-data.csv", dateKey, catID)
			file, errOpen := os.Open(filename)
			if errOpen != nil {
				err = errOpen
				log.Println(err)
				return
			}
			r := csv.NewReader(file)
			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if len(record) != 7 {
					invalidData++
				} else {
					productKey := fmt.Sprintf("%v-%v", record[0], catID)
					mapProd[productKey] = SimpleProduct{
						Itemid: record[0],
						Shopid: record[1],
						Catid:  catID,
						Name:   record[6],
					}
					pdKey := fmt.Sprintf("%v-%v", dateKey, productKey)
					mapProdDate[pdKey] = SimpleProduct{
						Itemid:         record[0],
						Price:          record[3],
						HistoricalSold: record[4],
					}
				}
			}
			file.Close()
		}
		readDur += time.Since(readStart).Milliseconds()
		writeStart := time.Now()
		// writing csv
		for k, v := range mapProd {
			var (
				totPrice, avgPrice, avail, minSold float64
				maxSold, diffSold, estGMV          float64
				prices, solds                      []string
				fSolds, extras                     []float64
			)

			for _, t := range dates {
				dateKey := getDate(t).Format("2006-01-02")
				pdKey := fmt.Sprintf("%v-%v", dateKey, k)
				var price, sold = "-", "-"
				if val, ok := mapProdDate[pdKey]; ok {
					price = val.Price
					sold = val.HistoricalSold

					fPrice, _ := strconv.ParseFloat(val.Price, 64)
					fSold, _ := strconv.ParseFloat(val.HistoricalSold, 64)

					fSolds = append(fSolds, fSold)
					totPrice += fPrice
					avail++
				}
				prices = append(prices, price)
				solds = append(solds, sold)
			}

			if len(fSolds) > 0 {
				sort.Slice(fSolds, func(i, j int) bool {
					return fSolds[i] < fSolds[j]
				})
				maxSold = fSolds[len(fSolds)-1]
				minSold = fSolds[0]
				diffSold = maxSold - minSold
			}

			if avail > 0 {
				avgPrice = totPrice / avail
			}
			estGMV = diffSold * avgPrice
			extras = []float64{avail, minSold, maxSold, diffSold, avgPrice, estGMV}
			row = append(row, formatCSVAGG(v, prices, solds, extras))
		}
		writeDur += time.Since(writeStart).Milliseconds()
		totProdCatDate += len(mapProdDate)
		totProdCat += len(mapProd)
		writeCSV("aggregated-"+catID, row, col)

	}

	readSec := int(float64(readDur) / 1000.0)
	writeSec := int(float64(writeDur) / 1000.0)

	ar.AggDuration = time.Since(prepareStart).String()
	ar.ReadDuration = strconv.Itoa(readSec)
	ar.CalcDuration = strconv.Itoa(writeSec)

	log.Println("dur", invalidData, totProdCatDate, ar)

	return
}
