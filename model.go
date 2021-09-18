package main

type SearchParam struct {
	Newest    int
	Keyword   string
	MinPrice  string
	Matchid   int
	SortBy    string
	CrawlType int
	PageType  string
	Limit     int
}

type workerDoCrawlReturn struct {
	Result  [][]string
	MatchID int
	Dur     string
	Err     error
	TotProd int
}

type BasicResp struct {
	Msg     string
	Data    interface{}
	Latency string
}

type CrawlCronResult struct {
	TotalCategories int
	TotalShopIDs    int
	TotalProduct    int
	AvgProductCount int
	CrawlDuration   string
	ErrCount        int
}

type AggCronResult struct {
	AggDuration  string
	ReadDuration string
	CalcDuration string
	TotalDays    int
}

type UserReview struct {
	Detail     UserDetail `json:"detail"`
	OrderID    int64      `json:"order_id"`
	ShopID     int        `json:"shop_id"`
	ItemID     int64      `json:"item_id"`
	RatingStar int        `json:"rating_star"`
	Comment    string     `json:"comment"`
}

type UserDetail struct {
	ShopID        int    `json:"shop_id"`
	Username      string `json:"username"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	City          string `json:"city"`
	District      string `json:"district"`
	Follower      int    `json:"follower"`
	Portrait      string `json:"portrait"`
	PhoneVerified bool   `json:"phone_verified"`
}
