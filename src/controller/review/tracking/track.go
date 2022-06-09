package tracking

import (
	// login 및 기본
	"cashApi/src/controller/cls"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func CrossTest(c echo.Context) error {
	// 172.30.1.237

	/*
	SELECT a.BIZ_NUM, a.REST_ID, b.USER_ID, a.COMP_NM
	FROM cc_comp_inf a
	LEFT JOIN priv_rest_user_info b ON a.REST_ID = b.REST_ID
	WHERE b.REST_AUTH = 0 AND a.biz_num IN(
		SELECT biz_num
		FROM cc_sync_inf
			WHERE bs_dt = '20210822'
			AND err_cd = '0000'
			AND biz_num IN(
				SELECT b.BIZ_NUM
				FROM priv_rest_info a INNER JOIN cc_comp_inf b ON a.REST_ID = b.REST_ID
				AND 1=1
	)) ORDER BY b.REG_DATE;
	 */

	return c.HTML(http.StatusFound, "http://www.naver.com")
}

func GetCustomerInfo(c echo.Context) error {

	m := make(map[string]interface{})

	params := cls.GetParamJsonMap(c)
	rst, b := GetBaeminCompInfo(params["customerId"], "50")
	if rst < 0{
		m["resultCode"] = "99"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = ""

		return c.JSON(http.StatusOK, m)
	}

	var totalCnt, totalAvg float64
	var w BaeminCustomerInfo

	for _,v := range b.Data.Reviews{
		totalAvg += v.Rating
		totalCnt ++

		if totalCnt <= 5{
			var r RecentReivew
			r.Rating = v.Rating
			r.Contents = v.Contents
			r.ShopName = v.Shop.Name
			r.Date = krToDate(v.Datetext)

			w.RecentReivews = append(w.RecentReivews, r)
		}
	}

	w.NicName = b.Data.Member.Nickname
	w.ReviewCount = b.Data.Reviewcount
	w.AvgRating = totalAvg/totalCnt

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = w

	return c.JSON(http.StatusOK, m)
}

func GetBaeminCompInfo(customerId, limit string) (int, BaeminCustomerReviews){
	var b BaeminCustomerReviews

	url := fmt.Sprintf("v1/members/%s/reviews?offset=0&limit=%s&appver=10.27.1&carrier=45008&site=7jWXRELC2e&deviceModel=SM-G906S&dvcid=OPUD5687e4685d245b9c&adid=NONE&sessionId=f0a64c6f0b15c3f27f6e38&osver=23&oscd=2", customerId, limit)

	resp, err := cls.HttpRequestDetail("HTTPS", "GET", "review-api.baemin.com", "443", url, nil, nil, "", false)
	if err != nil{
		lprintf(1, "[ERROR] baemin get comp addr err(%s)\n", err.Error())
		return -1, b
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		lprintf(1, "[ERROR] baemin get comp addr body read err(%s)\n", err.Error())
		return -1, b
	}

	err = json.Unmarshal(data, &b)
	if err != nil{
		lprintf(1, "[ERROR] baemin get comp addr parsing json err(%s)\n", err.Error())
		return -1, b
	}

	return 1, b
}

func krToDate(date string) string{

	n := time.Now()
	var td string

	if strings.Contains(date, "오늘"){
		td = n.Format("20060102")
	}else if strings.Contains(date, "어제"){
		td = n.AddDate(0,0,-1).Format("20060102")
	}else if strings.Contains(date, "그제"){
		td = n.AddDate(0,0,-2).Format("20060102")
	}else if strings.Contains(date, "이번 주"){
		td = n.AddDate(0,0,-3).Format("20060102")
	}else if strings.Contains(date, "지난 주"){
		td = n.AddDate(0,0,-7).Format("20060102")
	}else if strings.Contains(date, "이번 달"){
		td = n.AddDate(0,0,-8).Format("20060102")
	}else if strings.Contains(date, "지난 달"){
		td = n.AddDate(0,-1,0).Format("20060102")
	}else if strings.Contains(date, "2개월 전"){
		td = n.AddDate(0,-2,0).Format("20060102")
	}else if strings.Contains(date, "3개월 전"){
		td = n.AddDate(0,-3,0).Format("20060102")
	}else if strings.Contains(date, "4개월 전"){
		td = n.AddDate(0,-4,0).Format("20060102")
	}else{
		td = n.AddDate(-1,0,0).Format("20060102")
	}

	return td

}