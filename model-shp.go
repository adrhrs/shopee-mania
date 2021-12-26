package main

//search API
type ShpResV2 struct {
	TotalCount int `json:"total_count"`
	Items      []struct {
		ItemBasic ItemBasicNecessary `json:"item_basic"`
	} `json:"items"`
}

type ItemBasicNecessary struct {
	Itemid              int    `json:"itemid"`
	Shopid              int    `json:"shopid"`
	Name                string `json:"name"`
	Stock               int    `json:"stock"`
	HistoricalSold      int    `json:"historical_sold"`
	LikedCount          int    `json:"liked_count"`
	ViewCount           int    `json:"view_count"`
	Catid               int    `json:"catid"`
	Brand               string `json:"brand"`
	Price               int    `json:"price"`
	PriceBeforeDiscount int    `json:"price_before_discount"`
	RawDiscount         int    `json:"raw_discount"`
	ItemRating          Rating `json:"item_rating"`
	IsOfficialShop      bool   `json:"is_official_shop"`
}

type SimpleProduct struct {
	Itemid         string `json:"itemid"`
	Shopid         string `json:"shopid"`
	Catid          string `json:"catid"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	HistoricalSold string `json:"historical_sold"`
}

type FetchProductResponse struct {
	Total   int             `json:"total"`
	Product []CastedProduct `json:"products"`
	Dates   []string        `json:"dates"`
}

type CastedProduct struct {
	Itemid int64   `json:"item_id"`
	Shopid int64   `json:"shop_id"`
	Catid  int64   `json:"cat_id"`
	Name   string  `json:"name"`
	Prices []int64 `json:"prices"`
	Solds  []int64 `json:"solds"`

	Avail     int64 `json:"availibility"`
	MinSold   int64 `json:"min_sold"`
	MaxSold   int64 `json:"max_sold"`
	DeltaSold int64 `json:"delta_sold"`
	AvgPrice  int64 `json:"avg_price"`
	EstGMV    int64 `json:"est_gmv"`
}

type Rating struct {
	RatingStar        float64 `json:"rating_star"`
	RatingCount       []int   `json:"rating_count"`
	RcountWithContext int     `json:"rcount_with_context"`
	RcountWithImage   int     `json:"rcount_with_image"`
}

type ItemBasicAll struct {
	Itemid                  int         `json:"itemid"`
	Shopid                  int         `json:"shopid"`
	Name                    string      `json:"name"`
	LabelIds                []int       `json:"label_ids"`
	Image                   string      `json:"image"`
	Images                  []string    `json:"images"`
	Currency                string      `json:"currency"`
	Stock                   int         `json:"stock"`
	Status                  int         `json:"status"`
	Ctime                   int         `json:"ctime"`
	Sold                    int         `json:"sold"`
	HistoricalSold          int         `json:"historical_sold"`
	Liked                   bool        `json:"liked"`
	LikedCount              int         `json:"liked_count"`
	ViewCount               int         `json:"view_count"`
	Catid                   int         `json:"catid"`
	Brand                   string      `json:"brand"`
	CmtCount                int         `json:"cmt_count"`
	Flag                    int         `json:"flag"`
	CbOption                int         `json:"cb_option"`
	ItemStatus              string      `json:"item_status"`
	Price                   int         `json:"price"`
	PriceMin                int         `json:"price_min"`
	PriceMax                int         `json:"price_max"`
	PriceMinBeforeDiscount  int         `json:"price_min_before_discount"`
	PriceMaxBeforeDiscount  int         `json:"price_max_before_discount"`
	HiddenPriceDisplay      interface{} `json:"hidden_price_display"`
	PriceBeforeDiscount     int         `json:"price_before_discount"`
	HasLowestPriceGuarantee bool        `json:"has_lowest_price_guarantee"`
	ShowDiscount            int         `json:"show_discount"`
	RawDiscount             int         `json:"raw_discount"`
	Discount                interface{} `json:"discount"`
	IsCategoryFailed        bool        `json:"is_category_failed"`
	SizeChart               interface{} `json:"size_chart"`
	VideoInfoList           []struct {
		VideoID  string `json:"video_id"`
		ThumbURL string `json:"thumb_url"`
		Duration int    `json:"duration"`
		Version  int    `json:"version"`
		Vid      string `json:"vid"`
		Formats  []struct {
			Format  int    `json:"format"`
			Defn    string `json:"defn"`
			Profile string `json:"profile"`
			Path    string `json:"path"`
			URL     string `json:"url"`
			Width   int    `json:"width"`
			Height  int    `json:"height"`
		} `json:"formats"`
		DefaultFormat struct {
			Format  int    `json:"format"`
			Defn    string `json:"defn"`
			Profile string `json:"profile"`
			Path    string `json:"path"`
			URL     string `json:"url"`
			Width   int    `json:"width"`
			Height  int    `json:"height"`
		} `json:"default_format"`
	} `json:"video_info_list"`
	TierVariations []struct {
		Name       string        `json:"name"`
		Options    []string      `json:"options"`
		Images     []string      `json:"images"`
		Properties []interface{} `json:"properties"`
		Type       int           `json:"type"`
	} `json:"tier_variations"`
	ItemRating struct {
		RatingStar        float64 `json:"rating_star"`
		RatingCount       []int   `json:"rating_count"`
		RcountWithContext int     `json:"rcount_with_context"`
		RcountWithImage   int     `json:"rcount_with_image"`
	} `json:"item_rating"`
	ItemType                          int         `json:"item_type"`
	ReferenceItemID                   string      `json:"reference_item_id"`
	TransparentBackgroundImage        string      `json:"transparent_background_image"`
	IsAdult                           bool        `json:"is_adult"`
	BadgeIconType                     int         `json:"badge_icon_type"`
	ShopeeVerified                    bool        `json:"shopee_verified"`
	IsOfficialShop                    bool        `json:"is_official_shop"`
	ShowOfficialShopLabel             bool        `json:"show_official_shop_label"`
	ShowShopeeVerifiedLabel           bool        `json:"show_shopee_verified_label"`
	ShowOfficialShopLabelInTitle      bool        `json:"show_official_shop_label_in_title"`
	IsCcInstallmentPaymentEligible    bool        `json:"is_cc_installment_payment_eligible"`
	IsNonCcInstallmentPaymentEligible bool        `json:"is_non_cc_installment_payment_eligible"`
	CoinEarnLabel                     interface{} `json:"coin_earn_label"`
	ShowFreeShipping                  bool        `json:"show_free_shipping"`
	PreviewInfo                       interface{} `json:"preview_info"`
	CoinInfo                          interface{} `json:"coin_info"`
	ExclusivePriceInfo                interface{} `json:"exclusive_price_info"`
	BundleDealID                      int         `json:"bundle_deal_id"`
	CanUseBundleDeal                  bool        `json:"can_use_bundle_deal"`
	BundleDealInfo                    interface{} `json:"bundle_deal_info"`
	IsGroupBuyItem                    interface{} `json:"is_group_buy_item"`
	HasGroupBuyStock                  interface{} `json:"has_group_buy_stock"`
	GroupBuyInfo                      interface{} `json:"group_buy_info"`
	WelcomePackageType                int         `json:"welcome_package_type"`
	WelcomePackageInfo                interface{} `json:"welcome_package_info"`
	AddOnDealInfo                     interface{} `json:"add_on_deal_info"`
	CanUseWholesale                   bool        `json:"can_use_wholesale"`
	IsPreferredPlusSeller             bool        `json:"is_preferred_plus_seller"`
	ShopLocation                      string      `json:"shop_location"`
	HasModelWithAvailableShopeeStock  bool        `json:"has_model_with_available_shopee_stock"`
	VoucherInfo                       interface{} `json:"voucher_info"`
	CanUseCod                         bool        `json:"can_use_cod"`
	IsOnFlashSale                     bool        `json:"is_on_flash_sale"`
	SplInstallmentTenure              interface{} `json:"spl_installment_tenure"`
	IsLiveStreamingPrice              interface{} `json:"is_live_streaming_price"`
	IsMart                            bool        `json:"is_mart"`
	PackSize                          interface{} `json:"pack_size"`
}

//search API finish

//category API

type CatResp []struct {
	Main struct {
		DisplayName        string      `json:"display_name"`
		Name               string      `json:"name"`
		Catid              int         `json:"catid"`
		Image              string      `json:"image"`
		ParentCategory     int         `json:"parent_category"`
		IsAdult            int         `json:"is_adult"`
		BlockBuyerPlatform interface{} `json:"block_buyer_platform"`
		SortWeight         int         `json:"sort_weight"`
	} `json:"main"`
	Sub []struct {
		DisplayName        string      `json:"display_name"`
		Name               string      `json:"name"`
		Catid              int         `json:"catid"`
		Image              string      `json:"image"`
		ParentCategory     int         `json:"parent_category"`
		IsAdult            int         `json:"is_adult"`
		BlockBuyerPlatform interface{} `json:"block_buyer_platform"`
		SortWeight         int         `json:"sort_weight"`
		SubSub             []struct {
			Image              string      `json:"image"`
			BlockBuyerPlatform interface{} `json:"block_buyer_platform"`
			DisplayName        string      `json:"display_name"`
			Name               string      `json:"name"`
			Catid              int         `json:"catid"`
		} `json:"sub_sub"`
	} `json:"sub"`
}

//category API finish

//rating API
type RatingResponse struct {
	Data struct {
		Ratings           []RatingDetail `json:"ratings"`
		ItemRatingSummary struct {
			RatingTotal       int   `json:"rating_total"`
			RatingCount       []int `json:"rating_count"`
			RcountWithContext int   `json:"rcount_with_context"`
			RcountWithImage   int   `json:"rcount_with_image"`
			RcountWithMedia   int   `json:"rcount_with_media"`
			RcountLocalReview int   `json:"rcount_local_review"`
		} `json:"item_rating_summary"`
	} `json:"data"`
	Error    interface{} `json:"error"`
	ErrorMsg interface{} `json:"error_msg"`
}

type RatingDetail struct {
	Orderid        int64    `json:"orderid"`
	Itemid         int64    `json:"itemid"`
	Rating         int      `json:"rating"`
	Userid         int      `json:"userid"`
	Shopid         int      `json:"shopid"`
	Ctime          int64    `json:"ctime"`
	Comment        string   `json:"comment"`
	RatingStar     int      `json:"rating_star"`
	AuthorUsername string   `json:"author_username"`
	AuthorPortrait string   `json:"author_portrait"`
	AuthorShopid   int      `json:"author_shopid"`
	Anonymous      bool     `json:"anonymous"`
	Images         []string `json:"images"`
	ProductItems   []struct {
		Itemid    int64  `json:"itemid"`
		Shopid    int    `json:"shopid"`
		Name      string `json:"name"`
		Modelid   int64  `json:"modelid"`
		ModelName string `json:"model_name"`
	} `json:"product_items"`
	LikeCount interface{} `json:"like_count"`
}

//rating API finish

//shop info API
type ShopInfo struct {
	Error    int      `json:"error"`
	ErrorMsg string   `json:"error_msg"`
	Data     ShopData `json:"data"`
	Version  string   `json:"version"`
}

type ShopData struct {
	UserID        int    `json:"userid"`
	ShopID        int    `json:"shopid"`
	Name          string `json:"name"`
	ShopLocation  string `json:"shop_location"`
	ItemCount     int    `json:"item_count"`
	FollowerCount int    `json:"follower_count"`
	Account       struct {
		Username       string `json:"username"`
		FollowingCount int    `json:"following_count"`
		Portrait       string `json:"portrait"`
		PhoneVerified  bool   `json:"phone_verified"`
		EmailVerified  bool   `json:"email_verified"`
	} `json:"account"`

	Ctime          int64 `json:"ctime"`
	LastActiveTime int64 `json:"last_active_time"`
}

//shop info API finish

//pdp v4
type ItemDataV4 struct {
	ItemID              int64       `json:"itemid"`
	ShopID              int         `json:"shopid"`
	UserID              int         `json:"userid"`
	PriceBeforeDiscount int64       `json:"price_before_discount"`
	Price               int         `json:"price"`
	Stock               int         `json:"stock"`
	DiscountStock       int         `json:"discount_stock"`
	NormalStock         int         `json:"normal_stock"`
	Discount            string      `json:"discount"`
	HistoricalSold      int         `json:"historical_sold"`
	Sold                int         `json:"sold"`
	Name                string      `json:"name"`
	Ctime               int         `json:"ctime"`
	Condition           int         `json:"condition"`
	Catid               int         `json:"catid"`
	Description         string      `json:"description"`
	Brand               interface{} `json:"brand"`
	ItemRating          struct {
		RatingStar  float64 `json:"rating_star"`
		RatingCount []int   `json:"rating_count"`
	} `json:"item_rating"`
	LikedCount     int      `json:"liked_count"`
	Images         []string `json:"images"`
	Image          string   `json:"image"`
	IsOfficialShop bool     `json:"is_official_shop"`
	ShopLocation   string   `json:"shop_location"`
	Categories     []struct {
		Catid           int    `json:"catid"`
		DisplayName     string `json:"display_name"`
		NoSub           bool   `json:"no_sub"`
		IsDefaultSubcat bool   `json:"is_default_subcat"`
	} `json:"categories"`
}

type RespDetailV4 struct {
	Data ItemDataV4 `json:"data"`
}

//pdp v4 finish
