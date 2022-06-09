package datas

import (
	datasql "cashApi/query/datas"
	reviewsql "cashApi/query/reviews"
	"cashApi/src/controller"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	// login 및 기본

	"cashApi/src/controller/cls"
	"github.com/labstack/echo/v4"
)

/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
// darayo printf(로그레벨, 요청 컨텍스트, format, arg) => 무엇을(서비스, 요청), 어떻게(input), 왜(원인,조치)
var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func BaseUrl(c echo.Context) error {
	return c.JSONP(http.StatusOK, "", "homes")
}

func GetWeekData(c echo.Context) error {

	dprintf(4, c, "call GetWeekData\n")

	params := cls.GetParamJsonMap(c)
	// 1주간
	dt := time.Now().AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	endWeek := cls.GetEndOfWeek(dt)
	wStartDt := startWeek.Format("20060102")
	wEndDt := endWeek.Format("20060102")

	params["startDt"] = wStartDt
	params["endDt"] = wEndDt

	// 주간 매출 분석 (요일별 날자 입력, 지난 7일)
	weekSalesData, err := cls.GetSelectData(datasql.SelectWeekCash, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 주간 건수 분석 (요일별 날자 입력, 지난 7일)
	weekSalesCnt, err := cls.GetSelectData(datasql.SelectWeekCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난주 취소 분석
	cancleList, err := cls.GetSelectData(datasql.SelectLastCancleList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 가맹점 rating, keyword
	reviewOption, err := cls.GetSelectData2(datasql.SelectReivewOption, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난주 리뷰 분석
	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if len(deliveryInfo) == 0 {
		params["baeminId"] = "baeminId"
		params["yogiyoId"] = "yogiyoId"
		params["naverId"] = "naverId"
		params["coupangId"] = "coupangId"
	}else{
		if len(deliveryInfo[0]["baemin_id"]) == 0{
			params["baeminId"] = "baeminId"
		}else{
			params["baeminId"] = deliveryInfo[0]["baemin_id"]
		}

		if len(deliveryInfo[0]["yogiyo_id"]) == 0{
			params["yogiyoId"] = "yogiyoId"
		}else{
			params["yogiyoId"] = deliveryInfo[0]["yogiyo_id"]
		}

		if len(deliveryInfo[0]["naver_id"]) == 0{
			params["naverId"] = "naverId"
		}else{
			params["naverId"] = deliveryInfo[0]["naver_id"]
		}

		if len(deliveryInfo[0]["coupang_id"]) == 0 {
			params["coupangId"] = "coupangId"
		} else {
			params["coupangId"] = deliveryInfo[0]["coupang_id"]
		}
	}

	reviews, err := cls.GetSelectData2(datasql.SelectReviews, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 4주간
	params["startDt"] = endWeek.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = endWeek.Format("20060102")
	// 고객님 방문 분석 한달
	personData, err := cls.GetSelectData(datasql.SelectWeekPersonVisit, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 객 단가 예측
	personPrice, err := cls.GetSelectData(datasql.SelectAverageRevenuePerUser, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 주간 매출 분석 한달간 달아요 TIP
	dayAnalystic, err := cls.GetSelectDataUsingJson(datasql.SelectWeekAvgTime1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 이번달
	now := time.Now()
	params["bsDt"] = now.Format("200601")
	// web view 월간 컨텐츠
	webView, err := cls.GetSelectData2(datasql.SelectWebViewWeek, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if len(webView) == 0{
		webView, _ = cls.GetSelectData2(datasql.SelectWebViewWeekDefault, params)
	}

	/*
	// 주말 평일 비교
	dayAnalystic, err := cls.GetSelectData(datasql.SelectWeekdayAnalystic2, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 주말 바쁜 시간 도출
	busyTimeData, err := cls.GetSelectData(datasql.SelectWeeBusyTime, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 주중 바쁜 시간 도출
	busyTimeDataWork, err := cls.GetSelectData(datasql.SelectWeeBusyTimeWork, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	data := make(map[string]interface{})

	var avgAmt, maxAmt, minAmt int
	var minWeek, maxWeek, totCnt, avgCnt string
	if len(weekSalesData) > 0{
		avgAmt, _ = strconv.Atoi(weekSalesData[0]["avgAmt"])
		maxAmt, _ = strconv.Atoi(weekSalesData[0]["maxAmt"])
		minAmt, _ = strconv.Atoi(weekSalesData[0]["minAmt"])


		if weekSalesData[0]["minWeek"] != "%!s(<nil>)" {
			minWeek = weekSalesData[0]["minWeek"]
		}

		if weekSalesData[0]["maxWeek"] != "%!s(<nil>)" {
			maxWeek = weekSalesData[0]["maxWeek"]
		}

		if weekSalesData[0]["tot_cnt"] != "%!s(<nil>)" {
			totCnt = weekSalesData[0]["tot_cnt"]
		}

		if weekSalesData[0]["avg_cnt"] != "%!s(<nil>)" {
			avgCnt = weekSalesData[0]["avg_cnt"]
		}
	}

	var minh, mint, avgh, avgt int
	if maxAmt > 0{
		minh = int(((float64(minAmt) / float64(maxAmt)) * float64(100)))
		mint = 100 - minh

		avgh = int(((float64(avgAmt) / float64(maxAmt)) * float64(100)))
		avgt = 100 - avgh
	}

	data["avgAmt"] = avgAmt
	data["maxAmt"] = maxAmt
	data["minAmt"] = minAmt
	//data["maxWeek"] = dayNames(maxWeek)
	//data["minWeek"] = dayNames(minWeek)
	data["maxWeek"] = maxWeek
	data["minWeek"] = minWeek
	data["minh"] = minh
	data["mint"] = mint
	data["avgh"] = avgh
	data["avgt"] = avgt

	// 주간 건수 분석
	weekCnt := make(map[string]interface{})
	weekCnt["tot_cnt"] = totCnt
	weekCnt["avg_cnt"] = avgCnt
	tot_cnt, _ := strconv.Atoi(totCnt)
	var best int
	// 주간 건수 분석 베스트 1,2,3
	for i:=0; i<len(weekSalesCnt); i++{
		t0003,_ := strconv.Atoi(weekSalesCnt[i]["t0003"])
		t0306,_ := strconv.Atoi(weekSalesCnt[i]["t0306"])
		t0609,_ := strconv.Atoi(weekSalesCnt[i]["t0609"])
		t0912,_ := strconv.Atoi(weekSalesCnt[i]["t0912"])
		t1215,_ := strconv.Atoi(weekSalesCnt[i]["t1215"])
		t1518,_ := strconv.Atoi(weekSalesCnt[i]["t1518"])
		t1821,_ := strconv.Atoi(weekSalesCnt[i]["t1821"])
		t2124,_ := strconv.Atoi(weekSalesCnt[i]["t2124"])

		tTime := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}
		tm,cnt := FindBusyTime2(tTime)
		if best < cnt{
			best = cnt
		}

		weekCnt[fmt.Sprintf("bset%dtr_tm", i)] = fmt.Sprintf("%s시", tm)
		weekCnt[fmt.Sprintf("bset%dcnt", i)] = cnt
		weekCnt[fmt.Sprintf("bset%ddayName", i)] = weekSalesCnt[i]["day_name"]

		dayTotCnt, _ := strconv.Atoi(weekSalesCnt[i]["tot_cnt"])

		bestCntp := int(((float64(dayTotCnt) / float64(tot_cnt)) * float64(100)))
		weekCnt[fmt.Sprintf("bset%dcntp", i)] = bestCntp
	}
	weekCnt["len"] = len(weekSalesCnt)
	weekCnt["best"] = best
	data["weekCnt"] = weekCnt

	// 주간 매출 분석 제목
	var startMonth, endMonth string
	if weekSalesData[0]["minDt"] != "%!s(<nil>)" {
		startMonth = weekSalesData[0]["minDt"]
	}
	if weekSalesData[0]["maxDt"] != "%!s(<nil>)" {
		endMonth = weekSalesData[0]["maxDt"]
	}

	if len(startMonth) > 0 && len(endMonth) > 0{
		data["selectWeekTitle"] = fmt.Sprintf("주간 매출 (%s ~ %s)", dayFormat(wStartDt), dayFormat(wEndDt))
	}else{
		data["selectWeekTitle"] = fmt.Sprintf("주간 매출 (데이터 수집 중)")
	}

	var visitTotal, visit1, visit23, visit4 int
	var personPrices string

	if len(personData) > 0{
		visitTotal, _ = strconv.Atoi(personData[0]["visitTotal"])
		visit1, _ = strconv.Atoi(personData[0]["visit1"])
		visit23, _ = strconv.Atoi(personData[0]["visit23"])
		visit4, _ = strconv.Atoi(personData[0]["visit4"])

		if personPrice[0]["arpu"] != "%!s(<nil>)" {
			personPrices = personPrice[0]["arpu"]
		}
	}

	var visit1p, visit23p, visit4p int
	if visitTotal > 0{
		visit1p = int(((float64(visit1) / float64(visitTotal)) * float64(100)))
		visit23p = int(((float64(visit23) / float64(visitTotal)) * float64(100)))
		visit4p = int(((float64(visit4) / float64(visitTotal)) * float64(100)))
	}

	data["visitTotal"] = visitTotal
	data["visit1"] = visit1
	data["visit1p"] = visit1p
	data["visit23"] = visit23
	data["visit23p"] = visit23p
	data["visit4"] = visit4
	data["visit4p"] = visit4p
	data["personPrice"] = personPrices

	// 달아요팁
	/*
	dayAnalysticTotal, _ := strconv.Atoi(dayAnalystic[0]["total"])
	holy, _ := strconv.Atoi(dayAnalystic[0]["holy"])
	work, _ := strconv.Atoi(dayAnalystic[0]["work"])
	holyp := int(((float64(holy) / float64(dayAnalysticTotal)) * float64(100)))
	workp := int(((float64(work) / float64(dayAnalysticTotal)) * float64(100)))

	darayoTipMsg := ""
	darayoTipP := ""
	if holy > work {
		darayoTipP = strconv.Itoa(holyp - workp)
		darayoTipMsg = "주말에는 매출이 " + darayoTipP + "% 많으신 편이에요."
	} else {
		darayoTipP = strconv.Itoa(workp - holyp)
		darayoTipMsg = "평일에는 매출이 " + darayoTipP + "% 많으신 편이에요."
	}

	data["darayoTipMsg"] = darayoTipMsg

	data["holyBusyTime"] = busyTimeData[0]["atTime"]
	data["workBusyTime"] = busyTimeDataWork[0]["atTime"]
	 */

	// 요일 분석 팁
	//var weekAnalystic []map[string]interface{}
	var holySum, holy0003, holy0306, holy0609, holy0912, holy1215, holy1518, holy1821, holy2124 int
	var workSum, work0003, work0306, work0609, work0912, work1215, work1518, work1821, work2124 int
	for _, day := range dayAnalystic {
		week := make(map[string]interface{})
		week["trDt"] = day["trDt"]
		week["week"] = day["week"]
		week["weekNm"] = day["weekNm"]
		week["weekEnd"] = day["weekEnd"]

		totSum, _ := strconv.Atoi(day["totSum"])
		t0003, _ := strconv.Atoi(day["t0003"])
		t0306, _ := strconv.Atoi(day["t0306"])
		t0609, _ := strconv.Atoi(day["t0609"])
		t0912, _ := strconv.Atoi(day["t0912"])
		t1215, _ := strconv.Atoi(day["t1215"])
		t1518, _ := strconv.Atoi(day["t1518"])
		t1821, _ := strconv.Atoi(day["t1821"])
		t2124, _ := strconv.Atoi(day["t2124"])

		week["totSum"] = totSum
		week["t0003"] = t0003
		week["t0306"] = t0306
		week["t0609"] = t0609
		week["t0912"] = t0912
		week["t1215"] = t1215
		week["t1518"] = t1518
		week["t1821"] = t1821
		week["t2124"] = t2124

		if day["weekEnd"] == "HD" {
			holySum = holySum + totSum
			holy0003 = holy0003 + t0003
			holy0306 = holy0306 + t0306
			holy0609 = holy0609 + t0609
			holy0912 = holy0912 + t0912
			holy1215 = holy1215 + t1215
			holy1518 = holy1518 + t1518
			holy1821 = holy1821 + t1821
			holy2124 = holy2124 + t2124
		} else {
			workSum = workSum + totSum
			work0003 = work0003 + t0003
			work0306 = work0306 + t0306
			work0609 = work0609 + t0609
			work0912 = work0912 + t0912
			work1215 = work1215 + t1215
			work1518 = work1518 + t1518
			work1821 = work1821 + t1821
			work2124 = work2124 + t2124
		}
		//weekAnalystic = append(weekAnalystic, week)
	}

	//data["weekAnalystic"] = weekAnalystic

	// 달아요팁
	tipTotal := holySum + workSum
	holySump := int(((float64(holySum) / float64(tipTotal)) * float64(100)))
	workSump := int(((float64(workSum) / float64(tipTotal)) * float64(100)))

	darayoTipMsg := ""
	darayoTipP := ""
	if holySum > workSum {
		darayoTipP = strconv.Itoa(holySump - workSump)
		if holySum-workSum < 10 {
			darayoTipMsg = "지난 4주간 분석 결과 주말과 평일의 매출 차이가 많지 않아요.<br/><strong class='bl'>주말 매출 평균 " + darayoTipP + "% 정도 많았어요.</strong>"
		} else {
			darayoTipMsg = "지난 4주간 분석 결과 <strong class='bl'>주말 매출이 평균" + darayoTipP + "% 많았어요.</strong><br/>주말에는 영업 준비를 조금 더 많이 해보세요."
		}
	} else {
		darayoTipP = strconv.Itoa(workSump - holySump)
		if workSum-holySum < 10 {
			darayoTipMsg = "지난 4주간 분석 결과 주말과 평일의 매출 차이가 많지 않아요.<br/><strong class='bl'>평일 매출 평균 " + darayoTipP + "% 정도 많았어요.</strong>"
		} else {
			darayoTipMsg = "지난 4주간 분석 결과 <strong class='bl'>평일 매출이 평균" + darayoTipP + "% 많았어요.</strong><br/>평일 고객을 위해 다양한 준비를 해보세요."
		}
	}
	pageTip := make(map[string]interface{})
	pageTip["darayoTipMsg"] = darayoTipMsg

	holyTime := []int{holy0003, holy0306, holy0609, holy0912, holy1215, holy1518, holy1821, holy2124}
	workTime := []int{work0003, work0306, work0609, work0912, work1215, work1518, work1821, work2124}
	pageTip["holyBusyMsg"] = "<strong class='bl'>주말</strong>에는 평균 <strong class='bl'>" + findBusyTime(holyTime) + "시</strong>가 가장 <strong class='bl'>매출이 많아요.</strong>"
	pageTip["workBusyMsg"] = "<strong class='bl'>평일</strong>에는 평균 <strong class='bl'>" + findBusyTime(workTime) + "시</strong>가 가장 <strong class='bl'>매출이 많아요.</strong>"

	data["darayoTip"] = pageTip

	cancle := make(map[string]interface{})
	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int

	if len(cancleList) > 0{
		cancle["result"] = "y"

		for _,v := range cancleList{

			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle ++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0{
				noCancle ++ // 미 승인 취소
				continue
			}

			tr,_ := strconv.Atoi(v["tr_tm"][:2])
			otr,_ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10{
				nightCancle ++ // 심야 취소
			}else if v["tr_dt"] != cancleAprv[0]["tr_dt"]{
				dayCancle ++ // 일 취소
			}else if tr - otr > 3{
				timeCancle ++ // 시간 취소
			}else{
				okCancle ++ // 결제 취소
			}

		}

	}else{
		cancle["result"] = "n"
	}

	cancle["okCancle"] = okCancle
	cancle["timeCancle"] = timeCancle
	cancle["dayCancle"] = dayCancle
	cancle["nightCancle"] = nightCancle
	cancle["noCancle"] = noCancle

	data["cancle"] = cancle

	review := make(map[string]interface{})
	var reviewKeywords []string
	var reviewPoint []float64

	if len(reviewOption) > 0{
		ratings := strings.Split(reviewOption[0]["rating"], ",")
		keywords := strings.Split(reviewOption[0]["keyword"], "|")

		for _,v := range ratings{
			rating, err := strconv.ParseFloat(v, 64)
			if err != nil{
				continue
			}
			reviewPoint = append(reviewPoint, rating)
		}

		for _,v := range keywords{
			reviewKeywords = append(reviewKeywords, v)
		}

	}else{
		reviewKeywords = append(reviewKeywords, "맛있어요")
		reviewPoint = append(reviewPoint, 1)
	}

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0{
		review["result"] = "y"
	}else{
		review["result"] = "n"
	}

	Loop1:
	for _,v := range reviews{

		for _,key := range reviewKeywords{ // 키워드 포함 리뷰
			if strings.Contains(v["content"], key){
				keywordReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		for _,key := range reviewPoint{ // 평점 포함 리뷰
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r == key{
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		/*
		r, _ := strconv.ParseFloat(v["rating"], 64)
		if r <= reviewPoint{ // 평점 1점 이하 리뷰
			lowReview ++
			rst, tot := CheckReivewer(v["member_no"],v["rating"])
			if rst == 1{
				newCustomerTotal += tot
				newCustomer ++
			}else if rst == 2{
				oldCustomerTotal += tot
				oldCustomer ++
			}
			continue
		}
		 */

		okReview ++
	}

	if newCustomer+oldCustomer > 0 {
		review["result2"] = "y"
	}else{
		review["result2"] = "n"
	}

	review["okReview"] = okReview
	review["lowReview"] = lowReview
	review["keywordReview"] = keywordReview
	review["newCustomer"] = newCustomer
	review["oldCustomer"] = oldCustomer

	if newCustomer == 0{
		newCustomer = 1
	}

	if oldCustomer == 0{
		oldCustomer = 1
	}

	review["newCustomerTotal"] = strconv.FormatFloat(newCustomerTotal/float64(newCustomer), 'f', 2, 64)
	review["oldCustomerTotal"] = strconv.FormatFloat(oldCustomerTotal/float64(oldCustomer), 'f', 2, 64)

	/*
	var keyword string
	for _,key := range reviewKeywords{
		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0{
		//review["keyword"] = fmt.Sprintf("%s 키워드로 검색된 리뷰", keyword[:len(keyword)-1])
		review["keyword"] = fmt.Sprintf("키워드 : %s", keyword[:len(keyword)-1])
	}else{
		review["keyword"] = "설정된 키워드가 없습니다"
	}

	data["review"] = review
	 */

	var keyword string
	//var keyCnt int
	for _,key := range reviewKeywords{
		//keyCnt += utf8.RuneCountInString(key)
		//if keyCnt > 8{
		//	keyword = fmt.Sprintf("%s...",keyword[:len(keyword)-1])
		//	break
		//}

		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0{
		review["keyword"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	}else{
		review["keyword"] = "0"
	}

	data["review"] = review

	// 파트너 가입 여부
	data["partnerYN"] = "y"

	// web view content
	webViewContent := make(map[string]interface{})
	for _,v := range webView{
		webViewContent[fmt.Sprintf("conTitle%s", v["position"])] = v["title"]
		webViewContent[fmt.Sprintf("conBody%s", v["position"])] = v["content"]
	}
	data["webView"] = webViewContent

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

func CheckReivewer(customerId, rating string)(int, float64){

	if len(customerId) == 0{
		return -1,0
	}

	var params map[string]string
	params = make(map[string]string)
	params["customerId"] = customerId

	customer, err := cls.GetSelectData2(datasql.SelectBaeminCustomerCount, params)
	if err != nil {
		return -1,0
	}

	cnt,_ := strconv.Atoi(customer[0]["cnt"])
	if cnt > 1 { // 단골(재방문)
		tot, _ := strconv.ParseFloat(customer[0]["avg"], 64)
		return 2, tot
	}

	// 신규
	newCustomer, err := cls.GetSelectData2(datasql.SelectBaeminCustomerInfo, params)
	if err != nil {
		return -1,0
	}

	if len(newCustomer) > 0{
		tot, _ := strconv.ParseFloat(newCustomer[0]["avg"], 64)
		return 1, tot
	}

	tot, _ := strconv.ParseFloat(rating, 64)
	return 1, tot

/*
	rst, b := tracking.GetBaeminCompInfo(customerId, "50")
	if rst < 0{
		return -1,0
	}

	var total float64

	for _,v := range b.Data.Reviews{
		total += v.Rating
	}

	return 1, total/float64(len(b.Data.Reviews))
 */

}

func GetPosWeekData(c echo.Context) error {

	dprintf(4, c, "call GetPosWeekData\n")

	params := cls.GetParamJsonMap(c)
	t := time.Now()

	/* setting param */
	// 지난주 1주간 일~토
	dt := t.AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	endWeek := cls.GetEndOfWeek(dt)

	// 주간 12주 비교
	for i:=0; i<12; i++{
		//fmt.Printf("%s\t%s\n", startWeek.AddDate(0,0,-7*i).Format("20060102"), endWeek.AddDate(0,0,-7*i).Format("20060102"))

		startDt := startWeek.AddDate(0,0,-7*i)

		params[fmt.Sprintf("%dstartDt",i)] = startDt.Format("20060102")
		params[fmt.Sprintf("%dendDt",i)] = endWeek.AddDate(0,0,-7*i).Format("20060102")
		params[fmt.Sprintf("%dday1",i)] = fmt.Sprintf("%s01", startDt.Format("200601"))
	}

	wStartDt := startWeek.Format("20060102")
	wEndDt := endWeek.Format("20060102")
	params["startDt"] = wStartDt
	params["endDt"] = wEndDt
	//params["lStartDt"] = wStartDt
	//params["lEndDt"] = wEndDt

	// 어제부터 1주간
	//yStartDt := t.AddDate(0,0,-7)
	//yEndDt := t.AddDate(0,0,-1)
	//params["yStartDt"] = yStartDt.Format("20060102")
	//params["yEndDt"] = yEndDt.Format("20060102")

	// 단골 고객 비율 3개월
	//mDt1 := t.AddDate(0,-1,0)
	//mDt2 := t.AddDate(0,-2,0)
	//mDt3 := t.AddDate(0,-3,0)
	//params["mStartDt1"] = mDt1.Format("200601")
	//params["mStartDt2"] = mDt2.Format("200601")
	//params["mStartDt3"] = mDt3.Format("200601")

	// 지난주 취소 분석
	cancleList, err := cls.GetSelectData(datasql.SelectLastCancleList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 12주간 매출 비교
	allWeekAmt, err := cls.GetSelectData(datasql.Select12WeekAmt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}


	/* run sql */
	// 지난 주 매출 분석
	/*
	weekSalesData, err := cls.GetSelectData(datasql.SelectLastWeekAprv, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난 7일 매출 분석
	weekYeasterSalesData, err := cls.GetSelectData(datasql.SelectYeasterWeekAprv, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	 */

	// 지난주 리뷰 분석
	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 가맹점 rating, keyword
	reviewOption, err := cls.GetSelectData2(datasql.SelectReivewOption, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if len(deliveryInfo) == 0 {
		params["baeminId"] = "baeminId"
		params["yogiyoId"] = "yogiyoId"
		params["naverId"] = "naverId"
	}else{
		if len(deliveryInfo[0]["baemin_id"]) == 0{
			params["baeminId"] = "baeminId"
		}else{
			params["baeminId"] = deliveryInfo[0]["baemin_id"]
		}

		if len(deliveryInfo[0]["yogiyo_id"]) == 0{
			params["yogiyoId"] = "yogiyoId"
		}else{
			params["yogiyoId"] = deliveryInfo[0]["yogiyo_id"]
		}

		if len(deliveryInfo[0]["naver_id"]) == 0{
			params["naverId"] = "naverId"
		}else{
			params["naverId"] = deliveryInfo[0]["naver_id"]
		}
	}

	reviews, err := cls.GetSelectData2(datasql.SelectReviews, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난 3개월 단골 고객 비율 분석
	/*
	threeMonthSalesData, err := cls.GetSelectData(datasql.SelectThreeMonthAprv, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	/*
	threeMonthVisitData, err := cls.GetSelectData(datasql.SelectThreeMonthVisitAprv, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	threeMonthVisitData2, err := cls.GetSelectData(datasql.SelectThreeMonthVisitAprv2, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난 주 고객 매출분석 TIP (WD, HD)
	weekSalesDataTip, err := cls.GetSelectData(datasql.SelectLastWeekAprvTip, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	// 이번달
	params["bsDt"] = t.Format("200601")
	// web view 월간 컨텐츠
	webView, err := cls.GetSelectData2(datasql.SelectWebViewWeek, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}


	/* setting response */
	data := make(map[string]interface{})

	// 12주간 매출 비교
	allWeek := make(map[string]interface{})
	var weekMax, weekMin int
	var weekMaxIndex, weekMinIndex string
	var weekIndexArray []string
	var weekAmtArray []int

	for idx,v := range allWeekAmt{

		amt, _ := strconv.Atoi(v["amt"])

		//allWeek[fmt.Sprintf("week%dAmt", idx)] = amt
		weekAmtArray = append(weekAmtArray, amt)
		var dayIndex string

		if idx == 0{

			weekMin = amt

			dayIndex = "지난주 ("

			if wStartDt[4] == '0'{
				dayIndex += fmt.Sprintf("%s.%s",string(wStartDt[5]), wStartDt[6:])
			}else{
				dayIndex += fmt.Sprintf("%s.%s",string(wStartDt[4:6]), wStartDt[6:])
			}

			dayIndex += fmt.Sprintf(" ~ ")

			if wEndDt[4] == '0'{
				dayIndex += fmt.Sprintf("%s.%s",string(wEndDt[5]), wEndDt[6:])
			}else{
				dayIndex += fmt.Sprintf("%s.%s",string(wEndDt[4:6]), wEndDt[6:])
			}

			dayIndex += fmt.Sprintf(")")
		}else{
			dayIndex = v["dayIndex"]

			if dayIndex[0] == '0'{
				dayIndex = dayIndex[1:]
			}
		}

		weekIndexArray = append(weekIndexArray, dayIndex)

		if weekMax < amt{
			weekMaxIndex = fmt.Sprintf("week%dState", idx)
			weekMax = amt
		}

		if weekMin > amt{
			weekMinIndex = fmt.Sprintf("week%dState", idx)
			weekMin = amt
		}
	}

	allWeek["weekMaxIndex"] = weekMaxIndex
	allWeek["weekMinIndex"] = weekMinIndex
	allWeek["weekMax"] = weekMax

	allWeek["len"] = len(allWeekAmt)
	allWeek["weekAmtArray"] = weekAmtArray
	allWeek["weekIndexArray"] = weekIndexArray

	data["allWeek"] = allWeek

	// 지난 주 매출 분석
	/*
	mWeekSalesData := make(map[string]interface{})
	var sumWeekSalesData int

	for i:=0; i<len(weekSalesData); i++{
		mWeekSalesData["w"+weekSalesData[i]["day_index"]+"Cnt"] = weekSalesData[i]["tot_cnt"]
		mWeekSalesData["w"+weekSalesData[i]["day_index"]+"Amt"] = weekSalesData[i]["tot_amt"]
		mWeekSalesData["w"+weekSalesData[i]["day_index"]+"Dt"] = weekSalesData[i]["bs_dt"]
		mWeekSalesData["w"+weekSalesData[i]["day_index"]+"Kr"] = weekSalesData[i]["day_name"]
		tmpAmt, _ := strconv.Atoi(weekSalesData[i]["tot_amt"])
		sumWeekSalesData += tmpAmt
	}

	maxWeekSalesData, _ := strconv.Atoi(weekSalesData[0]["tot_amt"])
	minWeekSalesData, _ := strconv.Atoi(weekSalesData[len(weekSalesData)-1]["tot_amt"])
	mWeekSalesData["maxAmt"] = maxWeekSalesData
	mWeekSalesData["minAmt"] = minWeekSalesData
	mWeekSalesData["maxCnt"] = weekSalesData[0]["day_index"]
	mWeekSalesData["minCnt"] = weekSalesData[len(weekSalesData)-1]["day_index"]
	mWeekSalesData["startDt"] = mWeekSalesData["w1Dt"]
	mWeekSalesData["endDt"] = mWeekSalesData["w7Dt"]
	mWeekSalesData["avgAmt"] = int(sumWeekSalesData/len(weekSalesData))
	mWeekSalesData["totAmt"] = sumWeekSalesData

	data["mWeekSalesData"] = mWeekSalesData
	 */

	// 지난 7일 매출 분석
	/*
	mWeekYeasterSalesData := make(map[string]interface{})
	var sumWeekYeasterSalesData, maxWeekYeasterSalesData, minWeekYeasterSalesData int
	var maxIdx, minIdx int

	for i:=0; i<len(weekYeasterSalesData); i++{
		key := strconv.Itoa(i+1)
		mWeekYeasterSalesData["w"+key+"Cnt"] = weekYeasterSalesData[i]["tot_cnt"]
		mWeekYeasterSalesData["w"+key+"Amt"] = weekYeasterSalesData[i]["tot_amt"]
		mWeekYeasterSalesData["w"+key+"Dt"] = weekYeasterSalesData[i]["bs_dt"]
		mWeekYeasterSalesData["w"+key+"Kr"] = weekYeasterSalesData[i]["day_name"]
		tmpAmt, _ := strconv.Atoi(weekYeasterSalesData[i]["tot_amt"])
		sumWeekYeasterSalesData += tmpAmt

		if i==0	{
			minWeekYeasterSalesData = tmpAmt
		}

		if maxWeekYeasterSalesData < tmpAmt {
			maxWeekYeasterSalesData = tmpAmt
			maxIdx = i
		}

		if minWeekYeasterSalesData > tmpAmt {
			minWeekYeasterSalesData = tmpAmt
			minIdx = i
		}
	}

	// 10시 이전에 어제 데이터 수집 불가능
	if len(weekYeasterSalesData) == 6{
		mWeekYeasterSalesData["w7Amt"] = "0"
		mWeekYeasterSalesData["w7Cnt"] = "0"
		mWeekYeasterSalesData["w7Dt"] = ""
		mWeekYeasterSalesData["w7Kr"] = "데이터 수집중"
	}

	mWeekYeasterSalesData["maxAmt"] = maxWeekYeasterSalesData
	mWeekYeasterSalesData["minAmt"] = minWeekYeasterSalesData
	mWeekYeasterSalesData["maxCnt"] = maxIdx+1
	mWeekYeasterSalesData["minCnt"] = minIdx+1
	mWeekYeasterSalesData["startDt"] = weekYeasterSalesData[0]["bs_dt"]
	mWeekYeasterSalesData["endDt"] = weekYeasterSalesData[len(weekYeasterSalesData)-1]["bs_dt"]
	mWeekYeasterSalesData["avgAmt"] = int(sumWeekYeasterSalesData/len(weekYeasterSalesData))
	mWeekYeasterSalesData["totAmt"] = sumWeekYeasterSalesData

	data["mWeekYeasterSalesData"] = mWeekYeasterSalesData
	 */

	// 지난 3개월 단골 고객 비율 분석
	/*
	mThreeMonthSalesData := make(map[string]interface{})
	var threeMonthTotalPeople, threeMonthPayDay, threeMonthAvgAmt int
	var threeMonthVisit1, threeMonthVisit23, threeMonthVisit4 int

	for i:=0; i<len(threeMonthSalesData); i++{
		tmpTotCnt, _ := strconv.Atoi(threeMonthSalesData[i]["tot_cnt"])
		threeMonthTotalPeople += tmpTotCnt
		tmpDtCnt, _ := strconv.Atoi(threeMonthSalesData[i]["dt_cnt"])
		threeMonthPayDay += tmpDtCnt
		tmpAvgAmt, _ := strconv.Atoi(threeMonthSalesData[i]["avg_amt"])
		threeMonthAvgAmt += tmpAvgAmt
		tmpVisit1, _ := strconv.Atoi(threeMonthSalesData[i]["visit1"])
		threeMonthVisit1 += tmpVisit1
		tmpVisit23, _ := strconv.Atoi(threeMonthSalesData[i]["visit23"])
		threeMonthVisit23 += tmpVisit23
		tmpVisit4, _ := strconv.Atoi(threeMonthSalesData[i]["visit4"])
		threeMonthVisit4 += tmpVisit4
	}

	mThreeMonthSalesData["totalPeople"] = threeMonthTotalPeople
	mThreeMonthSalesData["onePeople"] = int(threeMonthTotalPeople/threeMonthPayDay)
	mThreeMonthSalesData["avgAmt"] = int(threeMonthAvgAmt/len(threeMonthSalesData))
	mThreeMonthSalesData["visit1"] = threeMonthVisit1
	mThreeMonthSalesData["visit1p"] = int(float64(threeMonthVisit1) / float64(threeMonthTotalPeople) * 100)
	mThreeMonthSalesData["visit23"] = threeMonthVisit23
	mThreeMonthSalesData["visit23p"] = int(float64(threeMonthVisit23) / float64(threeMonthTotalPeople) * 100)
	mThreeMonthSalesData["visit4"] = threeMonthVisit4
	mThreeMonthSalesData["visit4p"] = int(float64(threeMonthVisit4) / float64(threeMonthTotalPeople) * 100)

	data["mThreeMonthSalesData"] = mThreeMonthSalesData
	*/

	/*
	mThreeMonthVisitData := make(map[string]interface{})
	for i:=0; i<len(threeMonthVisitData); i++{
		mThreeMonthVisitData[fmt.Sprintf("visit%dDayName",i+1)] = threeMonthVisitData[i]["day_name"]

		t0003,_ := strconv.Atoi(threeMonthVisitData[i]["t0003"])
		t0306,_ := strconv.Atoi(threeMonthVisitData[i]["t0306"])
		t0609,_ := strconv.Atoi(threeMonthVisitData[i]["t0609"])
		t0912,_ := strconv.Atoi(threeMonthVisitData[i]["t0912"])
		t1215,_ := strconv.Atoi(threeMonthVisitData[i]["t1215"])
		t1518,_ := strconv.Atoi(threeMonthVisitData[i]["t1518"])
		t1821,_ := strconv.Atoi(threeMonthVisitData[i]["t1821"])
		t2124,_ := strconv.Atoi(threeMonthVisitData[i]["t2124"])

		tTime := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}
		tm,cnt := findBusyTime2(tTime)

		mThreeMonthVisitData[fmt.Sprintf("visit%dTotalCnt",i+1)] = threeMonthVisitData[i]["tot_cnt"]
		mThreeMonthVisitData[fmt.Sprintf("visit%dMent",i+1)] = fmt.Sprintf("%s시 : %d건",tm, cnt)
	}


	mThreeMonthVisitData["totCnt"] = threeMonthVisitData2[0]["tot_cnt"]
	mThreeMonthVisitData["totAvg"] = threeMonthVisitData2[0]["tot_avg"]

	data["mThreeMonthVisitData"] = mThreeMonthVisitData
	 */


	// 지난 주 고객 매출분석 TIP (WD, HD)
	/*
	mWeekSalesDataTip := make(map[string]interface{})

	wdAmt, _ := strconv.Atoi(weekSalesDataTip[0]["tot_amt"])
	hdAmt, _ := strconv.Atoi(weekSalesDataTip[1]["tot_amt"])
	totalAmt := wdAmt + hdAmt
	wdAmtP := int(float64(wdAmt) / float64(totalAmt) * 100)
	hdAmtP := 100 - wdAmtP

	wdCnt, _ := strconv.Atoi(weekSalesDataTip[0]["tot_cnt"])
	hdCnt, _ := strconv.Atoi(weekSalesDataTip[1]["tot_cnt"])
	totalCnt := wdCnt + hdCnt
	wdCntP := int(float64(wdCnt) / float64(totalCnt) * 100)
	hdCntP := 100 - wdCntP

	if wdAmt > hdAmt {
		mWeekSalesDataTip["amtTip"] = fmt.Sprintf("<strong class=\"bl\">평일 매출</strong>이 휴일 매출보다 <strong class=\"bl\">%d%% 많으신 편</strong>이에요.", wdAmtP-hdAmtP)
	}else{
		mWeekSalesDataTip["amtTip"] = fmt.Sprintf("<strong class=\"bl\">휴일 매출</strong>이 평일 매출보다 <strong class=\"bl\">%d%% 많으신 편</strong>이에요.", hdAmtP-wdAmtP)
	}

	if wdCnt > hdCnt {
		mWeekSalesDataTip["cntTip"] = fmt.Sprintf("고객님은 휴일 보다 <strong class=\"bl\">평일에 %d%% 더 많이 방문</strong>해요.", wdCntP - hdCntP)
	}else{
		mWeekSalesDataTip["cntTip"] = fmt.Sprintf("고객님은 평일 보다 <strong class=\"bl\">휴일에 %d%% 더 많이 방문</strong>해요.", hdCntP - wdCntP)
	}

	cnt03,_ := strconv.Atoi(weekSalesDataTip[2]["cnt03"])
	cnt36,_ := strconv.Atoi(weekSalesDataTip[2]["cnt36"])
	cnt69,_ := strconv.Atoi(weekSalesDataTip[2]["cnt69"])
	cnt912,_ := strconv.Atoi(weekSalesDataTip[2]["cnt912"])
	cnt1215,_ := strconv.Atoi(weekSalesDataTip[2]["cnt1215"])
	cnt1518,_ := strconv.Atoi(weekSalesDataTip[2]["cnt1518"])
	cnt1821,_ := strconv.Atoi(weekSalesDataTip[2]["cnt1821"])
	cnt2124,_ := strconv.Atoi(weekSalesDataTip[2]["cnt2124"])

	times := []string{"0시 ~ 3시", "3시 ~ 6시", "6시 ~ 9시", "9시 ~ 12시","12시 ~ 15시","15시 ~ 18시","18시 ~ 21시","21시 ~ 24시"}
	cnts := []int{cnt03, cnt36, cnt69, cnt912, cnt1215, cnt1518, cnt1821, cnt2124}
	timeIdx := sort.SearchInts(cnts, len(cnts))

	mWeekSalesDataTip["timeTip"] = fmt.Sprintf("보통 <strong class=\"bl\">%s</strong>사이에 가장 많이 방문해요.", times[timeIdx])
	mWeekSalesDataTip["tip"] = "<strong class=\"bl\">단골고객</strong>이 많이 방문하세요.<br />서비스를 통해 객단가를 높이시면 더욱 좋을 것 같아요 :)"


	data["mWeekSalesDataTip"] = mWeekSalesDataTip
	 */

	// web view content
	webViewContent := make(map[string]interface{})
	for _,v := range webView{
		webViewContent[fmt.Sprintf("conTitle%s", v["position"])] = v["title"]
		webViewContent[fmt.Sprintf("conBody%s", v["position"])] = v["content"]
	}
	data["webView"] = webViewContent

	// 지난주 취소 분석
	cancle := make(map[string]interface{})
	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int

	if len(cancleList) > 0{
		cancle["result"] = "y"

		for _,v := range cancleList{

			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle ++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0{
				noCancle ++ // 미 승인 취소
				continue
			}

			tr,_ := strconv.Atoi(v["tr_tm"][:2])
			otr,_ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10{
				nightCancle ++ // 심야 취소
			}else if v["tr_dt"] != cancleAprv[0]["tr_dt"]{
				dayCancle ++ // 일 취소
			}else if tr - otr > 3{
				timeCancle ++ // 시간 취소
			}else{
				okCancle ++ // 결제 취소
			}

		}

	}else{
		cancle["result"] = "n"
	}

	cancle["okCancle"] = okCancle
	cancle["timeCancle"] = timeCancle
	cancle["dayCancle"] = dayCancle
	cancle["nightCancle"] = nightCancle
	cancle["noCancle"] = noCancle

	data["cancle"] = cancle

	review := make(map[string]interface{})
	var reviewKeywords []string
	var reviewPoint []float64

	if len(reviewOption) > 0{
		ratings := strings.Split(reviewOption[0]["rating"], ",")
		keywords := strings.Split(reviewOption[0]["keyword"], "|")

		for _,v := range ratings{
			rating, err := strconv.ParseFloat(v, 64)
			if err != nil{
				continue
			}
			reviewPoint = append(reviewPoint, rating)
		}

		for _,v := range keywords{
			reviewKeywords = append(reviewKeywords, v)
		}

	}else{
		reviewKeywords = append(reviewKeywords, "맛있어요")
		reviewPoint = append(reviewPoint, 1)
	}

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0{
		review["result"] = "y"
	}else{
		review["result"] = "n"
	}

Loop1:
	for _,v := range reviews{

		for _,key := range reviewKeywords{ // 키워드 포함 리뷰
			if strings.Contains(v["content"], key){
				keywordReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		for _,key := range reviewPoint{ // 평점 포함 리뷰
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r == key{
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		/*
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r <= reviewPoint{ // 평점 1점 이하 리뷰
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue
			}
		*/

		okReview ++
	}

	if newCustomer+oldCustomer > 0 {
		review["result2"] = "y"
	}else{
		review["result2"] = "n"
	}

	review["okReview"] = okReview
	review["lowReview"] = lowReview
	review["keywordReview"] = keywordReview
	review["newCustomer"] = newCustomer
	review["oldCustomer"] = oldCustomer

	if newCustomer == 0{
		newCustomer = 1
	}

	if oldCustomer == 0{
		oldCustomer = 1
	}

	review["newCustomerTotal"] = strconv.FormatFloat(newCustomerTotal/float64(newCustomer), 'f', 2, 64)
	review["oldCustomerTotal"] = strconv.FormatFloat(oldCustomerTotal/float64(oldCustomer), 'f', 2, 64)

	var keyword string
	for _,key := range reviewKeywords{
		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0{
		review["keyword"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	}else{
		review["keyword"] = "0"
	}

	data["review"] = review

	// 파트너 가입 여부
	data["partnerYN"] = "y"

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

func GetPosMonthData(c echo.Context) error {

	dprintf(4, c, "call GetPosMonthData\n")
	params := cls.GetParamJsonMap(c)

	/* setting param */
	// 가입일 조회
	/*
	regDt, err := cls.GetSelectData(datasql.SelectRegistDate, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	 */

	now := time.Now()
	//thisYear := now.AddDate(0, 0, 0).Format("2006")
	//lastYear := now.AddDate(-1, 0, 0).Format("2006")

	/*
	// 월간분석 - 월 매출 분석(올해)
	if thisYear == regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 수집월
		params["startDt"] = regDt[0]["regDt"][:6]
	} else {
		// 작년 가입한 경우 1월
		params["startDt"] = thisYear + "01"
	}
	 */

	// 월간분석 - 월 매출 분석(올해)
	/*
	if thisYear == regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 수집월 (지난달 부터 시작)
		mon := regDt[0]["regDt"][4:6]
		monInt,_ :=strconv.Atoi(mon)
		monInt -= 1
		if monInt < 10{
			params["startDt"] = fmt.Sprintf("%s0%d", regDt[0]["regDt"][:4], monInt)
		}else{
			params["startDt"] = fmt.Sprintf("%s%d", regDt[0]["regDt"][:4], monInt)
		}
	} else {
		// 작년 가입한 경우 1월, 혹은 프리미엄
		params["startDt"] = thisYear + "01"
	}
	 */

	// 저번달
	date := now.AddDate(0, 0, -now.Day())
	params["lastDt"] = date.Format("200601")
	params["bsDt"] = date.Format("200601")

	// 월간 12달 비교
	for i:=0; i<12; i++{
		startDt := date.AddDate(0,-1*i,-3)

		params[fmt.Sprintf("%dbsDt",i)] = startDt.Format("200601")
		params[fmt.Sprintf("%dbsDt1",i)] = fmt.Sprintf("%s01", startDt.Format("200601"))
	}

	// 월간분석 - 월 매출 분석(올해)
	/*
	thisSalesData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthAprv, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 월간분석 - 월 매출 분석(작년)
	var lastSalesData []map[string]string
	if thisYear != regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 데이터 없음

		if lastYear == regDt[0]["regDt"][:4] {
			// 작년 가입한 경우 데이터 수집 시작 월
			params["startDt"] = regDt[0]["regDt"][:6]
		} else {
			// 재작년 가입한 경우 작년 1월
			params["startDt"] = lastYear + "01"
		}

		params["endDt"] = lastYear + "12"
		lastSalesData, err = cls.GetSelectDataUsingJson(datasql.SelectMonthAprv, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}
	}
	 */

	// 12주간 매출 비교
	allMonthAmt, err := cls.GetSelectData(datasql.Select12MonthAmtCard, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난달 요일별 평균 매출 분석
	lastDaySalesData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthDayAprv2, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 달아요팁
	darayoTipData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthDarayoTip1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난달 취소 분석
	cancleList, err := cls.GetSelectData(datasql.SelectLastMonthCancleList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난달 리뷰 분석
	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 가맹점 rating, keyword
	reviewOption, err := cls.GetSelectData2(datasql.SelectReivewOption, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if len(deliveryInfo) == 0 {
		params["baeminId"] = "baeminId"
		params["yogiyoId"] = "yogiyoId"
		params["naverId"] = "naverId"
	}else{
		if len(deliveryInfo[0]["baemin_id"]) == 0{
			params["baeminId"] = "baeminId"
		}else{
			params["baeminId"] = deliveryInfo[0]["baemin_id"]
		}

		if len(deliveryInfo[0]["yogiyo_id"]) == 0{
			params["yogiyoId"] = "yogiyoId"
		}else{
			params["yogiyoId"] = deliveryInfo[0]["yogiyo_id"]
		}

		if len(deliveryInfo[0]["naver_id"]) == 0{
			params["naverId"] = "naverId"
		}else{
			params["naverId"] = deliveryInfo[0]["naver_id"]
		}
	}

	reviews, err := cls.GetSelectData2(datasql.SelectMonthReviews, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 고객 방문 분석
	/*
	1. 결제가 가장 많이 발생한 날자 3
	2. 매출이 가장 많이 발생한 날자 3
	3. 결제 단가가 가장 높았던 날자 3
	4. 매출이 가장 적었던 날자 3
	 */
	monthTmCnt, err := cls.GetSelectDataUsingJson(datasql.SelectLastMonthAnal2, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 월 결제 금액, 결제 건수, 월 객 단가
	//monthTotalAmt, err := cls.GetSelectDataUsingJson(datasql.SelectMonthTotalAmt, params, c)
	//if err != nil {
	//	return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	//}


	// 주변 상점 -> 임시 데이터
	/*
	aroudVisitData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthAroundStroeVisit, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	// 이번달
	params["bsDt"] = now.Format("200601")
	// web view 월간 컨텐츠
	webView, err := cls.GetSelectData2(datasql.SelectWebViewWeek, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 금결원 가입여부
	/*
	kftcEnroll, err := cls.GetSelectDataUsingJson(datasql.SelectKFTCEnroll, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	/* setting response */
	data := make(map[string]interface{})

	/*
	mThisSalesData := make(map[string]interface{})
	var tsTotal, tsMax, tsMin int
	var tsMaxInedx, tsMinIndex int
	for i:=0; i<len(thisSalesData); i++{
		//mThisSalesData[fmt.Sprintf("dt%d",i+1)] = thisSalesData[i]["bs_dt"]
		tsAmt,_ := strconv.Atoi(thisSalesData[i]["tot_amt"])
		//mThisSalesData[fmt.Sprintf("amt%d",i+1)] = tsAmt

		if i==0{
			tsMax = tsAmt
			tsMin = tsAmt
			tsMaxInedx = i+1
			tsMinIndex = i+1
		}

		if tsMax < tsAmt{
			tsMax = tsAmt
			tsMaxInedx = i+1
		}

		if tsMin > tsAmt{
			tsMin = tsAmt
			tsMinIndex = i+1
		}

		tsTotal += tsAmt
	}

	mThisSalesData["totAmt"] = tsTotal
	mThisSalesData["avgAmt"] = int(tsTotal / len(thisSalesData))
	mThisSalesData["maxMonth"] = tsMaxInedx
	mThisSalesData["minMonth"] = tsMinIndex
	mThisSalesData["maxAmt"] = tsMax
	mThisSalesData["minAmt"] = tsMin
	mThisSalesData["monthCnt"] = len(thisSalesData)
	mThisSalesData["startMonth"] = thisSalesData[0]["bs_dt"]
	mThisSalesData["endMonth"] = thisSalesData[len(thisSalesData)-1]["bs_dt"]
	mThisSalesData["list"] = thisSalesData
	mThisSalesData["yaer"] = thisYear
	data["mThisSalesData"] = mThisSalesData
	 */

	/*
	if len(lastSalesData) > 0 {
		mLastSalesData := make(map[string]interface{})
		var lsTotal, lsMax, lsMin int
		var lsMaxInedx, lsMinIndex string

		for i:=0; i<len(lastSalesData); i++{
			//mLastSalesData[fmt.Sprintf("dt%s",lastSalesData[i]["bs_dt"])] = lastSalesData[i]["bs_dt"]
			tsAmt,_ := strconv.Atoi(lastSalesData[i]["tot_amt"])
			//mLastSalesData[fmt.Sprintf("amt%s",lastSalesData[i]["bs_dt"])] = tsAmt

			if i==0{
				lsMax = tsAmt
				lsMin = tsAmt
				lsMaxInedx = lastSalesData[i]["bs_dt"]
				lsMinIndex = lastSalesData[i]["bs_dt"]
			}

			if lsMax < tsAmt {
				lsMax = tsAmt
				lsMaxInedx = lastSalesData[i]["bs_dt"]
			}

			if lsMin > tsAmt {
				lsMin = tsAmt
				lsMinIndex = lastSalesData[i]["bs_dt"]
			}

			lsTotal += tsAmt
		}

		mLastSalesData["totAmt"] = lsTotal
		mLastSalesData["avgAmt"] = int(lsTotal / len(lastSalesData))
		mLastSalesData["maxMonth"] = lsMaxInedx
		mLastSalesData["minMonth"] = lsMinIndex
		mLastSalesData["maxAmt"] = lsMax
		mLastSalesData["minAmt"] = lsMin
		mLastSalesData["monthCnt"] = len(lastSalesData)
		mLastSalesData["startMonth"] = lastSalesData[0]["bs_dt"]
		mLastSalesData["endMonth"] = lastSalesData[len(lastSalesData)-1]["bs_dt"]
		mLastSalesData["list"] = lastSalesData
		mLastSalesData["yaer"] = lastYear
		data["mLastSalesData"] = mLastSalesData
	}else{
		data["mLastSalesData"] = []string{}
	}
	 */

	// 12월간 매출 비교
	allMonth := make(map[string]interface{})
	var monthMax, monthMin int
	var monthMaxIndex, monthMinIndex string
	var monthIndexArray []string
	var monthAmtArray []int

	for idx,v := range allMonthAmt{

		amt, _ := strconv.Atoi(v["amt"])

		//allMonth[fmt.Sprintf("month%dAmt", idx)] = amt
		monthAmtArray = append(monthAmtArray, amt)
		dayIndex := v["dayIndex"]

		if idx == 0 {
			monthMin = amt
		}else if idx != len(allMonthAmt)-1{
			dayIndex = fmt.Sprintf("%s", strings.TrimSpace(dayIndex[7:]))
		}

		monthIndexArray = append(monthIndexArray, dayIndex)

		if monthMax < amt{
			monthMaxIndex = fmt.Sprintf("month%dState", idx)
			monthMax = amt
		}

		if monthMin > amt{
			monthMinIndex = fmt.Sprintf("month%dState", idx)
			monthMin = amt
		}
	}

	allMonth["monthMaxIndex"] = monthMaxIndex
	allMonth["monthMinIndex"] = monthMinIndex
	allMonth["monthMax"] = monthMax
	allMonth["len"] = len(allMonthAmt)
	allMonth["monthAmtArray"] = monthAmtArray
	allMonth["monthIndexArray"] = monthIndexArray

	data["allMonth"] = allMonth

	// 월 매출 분석 결과 & 매장 관리 TIP
	var holySum, workSum int
	var holyMax, workMax string
	for _, tipData := range darayoTipData {
		t0003, _ := strconv.Atoi(tipData["t0003"])
		t0306, _ := strconv.Atoi(tipData["t0306"])
		t0609, _ := strconv.Atoi(tipData["t0609"])
		t0912, _ := strconv.Atoi(tipData["t0912"])
		t1215, _ := strconv.Atoi(tipData["t1215"])
		t1518, _ := strconv.Atoi(tipData["t1518"])
		t1821, _ := strconv.Atoi(tipData["t1821"])
		t2124, _ := strconv.Atoi(tipData["t2124"])
		timeList := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}

		strconv.Atoi(tipData["weekEnd"])
		if tipData["weekEnd"] == "HD" {
			holySum, _ = strconv.Atoi(tipData["total"])
			holyMax = findBusyTime(timeList)
		} else {
			workSum, _ = strconv.Atoi(tipData["total"])
			workMax = findBusyTime(timeList)
		}
	}

	tipTotal := holySum + workSum
	holySump := int(((float64(holySum) / float64(tipTotal)) * float64(100)))
	workSump := int(((float64(workSum) / float64(tipTotal)) * float64(100)))

	darayoTipMsg := ""
	darayoTipP := ""
	if holySum > workSum {
		darayoTipP = strconv.Itoa(holySump - workSump)
		if holySum-workSum < 10 {
			darayoTipMsg = "지난 달 분석 결과 주말과 평일의 매출 차이가 많지 않아요.<br>주말 매출 평균 " + darayoTipP + "% 정도 많았어요."
		} else {
			darayoTipMsg = "지난 달 분석 결과 주말 매출이 평균 " + darayoTipP + "% 많았어요.<br>주말에는 영업 준비를 조금 더 많이 해보세요."
		}
	} else {
		darayoTipP = strconv.Itoa(workSump - holySump)
		if workSum-holySum < 10 {
			darayoTipMsg = "지난 달 분석 결과 주말과 평일의 매출 차이가 많지 않아요.<br>평일 매출 평균 " + darayoTipP + "% 정도 많았어요."
		} else {
			darayoTipMsg = "지난 달 분석 결과 평일 매출이 평균 " + darayoTipP + "% 많았어요.<br>평일 고객을 위해 다양한 준비를 해보세요."
		}
	}
	pageTip := make(map[string]interface{})
	pageTip["darayoTipMsg"] = darayoTipMsg
	pageTip["holyBusyMsg"] = "<strong class='bl'>주말에는 " + holyMax + "시에 매출이 가장 많아요.</strong>"
	pageTip["workBusyMsg"] = "<strong class='bl'>평일에는 " + workMax + "시에 매출이 가장 많아요.</strong>"


	var pageTipHdAmt, pageTipWdAmt int
	var pageTipHdCnt, pageTipWdCnt int
	var hCnt, wCnt int

	for i:=0; i<len(lastDaySalesData); i++{

		tipAmt, _ := strconv.Atoi(lastDaySalesData[i]["tot_amt"])
		tipCnt, _ := strconv.Atoi(lastDaySalesData[i]["cnt"])

		// 휴일
		if lastDaySalesData[i]["day_index"] == "1" || lastDaySalesData[i]["day_index"] == "7"{
			if tipCnt > 0{

				pageTipHdAmt += tipAmt
				pageTipHdCnt += tipCnt

				hCnt++
			}
		}else{
			if tipCnt > 0{

				pageTipWdAmt += tipAmt
				pageTipWdCnt += tipCnt

				wCnt++
			}
		}
	}
	pageTip["avgHdAmt"] = int(pageTipHdAmt / hCnt)
	pageTip["cntHd"] = pageTipHdCnt
	pageTip["avgWdAmt"] = int(pageTipWdAmt / wCnt)
	pageTip["cntWd"] = pageTipWdCnt

	data["darayoSaleTip"] = pageTip

	// 지난달 요일 별 평균 분석
	mLastDaySalesData := make(map[string]interface{})
	mLastDaySalesData["list"] = lastDaySalesData

	lastDayMonth := params["lastDt"][5:]
	if lastDayMonth[0] == '0'{
		mLastDaySalesData["month"] = lastDayMonth[1]
	}else{
		mLastDaySalesData["month"] = lastDayMonth
	}
	data["mLastDaySalesData"] = mLastDaySalesData

	// 고객 방문 분석
	/*
		1. 결제가 가장 많이 발생한 날자 3
		2. 매출이 가장 많이 발생한 날자 3
		3. 결제 단가가 가장 높았던 날자 3
		4. 매출이 가장 적었던 날자 3
	*/
	mVisitData := make(map[string]interface{})

	if len(monthTmCnt) > 2{
		mVisitData["cnt"] = monthTmCnt[:3]
	}else{
		mVisitData["cnt"] = "n"
	}

	if len(monthTmCnt) > 5{
		mVisitData["maxAmt"] = monthTmCnt[3:6]
	}else{
		mVisitData["maxAmt"] = "n"
	}

	if len(monthTmCnt) > 8{
		mVisitData["avgAmt"] = monthTmCnt[3:6]
	}else{
		mVisitData["avgAmt"] = "n"
	}

	if len(monthTmCnt) > 11{
		mVisitData["minAmt"] = monthTmCnt[9:]
	}else {
		mVisitData["minAmt"] = "n"
	}

	// 월 결제 금액, 결제 건수, 월 객 단가
	//mVisitData["totAmt"] = monthTotalAmt[0]["total_amt"]
	//mVisitData["totCnt"] = monthTotalAmt[0]["tot_cnt"]
	//mVisitData["totAvg"] = monthTotalAmt[0]["tot_avg"]
	data["mVisitData"] = mVisitData

	// 우리 가게 방문고객이 자주가는 주변 상점
	/*
	mAroudVisitData := make(map[string]interface{})
	mAroudVisitData["visit1cd"] = aroudVisitData[0]["mv_cd_lv2_nm"]
	mAroudVisitData["visit1nm"] = aroudVisitData[0]["mv_rtl_name"]
	mAroudVisitData["visit1cnt"] = aroudVisitData[0]["mv_cnt"]
	mAroudVisitData["visit2cd"] = aroudVisitData[1]["mv_cd_lv2_nm"]
	mAroudVisitData["visit2nm"] = aroudVisitData[1]["mv_rtl_name"]
	mAroudVisitData["visit2cnt"] = aroudVisitData[1]["mv_cnt"]
	mAroudVisitData["visit3cd"] = aroudVisitData[2]["mv_cd_lv2_nm"]
	mAroudVisitData["visit3nm"] = aroudVisitData[2]["mv_rtl_name"]
	mAroudVisitData["visit3cnt"] = aroudVisitData[2]["mv_cnt"]
	mAroudVisitData["visit4cd"] = aroudVisitData[3]["mv_cd_lv2_nm"]
	mAroudVisitData["visit4nm"] = aroudVisitData[3]["mv_rtl_name"]
	mAroudVisitData["visit4cnt"] = aroudVisitData[3]["mv_cnt"]
	data["mAroudVisitData"] = mAroudVisitData
	 */

	// 동일 업종 상점의 지역별 매출 비교
	//var monthCompareMent string

	// 미가입 -> 금융결제원 회원 데이터가 부족하여 음식점(한식) 기준 데이터와 비교한 결과 입니다.
	/*
	if !(len(kftcEnroll) > 0 &&  kftcEnroll[0]["kftcStsCd"] == "1") {

		// 월 한식 가맹점 중 매출액 1위
		koreanData, err := cls.GetSelectData(datasql.SelectMonthCompareKorea, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		if len(koreanData) == 0{
			params["bsDt"] = "202103"
			// 월 한식 가맹점 중 매출액 1위
			koreanData, err = cls.GetSelectData(datasql.SelectMonthCompareKorea, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}
		}

		params["bizNum"] = koreanData[0]["reg_no"]
		monthCompareMent = "금융결제원 회원 데이터가 부족하여 음식점(한식) 기준 데이터와 비교한 결과 입니다."
	}

	// 동일업종 월 매출 비교
	salesData, err := cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	// 동일업종 월 매출 비교
	/*
	salesData, err := cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 미가입 -> 금융결제원 회원 데이터가 부족하여 음식점(한식) 기준 데이터와 비교한 결과 입니다.
	//if !(len(kftcEnroll) > 0 &&  kftcEnroll[0]["kftcStsCd"] == "1") {
	if len(salesData) == 0{

		// 대상 업체 코드 및 주소
		compInfo, err := cls.GetSelectData(datasql.SelectPrivRestInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		//params["buety"] = compInfo[0]["buety"]
		params["buety"] = "00"
		params["bsDt"] = "202103"

		// 지난 달 월매출 비교 가맹점 리스트
		compList, err := cls.GetSelectData(datasql.SelectMonthCompareList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		if len(compList) == 0{
			// 월 한식 가맹점 중 매출액 1위
			koreanData, err := cls.GetSelectData(datasql.SelectMonthCompareKorea, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}

			params["bizNum"] = koreanData[0]["reg_no"]
			monthCompareMent = "금융결제원 회원 데이터가 부족하여 음식점(한식) 기준 데이터와 비교한 결과 입니다."

			// 동일업종 월 매출 비교
			salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}
		}else if len(compList) == 1{
			params["bizNum"] = compList[0]["reg_no"]

			monthCompareMent = fmt.Sprintf("금융결제원 회원 데이터가 부족하여 동일업종(%s) 기준 데이터와 비교한 결과 입니다.", compList[0]["cd"])

			// 동일업종 월 매출 비교
			salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}
		}else{
			// 가맹점 주소 4글자
			addr := strings.TrimSpace(compInfo[0]["addr"])
			var newAddr string

			if strings.Contains(addr, "서울") || strings.Contains(addr, "서울시"){
				newAddr = "서울특별시"
			}else if strings.Contains(addr, "경기") || strings.Contains(addr, "경기도"){
				newAddr = "경기도"
			}else if strings.Contains(addr, "강원") || strings.Contains(addr, "강원도"){
				newAddr = "강원도"
			}else if strings.Contains(addr, "인천"){
				newAddr = "인천광역시"
			}else if strings.Contains(addr, "울산"){
				newAddr = "울산광역시"
			}else if strings.Contains(addr, "대전"){
				newAddr = "대전광역시"
			}else if strings.Contains(addr, "경남") || strings.Contains(addr, "경상남"){
				newAddr = "경상남도"
			}else if strings.Contains(addr, "경북") || strings.Contains(addr, "경상북"){
				newAddr = "경상북도"
			}else if strings.Contains(addr, "충북") || strings.Contains(addr, "충청북"){
				newAddr = "충청북도"
			}else if strings.Contains(addr, "충남") || strings.Contains(addr, "충청남"){
				newAddr = "충청남도"
			}else if strings.Contains(addr, "부산"){
				newAddr = "부산광역시"
			}else if strings.Contains(addr, "전북") || strings.Contains(addr, "전라북"){
				newAddr = "전라북도"
			}else if strings.Contains(addr, "전남") || strings.Contains(addr, "전라남"){
				newAddr = "전라남도"
			}else if strings.Contains(addr, "광주"){
				newAddr = "광주광역시"
			}else if strings.Contains(addr, "제주"){
				newAddr = "제주특별시"
			}

			for _, v := range compList{
				if v["sido_nm"] == newAddr{
					params["bizNum"] = v["reg_no"]

					monthCompareMent = fmt.Sprintf("금융결제원 회원 데이터가 부족하여 %s(%s) 기준 데이터와 비교한 결과 입니다.",newAddr, v["cd"])

					// 동일업종 월 매출 비교
					salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
					}

					break
				}
			}

			if len(salesData) == 0{
				params["bizNum"] = compList[0]["reg_no"]

				monthCompareMent = fmt.Sprintf("금융결제원 회원 데이터가 부족하여 %s(%s) 기준 데이터와 비교한 결과 입니다.",compList[0]["sido_nm"], compList[0]["cd"])

				// 동일업종 월 매출 비교
				salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
				}
			}
		}
	}
	 */

/*
	mSalesData := make(map[string]interface{})
	//myAmt, _ := strconv.Atoi(salesData[0]["my_amt"])
	myAmt, _ := strconv.Atoi(thisSalesData[len(thisSalesData)-1]["tot_amt"])
	sidoAmt, _ := strconv.Atoi(salesData[0]["sido_amt"])
	gunguAmt, _ := strconv.Atoi(salesData[0]["gungu_amt"])
	totalAmt, _ := strconv.Atoi(salesData[0]["total_amt"])

	if len(monthCompareMent) == 0{
		if myAmt > totalAmt{
			monthCompareMent = "주변 동종업종과 비교하여 매출을 확인해보세요.<br>"
			monthCompareMent += "지난달 전국적으로 상위 매출을 달성하셨어요!!<br>정말 수고 많으셨습니다. 이번달에는 더 많이 버세요 ^^<img src='/public/img/img_crown.png' alt='왕관'/>"
		}else if myAmt > sidoAmt{
			monthCompareMent = "동일한 업종의 지난 달 평균 매출을 지역별로 비교한 결과입니다.<br>"
			monthCompareMent += "지난달 도내(시군) 평균 이상 매출을 달성 하셨어요!<br>정말 수고 많으셨습니다. 이번달에는 더 많이 버세요 ^^<img src='/public/img/img_crown.png' alt='왕관'/>"
		}else if myAmt > gunguAmt{
			monthCompareMent = "동일한 업종의 지난 달 평균 매출을 지역별로 비교한 결과입니다.<br>"
			monthCompareMent += "지난달 도내(시군) 평균 이상 매출을 달성 하셨어요!<br>정말 수고 많으셨습니다. 이번달에는 더 많이 버세요 ^^<img src='/public/img/img_crown.png' alt='왕관'/>"
		}else{
			monthCompareMent = "동일한 업종을 대상으로 지난 달 매출을 지역별로 비교한 결과입니다.<br>"
			monthCompareMent += "업종에 따라 오차가 있을 수 있으며, 참고 자료로 보시기 바랍니다.<br>정말 수고 많으셨습니다. 이번달에는 더 번창하세요!"
		}
	}

	mSalesData["monthCompareMent"] = monthCompareMent
	mSalesData["myAmt"] = myAmt
	mSalesData["sidoAmt"] = sidoAmt
	mSalesData["gunguAmt"] = gunguAmt
	mSalesData["totalAmt"] = totalAmt
	//mSalesData["myCnt"] = salesData[0]["my_cnt"]
	mSalesData["myCnt"] = thisSalesData[len(thisSalesData)-1]["tot_cnt"]
	mSalesData["sidoCnt"] = salesData[0]["sido_cnt"]
	mSalesData["gunguCnt"] = salesData[0]["gungu_cnt"]
	mSalesData["totalCnt"] = salesData[0]["total_cnt"]
	//mSalesData["myAvg"] = salesData[0]["my_avgprice"]
	mSalesData["myAvg"] = thisSalesData[len(thisSalesData)-1]["tot_avg"]
	mSalesData["sidoAvg"] = salesData[0]["sido_avgprice"]
	mSalesData["gunguAvg"] = salesData[0]["gungu_avgprice"]
	mSalesData["totalAvg"] = salesData[0]["total_avgprice"]
	data["mSalesData"] = mSalesData

 */

	review := make(map[string]interface{})
	var reviewKeywords []string
	var reviewPoint []float64

	if len(reviewOption) > 0{
		ratings := strings.Split(reviewOption[0]["rating"], ",")
		keywords := strings.Split(reviewOption[0]["keyword"], "|")

		for _,v := range ratings{
			rating, err := strconv.ParseFloat(v, 64)
			if err != nil{
				continue
			}
			reviewPoint = append(reviewPoint, rating)
		}

		for _,v := range keywords{
			reviewKeywords = append(reviewKeywords, v)
		}

	}else{
		reviewKeywords = append(reviewKeywords, "맛있어요")
		reviewPoint = append(reviewPoint, 1)
	}

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0{
		review["result"] = "y"
	}else{
		review["result"] = "n"
	}

Loop1:
	for _,v := range reviews{

		for _,key := range reviewKeywords{ // 키워드 포함 리뷰
			if strings.Contains(v["content"], key){
				keywordReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		for _,key := range reviewPoint{ // 평점 포함 리뷰
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r == key{
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		/*
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r <= reviewPoint{ // 평점 1점 이하 리뷰
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue
			}
		*/

		okReview ++
	}

	if newCustomer+oldCustomer > 0 {
		review["result2"] = "y"
	}else{
		review["result2"] = "n"
	}

	review["okReview"] = okReview
	review["lowReview"] = lowReview
	review["keywordReview"] = keywordReview
	review["newCustomer"] = newCustomer
	review["oldCustomer"] = oldCustomer

	if newCustomer == 0{
		newCustomer = 1
	}

	if oldCustomer == 0{
		oldCustomer = 1
	}

	review["newCustomerTotal"] = strconv.FormatFloat(newCustomerTotal/float64(newCustomer), 'f', 2, 64)
	review["oldCustomerTotal"] = strconv.FormatFloat(oldCustomerTotal/float64(oldCustomer), 'f', 2, 64)

	var keyword string
	for _,key := range reviewKeywords{
		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0{
		review["keyword"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	}else{
		review["keyword"] = "0"
	}

	data["review"] = review


	// web view content
	webViewContent := make(map[string]interface{})
	for _,v := range webView{
		webViewContent[fmt.Sprintf("conTitle%s", v["position"])] = v["title"]
		webViewContent[fmt.Sprintf("conBody%s", v["position"])] = v["content"]
	}
	data["webView"] = webViewContent


	// 지난달 취소 분석
	cancle := make(map[string]interface{})
	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int

	if len(cancleList) > 0{
		cancle["result"] = "y"

		for _,v := range cancleList{

			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle ++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0{
				noCancle ++ // 미 승인 취소
				continue
			}

			tr,_ := strconv.Atoi(v["tr_tm"][:2])
			otr,_ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10{
				nightCancle ++ // 심야 취소
			}else if v["tr_dt"] != cancleAprv[0]["tr_dt"]{
				dayCancle ++ // 일 취소
			}else if tr - otr > 3{
				timeCancle ++ // 시간 취소
			}else{
				okCancle ++ // 결제 취소
			}

		}

	}else{
		cancle["result"] = "n"
	}

	cancle["okCancle"] = okCancle
	cancle["timeCancle"] = timeCancle
	cancle["dayCancle"] = dayCancle
	cancle["nightCancle"] = nightCancle
	cancle["noCancle"] = noCancle

	data["cancle"] = cancle

	// 파트너 가입 여부
	data["partnerYN"] = "y"

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

func GetPosWeekDataTimeAvg(c echo.Context) error {

	dprintf(4, c, "call GetPosWeekDataTimeAvg\n")

	params := cls.GetParamJsonMap(c)
	weekDay, err := strconv.Atoi(params["weekDay"])
	if err != nil{
		lprintf(1, "[ERROR] %s data fail \n", params["weekDay"])
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "parm fail"))
	}

	dt := time.Now().AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	startWeek = startWeek.AddDate(0,0,weekDay)
	params["weekDay"] = startWeek.Format("20060102")

	/* run sql */
	// 요일별/시간별 매출 분석
	weekSalesData, err := cls.GetSelectData(datasql.SelectLastWeek, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	data["tot_cnt"] = weekSalesData[0]["tot_cnt"]
	data["tot_amt"] = weekSalesData[0]["tot_amt"]
	data["cnt03"] = weekSalesData[0]["cnt03"]
	data["amt03"] = weekSalesData[0]["amt03"]
	data["cnt36"] = weekSalesData[0]["cnt36"]
	data["amt36"] = weekSalesData[0]["amt36"]
	data["cnt69"] = weekSalesData[0]["cnt69"]
	data["amt69"] = weekSalesData[0]["amt69"]
	data["cnt912"] = weekSalesData[0]["cnt912"]
	data["amt912"] = weekSalesData[0]["amt912"]
	data["cnt1215"] = weekSalesData[0]["cnt1215"]
	data["amt1215"] = weekSalesData[0]["amt1215"]
	data["cnt1518"] = weekSalesData[0]["cnt1518"]
	data["amt1518"] = weekSalesData[0]["amt1518"]
	data["cnt1821"] = weekSalesData[0]["cnt1821"]
	data["amt1821"] = weekSalesData[0]["amt1821"]
	data["cnt2124"] = weekSalesData[0]["cnt2124"]
	data["amt2124"] = weekSalesData[0]["amt2124"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

func GetWeekData1(c echo.Context) error {

	dprintf(4, c, "call GetWeekData1\n")

	params := cls.GetParamJsonMap(c)

	// 1주간
	dt := time.Now().AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	endWeek := cls.GetEndOfWeek(dt)
	params["startDt"] = startWeek.Format("20060102")
	params["endDt"] = endWeek.Format("20060102")
	// 주간 매출 분석
	weekSalesData, err := cls.GetSelectDataUsingJson(datasql.SelectWeekCash, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "startDt")
	delete(params, "endDt")

	// 4주간
	params["startDt"] = endWeek.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = endWeek.Format("20060102")
	// 월 평균고객 분석
	personData, err := cls.GetSelectDataUsingJson(datasql.SelectWeekPersonVisit, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 객단가 예측
	personPrice, err := cls.GetSelectDataUsingJson(datasql.SelectAverageRevenuePerUser, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 요일 분석 팁 (요일별 날자 입력, 지난 7일) -> 평일 주말 비교
	dayAnalystic, err := cls.GetSelectDataUsingJson(datasql.SelectWeekAvgTime1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})

	// 주간 매출 분석
	avgAmt, _ := strconv.Atoi(weekSalesData[0]["avgAmt"])
	maxAmt, _ := strconv.Atoi(weekSalesData[0]["maxAmt"])
	minAmt, _ := strconv.Atoi(weekSalesData[0]["minAmt"])

	minh := int(((float64(minAmt) / float64(maxAmt)) * float64(100)))
	mint := 100 - minh

	avgh := int(((float64(avgAmt) / float64(maxAmt)) * float64(100)))
	avgt := 100 - avgh

	pageSale := make(map[string]interface{})
	pageSale["avgAmt"] = avgAmt
	pageSale["maxAmt"] = maxAmt
	pageSale["minAmt"] = minAmt
	pageSale["minDt"] = weekSalesData[0]["minDt"]
	pageSale["maxDt"] = weekSalesData[0]["maxDt"]
	pageSale["minWeek"] = weekSalesData[0]["minWeek"]
	pageSale["maxWeek"] = weekSalesData[0]["maxWeek"]
	pageSale["minh"] = minh
	pageSale["mint"] = mint
	pageSale["avgh"] = avgh
	pageSale["avgt"] = avgt

	data["weekdaySaleAnalystic"] = pageSale

	// 방문
	visitTotal, _ := strconv.Atoi(personData[0]["visitTotal"])
	visit1, _ := strconv.Atoi(personData[0]["visit1"])
	visit23, _ := strconv.Atoi(personData[0]["visit23"])
	visit4, _ := strconv.Atoi(personData[0]["visit4"])

	visit1p := int(((float64(visit1) / float64(visitTotal)) * float64(100)))
	visit23p := int(((float64(visit23) / float64(visitTotal)) * float64(100)))
	visit4p := int(((float64(visit4) / float64(visitTotal)) * float64(100)))

	pageGuest := make(map[string]interface{})
	pageGuest["visitTotal"] = visitTotal
	pageGuest["visit1"] = visit1
	pageGuest["visit1p"] = visit1p
	pageGuest["visit23"] = visit23
	pageGuest["visit23p"] = visit23p
	pageGuest["visit4"] = visit4
	pageGuest["visit4p"] = visit4p
	pageGuest["arpu"] = personPrice[0]["arpu"]

	data["guestVisitAnalystic"] = pageGuest

	// 요일 분석 팁
	var weekAnalystic []map[string]interface{}
	var holySum, holy0003, holy0306, holy0609, holy0912, holy1215, holy1518, holy1821, holy2124 int
	var workSum, work0003, work0306, work0609, work0912, work1215, work1518, work1821, work2124 int
	for _, day := range dayAnalystic {
		week := make(map[string]interface{})
		week["trDt"] = day["trDt"]
		week["week"] = day["week"]
		week["weekNm"] = day["weekNm"]
		week["weekEnd"] = day["weekEnd"]

		totSum, _ := strconv.Atoi(day["totSum"])
		t0003, _ := strconv.Atoi(day["t0003"])
		t0306, _ := strconv.Atoi(day["t0306"])
		t0609, _ := strconv.Atoi(day["t0609"])
		t0912, _ := strconv.Atoi(day["t0912"])
		t1215, _ := strconv.Atoi(day["t1215"])
		t1518, _ := strconv.Atoi(day["t1518"])
		t1821, _ := strconv.Atoi(day["t1821"])
		t2124, _ := strconv.Atoi(day["t2124"])

		week["totSum"] = totSum
		week["t0003"] = t0003
		week["t0306"] = t0306
		week["t0609"] = t0609
		week["t0912"] = t0912
		week["t1215"] = t1215
		week["t1518"] = t1518
		week["t1821"] = t1821
		week["t2124"] = t2124

		if day["weekEnd"] == "HD" {
			holySum = holySum + totSum
			holy0003 = holy0003 + t0003
			holy0306 = holy0306 + t0306
			holy0609 = holy0609 + t0609
			holy0912 = holy0912 + t0912
			holy1215 = holy1215 + t1215
			holy1518 = holy1518 + t1518
			holy1821 = holy1821 + t1821
			holy2124 = holy2124 + t2124
		} else {
			workSum = workSum + totSum
			work0003 = work0003 + t0003
			work0306 = work0306 + t0306
			work0609 = work0609 + t0609
			work0912 = work0912 + t0912
			work1215 = work1215 + t1215
			work1518 = work1518 + t1518
			work1821 = work1821 + t1821
			work2124 = work2124 + t2124
		}
		weekAnalystic = append(weekAnalystic, week)
	}

	data["weekAnalystic"] = weekAnalystic

	// 달아요팁
	tipTotal := holySum + workSum
	holySump := int(((float64(holySum) / float64(tipTotal)) * float64(100)))
	workSump := int(((float64(workSum) / float64(tipTotal)) * float64(100)))

	darayoTipMsg := ""
	darayoTipP := ""
	if holySum > workSum {
		darayoTipP = strconv.Itoa(holySump - workSump)
		if holySum-workSum < 10 {
			darayoTipMsg = "지난 4주간 분석 결과 주말과 평일의 매출 차이가 많지 않아요.\n주말 매출 평균 " + darayoTipP + "% 정도 많았어요."
		} else {
			darayoTipMsg = "지난 4주간 분석 결과 주말 매출이 평균" + darayoTipP + "% 많았어요.\n주말에는 영업 준비를 조금 더 많이 해보세요."
		}
	} else {
		darayoTipP = strconv.Itoa(workSump - holySump)
		if workSum-holySum < 10 {
			darayoTipMsg = "지난 4주간 분석 결과 주말과 평일의 매출 차이가 많지 않아요.\n평일 매출 평균 " + darayoTipP + "% 정도 많았어요."
		} else {
			darayoTipMsg = "지난 4주간 분석 결과 평일 매출이 평균" + darayoTipP + "% 많았어요.\n평일 고객을 위해 다양한 준비를 해보세요."
		}
	}
	pageTip := make(map[string]interface{})
	pageTip["darayoTipMsg"] = darayoTipMsg

	holyTime := []int{holy0003, holy0306, holy0609, holy0912, holy1215, holy1518, holy1821, holy2124}
	workTime := []int{work0003, work0306, work0609, work0912, work1215, work1518, work1821, work2124}
	pageTip["holyBusyMsg"] = "주말에는 평균 " + findBusyTime(holyTime) + "시가 가장 매출이 많아요."
	pageTip["workBusyMsg"] = "평일에는 평균 " + findBusyTime(workTime) + "시가 가장 매출이 많아요."

	data["darayoTip"] = pageTip

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

func GetWeekAvgTimData(c echo.Context) error {

	dprintf(4, c, "call GetWeekAvgTimData\n")

	params := cls.GetParamJsonMap(c)
	//fmt.Println(params["weekDay"])
	//fmt.Println(params["bizNum"])
	weekDay, err := strconv.Atoi(params["weekDay"])
	if err != nil{
		lprintf(1, "[ERROR] %s data fail \n", params["weekDay"])
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "parm fail"))
	}

	dt := time.Now().AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	startWeek = startWeek.AddDate(0,0,weekDay)

	params["weekDay"] = startWeek.Format("20060102")
	//fmt.Println(params["weekDay"])

	avgTimeData, err := cls.GetSelectData(datasql.SelectWeekAvgTime, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	/*
	avgTimeData2, err := cls.GetSelectData(datasql.SelectWeekAvgTime2, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	var weekTotal, line0003, line0406, line0709, line1012, line1315, line1618, line1921, line2224 string
	if len(avgTimeData) > 0{
		weekTotal = avgTimeData[0]["weekTotal"]
		line0003 = avgTimeData[0]["t0003"]
		line0406 = avgTimeData[0]["t0306"]
		line0709 = avgTimeData[0]["t0609"]
		line1012 = avgTimeData[0]["t0912"]
		line1315 = avgTimeData[0]["t1215"]
		line1618 = avgTimeData[0]["t1518"]
		line1921 = avgTimeData[0]["t1821"]
		line2224 = avgTimeData[0]["t2124"]
	}

	data := make(map[string]interface{})
	// 선그래프

	if len(weekTotal) == 0{
		data["weekTotal"] = "0"
	}else{
		data["weekTotal"] = weekTotal
	}

	data["line0003"] = line0003
	data["line0406"] = line0406
	data["line0709"] = line0709
	data["line1012"] = line1012
	data["line1315"] = line1315
	data["line1618"] = line1618
	data["line1921"] = line1921
	data["line2224"] = line2224
	// 선그래프 끝

	/*
	if len(avgTimeData2) > 0{
		data["line0006"] = avgTimeData2[0]["t0006"]
		data["line0611"] = avgTimeData2[0]["t0611"]
		data["line1114"] = avgTimeData2[0]["t1114"]
		data["line1417"] = avgTimeData2[0]["t1417"]
		data["line1724"] = avgTimeData2[0]["t1724"]
	}else{
		data["line0006"] = ""
		data["line0611"] = ""
		data["line1114"] = ""
		data["line1417"] = ""
		data["line1724"] = ""
	}
	 */


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

func GetMonthData(c echo.Context) error {

	dprintf(4, c, "call GetMonthData\n")

	params := cls.GetParamJsonMap(c)

	weekSalesData, err := cls.GetSelectData(datasql.SelectMonthCash, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	MonthAvgWeekData, err := cls.GetSelectData(datasql.SelectMonthAvgWeek, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	darayoTipData, err := cls.GetSelectData(datasql.SelectMonthDarayoTip, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	visitData, err := cls.GetSelectData(datasql.SelectMonthVisit, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	aroudVisitData, err := cls.GetSelectData(datasql.SelectMonthAroundStroeVisit, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	monthBusyTimeData, err := cls.GetSelectData(datasql.SelectMonthBusyTime, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})

	//올해  매출 분석
	avgAmt, _ := strconv.Atoi(weekSalesData[0]["avgAmt"])
	maxAmt, _ := strconv.Atoi(weekSalesData[0]["maxAmt"])
	beforeMonthAmt, _ := strconv.Atoi(weekSalesData[0]["beforeMonthAmt"])
	bfmh := int(((float64(beforeMonthAmt) / float64(maxAmt)) * float64(100)))
	bfmt := 100 - bfmh
	avgh := int(((float64(avgAmt) / float64(maxAmt)) * float64(100)))
	avgt := 100 - avgh

	data["avgAmt"] = avgAmt
	data["maxAmt"] = maxAmt
	data["beforeMonthAmt"] = beforeMonthAmt
	data["bfmh"] = bfmh
	data["bfmt"] = bfmt
	data["avgh"] = avgh
	data["avgt"] = avgt

	data["maxMonth"] = weekSalesData[0]["maxMonth"]
	data["allAmt"] = weekSalesData[0]["allAmt"]

	// 작년 매출 분석

	l_avgAmt, _ := strconv.Atoi(weekSalesData[1]["avgAmt"])
	l_maxAmt, _ := strconv.Atoi(weekSalesData[1]["maxAmt"])
	l_beforeMonthAmt, _ := strconv.Atoi(weekSalesData[1]["beforeMonthAmt"])
	l_bfmh := int(((float64(l_beforeMonthAmt) / float64(l_maxAmt)) * float64(100)))
	l_bfmt := 100 - l_bfmh
	l_avgh := int(((float64(l_avgAmt) / float64(l_maxAmt)) * float64(100)))
	l_avgt := 100 - l_avgh

	data["l_avgAmt"] = l_avgAmt
	data["l_maxAmt"] = l_maxAmt
	data["l_beforeMonthAmt"] = l_beforeMonthAmt
	data["l_bfmh"] = l_bfmh
	data["l_bfmt"] = l_bfmt
	data["l_avgh"] = l_avgh
	data["l_avgt"] = l_avgt

	data["l_maxMonth"] = weekSalesData[1]["maxMonth"]
	data["l_allAmt"] = weekSalesData[1]["allAmt"]

	// 선그래프

	monthTotal, _ := strconv.Atoi(MonthAvgWeekData[0]["total"])
	mon, _ := strconv.Atoi(MonthAvgWeekData[0]["mon"])
	tue, _ := strconv.Atoi(MonthAvgWeekData[0]["tue"])
	wed, _ := strconv.Atoi(MonthAvgWeekData[0]["wed"])
	thr, _ := strconv.Atoi(MonthAvgWeekData[0]["thr"])
	fri, _ := strconv.Atoi(MonthAvgWeekData[0]["fri"])
	sat, _ := strconv.Atoi(MonthAvgWeekData[0]["sat"])
	sun, _ := strconv.Atoi(MonthAvgWeekData[0]["sun"])

	data["monthTotal"] = monthTotal
	data["mon"] = mon
	data["tue"] = tue
	data["wed"] = wed
	data["thr"] = thr
	data["fri"] = fri
	data["sat"] = sat
	data["sun"] = sun
	// 선그래프 끝

	// 달아요팁

	darayoTiptotal, _ := strconv.Atoi(darayoTipData[0]["darayoTiptotal"])
	holy, _ := strconv.Atoi(darayoTipData[0]["holy"])
	work, _ := strconv.Atoi(darayoTipData[0]["work"])
	holyp := int(((float64(holy) / float64(darayoTiptotal)) * float64(100)))
	workp := int(((float64(work) / float64(darayoTiptotal)) * float64(100)))

	darayoTipMsg := ""
	darayoTipP := ""
	if holy > work {
		darayoTipP = strconv.Itoa(holyp - workp)
		darayoTipMsg = "주말에는 매출이 " + darayoTipP + "% 많으신 편이에요."
	} else {
		darayoTipP = strconv.Itoa(workp - holyp)
		darayoTipMsg = "평일에는 매출이 " + darayoTipP + "% 많으신 편이에요."
	}

	data["darayoTipMsg"] = darayoTipMsg

	h0003, _ := strconv.Atoi(monthBusyTimeData[0]["0003"])
	h0406, _ := strconv.Atoi(monthBusyTimeData[0]["0406"])
	h0709, _ := strconv.Atoi(monthBusyTimeData[0]["0709"])
	h1012, _ := strconv.Atoi(monthBusyTimeData[0]["1012"])
	h1315, _ := strconv.Atoi(monthBusyTimeData[0]["1315"])
	h1618, _ := strconv.Atoi(monthBusyTimeData[0]["1618"])
	h1921, _ := strconv.Atoi(monthBusyTimeData[0]["1921"])
	h2224, _ := strconv.Atoi(monthBusyTimeData[0]["2224"])

	var hbusyData = []int{h0003, h0406, h0709, h1012, h1315, h1618, h1921, h2224}

	hTime := findBusyTime(hbusyData)

	w0003, _ := strconv.Atoi(monthBusyTimeData[1]["0003"])
	w0406, _ := strconv.Atoi(monthBusyTimeData[1]["0406"])
	w0709, _ := strconv.Atoi(monthBusyTimeData[1]["0709"])
	w1012, _ := strconv.Atoi(monthBusyTimeData[1]["1012"])
	w1315, _ := strconv.Atoi(monthBusyTimeData[1]["1315"])
	w1618, _ := strconv.Atoi(monthBusyTimeData[1]["1618"])
	w1921, _ := strconv.Atoi(monthBusyTimeData[1]["1921"])
	w2224, _ := strconv.Atoi(monthBusyTimeData[1]["2224"])

	var wbusyData = []int{w0003, w0406, w0709, w1012, w1315, w1618, w1921, w2224}

	wTime := findBusyTime(wbusyData)

	data["hTime"] = hTime
	data["wTime"] = wTime

	//방문분석
	visitTotal, _ := strconv.Atoi(visitData[0]["visitTotal"])
	visit01, _ := strconv.Atoi(visitData[0]["visit01"])
	visit25, _ := strconv.Atoi(visitData[0]["visit25"])
	visit69, _ := strconv.Atoi(visitData[0]["visit69"])
	visit10, _ := strconv.Atoi(visitData[0]["visit10"])

	var IntvisitData = []int{visit01, visit25, visit69, visit10}

	maxVisit := findMaxVist(IntvisitData)

	visit01p := int(((float64(visit01) / float64(visitTotal)) * float64(100)))
	visit25p := int(((float64(visit25) / float64(visitTotal)) * float64(100)))
	visit69p := int(((float64(visit69) / float64(visitTotal)) * float64(100)))
	visit10p := int(((float64(visit10) / float64(visitTotal)) * float64(100)))

	data["visitTotal"] = visitTotal
	data["visit01"] = visit01
	data["visit25"] = visit25
	data["visit69"] = visit69
	data["visit10"] = visit10
	data["visit01p"] = visit01p
	data["visit25p"] = visit25p
	data["visit69p"] = visit69p
	data["visit10p"] = visit10p
	data["maxVisit"] = maxVisit

	store1P, _ := strconv.Atoi(aroudVisitData[0]["mv_cnt"])
	store2P, _ := strconv.Atoi(aroudVisitData[1]["mv_cnt"])
	store3P, _ := strconv.Atoi(aroudVisitData[2]["mv_cnt"])
	store4P, _ := strconv.Atoi(aroudVisitData[3]["mv_cnt"])

	totalStore := store1P + store2P + store3P + store4P

	store1P = int(((float64(store1P) / float64(totalStore)) * float64(100)))
	store2P = int(((float64(store2P) / float64(totalStore)) * float64(100)))
	store3P = int(((float64(store3P) / float64(totalStore)) * float64(100)))
	store4P = int(((float64(store4P) / float64(totalStore)) * float64(100)))

	data["store1Nm"] = aroudVisitData[0]["mv_rtl_name"]
	data["store2Nm"] = aroudVisitData[1]["mv_rtl_name"]
	data["store3Nm"] = aroudVisitData[2]["mv_rtl_name"]
	data["store4Nm"] = "기타"

	data["store1P"] = store1P
	data["store2P"] = store1P
	data["store3P"] = store1P
	data["store4P"] = store1P

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

// 월간분석
func GetMonthData1(c echo.Context) error {
	dprintf(4, c, "call GetMonthData1\n")

	params := cls.GetParamJsonMap(c)

	// 가입일 조회
	regDt, err := cls.GetSelectData(datasql.SelectRegistDate, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if len(regDt) == 0{
		tmp := make(map[string]interface{})
		monthSaleAnalystic := make(map[string]interface{})
		data := make(map[string]interface{})
		monthSaleAnalystic["thisMonthTitle"] = "월간 매출 (데이터 수집 중)"
		data["monthSaleAnalystic"] = monthSaleAnalystic
		tmp["resultCode"] = "00"
		tmp["resultMsg"] = "응답 성공"
		tmp["resultData"] = data

		return c.JSON(http.StatusOK, tmp)
	}

	//params["restId"] = regDt[0]["restId"]

	now := time.Now()
	thisYear := now.AddDate(0, 0, 0).Format("2006")
	lastYear := now.AddDate(-1, 0, 0).Format("2006")

	params["lastDt"] = now.Format("200601")
	params["bsDt"] = now.Format("200601")

	// 월간 12달 비교
	for i:=0; i<6; i++{
		startDt := now.AddDate(0,-1*i,-3)

		params[fmt.Sprintf("%dbsDt",i)] = startDt.Format("200601")
		params[fmt.Sprintf("%dbsDt1",i)] = fmt.Sprintf("%s01", startDt.Format("200601"))
	}

	// 12주간 매출 비교
	allMonthAmt, err := cls.GetSelectData(datasql.Select6MonthAmtCard, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 월간분석 - 월 매출 분석(올해)
	if thisYear == regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 수집월 (지난달 부터 시작)
		mon := regDt[0]["regDt"][4:6]
		monInt,_ :=strconv.Atoi(mon)
		monInt -= 1
		if monInt < 10{
			params["startDt"] = fmt.Sprintf("%s0%d", regDt[0]["regDt"][:4], monInt)
		}else{
			params["startDt"] = fmt.Sprintf("%s%d", regDt[0]["regDt"][:4], monInt)
		}
	} else {
		// 작년 가입한 경우 1월, 혹은 프리미엄
		params["startDt"] = thisYear + "01"
	}

	// 저번달
	params["endDt"] = now.AddDate(0, 0, -now.Day()).Format("200601")

	// 월간분석 - 월 매출 분석(올해)
	thisSalesData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthCash1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if len(thisSalesData) == 0{
		tmp := make(map[string]interface{})
		monthSaleAnalystic := make(map[string]interface{})
		data := make(map[string]interface{})
		monthSaleAnalystic["thisMonthTitle"] = "월간 매출 (데이터 수집 중)"
		data["monthSaleAnalystic"] = monthSaleAnalystic
		tmp["resultCode"] = "00"
		tmp["resultMsg"] = "응답 성공"
		tmp["resultData"] = data

		return c.JSON(http.StatusOK, tmp)
	}
	delete(params, "startDt")
	delete(params, "endDt")

	// 월간분석 - 월 매출 분석(작년)
	var lastSalesData []map[string]string
	if thisYear != regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 데이터 없음

		if lastYear == regDt[0]["regDt"][:4] {
			// 작년 가입한 경우 데이터 수집 시작 월
			params["startDt"] = regDt[0]["regDt"][:6]
		} else {
			// 재작년 가입한 경우 작년 1월
			params["startDt"] = lastYear + "01"
		}

		params["endDt"] = lastYear + "12"
		lastSalesData, err = cls.GetSelectDataUsingJson(datasql.SelectMonthCash2, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}
		delete(params, "startDt")
		delete(params, "endDt")
	}


	// 저번달
	params["bsDt"] = now.AddDate(0, 0, -now.Day()).Format("200601")

	// 지난달 고객님 결제 건수 분석
	selectMonthCnt, err := cls.GetSelectDataUsingJson(datasql.SelectMonthCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난달 고객님 결제 건수 분석 디테일
	selectMonthCntDetail, err := cls.GetSelectDataUsingJson(datasql.SelectMonthCntDetail, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 지난주 취소 분석
	cancleList, err := cls.GetSelectData(datasql.SelectLastMonthCancleList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 월 평균 요일별 매출 분석
	MonthAvgWeekData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthAvgWeek1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 달아요팁
	darayoTipData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthDarayoTip1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 방문분석
	/*
	visitData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthVisit1, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "bsDt")
	 */

	// 주변 상점 -> 임시 데이터
	/*
	aroudVisitData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthAroundStroeVisit, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	// 금결원 가입여부
	/*
	kftcEnroll, err := cls.GetSelectDataUsingJson(datasql.SelectKFTCEnroll, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	*/

	// web view 월간 컨텐츠
	webView, err := cls.GetSelectData2(datasql.SelectWebViewMonth, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if len(webView) == 0{
		webView, _ = cls.GetSelectData2(datasql.SelectWebViewMonthDefault, params)
	}

	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryId, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 가맹점 rating, keyword
	reviewOption, err := cls.GetSelectData2(datasql.SelectReivewOption, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})

	// 고객님 방문 분석
	selectMonthAnalystic := make(map[string]interface{})
	selectMonthAnalystic["month"] = selectMonthCnt[0]["bs_dt"]
	selectMonthAnalystic["tot_cnt"] = selectMonthCnt[0]["tot_cnt"]
	selectMonthAnalystic["tot_avg"] = selectMonthCnt[0]["tot_avg"]
	selectMonthAnalystic["len"] = len(selectMonthCntDetail)
	tot_cnt, _ := strconv.Atoi(selectMonthCnt[0]["tot_cnt"])
	var best int

	for i:=0; i<len(selectMonthCntDetail); i++{

		t0003,_ := strconv.Atoi(selectMonthCntDetail[i]["t0003"])
		t0306,_ := strconv.Atoi(selectMonthCntDetail[i]["t0306"])
		t0609,_ := strconv.Atoi(selectMonthCntDetail[i]["t0609"])
		t0912,_ := strconv.Atoi(selectMonthCntDetail[i]["t0912"])
		t1215,_ := strconv.Atoi(selectMonthCntDetail[i]["t1215"])
		t1518,_ := strconv.Atoi(selectMonthCntDetail[i]["t1518"])
		t1821,_ := strconv.Atoi(selectMonthCntDetail[i]["t1821"])
		t2124,_ := strconv.Atoi(selectMonthCntDetail[i]["t2124"])

		tTime := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}
		tm,cnt := FindBusyTime2(tTime)
		if best < cnt{
			best = cnt
		}

		selectMonthAnalystic[fmt.Sprintf("monthDayName%d", i)] = selectMonthCntDetail[i]["day_name"]
		selectMonthAnalystic[fmt.Sprintf("monthTrTm%d", i)] = fmt.Sprintf("%s시", tm)
		selectMonthAnalystic[fmt.Sprintf("monthCnt%d", i)] = cnt

		dayTotCnt, _ := strconv.Atoi(selectMonthCntDetail[i]["tot_cnt"])

		bestCntp := int(((float64(dayTotCnt) / float64(tot_cnt)) * float64(100)))
		selectMonthAnalystic[fmt.Sprintf("monthCntp%d", i)] = bestCntp
	}
	selectMonthAnalystic["best"] = best
	data["selectMonthAnalystic"] = selectMonthAnalystic

	// 올해 월 매출 분석

	var maxAmt, allAmt, avgAmt int
	var maxMonthIndex string

	for _,v := range allMonthAmt{

		dayIndex := fmt.Sprintf("%s", strings.TrimSpace(v["dayIndex"][7:]))

		samt, _ := strconv.Atoi(v["amt"])
		if samt> maxAmt{
			maxAmt = samt
			maxMonthIndex = dayIndex
		}

		allAmt += samt
	}

	avgAmt = int(allAmt/6)
	beforeMonthAmt, _ := strconv.Atoi(allMonthAmt[0]["amt"])

	bfmh := int(((float64(beforeMonthAmt) / float64(maxAmt)) * float64(100)))
	bfmt := 100 - bfmh
	avgh := int(((float64(avgAmt) / float64(maxAmt)) * float64(100)))
	avgt := 100 - avgh

	pageSale := make(map[string]interface{})
	pageSale["fromDt"] = thisSalesData[0]["fromDt"]
	pageSale["toDt"] = thisSalesData[0]["toDt"]
	pageSale["year"] = thisYear
	pageSale["beforeMonthAmt"] = beforeMonthAmt
	pageSale["beforeMonth"] = thisSalesData[0]["beforeMonth"][4:]
	pageSale["avgAmt"] = avgAmt
	pageSale["maxAmt"] = maxAmt
	pageSale["maxMonth"] = maxMonthIndex
	pageSale["allAmt"] = allAmt
	pageSale["bfmh"] = bfmh
	pageSale["bfmt"] = bfmt
	pageSale["avgh"] = avgh
	pageSale["avgt"] = avgt

	thisMonthTitle := fmt.Sprintf("월간 매출 (%s월~%s월)",thisSalesData[0]["fromDt"],thisSalesData[0]["toDt"])
	pageSale["thisMonthTitle"] = thisMonthTitle

	// 작년 월 매출 분석
	if len(lastSalesData) > 0 {

		l_avgAmt, _ := strconv.Atoi(lastSalesData[0]["avgAmt"])
		l_maxAmt, _ := strconv.Atoi(lastSalesData[0]["maxAmt"])
		l_beforeMonthAmt, _ := strconv.Atoi(lastSalesData[0]["minMonthAmt"])
		l_allAmt, _ := strconv.Atoi(lastSalesData[0]["allAmt"])
		l_bfmh := int(((float64(l_beforeMonthAmt) / float64(l_maxAmt)) * float64(100)))
		l_bfmt := 100 - l_bfmh
		l_avgh := int(((float64(l_avgAmt) / float64(l_maxAmt)) * float64(100)))
		l_avgt := 100 - l_avgh

		pageSale["l_year"] = lastYear
		pageSale["l_minMonthAmt"] = l_beforeMonthAmt
		pageSale["l_minMonth"] = lastSalesData[0]["minMonth"][4:]
		pageSale["l_avgAmt"] = l_avgAmt
		pageSale["l_maxAmt"] = l_maxAmt
		pageSale["l_maxMonth"] = lastSalesData[0]["maxMonth"][4:]
		pageSale["l_allAmt"] = l_allAmt
		pageSale["l_bfmh"] = l_bfmh
		pageSale["l_bfmt"] = l_bfmt
		pageSale["l_avgh"] = l_avgh
		pageSale["l_avgt"] = l_avgt
		pageSale["l_flag"] = "y"
		pageSale["l_fromDt"] = lastSalesData[0]["fromDt"]
		pageSale["l_toDt"] = lastSalesData[0]["toDt"]
	} else {
		pageSale["l_year"] = lastYear
		pageSale["l_beforeMonthAmt"] = 0
		pageSale["l_beforeMonth"] = ""
		pageSale["l_avgAmt"] = 0
		pageSale["l_maxAmt"] = 0
		pageSale["l_maxMonth"] = ""
		pageSale["l_allAmt"] = 0
		pageSale["l_bfmh"] = 0
		pageSale["l_bfmt"] = 0
		pageSale["l_avgh"] = 0
		pageSale["l_avgt"] = 0
		pageSale["l_flag"] = "n"
		pageSale["l_fromDt"] = ""
		pageSale["l_toDt"] = ""
	}

	data["monthSaleAnalystic"] = pageSale

	// 월 평균 요일별 매출 분석
	var monthTotal, mon, tue, wed, thr, fri, sat, sun int
	for _, avgData := range MonthAvgWeekData {

		avg, _ := strconv.Atoi(avgData["totAvg"])

		//1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일'

		switch avgData["weekCd"] {
		case "2":
			mon = avg
		case "3":
			tue = avg
		case "4":
			wed = avg
		case "5":
			thr = avg
		case "6":
			fri = avg
		case "7":
			sat = avg
		default:
			sun = avg
		}

		/*
		total, _ := strconv.Atoi(avgData["total"])
		*/
		if monthTotal < avg {
			monthTotal = avg
		}

	}

	pageWeekAvg := make(map[string]interface{})
	pageWeekAvg["bestAVg"] = monthTotal
	pageWeekAvg["mon"] = mon
	pageWeekAvg["tue"] = tue
	pageWeekAvg["wed"] = wed
	pageWeekAvg["thr"] = thr
	pageWeekAvg["fri"] = fri
	pageWeekAvg["sat"] = sat
	pageWeekAvg["sun"] = sun

	data["monthWeekSaleAnalystic"] = pageWeekAvg
	// 선그래프 끝

	// 달아요팁
	var holySum, workSum int
	var holyMax, workMax string
	for _, tipData := range darayoTipData {
		t0003, _ := strconv.Atoi(tipData["t0003"])
		t0306, _ := strconv.Atoi(tipData["t0306"])
		t0609, _ := strconv.Atoi(tipData["t0609"])
		t0912, _ := strconv.Atoi(tipData["t0912"])
		t1215, _ := strconv.Atoi(tipData["t1215"])
		t1518, _ := strconv.Atoi(tipData["t1518"])
		t1821, _ := strconv.Atoi(tipData["t1821"])
		t2124, _ := strconv.Atoi(tipData["t2124"])
		timeList := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}

		strconv.Atoi(tipData["weekEnd"])
		if tipData["weekEnd"] == "HD" {
			holySum, _ = strconv.Atoi(tipData["total"])
			holyMax = findBusyTime(timeList)
		} else {
			workSum, _ = strconv.Atoi(tipData["total"])
			workMax = findBusyTime(timeList)
		}
	}

	tipTotal := holySum + workSum
	holySump := int(((float64(holySum) / float64(tipTotal)) * float64(100)))
	workSump := int(((float64(workSum) / float64(tipTotal)) * float64(100)))

	darayoTipMsg := ""
	darayoTipP := ""
	if holySum > workSum {
		darayoTipP = strconv.Itoa(holySump - workSump)
		if holySum-workSum < 10 {
			darayoTipMsg = "지난 달 분석 결과 주말과 평일의 매출 차이가 많지 않아요.<br>주말 매출 평균 " + darayoTipP + "% 정도 많았어요."
		} else {
			//darayoTipMsg = "지난 달 분석 결과 <strong class='bl'>주말 매출이 평균 " + darayoTipP + "%</strong> 많았어요.<br>주말에는 영업 준비를 조금 더 많이 해보세요."
			darayoTipMsg = "지난 달 분석 결과 <strong class='bl'>주말 매출이 평균 " + darayoTipP + "%</strong> 많았어요."
		}
	} else {
		darayoTipP = strconv.Itoa(workSump - holySump)
		if workSum-holySum < 10 {
			darayoTipMsg = "지난 달 분석 결과 주말과 평일의 매출 차이가 많지 않아요.<br>평일 매출 평균 " + darayoTipP + "% 정도 많았어요."
		} else {
			//darayoTipMsg = "지난 달 분석 결과 <strong class='bl'>평일 매출이 평균 " + darayoTipP + "%</strong> 많았어요.<br>평일 고객을 위해 다양한 준비를 해보세요."
			darayoTipMsg = "지난 달 분석 결과 <strong class='bl'>평일 매출이 평균 " + darayoTipP + "%</strong> 많았어요."
		}
	}
	pageTip := make(map[string]interface{})
	pageTip["darayoTipMsg"] = darayoTipMsg
	pageTip["holyBusyMsg"] = "<strong class='bl'>주말에는 " + holyMax + "</strong>시 매출이 가장 많아요."
	pageTip["workBusyMsg"] = "<strong class='bl'>평일에는 " + workMax + "</strong>시 매출이 가장 많아요."

	data["darayoSaleTip"] = pageTip

	//방문분석
	/*
	visitTotalAmt, _ := strconv.Atoi(visitData[0]["visitTotalAmt"])
	visitAmt1, _ := strconv.Atoi(visitData[0]["visitAmt1"])
	visitAmt23, _ := strconv.Atoi(visitData[0]["visitAmt23"])
	visitAmt49, _ := strconv.Atoi(visitData[0]["visitAmt49"])
	visitAmt10, _ := strconv.Atoi(visitData[0]["visitAmt10"])

	var IntvisitData = []int{visitAmt1, visitAmt23, visitAmt49, visitAmt10}

	maxVisit := findMaxVist(IntvisitData)

	visitAmt1p := int(((float64(visitAmt1) / float64(visitTotalAmt)) * float64(100)))
	visitAmt23p := int(((float64(visitAmt23) / float64(visitTotalAmt)) * float64(100)))
	visitAmt49p := int(((float64(visitAmt49) / float64(visitTotalAmt)) * float64(100)))
	visitAmt10p := int(((float64(visitAmt10) / float64(visitTotalAmt)) * float64(100)))

	pageVisit := make(map[string]interface{})
	pageVisit["visitTotal"] = visitTotalAmt
	pageVisit["visit1"] = visitAmt1
	pageVisit["visit23"] = visitAmt23
	pageVisit["visit49"] = visitAmt49
	pageVisit["visit10"] = visitAmt10
	pageVisit["visit1p"] = visitAmt1p
	pageVisit["visit23p"] = visitAmt23p
	pageVisit["visit49p"] = visitAmt49p
	pageVisit["visit10p"] = visitAmt10p
	pageVisit["maxVisit"] = maxVisit

	if maxVisit == "1" {
		pageVisit["visitTip"] = "<strong class='bl'>단골 보다는 신규 방문 고객의 비중이 높아요.</strong><br>신규 방문자를 더 끌어오는 홍보활동을 해보시거나 재방문을 유도하는 이벤트를 해보세요."
	} else {
		pageVisit["visitTip"] = "<strong class='bl'>단골 고객의 매출이 높은 편이에요.</strong><br>재방문 고객을 위한 이벤트나 서비스를 시도해 보세요." +
			"방문 횟수를 늘리면 매출이 더 증가할 수 있어요."
	}

	data["monthVisitAnalystic"] = pageVisit
	 */

	// 도넛 그래프
	/*
	var store1P, store2P, store3P, store4P, totalStore int
	pageVisitStore := make(map[string]interface{})

	switch len(aroudVisitData) {
		case 1:
			store1P, _ = strconv.Atoi(aroudVisitData[0]["mv_cnt"])
			store2P = 0
			store3P = 0
			store4P = 0

			store1P = 100

			pageVisitStore["store1Nm"] = aroudVisitData[0]["mv_rtl_name"]
			pageVisitStore["store1P"] = store1P
		case 2:
			store1P, _ = strconv.Atoi(aroudVisitData[0]["mv_cnt"])
			store2P, _ = strconv.Atoi(aroudVisitData[1]["mv_cnt"])
			store3P = 0
			store4P = 0

			totalStore = store1P + store2P
			store1P = int(((float64(store1P) / float64(totalStore)) * float64(100)))
			store2P = int(((float64(store2P) / float64(totalStore)) * float64(100)))


			pageVisitStore["store1Nm"] = aroudVisitData[0]["mv_rtl_name"]
			pageVisitStore["store1P"] = store1P
			pageVisitStore["store2Nm"] = aroudVisitData[1]["mv_rtl_name"]
			pageVisitStore["store2P"] = store1P
		case 3:
			store1P, _ = strconv.Atoi(aroudVisitData[0]["mv_cnt"])
			store2P, _ = strconv.Atoi(aroudVisitData[1]["mv_cnt"])
			store3P, _ = strconv.Atoi(aroudVisitData[2]["mv_cnt"])
			store4P = 0

			totalStore = store1P + store2P + store3P
			store1P = int(((float64(store1P) / float64(totalStore)) * float64(100)))
			store2P = int(((float64(store2P) / float64(totalStore)) * float64(100)))
			store3P = int(((float64(store3P) / float64(totalStore)) * float64(100)))

			pageVisitStore["store1Nm"] = aroudVisitData[0]["mv_rtl_name"]
			pageVisitStore["store1P"] = store1P
			pageVisitStore["store2Nm"] = aroudVisitData[1]["mv_rtl_name"]
			pageVisitStore["store2P"] = store2P
			pageVisitStore["store3Nm"] = aroudVisitData[2]["mv_rtl_name"]
			pageVisitStore["store3P"] = store3P

		case 4:
			store1P, _ = strconv.Atoi(aroudVisitData[0]["mv_cnt"])
			store2P, _ = strconv.Atoi(aroudVisitData[1]["mv_cnt"])
			store3P, _ = strconv.Atoi(aroudVisitData[2]["mv_cnt"])
			store4P, _ = strconv.Atoi(aroudVisitData[3]["mv_cnt"])

			totalStore = store1P + store2P + store3P + store4P
			store1P = int(((float64(store1P) / float64(totalStore)) * float64(100)))
			store2P = int(((float64(store2P) / float64(totalStore)) * float64(100)))
			store3P = int(((float64(store3P) / float64(totalStore)) * float64(100)))
			store4P = int(((float64(store4P) / float64(totalStore)) * float64(100)))

			pageVisitStore["store1Nm"] = aroudVisitData[0]["mv_rtl_name"]
			pageVisitStore["store1P"] = store1P
			pageVisitStore["store2Nm"] = aroudVisitData[1]["mv_rtl_name"]
			pageVisitStore["store2P"] = store2P
			pageVisitStore["store3Nm"] = aroudVisitData[2]["mv_rtl_name"]
			pageVisitStore["store3P"] = store3P
			pageVisitStore["store4Nm"] = aroudVisitData[3]["mv_rtl_name"]
			pageVisitStore["store4P"] = store4P
	default:
			pageVisitStore["store1Nm"] = "예시데이터"
			pageVisitStore["store1P"] = 25
			pageVisitStore["store2Nm"] = "금융결제원"
			pageVisitStore["store2P"] = 25
			pageVisitStore["store3Nm"] = "회원에게만"
			pageVisitStore["store3P"] = 25
			pageVisitStore["store4Nm"] = "적용됩니다"
			pageVisitStore["store4P"] = 25
	}
	pageVisitStore["visitLen"] = len(aroudVisitData)
	data["visitStores"] = pageVisitStore
	 */

	// 달아요팁2
	/*
	var VisitStoreTipMsg string
	dprintf(4, c, "call kftcEnroll=%v\n", kftcEnroll)
	if len(kftcEnroll) > 0 &&  kftcEnroll[0]["kftcStsCd"] == "1" {
		VisitStoreTipMsg = "동일한 업종을 방문하는 고객이 많으시다면<br>경쟁 상품이 무엇인지 파악하고 차별성을 높여 보세요.<br><br>" +
				"상호 보완 업종을 방문하는 고객을 위해<br>협력 이벤트를 진행하면, 단골이 늘고 고정 매출이 오를 수 있어요."
	} else {
		VisitStoreTipMsg = "금육결제원 회원사에게만 제공됩니다.<br>주변 가맹점을 대상으로 매출 및 방문 데이터를 분석합니다.<br><br>" +
			"우리 가게 고객이 방문하는 가게를 알려드려 협력 할 가게나 경쟁 가게를 파악할 수 있습니다."
	}
	data["darayoStoreTip"] = VisitStoreTipMsg
	 */

	// 배달업체 만족도 분석
	var bId, yId, nId, dTip string
	delivery := make(map[string]interface{})

	if len(deliveryInfo) != 0{
		bId = deliveryInfo[0]["baemin_id"]
		yId = deliveryInfo[0]["yogiyo_id"]
		nId = deliveryInfo[0]["naver_id"]
	}

	if len(bId) > 0{
		params["baeminId"] = bId
	}else{
		params["baeminId"] = "baeminId"
	}

	if len(yId) > 0{
		params["yogiyoId"] = yId
	}else{
		params["yogiyoId"] = "yogiyoId"
	}

	if len(nId) > 0{
		params["naverId"] = nId
	}else{
		params["naverId"] = "naverId"
	}

	reviews, err := cls.GetSelectData2(datasql.SelectMonthReviews, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	params["startDt"] = now.AddDate(0, -3, -now.Day()).Format("200601")
	params["endDt"] = now.AddDate(0, 0, -now.Day()).Format("200601")

	// 리뷰 사용 안할경우 배달 리뷰 평점 및 워드 클라우드 샘플
	if len(bId) == 0 && len(yId) == 0 && len(nId) == 0{
		params["baeminId"] = "10797768"
		params["naverId"] = "38231385"
		params["yogiyoId"] = "287066"

		dTip = "<span>배달 서비스 업체 분석 예시 데이터입니다.</span><br>" +
			"<span>배달 서비스 업체 이용시 리뷰 평점 변화를 한눈에 볼 수 있어요!</span>"
	}else{
		dTip = "<span>만족도가 낮은 배달 업체와 높은 배달 업체를 비교해 보세요.</span></br>"+
			"<span>배달 업체 맛집 랭킹을 높여 보세요! </span>"
	}

	baeminReview, err := cls.GetSelectData2(datasql.SelectBaeminReview, params)
	if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

	if len(baeminReview) > 0{
		delivery["baeminYn"] = "y"
	}else{
		delivery["baeminYn"] = "n"
	}

	for idx, v := range baeminReview{
		delivery[fmt.Sprintf("baemin%davg", idx+1)], _ = strconv.ParseFloat(v["rating_avg"], 8)
		delivery[fmt.Sprintf("baemin%dcnt", idx+1)],_ = strconv.Atoi(v["cnt"])
		delivery[fmt.Sprintf("baemin%ddate", idx+1)]= v["date"]
	}

	naverReview, err := cls.GetSelectData2(datasql.SelectNaverReview, params)
	if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

	if len(naverReview) > 0{
		delivery["naverYn"] = "y"
	}else{
		delivery["naverYn"] = "n"
	}

	for idx, v := range naverReview{
		delivery[fmt.Sprintf("naver%davg", idx+1)], _ = strconv.ParseFloat(v["rating_avg"], 8)
		delivery[fmt.Sprintf("naver%dcnt", idx+1)],_ = strconv.Atoi(v["cnt"])
		delivery[fmt.Sprintf("naver%ddate", idx+1)]= v["date"]
	}

	yogiyoReview, err := cls.GetSelectData2(datasql.SelectYogiyoReview, params)
	if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

	if len(yogiyoReview) > 0{
		delivery["yogiyoYn"] = "y"
	}else{
		delivery["yogiyoYn"] = "n"
	}

	for idx, v := range yogiyoReview{
		delivery[fmt.Sprintf("yogiyo%davg", idx+1)], _ = strconv.ParseFloat(v["rating_avg"], 8)
		delivery[fmt.Sprintf("yogiyo%dcnt", idx+1)],_ = strconv.Atoi(v["cnt"])
		delivery[fmt.Sprintf("yogiyo%ddate", idx+1)]= v["date"]
	}

	delivery["deliveryTip"] = dTip
	data["delivery"] = delivery

	cancle := make(map[string]interface{})
	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int

	if len(cancleList) > 0{
		cancle["result"] = "y"

		for _,v := range cancleList{

			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle ++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0{
				noCancle ++ // 미 승인 취소
				continue
			}

			tr,_ := strconv.Atoi(v["tr_tm"][:2])
			otr,_ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10{
				nightCancle ++ // 심야 취소
			}else if v["tr_dt"] != cancleAprv[0]["tr_dt"]{
				dayCancle ++ // 일 취소
			}else if tr - otr > 3{
				timeCancle ++ // 시간 취소
			}else{
				okCancle ++ // 결제 취소
			}

		}

	}else{
		cancle["result"] = "n"
	}

	cancle["okCancle"] = okCancle
	cancle["timeCancle"] = timeCancle
	cancle["dayCancle"] = dayCancle
	cancle["nightCancle"] = nightCancle
	cancle["noCancle"] = noCancle

	data["cancle"] = cancle

	review := make(map[string]interface{})
	var reviewKeywords []string
	var reviewPoint []float64

	if len(reviewOption) > 0{
		ratings := strings.Split(reviewOption[0]["rating"], ",")
		keywords := strings.Split(reviewOption[0]["keyword"], "|")

		for _,v := range ratings{
			rating, err := strconv.ParseFloat(v, 64)
			if err != nil{
				continue
			}
			reviewPoint = append(reviewPoint, rating)
		}

		for _,v := range keywords{
			reviewKeywords = append(reviewKeywords, v)
		}

	}else{
		reviewKeywords = append(reviewKeywords, "맛있어요")
		reviewPoint = append(reviewPoint, 1)
	}

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0{
		review["result"] = "y"
	}else{
		review["result"] = "n"
	}

Loop1:
	for _,v := range reviews{
		for _,key := range reviewKeywords{ // 키워드 포함 리뷰
			if strings.Contains(v["content"], key){
				keywordReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		for _,key := range reviewPoint{ // 평점 포함 리뷰
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r == key{
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue Loop1
			}
		}

		/*
			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r <= reviewPoint{ // 평점 1점 이하 리뷰
				lowReview ++
				rst, tot := CheckReivewer(v["member_no"],v["rating"])
				if rst == 1{
					newCustomerTotal += tot
					newCustomer ++
				}else if rst == 2{
					oldCustomerTotal += tot
					oldCustomer ++
				}
				continue
			}
		*/

		okReview ++
	}

	if newCustomer+oldCustomer > 0 {
		review["result2"] = "y"
	}else{
		review["result2"] = "n"
	}

	review["okReview"] = okReview
	review["lowReview"] = lowReview
	review["keywordReview"] = keywordReview
	review["newCustomer"] = newCustomer
	review["oldCustomer"] = oldCustomer

	if newCustomer == 0{
		newCustomer = 1
	}

	if oldCustomer == 0{
		oldCustomer = 1
	}

	review["newCustomerTotal"] = strconv.FormatFloat(newCustomerTotal/float64(newCustomer), 'f', 2, 64)
	review["oldCustomerTotal"] = strconv.FormatFloat(oldCustomerTotal/float64(oldCustomer), 'f', 2, 64)

	var keyword string
	//var keyCnt int
	for _,key := range reviewKeywords{

		//keyCnt += utf8.RuneCountInString(key)
		//if keyCnt > 8{
		//	keyword = fmt.Sprintf("%s...",keyword[:len(keyword)-1])
		//	break
		//}

		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0{
		review["keyword"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	}else{
		review["keyword"] = "0"
	}

	data["review"] = review

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

	// 파트너 가입 여부
	data["partnerYN"] = "y"

	// web view content
	webViewContent := make(map[string]interface{})
	for _,v := range webView{
		webViewContent[fmt.Sprintf("conTitle%s", v["position"])] = v["title"]
		webViewContent[fmt.Sprintf("conBody%s", v["position"])] = v["content"]
	}
	data["webView"] = webViewContent

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}
/*
func GetWordCloud(c echo.Context) error {


	//data := "사장님..일부러 김치 받을라고 대자 시켰는데... 별다른 대처가 이틀째 없으셔서 별은 하나만 드립니다."

	textList := []string{}
	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	//algorithmDef := textrank.NewDefaultAlgorithm()
	algorithmDef :=textrank.NewChainAlgorithm()

	//tr.Populate(data, language, rule)
	//tr.Ranking(algorithmDef)

	data := "아 진짜 매일 먹고싶어요ㅠㅠㅠㅠ 너무 맛있어요 월급받고 또 먹어야징\n" +
		"오랜만에 시켜먹었는데 역시 맛있네요\n" +
		"잘먹었습니다 고기가 퍽퍽했어요\n" +
		"고기가 맛있네요 ㅎㅎ"



	tr.Populate(data, language, rule)
	tr.Ranking(algorithmDef)
	//rankedPhrases := textrank.FindPhrases(tr)

	words := textrank.FindSingleWords(tr)
	fmt.Println(len(words))
	//fmt.Println(len(rankedPhrases))

	//if len(rankedPhrases) == 0{

	//	return c.JSON(http.StatusOK, "0")
	//}

	//for _, v := range rankedPhrases{

	//	fmt.Println(v.Right, " ", v.Left, " ", v.Weight, " ", v.Qty)

			//if len(idMap) == 20{
			//	break
			//}
	//}

	var testData map[string]int
	testData = make(map[string]int)

	for _, v := range words{
		fmt.Println(v.Word , " ", v.Weight, " ", v.Qty)

		testData[v.Word] = v.Qty
	}

	w := wordclouds.NewWordcloud(testData, wordclouds.FontFile("./gulim.ttc"),wordclouds.Height(1024), wordclouds.Width(1024))
	img := w.Draw()

	outFile,_ := os.Create("./out1.png")
	png.Encode(outFile, img)
	outFile.Close()

	return c.JSON(http.StatusOK, "kk")

	//return c.JSON(http.StatusOK, "0")

	//for key, _ := range idMap{
	//	textList = append(textList, key)
	//}

	//textList := []string{"恭喜", "发财", "万事", "如意"}
	//angles := []int{0, 15, -15, 90} // 여러 방향
	angles := []int{0, 90} // 여러 방향
	//angles := []int{0} // 정방향
	// <img src='/public/img/img_crown.png' alt='왕관'/>


	colors := []*wordcloud.Color{
		&wordcloud.Color{0x80,0x00,0x00}, // 다크레드
		&wordcloud.Color{0x60,0x30,0x00}, // 다크초록
		&wordcloud.Color{0x00,0x00,0x80}, // 다크블루
		&wordcloud.Color{0x0, 0x60, 0x30}, // 라이트 블루
		&wordcloud.Color{0x60, 0x0, 0x0}, // 와인
		&wordcloud.Color{0x73, 0x73, 0x0}, //회색그린
	}

	render := wordcloud.NewWordCloudRender(80, 5,
		"./gulim.ttc",
		"./water.png", textList, angles, colors, "./out.png")


	render.Render()

	return c.JSON(http.StatusOK, "ok")
}
*/

func GetMonthCompareData(c echo.Context) error {

	dprintf(4, c, "call GetMonthCompareData\n")

	params := cls.GetParamJsonMap(c)

	// 저번 달
	now := time.Now()
	params["bsDt"] = now.AddDate(0, 0, -now.Day()).Format("200601")
	var monthCompareMent string

	// 가맹점 지난달 매출금액
	compAprv, err := cls.GetSelectData(datasql.SelectCompAprv, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 동일업종 월 매출 비교
	salesData, err := cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 미가입 -> 금융결제원 회원 데이터가 부족하여 음식점(한식) 기준 데이터와 비교한 결과 입니다.
	//if !(len(kftcEnroll) > 0 &&  kftcEnroll[0]["kftcStsCd"] == "1") {
	if len(salesData) == 0{

		// 대상 업체 코드 및 주소
		compInfo, err := cls.GetSelectData(datasql.SelectPrivRestInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		params["buety"] = compInfo[0]["buety"]

		// 지난 달 월매출 비교 가맹점 리스트
		compList, err := cls.GetSelectData(datasql.SelectMonthCompareList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		if len(compList) == 0{
			// 월 한식 가맹점 중 매출액 1위
			koreanData, err := cls.GetSelectData(datasql.SelectMonthCompareKorea, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}

			if len(koreanData) == 0{
				now := time.Now()
				params["bsDt"] = now.AddDate(0, -1, -now.Day()).Format("200601")

				koreanData, err = cls.GetSelectData(datasql.SelectMonthCompareKorea, params, c)
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
				}
			}

			params["bizNum"] = koreanData[0]["reg_no"]
			monthCompareMent = "금융결제원 회원 데이터가 부족하여 음식점(한식) 기준 데이터와 비교한 결과 입니다."

			// 동일업종 월 매출 비교
			salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}
		}else if len(compList) == 1{
			params["bizNum"] = compList[0]["reg_no"]

			monthCompareMent = fmt.Sprintf("금융결제원 회원 데이터가 부족하여 동일업종(%s) 기준 데이터와 비교한 결과 입니다.", compList[0]["cd"])

			// 동일업종 월 매출 비교
			salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
			if err != nil {
				return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
			}
		}else{
			// 가맹점 주소 4글자
			addr := strings.TrimSpace(compInfo[0]["addr"])
			var newAddr string

			if strings.Contains(addr, "서울") || strings.Contains(addr, "서울시"){
				newAddr = "서울특별시"
			}else if strings.Contains(addr, "경기") || strings.Contains(addr, "경기도"){
				newAddr = "경기도"
			}else if strings.Contains(addr, "강원") || strings.Contains(addr, "강원도"){
				newAddr = "강원도"
			}else if strings.Contains(addr, "인천"){
				newAddr = "인천광역시"
			}else if strings.Contains(addr, "울산"){
				newAddr = "울산광역시"
			}else if strings.Contains(addr, "대전"){
				newAddr = "대전광역시"
			}else if strings.Contains(addr, "경남") || strings.Contains(addr, "경상남"){
				newAddr = "경상남도"
			}else if strings.Contains(addr, "경북") || strings.Contains(addr, "경상북"){
				newAddr = "경상북도"
			}else if strings.Contains(addr, "충북") || strings.Contains(addr, "충청북"){
				newAddr = "충청북도"
			}else if strings.Contains(addr, "충남") || strings.Contains(addr, "충청남"){
				newAddr = "충청남도"
			}else if strings.Contains(addr, "부산"){
				newAddr = "부산광역시"
			}else if strings.Contains(addr, "전북") || strings.Contains(addr, "전라북"){
				newAddr = "전라북도"
			}else if strings.Contains(addr, "전남") || strings.Contains(addr, "전라남"){
				newAddr = "전라남도"
			}else if strings.Contains(addr, "광주"){
				newAddr = "광주광역시"
			}else if strings.Contains(addr, "제주"){
				newAddr = "제주특별시"
			}

			for _, v := range compList{
				if v["sido_nm"] == newAddr{
					params["bizNum"] = v["reg_no"]

					monthCompareMent = fmt.Sprintf("금융결제원 회원 데이터가 부족하여 %s(%s) 기준 데이터와 비교한 결과 입니다.",newAddr, v["cd"])

					// 동일업종 월 매출 비교
					salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
					if err != nil {
						return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
					}

					break
				}
			}

			if len(salesData) == 0{
				params["bizNum"] = compList[0]["reg_no"]

				monthCompareMent = fmt.Sprintf("금융결제원 회원 데이터가 부족하여 %s(%s) 기준 데이터와 비교한 결과 입니다.",compList[0]["sido_nm"], compList[0]["cd"])

				// 동일업종 월 매출 비교
				salesData, err = cls.GetSelectData(datasql.SelectMonthCompareAll, params, c)
				if err != nil {
					return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
				}
			}


		}
	}

	data := make(map[string]interface{})

	samecnt, _ := strconv.Atoi(compAprv[0]["tot_cnt"])
	sameAmt, _ := strconv.Atoi(compAprv[0]["tot_amt"])
	a_sameAmt, _ := strconv.Atoi(compAprv[0]["tot_avg"])

	var scString, saString, sapString string

	switch params["qtype"] {
		// 군
		case "0":
			scString = "gungu_cnt"
			saString = "gungu_amt"
			sapString = "gungu_avgprice"
		// 시도
		case "1":
			scString = "sido_cnt"
			saString = "sido_amt"
			sapString = "sido_avgprice"
		// 전국(2)
		default:
			scString = "total_cnt"
			saString = "total_amt"
			sapString = "total_avgprice"
	}

	n_samecnt, _ := strconv.Atoi(salesData[0][scString])
	n_sameAmt, _ := strconv.Atoi(salesData[0][saString])
	n_a_sameAmt, _ := strconv.Atoi(salesData[0][sapString])

	// 가맹점 매출건수, 매출액, 객단가
	data["samecnt"] = samecnt
	data["sameAmt"] = sameAmt
	data["a_sameAmt"] = a_sameAmt
	// 시군구, 시도, 전체
	data["n_samecnt"] = n_samecnt
	data["n_sameAmt"] = n_sameAmt
	data["n_a_sameAmt"] = n_a_sameAmt
	// 월 비교 멘트

	same1max := n_samecnt + (n_samecnt / 3)
	same2max := n_sameAmt + (n_sameAmt / 3)
	same3max := n_a_sameAmt + (n_a_sameAmt / 3)
	if samecnt > n_samecnt {
		same1max = samecnt + (samecnt / 3)
	}
	if sameAmt > n_sameAmt {
		same2max = sameAmt + (sameAmt / 3)
	}
	if a_sameAmt > n_a_sameAmt {
		same3max = a_sameAmt + (a_sameAmt / 3)
	}

	same1h := int(((float64(samecnt) / float64(same1max)) * float64(100)))
	same1t := 100 - same1h
	same2h := int(((float64(n_samecnt) / float64(same1max)) * float64(100)))
	same2t := 100 - same2h

	same3h := int(((float64(sameAmt) / float64(same2max)) * float64(100)))
	same3t := 100 - same3h
	same4h := int(((float64(n_sameAmt) / float64(same2max)) * float64(100)))
	same4t := 100 - same4h

	same5h := int(((float64(a_sameAmt) / float64(same3max)) * float64(100)))
	same5t := 100 - same5h
	same6h := int(((float64(n_a_sameAmt) / float64(same3max)) * float64(100)))
	same6t := 100 - same6h

	data["same1h"] = same1h
	data["same1t"] = same1t
	data["same2h"] = same2h
	data["same2t"] = same2t
	data["same3h"] = same3h
	data["same3t"] = same3t
	data["same4h"] = same4h
	data["same4t"] = same4t
	data["same5h"] = same5h
	data["same5t"] = same5t
	data["same6h"] = same6h
	data["same6t"] = same6t

	// 가맹점 vs 전국, 도내, 시군 매출 비교
	sidoAmt, _ := strconv.Atoi(salesData[0]["sido_amt"])
	gunguAmt, _ := strconv.Atoi(salesData[0]["gungu_amt"])
	totalAmt, _ := strconv.Atoi(salesData[0]["total_amt"])

	if len(monthCompareMent) == 0{

		monthCompareMent = "동일한 업종의 지난달 평균 매출을 지역별로 비교한 결과입니다.<br>"

		if sameAmt > totalAmt{
			monthCompareMent += "지난달 전국적으로 상위 매출을 달성하셨어요!!<br>정말 수고 많으셨습니다. 이번달에는 더 많이 버세요 ^^<img src='/public/img/img_crown.png' alt='왕관'/>"
		}else if sameAmt > sidoAmt{
			monthCompareMent += "지난달 도내(시군) 평균 이상 매출을 달성 하셨어요!<br>정말 수고 많으셨습니다. 이번달에는 더 많이 버세요 ^^<img src='/public/img/img_crown.png' alt='왕관'/>"
		}else if sameAmt > gunguAmt{
			monthCompareMent += "지난달 도내(시군) 평균 이상 매출을 달성 하셨어요!<br>정말 수고 많으셨습니다. 이번달에는 더 많이 버세요 ^^<img src='/public/img/img_crown.png' alt='왕관'/>"
		}else{
			monthCompareMent += "업종에 따라 오차가 있을 수 있으며, 참고 자료로 보시기 바랍니다.<br>정말 수고 많으셨습니다. 이번달에는 더 번창하세요!"
		}

	}

	data["monthCompareMent"] = monthCompareMent

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

func dayNames(num string) string {
	day := ""

	switch num {
	case "1":
		day = "일요일"
	case "2":
		day = "월요일"
	case "3":
		day = "화요일"
	case "4":
		day = "수요일"
	case "5":
		day = "목요일"
	case "6":
		day = "금요일"
	default:
		day = "토요일"
	}
	return day

}

func findBusyTime(arr []int) string {
	var max = arr[0]
	var index = 0
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
			index = i
		}
	}
	var stime = ""
	switch index {
	case 0:
		stime = "0~3"
	case 1:
		stime = "3~6"
	case 2:
		stime = "6~9"
	case 3:
		stime = "9~12"
	case 4:
		stime = "12~15"
	case 5:
		stime = "15~18"
	case 6:
		stime = "18~21"
	case 7:
		stime = "21~24"
	}

	return stime
}

func FindBusyTime2(arr []int) (string,int) {
	var max = arr[0]
	var index = 0
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
			index = i
		}
	}
	var stime = ""
	switch index {
	case 0:
		stime = "0~3"
	case 1:
		stime = "3~6"
	case 2:
		stime = "6~9"
	case 3:
		stime = "9~12"
	case 4:
		stime = "12~15"
	case 5:
		stime = "15~18"
	case 6:
		stime = "18~21"
	case 7:
		stime = "21~24"
	}

	return stime,max
}

func findMaxVist(arr []int) string {
	var max = arr[0]
	var index = 0
	for i := range arr {
		if arr[i] > max {
			max = arr[i]
			index = i
		}
	}
	var stime = ""
	switch index {
	case 0:
		stime = "1"
	case 1:
		stime = "2~5"
	case 2:
		stime = "6~9"
	case 3:
		stime = "10"
	}

	return stime
}

func dayFormat(date string) string{

	// date - 20210723
	// return 7/23

	if date[4] == '0'{
		return fmt.Sprintf("%s/%s", string(date[5]), date[6:])
	}

	return fmt.Sprintf("%s/%s", date[4:6], date[6:])
}



/// 월간 보고서
func GetMonthly(c echo.Context) error {
	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})


	endBsDt := time.Now().AddDate(0, -1, 0).Format("200601")
	startBsDt :=time.Now().AddDate(0, -6, 0).Format("200601")

	params["startBsDt"]=startBsDt
	params["endBsDt"]=endBsDt

	bsDt :=time.Now().AddDate(0, -1, 0).Format("200601")
	params["bsDt"]=bsDt

	deliveryInfo, err := cls.GetSelectData2(reviewsql.SelectStoreInfo_v2, params)
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

	restNm :=deliveryInfo[0]["comp_nm"]
	restId :=deliveryInfo[0]["rest_id"]

	salesInfo, err := cls.GetSelectData(datasql.SelectLastSaleInfo, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	sixMonthSales, err := cls.GetSelectData(datasql.Select6MonthSales, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	contentsInfo, err := cls.GetSelectType(reviewsql.SelectContentList, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	params["endDt"]=bsDt

	lastMonthPayAmt, err := cls.GetSelectData(datasql.SelectMonthPayAmt, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}



	// 지난달 취소 분석


	cancleList, err := cls.GetSelectData(datasql.SelectLastMonthCancleList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	cancle := make(map[string]interface{})
	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int

	if len(cancleList) > 0{
		cancle["result"] = "y"

		for _,v := range cancleList{

			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle ++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0{
				noCancle ++ // 미 승인 취소
				continue
			}

			tr,_ := strconv.Atoi(v["tr_tm"][:2])
			otr,_ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10{
				nightCancle ++ // 심야 취소
			}else if v["tr_dt"] != cancleAprv[0]["tr_dt"]{
				dayCancle ++ // 일 취소
			}else if tr - otr > 3{
				timeCancle ++ // 시간 취소
			}else{
				okCancle ++ // 결제 취소
			}

		}

	}else{
		cancle["result"] = "n"
	}

	cancle["okCancle"] = okCancle
	cancle["timeCancle"] = timeCancle
	cancle["dayCancle"] = dayCancle
	cancle["nightCancle"] = nightCancle
	cancle["noCancle"] = noCancle






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


	reviewInfo, err := cls.GetSelectData2(reviewsql.SelectReivewRatingMonth_v2, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}





	orderInfo, err := cls.GetSelectData(datasql.SelectOrderAmt, params,c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}




	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["salesInfo"] = salesInfo[0]
	m["sixMonthSales"] = sixMonthSales
	m["contentsInfo"] = contentsInfo
	m["lastMonthPayAmt"] = lastMonthPayAmt[0]
	m["cancle"] = cancle
	m["reviewInfo"] = reviewInfo[0]
	m["orderInfo"] = orderInfo[0]
	m["restNm"] =restNm
	m["restId"] =restId

	return c.JSON(http.StatusOK, m)

}