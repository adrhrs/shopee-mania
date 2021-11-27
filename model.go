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

type GetReviews struct {
	ProductID string
	ShopID    string
	Offset    int

	Reviews []RatingDetail
	Latency string
	Err     error
}

type GetUser struct {
	UserShopID  string
	ReviewCtime int64

	UserData ShopData
	Latency  string
	Err      error
}

type Buyer struct {
	UserShopID     int64  `json:"user_shop_id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	Location       string `json:"location"`
	FollowingCount int    `json:"following_count"`
	Portrait       string `json:"portrait"`

	BuyingTimeFormatted     string `json:"buying_time"`
	UserCreateTimeFormatted string `json:"user_create_time"`

	AccountAge int `json:"account_age"`
}
