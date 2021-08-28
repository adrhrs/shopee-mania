package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func buildSearchURL(p SearchParam) (URL, param string) {

	newest := strconv.Itoa(p.Newest)
	// keyword := url.QueryEscape(p.Data.Keyword)

	matchID := ""
	keywordParam := ""
	pageTypeParam := ""
	priceParam := ""
	pageTypeParam = "&page_type=search&version=2"
	sortBy := "by=sales"

	matchID = strconv.Itoa(p.Matchid)
	URL = "https://shopee.co.id/api/v4/search/search_items?" + sortBy + "&limit=60&order=desc&newest=" +
		newest + "&match_id=" + matchID + keywordParam + priceParam + pageTypeParam + "&scenario=PAGE_OTHERS"

	param = sortBy + "&limit=60&order=desc&newest=" +
		newest + matchID + keywordParam + priceParam + pageTypeParam

	return
}

func doSearchCrawl(p SearchParam) (resp ShpResV2, dur string, err error) {

	t := time.Now()
	url, par := buildSearchURL(p)
	req, err := http.NewRequest("GET", url, nil)
	req = addReqHeader(req, generateIfNoneMatch(par))
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{}

	res, err := client.Do(req)
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println(err)
	}

	dur = time.Since(t).String()
	defer res.Body.Close()

	return
}

func doGetCategories() (resp CatResp, err error) {
	url := "https://shopee.co.id/api/v1/category_list/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "SPC_IA=-1; SPC_EC=-; SPC_F=HIv65zB72GPw8zwivK12luZIl6ehh1jx; REC_T_ID=e8f38444-5609-11ea-93cd-ccbbfe5d5cda; SPC_U=-; SPC_R_T_ID=\"w8nBNM3NDZ7/Tv5q9Otd+w4vHuv0iOpo9Y9hNnrcwUQUtLGsOZmEswbrbbq3yNQr4cNQR5prO2+5jEGVYeRtWj3Gaa5sOqJGkhzffDYxaRw=\"; SPC_T_ID=\"w8nBNM3NDZ7/Tv5q9Otd+w4vHuv0iOpo9Y9hNnrcwUQUtLGsOZmEswbrbbq3yNQr4cNQR5prO2+5jEGVYeRtWj3Gaa5sOqJGkhzffDYxaRw=\"; SPC_SI=mall.EY1oB2HNr2gFAqJeF35kLABpd4YEPHPv; SPC_R_T_IV=\"h5ydLCiRp02IRGRK4M9CWQ==\"; SPC_T_IV=\"h5ydLCiRp02IRGRK4M9CWQ==\"")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println(err)
	}
	return
}

func fixProduct(pr ItemBasicNecessary) (fixed ItemBasicNecessary) {
	fixed = pr
	fixed.Price = pr.Price / 100000
	return
}

func hitDetailAPI(itemID, shopID string) (resp Detail, err error) {
	url := "https://shopee.co.id/api/v2/item/get?itemid=" + itemID + "&shopid=" + shopID
	par := "itemid=" + itemID + "&shopid=" + shopID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	req = addReqHeader(req, generateIfNoneMatch(par))
	client := &http.Client{}

	res, err := client.Do(req)
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	return
}

func hitRating(itemID, shopID, typeRating, offset string) (resp RatingResponse, err error) {
	url := "https://shopee.co.id/api/v2/item/get_ratings?filter=0&flag=1&limit=6&offset=" + offset + "&type=" + typeRating + "&itemid=" + itemID + "&shopid=" + shopID
	par := "itemid=" + itemID + "&shopid=" + shopID + "&offset=" + offset + "&filter=0&flag=1&limit=6" + "&type=" + typeRating

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	req = addReqHeader(req, generateIfNoneMatch(par))
	client := &http.Client{}

	res, err := client.Do(req)
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	return
}

func hitShopInfo(shopID, username string, isAnonymous bool) (resp ShopInfo, err error) {
	var url, par string

	if isAnonymous {
		url = "https://shopee.co.id/api/v4/shop/get_shop_detail?shopid=" + shopID
		par = "shopid=" + shopID
	} else {
		url = "https://shopee.co.id/api/v4/shop/get_shop_detail?username=" + username
		par = "username=" + username
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	req = addReqHeader(req, generateIfNoneMatch(par))
	client := &http.Client{}

	res, err := client.Do(req)
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()

	return
}
