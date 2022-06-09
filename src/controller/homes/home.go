package homes

import (
	commonsql "cashApi/query/commons"
	datasql "cashApi/query/datas"
	salesql "cashApi/query/sales"
	"cashApi/src/controller"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	homesql "cashApi/query/homes"
	users "cashApi/src/controller/users"

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

// 홈화면 데이터
func GetHomeData(c echo.Context) error {

	dprintf(4, c, "call GetHomeData\n")

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectData(homesql.SelectIntroMsg, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	storeData, err := cls.GetSelectData(homesql.SelectStoreService, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if storeData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}



	introMsg := resultData[0]["introMsg"]

	billingYn := storeData[0]["billingYn"]
	cardSalesYn := storeData[0]["cardSalesYn"]
	homeTaxYn := storeData[0]["homeTaxYn"]

	restAuth :=  storeData[0]["REST_AUTH"]





	addr := storeData[0]["ADDR"]

	addrStr := strings.Split(addr, " ")
	addrShort := ""
	if len(addrStr) > 1 {
		addrShort = addrStr[0] + " " + addrStr[1]
	}

	if billingYn != "N" {
		billingYn = "Y"
	}

	if cardSalesYn != "N" {
		cardSalesYn = "Y"
	}

	if homeTaxYn != "N" {
		homeTaxYn = "Y"
	}

	if restAuth !="0" {
		billingYn ="N"
		cardSalesYn ="N"
		homeTaxYn ="N"
	}

	//billingYn = "Y"
	//cardSalesYn = "Y"
	//homeTaxYn = "Y"

	homeData := make(map[string]interface{})
	homeData["introMsg"] = introMsg
	homeData["billingYn"] = billingYn
	homeData["cardSalesYn"] = cardSalesYn
	homeData["addrShort"] = addrShort
	homeData["homeTaxYn"] = homeTaxYn

	imgUrl :="/public/img/banner/darayo_20210409_1.png"
	if cardSalesYn=="Y"{
		imgUrl = "/public/img/banner/darayo_authY.png"
	}

	homeData["bannerImgUrl"] = imgUrl
	homeData["bannerLink"] = "https://blog.naver.com/darayocash/222312331179"

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = homeData

	return c.JSON(http.StatusOK, m)

}

func GetBoardList(c echo.Context) error {

	dprintf(4, c, "call GetBoardList\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectType(commonsql.SelectBoardList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)

}

func GetPayDayList(c echo.Context) error{

	dprintf(4, c, "call GetSalesList\n")

	params := cls.GetParamJsonMap(c)
	data := make(map[string]interface{})
	tn := time.Now()
	today := tn.Format("20060102")
	//today := tn.AddDate(0,0,3).Format("20060102")

	params["bsDt"] = today
	// 오늘 입금 예정 금액
	payDailyList, err := cls.GetSelectDataUsingJson(salesql.SelectPayDailyListHome, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 합계건수, 금액 생성
	var payList []map[string]interface{}
	var pcaCnt, pcaAmt2, totFee, outpExptAmt int
	for _, payData := range payDailyList {
		data := make(map[string]interface{})
		//data["rNum"] = idx + 2
		data["cardCd"] = payData["cardCd"]
		data["cardNm"] = payData["cardNm"]
		data["merNo"] = payData["merNo"]

		tmp, _ := strconv.Atoi(payData["pcaCnt"])
		pcaCnt = pcaCnt + tmp
		data["pcaCnt"] = tmp
		tmp, _ = strconv.Atoi(payData["pcaAmt"])
		pcaAmt2 = pcaAmt2 + tmp
		data["pcaAmt"] = tmp
		tmp, _ = strconv.Atoi(payData["totFee"])
		totFee = totFee + tmp
		data["totFee"] = tmp
		//tmp, _ = strconv.Atoi(payData["vatAmt"])
		//vatAmt = vatAmt + tmp
		//data["vatAmt"] = tmp
		tmp, _ = strconv.Atoi(payData["outpExptAmt"])
		outpExptAmt = outpExptAmt + tmp
		data["outpExptAmt"] = tmp
		//tmp, _ = strconv.Atoi(payData["realInAmt"])
		//realInAmt = realInAmt + tmp
		//data["realInAmt"] = tmp
		//tmp, _ = strconv.Atoi(payData["diffAmt"])
		//diffAmt = diffAmt + tmp
		//data["diffAmt"] = tmp

		//data["diffNm"] = payData["diffNm"]
		//data["diffColor"] = payData["diffColor"]

		payList = append(payList, data)
	}

	var payAllList []map[string]interface{}

	// 합계
	sum := make(map[string]interface{})
	sum["cardCd"] = "99"
	sum["cardNm"] = "합계"
	sum["pcaCnt"] = pcaCnt
	sum["pcaAmt"] = pcaAmt2
	sum["totFee"] = totFee
	sum["outpExptAmt"] = outpExptAmt


	/*
	sum["realInAmt"] = realInAmt
	sum["diffAmt"] = diffAmt

	var diffNm, diffColor string
	if realInAmt == outpExptAmt {
		diffNm = "일치"
	} else if realInAmt < outpExptAmt {
		diffNm = "일부입금"
	} else {
		diffNm = "초과입금"
	}
	sum["diffNm"] = diffNm

	if realInAmt == outpExptAmt {
		diffColor = "0"
	} else if realInAmt < outpExptAmt {
		diffColor = "1"
	} else {
		diffColor = "2"
	}
	sum["diffColor"] = diffColor
	 */
	data["sum"] = sum

	// 신용카드
	payAllList = append(payAllList, payList...)
	if len(payAllList) == 0{
		s := []string{}
		data["payAllList"] = s
	}else{
		data["payAllList"] = payAllList
	}


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultYn"] = "Y"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 홈화면 매출/빅데이터
func GetSalesInfo(c echo.Context) error {

	dprintf(4, c, "call GetSalesList\n")

	params := cls.GetParamJsonMap(c)
	data := make(map[string]interface{})
	tn := time.Now()

	today := tn.Format("20060102")
	yesterday := tn.AddDate(0, 0, -1).Format("20060102")

	// 가입일 체크
	regInfo, err := cls.GetSelectData(homesql.SelectCashRegInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if regInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	/*
	regDt : 등록 일 시 (cash 가입 날짜)
	lnOpenDt : 여신 협회 개시일
	hometaxOpenDt : 홈텍스 개시일
	hometaxJoinStsDt : 홈텍스 가입상태
	lnJoinStsDt : 여신 가입상태
	 */

	// 해당 api 호출 전 여신 인증 / hometax 인증 여부를 확인 함
	/*
	if regInfo[0]["regDt"] == today ||
		(regInfo[0]["lnOpenDt"] == "" && regInfo[0]["hometaxOpenDt"] == "") ||
		(regInfo[0]["lnOpenDt"] == today && regInfo[0]["hometaxOpenDt"] == "") ||
		(regInfo[0]["lnOpenDt"] == "" && regInfo[0]["hometaxOpenDt"] == today) ||
		(regInfo[0]["lnOpenDt"] == today && regInfo[0]["hometaxOpenDt"] == today) ||
		(regInfo[0]["lnJoinStsDt"] != "1" && regInfo[0]["hometaxJoinStsDt"] != "1") {
	*/

	//if regInfo[0]["regDt"] == today || len(regInfo) == 0 {
	if len(regInfo) == 0 {
		data["aprvMonthSum"] = 0
		data["aprvCnt"] = 0
		data["arpuAmt"] = 0
		data["busyTime"] = ""
		data["cardAmt"] = 0
		data["cashAmt"] = 0
		data["diffCardAmt"] = 0
		data["expectMsg"] = "수집전 예상 매출"
		data["expectArrowMsg"] = "어제실적 수집전 이에요."
		data["lastAmt"] = 0
		data["lastMsg"] = "데이터 수집 중"
		data["lastArrowMsg"] = "수집완료 후 확인 가능"
		data["month"] = ""
		data["payAmt"] = 0
		data["payMonthSum"] = 0
		data["todayExpectAmt"] = 0
		data["todayPay"] = "매출내역 수집 중"
		data["visitRate"] = ""
		data["todayLuckyMsg"] = "소띠 보러가기"
		data["todayLuckyUrl"] = "https://m.fortune.nate.com/today/todayOriental.nate?contsCd=CT000136&day=0&tti=01"

		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultYn"] = "N"
		//m["resultYn"] = "Y"
		m["resultData"] = data

		return c.JSON(http.StatusOK, m)
	}

	params["restId"] = regInfo[0]["restId"]
	// get user id
	rUserInfo, err := cls.GetSelectData(homesql.SelectRestUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.GetSelectData(homesql.SelectUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	birth := userInfo[0]["birth"]

	var ceoTti string
	if len(birth) > 0{
		lprintf(4, "[INFO] user birth(%s) \n", birth)
		ceoTti = users.GetBrithdayTti(userInfo[0]["birth"])

		data["todayLuckyMsg"] = TtiName(ceoTti)
		data["todayLuckyUrl"] = fmt.Sprintf("https://m.fortune.nate.com/today/todayOriental.nate?contsCd=CT000136&day=0&tti=%s", ceoTti)
	}else{
		data["todayLuckyMsg"] = "소띠 보러가기"
		data["todayLuckyUrl"] = "https://m.fortune.nate.com/today/todayOriental.nate?contsCd=CT000136&day=0&tti=01"
	}

	// 전영업일 매출확인
	params["trDt"] = yesterday
	lastSales, err := cls.GetSelectData(homesql.SelectDaySaleData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	/*
	if lastSales == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}
	*/

	//lastAmt, _ := strconv.Atoi(lastSales[0]["realAmt"])
	var beforeLastExpectAmt int
	if len(lastSales) > 0{
		beforeLastExpectAmt, _ = strconv.Atoi(lastSales[0]["expectAmt"])
	}else if regInfo[0]["lnOpenDt"] == today{
		data["aprvMonthSum"] = 0
		data["aprvCnt"] = 0
		data["arpuAmt"] = 0
		data["busyTime"] = ""
		data["cardAmt"] = 0
		data["cashAmt"] = 0
		data["diffCardAmt"] = 0
		data["expectMsg"] = "수집전 예상 매출"
		data["expectArrowMsg"] = "어제실적 수집전 이에요."
		data["lastAmt"] = 0
		data["lastMsg"] = "데이터 수집 중"
		data["lastArrowMsg"] = "수집완료 후 확인 가능"
		data["month"] = ""
		data["payAmt"] = 0
		data["todayPay"] = "매출내역 수집 중"
		data["payMonthSum"] = 0
		data["todayExpectAmt"] = 0
		data["visitRate"] = ""

		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultYn"] = "N"
		//m["resultYn"] = "Y"
		m["resultData"] = data

		return c.JSON(http.StatusOK, m)
	}

	// 매출 조회
	cardSale, err := cls.GetSelectData(homesql.SelectCardAmt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	cashSale, err := cls.GetSelectData(homesql.SelectCashAmt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	pcaSale, err := cls.GetSelectData(homesql.SelectPrePcaAmt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	payASale, err := cls.GetSelectData(homesql.SelectPayAmt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "trDt")

	var cardAmt, cashAmt, pcaAmt, payAmt int
	if len(cardSale) > 0{
		cardAmt, _ = strconv.Atoi(cardSale[0]["aprvAmt"])
	}
	if len(cashSale) > 0{
		cashAmt, _ = strconv.Atoi(cashSale[0]["cashAmt"])
	}
	if len(pcaSale) > 0{
		pcaAmt, _ = strconv.Atoi(pcaSale[0]["pcaAmt"])
	}
	if len(payASale) > 0{
		payAmt, _ = strconv.Atoi(payASale[0]["payAmt"])
	}

	lastAmt := cardAmt+cashAmt

	// 어제 실제 매출액
	data["lastAmt"] = lastAmt
	//if lastAmt != 0 && lastAmt > beforeLastExpectAmt {
	if lastAmt > 0 {
		//lastDiff := (float64(lastAmt) - float64(beforeLastExpectAmt)) / float64(beforeLastExpectAmt) * 100
		//data["lastMsg"] = fmt.Sprintf("잘했어요! 예상보다 %0.f%%높아요", math.Round(lastDiff))

		if lastAmt < int(float64(beforeLastExpectAmt) / 100 * 15){
			data["lastMsg"] = "어제 휴일이셨나요?"
			data["lastArrowMsg"] = fmt.Sprintf("%s 매출", monthDayName(tn))
		}else{
			if cardAmt > 0 && cashAmt > 0{
				data["lastMsg"] = "카드,현금 매출"
			}else if cashAmt > 0{
				data["lastMsg"] = "현금 매출"
			}else if cardAmt > 0{
				data["lastMsg"] = "카드 매출"
			}

			data["lastArrowMsg"] = fmt.Sprintf("%s 매출", monthDayName(tn))
		}
	} else {
		// 여신데이터 수집 했는데 0원인 경우
		if tn.Hour() < 11 || regInfo[0]["regDt"] == today{
			data["lastMsg"] = "데이터 수집 중"
			data["lastArrowMsg"] = "수집완료 후 확인 가능"
		}else{
			data["lastMsg"] = "어제 휴일이셨나요?"
			data["lastArrowMsg"] = fmt.Sprintf("%s 매출", monthDayName(tn))
		}
	}
	/*
	else{
		// lastAmt(어제 매출)이 beforeLastExpectAmt(예상 매출)보다 현저히(15% 미만?) 낮을경우 휴일로 판단?
		if lastAmt < int(float64(beforeLastExpectAmt) / 100 * 15){
			data["lastMsg"] = "어제 휴일이셨나요?"
		}else{
			data["lastMsg"] = fmt.Sprintf("오늘은 좋은일만 가득하세요!")
		}

		// 6월 2일 수요일 매출
		data["lastArrowMsg"] = fmt.Sprintf("%s 매출", monthDayName(tn))
	}
	 */

	data["cardAmt"] = cardAmt
	data["cashAmt"] = cashAmt
	data["payAmt"] = payAmt
	// 입금 - 매입
	//data["diffCardAmt"] = pcaAmt - payAmt
	data["diffCardAmt"] = payAmt - pcaAmt

	// 오늘 매출 예측
	params["trDt"] = today
	todaySales, err := cls.GetSelectData(homesql.SelectDaySaleData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	/*
	if todaySales == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}
	*/

	todayDateName, err := cls.GetSelectData(homesql.SelectDayName, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	var todayExpectAmt int
	if len(todaySales) > 0{
		todayExpectAmt, _ = strconv.Atoi(todaySales[0]["expectAmt"])
	}

	todayPay, err := cls.GetSelectData(homesql.SelectTodayPay, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	tPay, _ := strconv.Atoi(todayPay[0]["amt"])
	if tPay > 0{
		data["todayPay"] = fmt.Sprintf("%d",tPay)
	}else{
		if tn.Hour() < 11{
			data["todayPay"] = "매출내역 수집 중"
		}else{
			data["todayPay"] = "0"
		}
	}

	data["todayExpectAmt"] = todayExpectAmt
	if lastAmt > 0 {
		//todayDiff := (float64(todayExpectAmt) - float64(lastAmt)) / float64(lastAmt) * 100
		//data["expectMsg"] = fmt.Sprintf("전일보다 %0.f%% 예상됩니다.", math.Round(todayDiff))

		// 오늘 휴일이면?
		if todayExpectAmt == 0{
			data["expectMsg"] = "오늘 휴일인가요?"
			data["expectArrowMsg"] = "오늘 영업하신다면..."
		}else{
			var dayName string
			switch tn.Weekday() {
			case 0:
				dayName = "일요일"
			case 1:
				dayName = "월요일"
			case 2:
				dayName = "화요일"
			case 3:
				dayName = "수요일"
			case 4:
				dayName = "목요일"
			case 5:
				dayName = "금요일"
			case 6:
				dayName = "토요일"
			}

			data["expectMsg"] = fmt.Sprintf("%s 예상 매출", dayName)
			if len(todayDateName) > 0 && len(todayDateName[0]["datename"]) > 0{
				data["expectArrowMsg"] = fmt.Sprintf("오늘은 %s!", todayDateName[0]["datename"])
			}else{
				if todayExpectAmt > lastAmt{
					data["expectArrowMsg"] = "어제보다 바쁠것 같아요."
					//data["expectArrowMsg"] = "전국동시지방선거에는 어제보다 바쁠것 같아요." -> 너무 김 화면 꽉참
					//data["expectArrowMsg"] = "작년 전국동시지방선거"
				}else{
					data["expectArrowMsg"] = "어제보다 여유가 있어요."
				}
			}
		}
	} else {
		if tn.Hour() < 11{
			data["expectMsg"] = "수집전 예상 매출"
			data["expectArrowMsg"] = "어제실적 수집전 이에요."
		}else{
			// 어제가 0원 혹은 휴일
			data["expectMsg"] = "오늘의 매출 예측"
			data["expectArrowMsg"] = "희망찬 하루 되세요!"
		}
	}

	// 월 누적 매출/입금 조회
	aprvSaleSum, err := cls.GetSelectData(homesql.SelectSaleSum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	aprvCashSum, err := cls.GetSelectData(homesql.SelectCashSum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	paySaleSum, err := cls.GetSelectData(homesql.SelectPaySum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "trDt")

	var aprvCardMonthSum, aprvCashMonthSum, payMonthSum int
	if len(aprvSaleSum) > 0{
		aprvCardMonthSum, _ = strconv.Atoi(aprvSaleSum[0]["aprvSum"])
	}
	if len(aprvCashSum) > 0{
		aprvCashMonthSum, _ = strconv.Atoi(aprvCashSum[0]["aprvSum"])
	}
	if len(paySaleSum) > 0{
		payMonthSum, _ = strconv.Atoi(paySaleSum[0]["paySum"])
	}

	data["month"] = today[:6]
	data["aprvMonthSum"] = aprvCardMonthSum + aprvCashMonthSum
	data["payMonthSum"] = payMonthSum

	// 바쁜 시간 예측
	params["startDt"] = time.Now().AddDate(0, 0, -29).Format("20060102") // 어제 기준 4주전
	params["endDt"] = yesterday
	expectBusyTime, err := cls.GetSelectData(homesql.SelectExpectBusyTime, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if len(expectBusyTime) > 0{
		var maxTime string
		var maxValue int
		for key, value := range expectBusyTime[0] {
			if maxTime == "" {
				maxTime = key
				maxValue, _ = strconv.Atoi(value)
			} else {
				newValue, _ := strconv.Atoi(value)
				if maxValue < newValue {
					maxTime = key
					maxValue = newValue
				}
			}
		}
		bytes := []byte(maxTime)
		busyTime := fmt.Sprintf("%s:00~%s:00", bytes[1:3], bytes[3:])
		data["busyTime"] = busyTime
	}else{
		data["busyTime"] = ""
	}


	// 객단가 예측
	arpuData, err := cls.GetSelectData(datasql.SelectAverageRevenuePerUser, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	cashArpuData, err := cls.GetSelectData(datasql.SelectCashAverageRevenuePerUser, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "startDt")
	delete(params, "endDt")

	var arpuAmt, cashArpuAmt int
	if len(arpuData) > 0{
		arpuAmt, _ = strconv.Atoi(arpuData[0]["arpu"])
	}
	if len(cashArpuData) > 0{
		cashArpuAmt, _ = strconv.Atoi(cashArpuData[0]["arpu"])
	}
	data["arpuAmt"] = arpuAmt + cashArpuAmt

	// 단골 비율 예측 -> 결제 건수 예측 변경
	//params["startDt"] = time.Now().AddDate(0, 0, -113).Format("20060102") // 어제 기준 12주전
	params["startDt"] = time.Now().AddDate(0, 0, -28).Format("20060102") // 어제 기준 12주전
	params["endDt"] = yesterday
	visitData, err := cls.GetSelectData(homesql.SelectExpectCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "startDt")
	delete(params, "endDt")

	//visitTotal, _ := strconv.Atoi(visitData[0]["visitTotal"])
	//visit2, _ := strconv.Atoi(visitData[0]["visit2"])
	//data["visitRate"] = fmt.Sprintf("%.0f%%(2회 이상)", (float32(visit2)/float32(visitTotal))*100)
	data["visitRate"] = ""
	if len(visitData) > 0{
		avgCnt,_ := strconv.Atoi(visitData[0]["avg_cnt"])
		data["aprvCnt"] = avgCnt
	}else{
		data["aprvCnt"] = 0
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultYn"] = "Y"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// day format
func monthDayName(tn time.Time) string{

	yesterdayTn := tn.AddDate(0, 0, -1)
	fullDay := yesterdayTn.Format("20060102")

	month := fullDay[4:6]
	day := fullDay[6:]

	var dayName string
	switch yesterdayTn.Weekday() {
		case 0:
			dayName = "일요일"
		case 1:
			dayName = "월요일"
		case 2:
			dayName = "화요일"
		case 3:
			dayName = "수요일"
		case 4:
			dayName = "목요일"
		case 5:
			dayName = "금요일"
		case 6:
			dayName = "토요일"
	}

	lprintf(4, "[INFO] month(%s), day(%s), dayIndex(%d), dayName(%s)", month, day, tn.Weekday(), dayName)

	var returnDay string
	if month[0] == '0'{
		returnDay = fmt.Sprintf("%v월 ", month[1:])
	}else{
		returnDay = fmt.Sprintf("%v월 ", month)
	}

	if day[0] == '0'{
		returnDay += fmt.Sprintf("%v일 ", day[1:])
	}else{
		returnDay += fmt.Sprintf("%v일 ", day)
	}

	returnDay += dayName

	// 6월 2일 수요일
	return returnDay
}

func TtiName(tti string) string{

	lprintf(4, "[INFO] tti(%s) \n", tti)

	var ttiName string

	switch tti {
	case "00":
		ttiName = "쥐"
	case "01":
		ttiName = "소"
	case "02":
		ttiName = "호랑이"
	case "03":
		ttiName = "토끼"
	case "04":
		ttiName = "용"
	case "05":
		ttiName = "뱀"
	case "06":
		ttiName = "말"
	case "07":
		ttiName = "양"
	case "08":
		ttiName = "원숭이"
	case "09":
		ttiName = "닭"
	case "10":
		ttiName = "개"
	case "11":
		ttiName = "돼지"
	}

	return fmt.Sprintf("%s띠 보러가기", ttiName)
}