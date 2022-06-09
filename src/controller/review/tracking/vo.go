package tracking

type BaeminCustomerReviews struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	Serverdatetime string `json:"serverDatetime"`
	Data           struct {
		Member struct {
			Nickname string `json:"nickname"`
			Grade    string `json:"grade"`
			Imageurl string `json:"imageUrl"`
		} `json:"member"`
		Reviewcount int `json:"reviewCount"`
		Reviews     []struct {
			ID                      int64   `json:"id"`
			Rating                  float64 `json:"rating"`
			Ceoonlymessage          string  `json:"ceoOnlyMessage"`
			Abusingsuspectedmessage string  `json:"abusingSuspectedMessage"`
			Blockmessage            string  `json:"blockMessage"`
			Contents                string  `json:"contents"`
			Modifiable              bool    `json:"modifiable"`
			Deletable               bool    `json:"deletable"`
			Displaytype             string  `json:"displayType"`
			Displaystatus           string  `json:"displayStatus"`
			Menudisplaytype         string  `json:"menuDisplayType"`
			Menus                   []struct {
				Menuid         int    `json:"menuId"`
				Reviewmenuid   int64  `json:"reviewMenuId"`
				Name           string `json:"name"`
				Recommendation string `json:"recommendation"`
				Contents       string `json:"contents"`
			} `json:"menus"`
			Comments []interface{} `json:"comments"`
			Images   []struct {
				ID  int64  `json:"id"`
				URL string `json:"url"`
			} `json:"images"`
			Shop struct {
				No          int    `json:"no"`
				Name        string `json:"name"`
				Servicetype string `json:"serviceType"`
			} `json:"shop"`
			Datetext string `json:"dateText"`
		} `json:"reviews"`
	} `json:"data"`
}

type BaeminCustomerInfo struct{
	NicName 	string 	`json:"nicName"`
	ReviewCount int 	`json:"reviewCount"`
	AvgRating 	float64 `json:"avgRating"`

	RecentReivews []RecentReivew `json:"recentReviews"`
}

type RecentReivew struct{
	Rating 		float64 `json:"rating"`
	Contents 	string 	`json:"contents"`
	ShopName 	string 	`json:"shopName"`
	Date		string 	`json:"date"`
}