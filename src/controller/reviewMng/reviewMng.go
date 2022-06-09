package reviewMng

import (
	reviewsql "cashApi/query/reviews"
	storesql "cashApi/query/stores"
	"cashApi/src/controller"
	"cashApi/src/controller/cls"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	// login 및 기본

)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func TestPage(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "reviewMng/testPage.htm", m)
}


func ReviewSetting(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["restId"] = params["restId"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "reviewMng/reviewSetting.htm", m)
}

func ReviewList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["restId"] = params["restId"]
	m["startDt"] = params["startDt"]
	m["endDt"] = params["endDt"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "reviewMng/reviewList.htm", m)
}

func CustomReviewList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["restId"] = params["restId"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "reviewMng/customReviewList.htm", m)
}


func ReviewMain(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["restId"] = params["restId"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "reviewMng/reviewMain.htm", m)
}




func GetReviewMain(c echo.Context) error {
	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})


	endDate := time.Now().Format("20060102")
	startDate :=time.Now().AddDate(0, -5, 0).Format("20060102")
	params["startDate"]=startDate
	params["endDate"]=endDate

	deliveryInfo, err := cls.GetSelectData2(reviewsql.SelectStoreInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	delivery := make(map[string]interface{})
	if deliveryInfo == nil{
		m["resultCode"] = "01"
		m["resultMsg"] = "가맹점 정보가 없습니다."
		m["resultData"] = delivery
		return c.JSON(http.StatusOK, m)
	}

	b := deliveryInfo[0]["baeminId"]
	n := deliveryInfo[0]["naverId"]
	y := deliveryInfo[0]["yogiyoId"]
	cp := deliveryInfo[0]["coupangId"]

	if len(b) > 0{
		params["baeminId"] = b
	}else{
		params["baeminId"] = "baemin"
	}

	if len(n) > 0{
		params["naverId"] = n
	}else{
		params["naverId"] = "naver"
	}

	if len(y) > 0{
		params["yogiyoId"] = y
	}else{
		params["yogiyoId"] = "yogiyo"
	}

	if len(cp) > 0{
		params["coupangId"] =cp
	}else{
		params["coupangId"] = "coupang"
	}

	ratingInfo, err := cls.GetSelectData2(reviewsql.SelectReivewRating, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	allCnt :=0
	baminCnt :=0
	naverCnt :=0
	yogiyoCnt :=0
	coupangCnt :=0

	allPoint :=0
	baminPoint :=0
	naverPoint :=0
	yogiyoPoint :=0
	coupangPoint :=0

	for _,v := range ratingInfo{
		cnt := v["tot_cnt"]
		bcnt := v["bamin_cnt"]
		ycnt := v["yogiyo_cnt"]
		ncnt := v["naver_cnt"]
		ccnt := v["coupang_cnt"]
		t, _ := strconv.Atoi(cnt)
		b, _ := strconv.Atoi(bcnt)
		y, _ := strconv.Atoi(ycnt)
		n, _ := strconv.Atoi(ncnt)
		c, _ := strconv.Atoi(ccnt)

		star, _ := strconv.Atoi(v["rating"])
		allPoint += star*t
		allCnt += t

		baminCnt += b
		baminPoint += star*b

		naverCnt += n
		naverPoint += star*n

		yogiyoCnt += y
		yogiyoPoint += star*y

		coupangCnt += c
		coupangPoint += star*c
	}

	var allAvg, baminAvg, naverAvg, yogiyoAvg, coupangAvg float32
	if (allCnt > 0) {
		allAvg = float32(allPoint) /float32(allCnt)
	}
	if (baminCnt > 0) {
		baminAvg = float32(baminPoint) /float32(baminCnt)
	}
	if (naverCnt > 0) {
		naverAvg = float32(naverPoint) /float32(naverCnt)
	}
	if (yogiyoCnt > 0) {
		yogiyoAvg = float32(yogiyoPoint) /float32(yogiyoCnt)
	}
	if (coupangCnt > 0) {
		coupangAvg = float32(coupangPoint) /float32(coupangCnt)
	}



	monthRatingInfo, err := cls.GetSelectData2(reviewsql.SelectReivewRatingMonth, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	ratingCnt := make(map[string]interface{})
	ratingCnt["allCnt"] = allCnt
	ratingCnt["allAvg"] = allAvg
	ratingCnt["baminCnt"] = baminCnt
	ratingCnt["naverCnt"] = naverCnt
	ratingCnt["yogiyoCnt"] = yogiyoCnt
	ratingCnt["coupangCnt"] = coupangCnt
	ratingCnt["baminAvg"] = baminAvg
	ratingCnt["naverAvg"] = naverAvg
	ratingCnt["yogiyoAvg"] = yogiyoAvg
	ratingCnt["coupangAvg"] = coupangAvg


	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["ratingCnt"] = ratingCnt
	m["monthRatingInfo"] = monthRatingInfo
	m["resultData"] = ratingInfo
	return c.JSON(http.StatusOK, m)

}

func GetReviewMain2(c echo.Context) error {

	data := make(map[string]interface{})
	params := cls.GetParamJsonMap(c)

	contentsInfo, err := cls.GetSelectType(reviewsql.SelectContentList, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	tipInfo, err := cls.GetSelectType(reviewsql.SelectTipList, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	// 가맹점 wordcloud image 여부 판단
	fPath := fmt.Sprintf("/app/SharedStorage/wordCloud/%s/old_%s.png", params["restId"], params["restId"])
	lprintf(4, "[INFO] fPath(%s) \n", fPath)

	fInfo, err := os.Stat(fPath)
	if err != nil{
		data["wordCloudResult"] = "n"
	}else{
		if fInfo.Size() < 250000 {
			data["wordCloudResult"] = "n"
		}else{
			data["wordCloudResult"] = "y"
		}
	}

	billingInfo, err := cls.GetSelectData(reviewsql.SelectBillingInfo, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if billingInfo == nil {
		data["billingInfo"] = ""
	}else{
		data["billingInfo"] =billingInfo[0]
	}

	data["contents"] = contentsInfo
	data["tips"] = tipInfo

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data
	return c.JSON(http.StatusOK, m)

}



func GetReviewSetting(c echo.Context) error {
	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})

	deliveryInfo, err := cls.GetSelectData2(reviewsql.SelectStoreInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	delivery := make(map[string]interface{})


	storeInfo, err := cls.GetSelectData(storesql.SelectStoreInfo, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	delivery["addr"] = storeInfo[0]["ADDR"]
	delivery["addr2"] = storeInfo[0]["ADDR2"]
	delivery["bizNum"] = storeInfo[0]["BUSID"]


	if deliveryInfo == nil{

		m["resultCode"] = "01"
		m["resultMsg"] = "가맹점 정보가 없습니다."
		m["resultData"] = delivery
		return c.JSON(http.StatusOK, m)
	}

	b := deliveryInfo[0]["baeminId"]
	n := deliveryInfo[0]["naverId"]
	y := deliveryInfo[0]["yogiyoId"]
	cp := deliveryInfo[0]["coupangId"]


	if len(b) > 0{
		params["baeminId"] = b
	}else{
		params["baeminId"] = "baemin"
	}

	if len(n) > 0{
		params["naverId"] = n
	}else{
		params["naverId"] = "naver"
	}

	if len(y) > 0{
		params["yogiyoId"] = y
	}else{
		params["yogiyoId"] = "yogiyo"
	}

	if len(cp) > 0{
		params["coupangId"] =cp
	}else{
		params["coupangId"] = "coupang"
	}

	baeminInfo, err := cls.GetSelectData2(reviewsql.SelectBaeminInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	naverInfo, err := cls.GetSelectData2(reviewsql.SelectNaverInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	yogiyoInfo, err := cls.GetSelectData2(reviewsql.SelectYogiyoInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	coupangInfo, err := cls.GetSelectData2(reviewsql.SelectCoupangInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}



	if len(baeminInfo) > 0{
		delivery["baemin"] = baeminInfo[0]
	}else{
		delivery["baemin"] = "n"
	}

	if len(naverInfo) > 0{
		delivery["naverInfo"] = naverInfo[0]
	}else{
		delivery["naverInfo"] = "n"
	}

	if len(yogiyoInfo) > 0{
		delivery["yogiyoInfo"] = yogiyoInfo[0]
	}else{
		delivery["yogiyoInfo"] = "n"
	}

	if len(coupangInfo) > 0{
		delivery["coupangInfo"] = coupangInfo[0]
	}else{
		delivery["coupangInfo"] = "n"
	}


	filterData, err := cls.GetSelectData2(reviewsql.SelectBlackFilter, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	delivery["filterData"] = filterData

	billingInfo, err := cls.GetSelectData(reviewsql.SelectBillingInfo, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if billingInfo == nil {
		delivery["billingInfo"] = ""
	}else{
		delivery["billingInfo"] =billingInfo[0]
	}


	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = delivery
	return c.JSON(http.StatusOK, m)

}




func SetKeywordSetup(c echo.Context) error {

	dprintf(4, c, "call SetKeywordSetup\n")
	//
	params := cls.GetParamJsonMap(c)

	runType:="U"
	selectStoreCCompChk, err := cls.GetSelectData2(reviewsql.SelectBlackFilter, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if selectStoreCCompChk == nil{
		runType="I"
	}


	tx, err := cls.DBc2.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -필터 설정 SetBlackSetup)  \n")
			tx.Rollback()
		}
	}()
	rQuery :=""
	if runType=="U"{
		rQuery = reviewsql.UpdateStoreBlackFilter
	}else{
		rQuery = reviewsql.InsertStoreBlackFilter
	}

	StoreBlackFilterQuery, err := cls.GetQueryJson(rQuery, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UserCreateQuery parameter fail"))
	}

	_, err = tx.Exec(StoreBlackFilterQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", StoreBlackFilterQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"


	return c.JSON(http.StatusOK, m)

}




func SetStoreReviewInfo(c echo.Context) error {

	dprintf(4, c, "call SetStoreReviewInfo\n")
	//
	params := cls.GetParamJsonMap(c)

	selectQuery, err := cls.SetUpdateParam(storesql.UpdateStore, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}



	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"


	return c.JSON(http.StatusOK, m)

}




// 가맹점 현재 배달업체 셋팅값
func SetCompSetting(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	t := params["type"]

	var query string

	if t == "baemin" && len(params["baeminId"]) > 0{
		query = "UPDATE b_store SET baemin_id = '"+params["baeminId"]+"' WHERE REST_ID = '"+params["restId"]+"' ;"
	}else if t == "naver" && len(params["naverId"]) > 0{
		query = "UPDATE b_store SET naver_id = '"+params["naverId"]+"' WHERE REST_ID = '"+params["restId"]+"' ;"
	}else if t == "yogiyo" && len(params["yogiyoId"]) > 0{
		query = "UPDATE b_store SET yogiyo_id = '"+params["yogiyoId"]+"' WHERE REST_ID = '"+params["restId"]+"' ;"
	}else if t == "coupang" && len(params["coupangId"]) > 0{
		query = "UPDATE b_store SET COUPANG_ID = '"+params["coupangId"]+"' WHERE REST_ID = '"+params["restId"]+"' ;"
	}else{
		m["resultCode"] = "99"
		m["resultMsg"] = "응답 실패"

		return c.JSON(http.StatusOK, m)
	}


	_, err := cls.QueryDB2(query)
	if err != nil {
		lprintf(1, "[ERROR] cls.ExecDBbyParam error(%s) \n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}



// 가맹점 현재 배달업체 셋팅값
func GetDeliveryList(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	keyword := params["keyword"]

	var page int
	tmp, err := strconv.Atoi(params["page"])
	if err != nil{
		page = tmp
	}


	comp, err := cls.GetSelectData(reviewsql.SelectCompInfo, params,c)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	if len(comp) == 0{
		m["resultCode"] = "99"
		m["resultMsg"] = "응답 실패"

		return c.JSON(http.StatusOK, m)
	}

	lat := comp[0]["lat"]
	lng := comp[0]["lng"]
	compNm := comp[0]["comp_nm"]


	if len(keyword) == 0{
		keyword = compNm
	}

	t := params["type"]


	println(t)

	if t == "baemin"{

		rst, b := BaeminCompList(lat, lng, keyword)
		if rst > 0{
			m["resultData"] = b
		}

	}else if t == "naver"{

		rst, n := NaverCompList(lat, lng, keyword)
		if rst > 0{
			m["resultData"] = n
		}

	}else if t == "yogiyo"{

		rst, y := YogiyoCompList(lat, lng, keyword, page)
		if rst > 0{
			m["resultData"] = y
		}

	}else if t == "coupang"{
		rst, c := CoupangCompList(lat, lng, keyword)
		if rst > 0{
			m["resultData"] = c
		}
	}else{
		m["resultCode"] = "99"
		m["resultMsg"] = "응답 실패"

		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func CoupangCompList(lat, lng, compNm string) (int, []map[string]interface{}){

	rst, cs := GetCoupangComp(lat, lng, compNm)
	if rst < 0{
		return -1, nil
	}

	go SetCoupangComp(cs)

	deliverys := []map[string]interface{}{}

	for _,v := range cs.Data.EntityList{

		c := v.Entity.Data

		delivery := make(map[string]interface{})

		delivery["name"] = c.Name

		var tmp string
		for idx, ctg := range c.Categories{
			if idx == 0 {
				tmp = ctg
			}else{
				tmp += ","+ctg
			}
		}

		delivery["category"] = tmp
		delivery["addr"] = c.Address
		if len(c.ImagePaths) > 0{
			delivery["logo"] = c.ImagePaths[0]
		}else{
			delivery["logo"] = c.BrandLogoPath
		}
		delivery["coupangId"] = c.ID

		deliverys = append(deliverys, delivery)
	}

	return 1, deliverys

}

func SetCoupangComp(c CoupangComps){

	var query string

	for _, v := range c.Data.EntityList{

		c := v.Entity.Data

		if c.ID == 0{
			continue
		}

		var tmp string
		query = "REPLACE INTO a_coupang(COUPANG_ID, CATEGORIES, NAME, PAYMENT_STORE_ID, MERCHANT_ID, DESCRIPTION, TEL_NO, BIZNUM, ZIP_NO, ADDRESS, ADDRESS_DETAIL, " +
			"LAT, LNG, SERVICE_FEE_RATIO, MENUS, IMAGEPATH, TOPDISH_IMAGEPATH, BRANDLOGO_PATH) " +
			"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"

		var params []interface{}
		params = append(params, strconv.Itoa(c.ID))

		for idx, ctg := range c.Categories{
			if idx == 0 {
				tmp = ctg
			}else{
				tmp += ","+ctg
			}
		}
		params = append(params, tmp)

		params = append(params, c.Name)
		params = append(params, c.PaymentStoreID)
		params = append(params, c.MerchantID)
		params = append(params, c.Description)
		params = append(params, c.TelNo)

		bizNo := strings.ReplaceAll(c.BizNo, "-", "")

		params = append(params, bizNo)
		params = append(params, c.ZipNo)
		params = append(params, c.Address)
		params = append(params, c.AddressDetail)
		params = append(params, c.Latitude)
		params = append(params, c.Longitude)
		params = append(params, c.ServiceFeeRatio)

		tmp = ""
		for idx, menu := range c.Menus{
			if idx == 0 {
				tmp = menu
			}else{
				tmp += ","+menu
			}
		}
		params = append(params, tmp)

		if len(c.ImagePaths) > 0{
			params = append(params, c.ImagePaths[0])
		}else{
			params = append(params, "")
		}

		if len(c.TopDishImagePaths) > 0{
			params = append(params, c.TopDishImagePaths[0])
		}else{
			params = append(params, "")
		}

		params = append(params, c.BrandLogoPath)

		_, err := cls.ExecDBbyParam2(query, params)
		if err != nil {
			lprintf(1, "[ERROR] cls.ExecDBbyParam error(%s) \n", err.Error())
			continue
		}
	}

	return
}

func GetCoupangComp(lat, lng, compNm string) (int, CoupangComps){
	var c CoupangComps

	url := fmt.Sprintf("endpoint/store.get_search?keyWord=%s&sort=nearby", url.QueryEscape(compNm))

	var httpHeader map[string]string
	httpHeader = make(map[string]string)
	httpHeader["X-EATS-LOCALE"] = "ko-KR"
	httpHeader["X-EATS-LOCATION"] = fmt.Sprintf("{\"addressId\":0,\"latitude\":%s,\"longitude\":%s,\"regionId\":23}", lat, lng)
	httpHeader["X-EATS-APP-VERSION"] = "1.3.8"
	httpHeader["X-EATS-PCID"] = "03e2556d-8f3b-3f74-982e-5264228f0f78"
	httpHeader["X-EATS-DEVICE-ID"] = "03e2556d-8f3b-3f74-982e-5264228f0f78"
	httpHeader["X-EATS-OS-TYPE"] = "ANDROID"

	resp, err := cls.HttpRequestDetail("HTTPS", "GET", "api.coupangeats.com", "443", url, nil, httpHeader, "", false)
	if err != nil{
		lprintf(1, "[ERROR] coupang get comp addr err(%s)\n", err.Error())
		return -1, c
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		lprintf(1, "[ERROR] coupang get comp addr body read err(%s)\n", err.Error())
		return -1, c
	}

	err = json.Unmarshal(data, &c)
	if err != nil{
		lprintf(1, "[ERROR] coupang get comp addr parsing json err(%s)\n", err.Error())
		return -1, c
	}

	return 1, c
}

func YogiyoCompList(lat, lng, keyword string, page int) (int, []map[string]interface{}){

	rst, y := GetYogiyoComp(lat, lng, keyword,page)
	if rst < 0{
		return -1, nil
	}

	go SetYogiyoComp(y)

	deliverys := []map[string]interface{}{}

	for _,v := range y.Restaurants{

		delivery := make(map[string]interface{})

		delivery["name"] = v.Name

		var tmp string
		for idx, ctg := range v.Categories{
			if idx == 0 {
				tmp = ctg
			}else{
				tmp += ","+ctg
			}
		}

		delivery["category"] = tmp
		delivery["addr"] = v.Address
		delivery["logo"] = v.ThumbnailURL
		delivery["yogiyoId"] = v.ID

		deliverys = append(deliverys, delivery)
	}

	return 1, deliverys
}


func SetYogiyoComp(y YogiyoComps) int {

	var query, tmp string

	for _, v := range y.Restaurants{
		query = "REPLACE INTO a_yogiyo(YOGIYO_ID, COMP_NM, CATEGORIES, THUMBNAIL_URL, BEGIN, END, ADDRESS, YOGIYO_TYPE, REVIEW_AVG, CITY, " +
			"DELIVERY_FEE_EXPLANATION, FRANCHISE_ID, FRANCHISE_NAME, REVIEW_IMAGE_COUNT, OWNER_REPLY_COUNT, LAT, LNG, MIN_ORDER_AMOUNT, " +
			"ADDITIONAL_DISCOUNT) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"

		var params []interface{}
		params = append(params, v.ID)
		params = append(params, v.Name)

		for idx, ctg := range v.Categories{
			if idx == 0 {
				tmp = ctg
			}else{
				tmp += ","+ctg
			}
		}
		params = append(params, tmp)

		params = append(params, v.ThumbnailURL)
		params = append(params, v.Begin)
		params = append(params, v.End)
		params = append(params, v.Address)
		params = append(params, v.RestaurantType)
		params = append(params, v.ReviewAvg)
		params = append(params, v.City)
		params = append(params, v.DeliveryFeeExplanation)
		params = append(params, v.FranchiseID)
		params = append(params, v.FranchiseName)
		params = append(params, v.ReviewImageCount)
		params = append(params, v.OwnerReplyCount)
		params = append(params, v.Lat)
		params = append(params, v.Lng)
		params = append(params, v.MinOrderAmount)
		params = append(params, v.AdditionalDiscount)

		_, err := cls.ExecDBbyParam2(query, params)
		if err != nil {
			lprintf(1, "[ERROR] cls.ExecDBbyParam error(%s) \n", err.Error())
			continue
		}
	}

	return 1
}



func GetYogiyoComp(lat, lng, keyword string, page int) (int, YogiyoComps){
	var y YogiyoComps

	url := fmt.Sprintf("api/v1/restaurants-geo/search?items=100&lat=%s&lng=%s&order=distance&page=%d&search=%s", lat, lng, page, url.QueryEscape(keyword))

	var httpHeader map[string]string
	httpHeader = make(map[string]string)
	httpHeader["x-apisecret"] = "fe5183cc3dea12bd0ce299cf110a75a2"
	httpHeader["x-apikey"] = "iphoneap"

	resp, err := cls.HttpRequestDetail("HTTPS", "GET", "www.yogiyo.co.kr", "443", url, nil, httpHeader, "", false)
	if err != nil{
		lprintf(1, "[ERROR] yogiyo get comp addr err(%s)\n", err.Error())
		return -1, y
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		lprintf(1, "[ERROR] yogiyo get comp addr body read err(%s)\n", err.Error())
		return -1, y
	}

	err = json.Unmarshal(data, &y)
	if err != nil{
		lprintf(1, "[ERROR] yogiyo get comp addr parsing json err(%s)\n", err.Error())
		return -1, y
	}

	return 1, y
}


// 가맹점 리스트 받기
type YogiyoComps struct {
	Pagination struct {
		PerPage      int `json:"per_page"`
		TotalObjects int `json:"total_objects"`
		CurrentPage  int `json:"current_page"`
		TotalPages   int `json:"total_pages"`
	} `json:"pagination"`
	Restaurants []struct {
		ID                                string   `json:"id"`
		Name                              string   `json:"name"`
		Rating                            float64  `json:"rating"`
		SmsBonus                          bool     `json:"sms_bonus"`
		RelayMethods                      []string `json:"relay_methods"`
		CentralBilling                    bool     `json:"central_billing"`
		Reachable                         bool     `json:"reachable"`
		Phone                             string   `json:"phone"`
		PhoneOrder                        bool     `json:"phone_order"`
		AppOrder                          bool     `json:"app_order"`
		Address                           string   `json:"address"`
		Slug                              string   `json:"slug"`
		Takeout                           bool     `json:"takeout"`
		IsAvailableDelivery               bool     `json:"is_available_delivery"`
		IsAvailablePickup                 bool     `json:"is_available_pickup"`
		MinimumPickupMinutes              int      `json:"minimum_pickup_minutes"`
		IsFoodfly                         bool     `json:"is_foodfly"`
		NewUntil                          string   `json:"new_until"`
		ExceptCash                        bool     `json:"except_cash"`
		PhoneDownlisted                   bool     `json:"phone_downlisted"`
		City                              string   `json:"city"`
		RestaurantType                    string   `json:"restaurant_type"`
		DiscountPercent                   int      `json:"discount_percent"`
		ReviewCount                       int      `json:"review_count"`
		OwnerReplyCount                   int      `json:"owner_reply_count"`
		ReviewAvg                         float64  `json:"review_avg"`
		Top28                             bool     `json:"top28"`
		FranchiseID                       int      `json:"franchise_id"`
		FranchiseName                     string   `json:"franchise_name"`
		PaymentMethods                    []string `json:"payment_methods"`
		Tags                              []string `json:"tags"`
		AdditionalDiscountCurrentlyActive bool     `json:"additional_discount_currently_active"`
		Threshold                         int      `json:"threshold"`
		ReviewImageCount                  int      `json:"review_image_count"`
		Lat                               float64  `json:"lat"`
		Lng                               float64  `json:"lng"`
		EstimatedDeliveryTime             string   `json:"estimated_delivery_time"`
		EstimatedDeliveryTimeKey          int      `json:"estimated_delivery_time_key"`
		Categories                        []string `json:"categories"`
		HasTerminal                       bool     `json:"has_terminal"`
		LogoURL                           string   `json:"logo_url"`
		ThumbnailURL                      string   `json:"thumbnail_url"`
		CanReview                         int      `json:"can_review"`
		HasLoyaltySupport                 bool     `json:"has_loyalty_support"`
		AdditionalDiscount                int      `json:"additional_discount"`
		PhoneRating                       int      `json:"phone_rating"`
		Section                           string   `json:"section"`
		AdvertisementType                 string   `json:"advertisement_type"`
		AdvertisementRank                 int      `json:"advertisement_rank"`
		AdvDistance                       int      `json:"adv_distance"`
		Premium                           struct {
		} `json:"premium"`
		ListPos                int      `json:"list_pos"`
		Begin                  string   `json:"begin"`
		End                    string   `json:"end"`
		DeliveryFee            int      `json:"delivery_fee"`
		FreeDeliveryThreshold  float64  `json:"free_delivery_threshold"`
		MinOrderAmount         int      `json:"min_order_amount"`
		DeliveryFeeExplanation string   `json:"delivery_fee_explanation"`
		Open                   bool     `json:"open"`
		IsDeliverable          float64  `json:"is_deliverable"`
		Score                  float64  `json:"score"`
		Menus                  []string `json:"menus"`
		SearchScore            float64  `json:"search_score"`
		Distance               float64  `json:"distance"`
		OpenTimeDescription    string   `json:"open_time_description"`
		New                    bool     `json:"new"`
		ThumbnailMessage       string   `json:"thumbnail_message"`
		DiscountUntil          string   `json:"discount_until,omitempty"`
		DiscountFrom           string   `json:"discount_from,omitempty"`
		AdjustedDeliveryFee    int      `json:"adjusted_delivery_fee"`
		DiscountedDeliveryFee  int      `json:"discounted_delivery_fee"`
		DeliveryFeeToDisplay   struct {
			Basic string `json:"basic"`
		} `json:"delivery_fee_to_display"`
		DeliveryMethod string `json:"delivery_method"`
		SectionPos     int    `json:"section_pos"`
		MenusInfo      []struct {
			MenuName string `json:"menu_name"`
			MenuID   string `json:"menu_id"`
		} `json:"menus_info"`
		RepresentativeMenus string `json:"representative_menus"`
	} `json:"restaurants"`
}



func BaeminCompList(lat, lng, compNm string) (int, []map[string]interface{}){

	rst, b := GetBaeminComp(lat, lng, compNm)
	if rst < 0{
		return -1, nil
	}

	go SetBaeminComp(b)

	deliverys := []map[string]interface{}{}

	for _,vShop := range b.Data.Shops{

		v := vShop.Shopinfo
		delivery := make(map[string]interface{})

		delivery["name"] = v.Shopname
		delivery["category"] = v.Categorynamekor
		delivery["addr"] = v.Address
		delivery["logo"] = v.Logourl
		delivery["baeminId"] = v.Shopnumber

		deliverys = append(deliverys, delivery)
	}

	return 1, deliverys
}

func SetBaeminComp(b BaeminComps) int {

	var query string

	for _, vShop := range b.Data.Shops{
		query = "REPLACE INTO a_baemin(BAEMIN_ID, COMP_NM, BAEMIN_TYPE, CATEGORY_NAME_KR, CATEGORY_NAME_ENG, LOGO_URL, INTRO_TEXT, CLOSE_DAY_TEXT, ADDRESS, TEL, TEL_VIRTUAL, " +
			"FRANCHISE_NUMBER, FRANCHISE_TEL_NUMBER, REPRESENTATION_MENU, DELIVERY_AREA_TEXT, MINIMUM_ORDER_PRICE, DELIVERY_TIP_PHRASE, EXPECTED_DELIVERY_TIME_PHRASE, DISTANCE_PHRASE, " +
			"DELIVERY_TIP_DISCOUNT, DELIVERY_TIP_ZERO, FASE_DELIVERY) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"

		v := vShop.Shopinfo

		var params []interface{}
		params = append(params, strconv.Itoa(v.Shopnumber))
		params = append(params, v.Shopname)
		params = append(params, v.Servicetype)
		params = append(params, v.Categorynamekor)
		params = append(params, v.Categorynameeng)
		params = append(params, v.Logourl)
		params = append(params, v.Introtext)
		params = append(params, v.Closedaytext)
		params = append(params, v.Address)
		params = append(params, v.Telnumber)
		params = append(params, v.Virtualtelnumber)
		params = append(params, v.Franchisenumber)
		params = append(params, v.Franchisetelnumber)
		params = append(params, v.Representationmenu)
		params = append(params, vShop.Deliveryinfo.Deliveryareatext)
		params = append(params, vShop.Deliveryinfo.Minimumorderprice)
		params = append(params, vShop.Deliveryinfo.Deliverytipphrase)
		params = append(params, vShop.Deliveryinfo.Expecteddeliverytimephrase)
		params = append(params, vShop.Deliveryinfo.Distancephrase)
		params = append(params, vShop.Deliveryinfo.Deliverytipdiscount)
		params = append(params, vShop.Deliveryinfo.Deliverytipzero)
		params = append(params, vShop.Deliveryinfo.Fastdelivery)

		_, err := cls.ExecDBbyParam2(query, params)
		if err != nil {
			lprintf(1, "[ERROR] cls.ExecDBbyParam error(%s) \n", err.Error())
			continue
		}
	}

	return 1
}

func GetBaeminComp(lat, lng, compNm string) (int, BaeminComps){
	var b BaeminComps

	url := fmt.Sprintf("v2/SEARCH/shops?keyword=%s&filter=&sort=SORT__DISTANCE&referral=Search&kind=DEFAULT&offset=0&limit=30&latitude=%s&longitude=%s&extension=&appver=10.27.1&carrier=45008&site=7jWXRELC2e&deviceModel=SM-G906S&dvcid=OPUD5687e4685d245b9c&adid=NONE&sessionId=91f4e62a5d1d615dfb1865&osver=23&oscd=2", url.QueryEscape(compNm), lat, lng)

	resp, err := cls.HttpRequestDetail("HTTPS", "GET", "shopdp-api.baemin.com", "443", url, nil, nil, "", false)
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

// 가맹점 리스트 얻기
type BaeminComps struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	Serverdatetime string `json:"serverDatetime"`
	Data           struct {
		Totalcount  int `json:"totalCount"`
		Serviceinfo struct {
			Existresult          bool   `json:"existResult"`
			Requestid            string `json:"requestId"`
			Resultkind           string `json:"resultKind"`
			Noresultimageurl     string `json:"noResultImageUrl"`
			Noresulttext         string `json:"noResultText"`
			Correctedkeywordtext string `json:"correctedKeywordText"`
			Correctedkeyword     string `json:"correctedKeyword"`
			Blocks               []struct {
				Blocktype   string `json:"blockType"`
				Description string `json:"description"`
				Tooltiptext string `json:"tooltipText"`
				Existbar    bool   `json:"existBar"`
				Ad          bool   `json:"ad"`
			} `json:"blocks"`
			Recommend struct {
				Text     string `json:"text"`
				Subtext  string `json:"subText"`
				Category string `json:"category"`
			} `json:"recommend"`
		} `json:"serviceInfo"`
		Shops []struct {
			Blocktype   string `json:"blockType"`
			Landingtype string `json:"landingType"`
			Shopinfo    struct {
				Shopnumber         int    `json:"shopNumber"`
				Shopname           string `json:"shopName"`
				Servicetype        string `json:"serviceType"`
				Categorycode       string `json:"categoryCode"`
				Categorynamekor    string `json:"categoryNameKor"`
				Categorynameeng    string `json:"categoryNameEng"`
				Logourl            string `json:"logoUrl"`
				Introtext          string `json:"introText"`
				Closedaytext       string `json:"closeDayText"`
				Address            string `json:"address"`
				Telnumber          string `json:"telNumber"`
				Virtualtelnumber   string `json:"virtualTelNumber"`
				Franchisenumber    int    `json:"franchiseNumber"`
				Franchisetelnumber string `json:"franchiseTelNumber"`
				Representationmenu string `json:"representationMenu"`
			} `json:"shopInfo"`
			Shopstatus struct {
				Inoperation bool `json:"inOperation"`
			} `json:"shopStatus"`
			Deliveryinfo struct {
				Deliveryareatext           string `json:"deliveryAreaText"`
				Minimumorderprice          int    `json:"minimumOrderPrice"`
				Deliverytipphrase          string `json:"deliveryTipPhrase"`
				Expecteddeliverytimephrase string `json:"expectedDeliveryTimePhrase"`
				Distancephrase             string `json:"distancePhrase"`
				Deliverytipdiscount        bool   `json:"deliveryTipDiscount"`
				Deliverytipzero            bool   `json:"deliveryTipZero"`
				Fastdelivery               bool   `json:"fastDelivery"`
			} `json:"deliveryInfo"`
			Adinfo struct {
				Campaignid string `json:"campaignId"`
			} `json:"adInfo"`
			Shopstatistics struct {
				Averagestarscore      float64 `json:"averageStarScore"`
				Favoritecount         int     `json:"favoriteCount"`
				Latestreviewcount     int     `json:"latestReviewCount"`
				Latestceocommentcount int     `json:"latestCeoCommentCount"`
				Latestordercount      int     `json:"latestOrderCount"`
			} `json:"shopStatistics"`
			Decoinfo struct {
				Thumbnail      bool `json:"thumbnail"`
				Backgrounddeco bool `json:"backgroundDeco"`
				Addonbadges    []struct {
					Type       string `json:"type"`
					Text       string `json:"text"`
					Background struct {
						Color string `json:"color"`
						Alpha string `json:"alpha"`
					} `json:"background"`
					Border struct {
						Color string `json:"color"`
						Alpha string `json:"alpha"`
					} `json:"border"`
					Font struct {
						Color string `json:"color"`
						Alpha string `json:"alpha"`
					} `json:"font"`
				} `json:"addonBadges"`
				Servicebadges []interface{} `json:"serviceBadges"`
				Shopbadges    []interface{} `json:"shopBadges"`
			} `json:"decoInfo"`
			Loginfo struct {
				Displaymenus          []string `json:"displayMenus"`
				Deliverytips          []int    `json:"deliveryTips"`
				Expecteddeliverytimes []int    `json:"expectedDeliveryTimes"`
				Trackinglog           struct {
					Fastdelivery                 bool `json:"fastDelivery"`
					Onepickdelivery              bool `json:"onePickDelivery"`
					Fastdeliverydeliverytiplimit int  `json:"fastDeliveryDeliveryTipLimit"`
				} `json:"trackingLog"`
			} `json:"logInfo"`
		} `json:"shops"`
		Resulttype string `json:"resultType"`
		Extension  string `json:"extension"`
	} `json:"data"`
}




func NaverCompList(lat, lng, compNm string) (int, []map[string]interface{}){


	println(compNm)

	rst, n := GetNaverCompWithName(lat, lng, compNm)
	if rst < 0{
		return -1, nil
	}

	go SetNaverCompWithName(n)

	deliverys := []map[string]interface{}{}

	for _,v := range n.Place{

		delivery := make(map[string]interface{})

		delivery["name"] = v.Title
		delivery["category"] = v.Ctg
		delivery["addr"] = v.Jibunaddress
		delivery["addr_road"] = v.Roadaddress
		delivery["naverId"] = v.ID

		deliverys = append(deliverys, delivery)
	}

	return 1, deliverys
}



func SetNaverCompWithName(n NaverCompsWithName) int {

	var query string

	for _, v := range n.Place{
		query = "REPLACE INTO a_naver(NAVER_ID, NAVER_TYPE, COMP_NM, LAT, LNG, DIST, TOTAL_SCORE, CTG, CID, JIBUN_ADDRESS, ROAD_ADDRESS) VALUES(?,?,?,?,?,?,?,?,?,?,?);"

		var params []interface{}
		params = append(params, v.ID)
		params = append(params, v.Type)
		params = append(params, v.Title)
		params = append(params, v.Y)
		params = append(params, v.X)
		params = append(params, v.Dist)
		params = append(params, v.Totalscore)
		params = append(params, v.Ctg)
		params = append(params, v.Cid)
		params = append(params, v.Jibunaddress)
		params = append(params, v.Roadaddress)

		_, err := cls.ExecDBbyParam2(query, params)
		if err != nil {
			lprintf(1, "[ERROR] cls.ExecDBbyParam error(%s) \n", err.Error())
			continue
		}
	}

	return 1
}




func GetNaverCompWithName(lat, lng, compNm string) (int, NaverCompsWithName){

	var n NaverCompsWithName

	url := fmt.Sprintf("v5/api/instantSearch?lang=ko&caller=pcweb&types=place,address,bus&coords=%s,%s&query=%s", lat, lng, url.QueryEscape(compNm))
	//url = fmt.Sprintf("v5/api/instantSearch?lang=ko&caller=pcweb&types=place,address,bus&coords=%s,%s&query=%s", lat, lng, url.QueryEscape(compNm))

	resp, err := cls.HttpRequestDetail("HTTPS", "GET", "map.naver.com", "443", url, nil, nil, "", false)
	if err != nil{
		lprintf(1, "[ERROR] naver get comp addr err(%s)\n", err.Error())
		return -1, n
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		lprintf(1, "[ERROR] naver get comp addr body read err(%s)\n", err.Error())
		return -1, n
	}

	err = json.Unmarshal(data, &n)
	if err != nil{
		lprintf(1, "[ERROR] naver get comp addr parsing json err(%s)\n", err.Error())
		return -1, n
	}

	return 1, n
}

// 가맹점 이름으로 아이디 얻기
type NaverCompsWithName struct {
	Meta struct {
		Model     string `json:"model"`
		Query     string `json:"query"`
		Requestid string `json:"requestId"`
	} `json:"meta"`
	Ac    []string `json:"ac"`
	Place []struct {
		Type         string  `json:"type"`
		ID           string  `json:"id"`
		Title        string  `json:"title"`
		X            string  `json:"x"`
		Y            string  `json:"y"`
		Dist         float64 `json:"dist"`
		Totalscore   float64 `json:"totalScore"`
		Sid          string  `json:"sid"`
		Ctg          string  `json:"ctg"`
		Cid          string  `json:"cid"`
		Jibunaddress string  `json:"jibunAddress"`
		Roadaddress  string  `json:"roadAddress"`
		Review       struct {
			Count string `json:"count"`
		} `json:"review"`
	} `json:"place"`
	Address []interface{} `json:"address"`
	Bus     []interface{} `json:"bus"`
	Menu    []interface{} `json:"menu"`
	All     []struct {
		Place struct {
			Type         string  `json:"type"`
			ID           string  `json:"id"`
			Title        string  `json:"title"`
			X            string  `json:"x"`
			Y            string  `json:"y"`
			Dist         float64 `json:"dist"`
			Totalscore   float64 `json:"totalScore"`
			Sid          string  `json:"sid"`
			Ctg          string  `json:"ctg"`
			Cid          string  `json:"cid"`
			Jibunaddress string  `json:"jibunAddress"`
			Roadaddress  string  `json:"roadAddress"`
			Review       struct {
				Count string `json:"count"`
			} `json:"review"`
		} `json:"place"`
		Address interface{} `json:"address"`
		Bus     interface{} `json:"bus"`
	} `json:"all"`
}

// 가맹점 리스트 받기
type CoupangComps struct {
	Data struct {
		MappedKeyword interface{} `json:"mappedKeyword"`
		NextToken     interface{} `json:"nextToken"`
		EntityList    []struct {
			Entity struct {
				Data struct {
					OpenStatus            string      `json:"openStatus"`
					OpenStatusText        string      `json:"openStatusText"`
					NextOpenAt            string      `json:"nextOpenAt"`
					RemainingTime         interface{} `json:"remainingTime"`
					Distance              string      `json:"distance"`
					EstimatedDeliveryTime string      `json:"estimatedDeliveryTime"`
					Shareable             bool        `json:"shareable"`
					Benefit               interface{} `json:"benefit"`
					Favorite              bool        `json:"favorite"`
					ReviewRating          float64     `json:"reviewRating"`
					ReviewCount           int         `json:"reviewCount"`
					ReviewCountTexts      []struct {
						Text          string `json:"text"`
						Color         string `json:"color"`
						Size          int    `json:"size"`
						Bold          bool   `json:"bold"`
						StrikeThrough bool   `json:"strikeThrough"`
						Underline     bool   `json:"underline"`
					} `json:"reviewCountTexts"`
					OrderCountText  interface{} `json:"orderCountText"`
					ShowOrderCount  bool        `json:"showOrderCount"`
					DeliveryFeeInfo string      `json:"deliveryFeeInfo"`
					DeliveryFeeText struct {
						Text          string `json:"text"`
						Color         string `json:"color"`
						Size          int    `json:"size"`
						Bold          bool   `json:"bold"`
						StrikeThrough bool   `json:"strikeThrough"`
						Underline     bool   `json:"underline"`
					} `json:"deliveryFeeText"`
					ServiceFeeInfo []struct {
						Text          string `json:"text"`
						Color         string `json:"color"`
						Size          int    `json:"size"`
						Bold          bool   `json:"bold"`
						StrikeThrough bool   `json:"strikeThrough"`
						Underline     bool   `json:"underline"`
					} `json:"serviceFeeInfo"`
					ExtraDocuments interface{} `json:"extraDocuments"`
					Badges         struct {
					} `json:"badges"`
					SearchID             string `json:"searchId"`
					ExpressDeliveryBadge struct {
						TierName     string `json:"tierName"`
						ImagePath    string `json:"imagePath"`
						BadgeToolTip string `json:"badgeToolTip"`
						ExpressBadge bool   `json:"expressBadge"`
					} `json:"expressDeliveryBadge"`
					BenefitCouponInfo      interface{}   `json:"benefitCouponInfo"`
					PreviouslyOrderedInfo  interface{}   `json:"previouslyOrderedInfo"`
					CuratedType            interface{}   `json:"curatedType"`
					NewStoreBadge          interface{}   `json:"newStoreBadge"`
					NameTexts              interface{}   `json:"nameTexts"`
					MatchedDishes          interface{}   `json:"matchedDishes"`
					ReviewCarousels        interface{}   `json:"reviewCarousels"`
					ReviewShortCuts        interface{}   `json:"reviewShortCuts"`
					Emphasized             bool          `json:"emphasized"`
					MinimumOrderThresholds interface{}   `json:"minimumOrderThresholds"`
					AdProperty             interface{}   `json:"adProperty"`
					ID                     int           `json:"id"`
					MerchantID             int           `json:"merchantId"`
					Categories             []string      `json:"categories"`
					PaymentStoreID         string        `json:"paymentStoreId"`
					Name                   string        `json:"name"`
					Description            interface{}   `json:"description"`
					TelNo                  string        `json:"telNo"`
					BizNo                  string        `json:"bizNo"`
					ApprovalStatus         string        `json:"approvalStatus"`
					ZipNo                  string        `json:"zipNo"`
					Address                string        `json:"address"`
					AddressDetail          string        `json:"addressDetail"`
					Latitude               float64       `json:"latitude"`
					Longitude              float64       `json:"longitude"`
					ServiceFeeRatio        float64       `json:"serviceFeeRatio"`
					Menus                  []string      `json:"menus"`
					MenuSource             interface{}   `json:"menuSource"`
					ImagePaths             []string      `json:"imagePaths"`
					TopDishImagePaths      []string      `json:"topDishImagePaths"`
					ImageHeightRatio       float64       `json:"imageHeightRatio"`
					TaxBaseType            string        `json:"taxBaseType"`
					StoreLevelInfoID       int           `json:"storeLevelInfoId"`
					ManuallyShutdown       bool          `json:"manuallyShutdown"`
					Deleted                bool          `json:"deleted"`
					BrandLogoPath          string        `json:"brandLogoPath"`
				} `json:"data"`
			} `json:"entity"`
			ViewType string `json:"viewType"`
		} `json:"entityList"`
	} `json:"data"`
	Error interface{} `json:"error"`
}





func GetReviewList(c echo.Context) error {

	dprintf(4, c, "call GetReviewList\n")

	params := cls.GetParamJsonMap(c)

	/* param
	restId - id
	reviewType - all, baemin, naver, yogiyo
	rating - all, 5, 4, 3, 2, 1
	startDt - 시작
	endDt - 끝
	*/


	dType := params["dType"]

	m := make(map[string]interface{})
	m["resultCode"] = "98"
	m["resultMsg"] = "DB fail"

	deliveryInfo, err := cls.GetSelectData2(reviewsql.SelectStoreInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if deliveryInfo == nil{

		m["resultCode"] = "01"
		m["resultMsg"] = "리뷰 정보가 없습니다."
		m["resultData"] = []string{}
		return c.JSON(http.StatusOK, m)
	}

	b := deliveryInfo[0]["baeminId"]
	n := deliveryInfo[0]["naverId"]
	y := deliveryInfo[0]["yogiyoId"]
	cp := deliveryInfo[0]["coupangId"]



	if len(b) > 0   {
		params["baeminId"] = b
	}else{
		params["baeminId"] = "baemin"
	}

	if len(n) > 0{
		params["naverId"] = n
	}else{
		params["naverId"] = "naver"
	}

	if len(y) > 0{
		params["yogiyoId"] = y
	}else{
		params["yogiyoId"] = "yogiyo"
	}

	if len(cp) > 0{
		params["coupangId"] = cp
	}else{
		params["coupangId"] = "coupang"
	}


	if dType=="naver"{
		params["yogiyoId"] = "yogiyo"
		params["baeminId"] = "baemin"
		params["coupangId"] = "coupang"
	}else if dType=="baemin"{
		params["yogiyoId"] = "yogiyo"
		params["naverId"] = "naver"
		params["coupangId"] = "coupang"
	}else if dType=="yogiyo"{
		params["naverId"] = "naver"
		params["baeminId"] = "baemin"
		params["coupangId"] = "coupang"
	}else if dType=="coupang"{
		params["naverId"] = "naver"
		params["baeminId"] = "baemin"
		params["yogiyoId"] = "yogiyo"
	}



	resultCnt, err := cls.GetSelectData2(reviewsql.SelectReviewRatingContentCnt, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}


	review, err := cls.GetSelectData2(reviewsql.SelectReviewRatingContent, params)
	if err != nil {
		return c.JSON(http.StatusOK, m)
	}


	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] =resultCnt[0]["TOTAL_COUNT"]
	m["resultList"] = review


	return c.JSON(http.StatusOK, m)
}


func GetCustomReviewList(c echo.Context) error {

	dprintf(4, c, "call GetCustomReviewList\n")

	params := cls.GetParamJsonMap(c)

	dType := params["dType"]

	m := make(map[string]interface{})
	m["resultCode"] = "98"
	m["resultMsg"] = "DB fail"

	filterData, err := cls.GetSelectData2(reviewsql.SelectBlackFilter, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	rating := ""
	if filterData == nil{
		rating = ""
		params["rating"] =rating
		params["keyword"]=""

	}else{
		rating = filterData[0]["rating"]
		params["rating"] =rating
		params["keyword"]=filterData[0]["keyword"]
	}



	deliveryInfo, err := cls.GetSelectData2(reviewsql.SelectStoreInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if deliveryInfo == nil{

		m["resultCode"] = "01"
		m["resultMsg"] = "리뷰 정보가 없습니다."
		m["filterData"] = filterData
		m["resultData"] = []string{}
		return c.JSON(http.StatusOK, m)
	}

	b := deliveryInfo[0]["baeminId"]
	n := deliveryInfo[0]["naverId"]
	y := deliveryInfo[0]["yogiyoId"]
	cp := deliveryInfo[0]["coupangId"]



	if len(b) > 0   {
		params["baeminId"] = b
	}else{
		params["baeminId"] = "baemin"
	}

	if len(n) > 0{
		params["naverId"] = n
	}else{
		params["naverId"] = "naver"
	}

	if len(y) > 0{
		params["yogiyoId"] = y
	}else{
		params["yogiyoId"] = "yogiyo"
	}

	if len(cp) > 0{
		params["coupangId"] = cp
	}else{
		params["coupangId"] = "coupang"
	}


	if dType=="naver"{
		params["yogiyoId"] = "yogiyo"
		params["baeminId"] = "baemin"
		params["coupangId"] = "coupang"
	}else if dType=="baemin"{
		params["yogiyoId"] = "yogiyo"
		params["naverId"] = "naver"
		params["coupangId"] = "coupang"
	}else if dType=="yogiyo"{
		params["naverId"] = "naver"
		params["baeminId"] = "baemin"
		params["coupangId"] = "coupang"
	}else if dType=="coupang"{
		params["naverId"] = "naver"
		params["baeminId"] = "baemin"
		params["yogiyoId"] = "yogiyo"
	}



	cntQuery :=reviewsql.SelectCustomReviewRatingContentCnt_Norating
	seectQuery :=reviewsql.SelectCustomReviewRatingContent_Norating

	if len(rating) > 0 {
		cntQuery =reviewsql.SelectCustomReviewRatingContentCnt
		seectQuery =reviewsql.SelectCustomReviewRatingContent
	}

	resultCnt, err := cls.GetSelectData2(cntQuery, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	review, err := cls.GetSelectData2(seectQuery, params)
	if err != nil {
		return c.JSON(http.StatusOK, m)
	}

	billingInfo, err := cls.GetSelectData(reviewsql.SelectBillingInfo, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if billingInfo == nil {
		m["billingInfo"] = ""
	}else{
		m["billingInfo"] =billingInfo[0]
	}


	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultCnt"] =resultCnt[0]["TOTAL_COUNT"]
	m["resultList"] = review
	m["filterData"] = filterData


	return c.JSON(http.StatusOK, m)
}



// 리뷰 작성자 정보
func GetReviewWriter(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})

	memberNo := params["memberNo"]

	rst, b := httpBaeminCustomerInfo(memberNo, "2000")
	if rst < 0{
		m["resultCode"] = "99"
		m["resultMsg"] = "응답 실패"
		return c.JSON(http.StatusOK, m)
	}

	go SetBaeminCustomerInfo(b, memberNo)

	var totRating float64
	writer := make(map[string]interface{})
	reviews := []map[string]interface{}{}

	for _,v := range b.Data.Reviews{
		totRating += v.Rating

		review := make(map[string]interface{})
		review["id"] = v.ID
		review["rating"] = v.Rating
		review["review"] = v.Contents
		review["shopName"] = v.Shop.Name

		// 주문 메뉴
		var menu string
		for _,m := range v.Menus{
			menu += fmt.Sprintf("%s,",m.Name)
		}

		if len(menu) > 0{
			review["menu"] = menu[:len(menu)-1]
		}else{
			review["menu"] = ""
		}

		// 사장님 답변, 닉네임
		if len(v.Comments) > 0{
			review["commentYn"] = "y"
			review["commentReview"] = v.Comments[0].Contents
			review["commentName"] = v.Comments[0].Nickname
		}else{
			review["commentYn"] = "n"
		}

		review["date"] = v.DateText

		reviews = append(reviews, review)
	}

	writer["name"] = b.Data.Member.Nickname
	writer["reviewCnt"] = b.Data.ReviewCount
	writer["reviewRatingAvg"] = fmt.Sprintf("%.2f", totRating/float64(b.Data.ReviewCount))
	writer["reviews"] = reviews

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = writer

	return c.JSON(http.StatusOK, m)
}


func httpBaeminCustomerInfo(customerId, limit string) (int, BaeminCustomerInfo){
	var b BaeminCustomerInfo

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


type BaeminCustomerInfo struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	ServerDatetime string `json:"serverDatetime"`
	Data           struct {
		Member struct {
			Nickname string `json:"nickname"`
			Grade    string `json:"grade"`
			ImageURL string `json:"imageUrl"`
		} `json:"member"`
		ReviewCount int `json:"reviewCount"`
		Reviews     []struct {
			ID                      int64   `json:"id"`
			Rating                  float64 `json:"rating"`
			CeoOnlyMessage          string  `json:"ceoOnlyMessage"`
			AbusingSuspectedMessage string  `json:"abusingSuspectedMessage"`
			BlockMessage            string  `json:"blockMessage"`
			Contents                string  `json:"contents"`
			Modifiable              bool    `json:"modifiable"`
			Deletable               bool    `json:"deletable"`
			DisplayType             string  `json:"displayType"`
			DisplayStatus           string  `json:"displayStatus"`
			MenuDisplayType         string  `json:"menuDisplayType"`
			Menus                   []struct {
				MenuID         int    `json:"menuId"`
				ReviewMenuID   int64  `json:"reviewMenuId"`
				Name           string `json:"name"`
				Recommendation string `json:"recommendation"`
				Contents       string `json:"contents"`
			} `json:"menus"`
			Comments []struct {
				ID             int64  `json:"id"`
				Nickname       string `json:"nickname"`
				ImageURL       string `json:"imageUrl"`
				Contents       string `json:"contents"`
				DisplayStatus  string `json:"displayStatus"`
				BlockMessage   string `json:"blockMessage"`
				CeoOnlyMessage string `json:"ceoOnlyMessage"`
				DateText       string `json:"dateText"`
			} `json:"comments"`
			Images []struct {
				ID  int64  `json:"id"`
				URL string `json:"url"`
			} `json:"images"`
			Shop struct {
				No          int    `json:"no"`
				Name        string `json:"name"`
				ServiceType string `json:"serviceType"`
			} `json:"shop"`
			DateText string `json:"dateText"`
		} `json:"reviews"`
	} `json:"data"`
}


func SetBaeminCustomerInfo(b BaeminCustomerInfo, memberNo string){

	query := "REPLACE INTO a_baemin_customer(REVIEW_ID, MEMBER_NO, MEMBER_NAME, RATING, CONTENTS, MENUS, BAEMIN_ID, BAEMIN_NAME, " +
		"BOSS_NICNAME, BOSS_CONTENTS, DATETEXT, DATE) " +
		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?);"

	// 이모티콘 제거
	// mysql varchar2 -> uft8 기준(3비트)
	// 이모티콘 utf16? -> 여튼 4비트라서 mysql에 안들어감
	reg,regErr := regexp.Compile("[^\u0000-\uFFFF]")
	if regErr != nil{
		lprintf(1, "[ERROR] regexp compile err(%s)\n")
	}

	for _, review := range b.Data.Reviews{

		var params []interface{}

		params = append(params, strconv.FormatInt(review.ID, 10))
		params = append(params, memberNo)
		params = append(params, b.Data.Member.Nickname)
		params = append(params, review.Rating)

		if regErr == nil{
			params = append(params, reg.ReplaceAllString(review.Contents, ""))
		}else{
			params = append(params, review.Contents)
		}

		// 메뉴
		var tmp string
		for _,menu := range review.Menus{
			tmp += fmt.Sprintf("%s,",menu.Name)
		}

		if len(tmp) > 0{
			params = append(params, tmp)
		}else{
			params = append(params, "")
		}

		params = append(params, fmt.Sprintf("%d",review.Shop.No))
		params = append(params, review.Shop.Name)

		if len(review.Comments) > 0{
			params = append(params, review.Comments[0].Nickname)

			if regErr == nil{
				params = append(params, reg.ReplaceAllString(review.Comments[0].Contents, ""))
			}else{
				params = append(params, review.Comments[0].Contents)
			}

		}else{
			params = append(params, "사장님")
			params = append(params, "")
		}

		params = append(params, review.DateText)
		params = append(params, krToDate(review.DateText))

		_, err := cls.ExecDBbyParam2(query, params)
		if err != nil {
			lprintf(1, "[ERROR] cls.ExecDBbyParam error(%s) \n", err.Error())
			continue
		}
	}
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
