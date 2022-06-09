package sales

import (
	"cashApi/src/controller"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"strconv"
	"strings"
	"time"

	salesql "cashApi/query/sales"

	// login 및 기본

	"cashApi/src/controller/cls"

	"github.com/labstack/echo/v4"
)

var RedisAddr []string
var PartnerUrl string

func RedisConfig(fname string) {

	v, r := cls.GetTokenValue("REDIS_INFO", fname)
	if r != cls.CONF_ERR {
		rCnt, err := strconv.Atoi(v)
		if err == nil && rCnt > 0 {
			for i := 0; i < rCnt; i++ {
				v, r = cls.GetTokenValue(fmt.Sprintf("REDIS_INF0%d", i), fname)
				if r != cls.CONF_ERR {
					rConfig := strings.Split(v, ",")
					if len(rConfig) == 2 {
						RedisAddr = append(RedisAddr, fmt.Sprintf("%s:%s", rConfig[1], rConfig[0]))
					}
				}
			}
		}
	}

	PartnerUrl = "https://partner.darayo.com/redis/update"

	v, r = cls.GetTokenValue("PARTNER_URL", fname)
	if r != cls.CONF_ERR {
		PartnerUrl = v
	}

}

/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
// darayo printf(로그레벨, 요청 컨텍스트, format, arg) => 무엇을(서비스, 요청), 어떻게(input), 왜(원인,조치)
var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func BaseUrl(c echo.Context) error {
	return c.JSONP(http.StatusOK, "", "homes")
}

// 카드 승인 합계
func GetCardSum(c echo.Context) error {

	dprintf(4, c, "call GetCardSum\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectData(salesql.SelectDayCardSum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	var allCnt, allAmt int

	if resultList == nil {
		s := []string{}
		data["list"] = s
		data["sumCnt"] = 0
		data["sumAmt"] = 0
	} else {
		var list []map[string]interface{}
		for _, result := range resultList {
			totCnt, _ := strconv.Atoi(result["totCnt"])
			totAmt, _ := strconv.Atoi(result["totAmt"])
			aprvCnt, _ := strconv.Atoi(result["aprvCnt"])
			aprvAmt, _ := strconv.Atoi(result["aprvAmt"])
			canCnt, _ := strconv.Atoi(result["canCnt"])
			canAmt, _ := strconv.Atoi(result["canAmt"])

			allCnt = allCnt + totCnt
			allAmt = allAmt + totAmt

			tmp := make(map[string]interface{})
			tmp["trDt"] = result["trDt"]
			tmp["totCnt"] = totCnt
			tmp["totAmt"] = totAmt
			tmp["aprvCnt"] = aprvCnt
			tmp["aprvAmt"] = aprvAmt
			tmp["canCnt"] = canCnt
			tmp["canAmt"] = canAmt

			list = append(list, tmp)
		}
		data["list"] = list
		data["sumCnt"] = allCnt
		data["sumAmt"] = allAmt
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 카드 승인 상세
func GetCardList(c echo.Context) error {

	dprintf(4, c, "call GetCardList\n")

	params := cls.GetParamJsonMap(c)
	pagePerRow, _ := strconv.Atoi(params["pagePerRow"])
	pageNo, _ := strconv.Atoi(params["pageNo"])
	params["startNum"] = strconv.Itoa((pageNo-1)*pagePerRow + 1)
	params["endNum"] = strconv.Itoa((pageNo-1)*pagePerRow + 1 + pagePerRow)

	resultList, err := cls.GetSelectType(salesql.SelectDayCardDetail, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	resultListCnt, err := cls.GetSelectData(salesql.SelectDayCardDetailCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	if resultList == nil {
		s := []string{}
		data["list"] = s
	} else {
		data["list"] = resultList
	}

	// page 처리
	paging := make(map[string]interface{})
	totCnt, _ := strconv.Atoi(resultListCnt[0]["totCnt"])
	paging["totalPage"] = strconv.Itoa((totCnt / pagePerRow) + 1)
	paging["currentPage"] = strconv.Itoa(pageNo)
	data["paging"] = paging

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 현금영수증 승인 합계
func GetCashSum(c echo.Context) error {

	dprintf(4, c, "call GetCashSum\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectData(salesql.SelectDayCashSum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	var allCnt, allAmt int

	if resultList == nil {
		s := []string{}
		data["list"] = s
		data["sumCnt"] = 0
		data["sumAmt"] = 0
	} else {
		var list []map[string]interface{}
		for _, result := range resultList {
			totCnt, _ := strconv.Atoi(result["totCnt"])
			totAmt, _ := strconv.Atoi(result["totAmt"])
			aprvCnt, _ := strconv.Atoi(result["aprvCnt"])
			aprvAmt, _ := strconv.Atoi(result["aprvAmt"])
			canCnt, _ := strconv.Atoi(result["canCnt"])
			canAmt, _ := strconv.Atoi(result["canAmt"])

			allCnt = allCnt + totCnt
			allAmt = allAmt + totAmt

			tmp := make(map[string]interface{})
			tmp["trDt"] = result["trDt"]
			tmp["totCnt"] = totCnt

			tmp["totAmt"] = totAmt
			tmp["aprvCnt"] = aprvCnt
			tmp["aprvAmt"] = aprvAmt
			tmp["canCnt"] = canCnt
			tmp["canAmt"] = canAmt

			list = append(list, tmp)
		}
		data["list"] = list
		data["sumCnt"] = allCnt
		data["sumAmt"] = allAmt
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 현금영수증 승인 상세
func GetCashList(c echo.Context) error {

	dprintf(4, c, "call GetCashList\n")

	params := cls.GetParamJsonMap(c)
	pagePerRow, _ := strconv.Atoi(params["pagePerRow"])
	pageNo, _ := strconv.Atoi(params["pageNo"])
	params["startNum"] = strconv.Itoa((pageNo-1)*pagePerRow + 1)
	params["endNum"] = strconv.Itoa((pageNo-1)*pagePerRow + 1 + pagePerRow)

	resultList, err := cls.GetSelectType(salesql.SelectDayCashDetail, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	resultListCnt, err := cls.GetSelectData(salesql.SelectDayCashDetailCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	if resultList == nil {
		s := []string{}
		data["list"] = s
	} else {
		data["list"] = resultList
	}

	// page 처리
	paging := make(map[string]interface{})
	totCnt, _ := strconv.Atoi(resultListCnt[0]["totCnt"])
	paging["totalPage"] = strconv.Itoa((totCnt / pagePerRow) + 1)
	paging["currentPage"] = strconv.Itoa(pageNo)
	data["paging"] = paging

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 입금정보 합계
func GetPaySum(c echo.Context) error {

	dprintf(4, c, "call GetPaySum\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectData(salesql.SelectDayPaySum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	var allAmt int

	if resultList == nil {
		s := []string{}
		data["list"] = s
		data["sumCnt"] = 0
		data["sumAmt"] = 0
	} else {
		var list []map[string]interface{}
		for _, result := range resultList {
			pcaCnt, _ := strconv.Atoi(result["pcaCnt"])
			pcaAmt, _ := strconv.Atoi(result["pcaAmt"])
			realPayAmt, _ := strconv.Atoi(result["realPayAmt"])
			delayAmt, _ := strconv.Atoi(result["delayAmt"])

			allAmt = allAmt + realPayAmt

			tmp := make(map[string]interface{})
			tmp["trDt"] = result["trDt"]
			tmp["pcaCnt"] = pcaCnt
			tmp["pcaAmt"] = pcaAmt
			tmp["realPayAmt"] = realPayAmt
			tmp["delayAmt"] = delayAmt
			list = append(list, tmp)
		}
		data["list"] = list
		data["sumAmt"] = allAmt
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 입금정보 상세
func GetPayList(c echo.Context) error {

	dprintf(4, c, "call GetPayList\n")

	params := cls.GetParamJsonMap(c)
	pagePerRow, _ := strconv.Atoi(params["pagePerRow"])
	pageNo, _ := strconv.Atoi(params["pageNo"])
	params["startNum"] = strconv.Itoa((pageNo-1)*pagePerRow + 1)
	params["endNum"] = strconv.Itoa((pageNo-1)*pagePerRow + 1 + pagePerRow)

	resultList, err := cls.GetSelectType(salesql.SelectDayPayDetail, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	resultListCnt, err := cls.GetSelectData(salesql.SelectDayPayDetailCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	if resultList == nil {
		s := []string{}
		data["list"] = s
	} else {
		data["list"] = resultList
	}

	// page 처리
	paging := make(map[string]interface{})
	totCnt, _ := strconv.Atoi(resultListCnt[0]["totCnt"])
	paging["totalPage"] = strconv.Itoa((totCnt / pagePerRow) + 1)
	paging["currentPage"] = strconv.Itoa(pageNo)
	data["paging"] = paging

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

func redisConnect(addr string) redis.Conn {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		lprintf(1, "[ERROR] redis con err(%s) \n", err.Error())
		return nil
	}

	pong, err := redis.String(c.Do("PING"))
	if err != nil {
		lprintf(1, "[ERROR] redis ping pong err(%s) \n", err.Error())
		c.Close()
		return nil
	}

	lprintf(4, "[INFO] redis con ping pong(%s)\n", pong)
	return c
}

func RedisGet(key string) (int, string) {

	for _, addr := range RedisAddr {
		c := redisConnect(addr)
		if c == nil {
			continue
		}

		reply, err := redis.String(c.Do("GET", key))
		if err != nil {
			lprintf(1, "[ERROR] redis get err(%s)\n", err.Error())
			c.Close()
			continue
		}

		c.Close()
		return 1, reply
	}

	return -1, ""
}

// 매출캘린더
func GetAprvCalendar(c echo.Context) error {

	dprintf(4, c, "call GetAprvCalendar\n")

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})

	aprvKey := fmt.Sprintf("%saprvCalender", params["bizNum"])
	rst, aprvCalender := RedisGet(aprvKey)
	if rst > 0 {
		data := make(map[string]interface{})
		if err := json.Unmarshal([]byte(aprvCalender), &data); err == nil {
			m["resultCode"] = "00"
			m["resultMsg"] = "응답 성공"
			m["resultData"] = data

			return c.JSON(http.StatusOK, m)
		}
	}

	if len(PartnerUrl) > 0 {
		go func() {
			resp, _ := http.Get(fmt.Sprintf("%s?bizNum=%s", PartnerUrl, params["bizNum"]))
			resp.Body.Close()
		}()
	}

	t := time.Now()
	params["startDt"] = fmt.Sprintf("%s01", t.AddDate(0, -8, 0).Format("200601"))
	//params["startDt"] = fmt.Sprintf("%s01", t.AddDate(0, -7, 0).Format("200601"))
	//params["endDt"] = t.Format("20060102")
	params["endDt"] = t.AddDate(0, 0, -1).Format("20060102")

	aprvSumList, err := cls.GetSelectDataUsingJson(salesql.SelectAprvCalendarSumList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "startDt")
	delete(params, "endDt")

	data := make(map[string]interface{})

	var summary []map[string]interface{}
	for _, sumData := range aprvSumList {
		tmp := make(map[string]interface{})
		aprv, _ := strconv.Atoi(sumData["aprvAmt"])
		cash, _ := strconv.Atoi(sumData["cashAmt"])
		pca, _ := strconv.Atoi(sumData["pcaAmt"])
		tot, _ := strconv.Atoi(sumData["totAmt"])

		tmp["trMonth"] = sumData["trMonth"][4:]
		tmp["aprvAmt"] = aprv
		tmp["cashAmt"] = cash
		tmp["pcaAmt"] = pca
		tmp["totAmt"] = tot

		summary = append(summary, tmp)
	}
	data["summary"] = summary

	var monthList []map[string]interface{}
	for idx, sumData := range aprvSumList {
		if idx == len(aprvSumList)-2 {
			break
		}

		dprintf(4, c, "trMonth=%s\n", sumData["trMonth"])
		// 날짜변경을 위해 Time 값으로 변경
		timeTrDt, err := time.Parse("20060102", fmt.Sprintf("%s01", sumData["trMonth"]))
		timeFirst, timeLast := cls.GetFirstAndLastOfMonth(timeTrDt)
		firstDay := cls.GetFirstOfWeek(timeFirst).Format("20060102")
		lastDay := cls.GetEndOfWeek(timeLast).Format("20060102")

		params["startDt"] = firstDay
		params["endDt"] = lastDay
		aprvList, err := cls.GetSelectDataUsingJson(salesql.SelectAprvCalendarList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}
		delete(params, "startDt")
		delete(params, "endDt")

		var dayList []map[string]interface{}
		for _, aprvData := range aprvList {
			tmp := make(map[string]interface{})
			row, _ := strconv.Atoi(aprvData["rNum"])
			aprv, _ := strconv.Atoi(aprvData["aprvAmt"])
			cash, _ := strconv.Atoi(aprvData["cashAmt"])
			pca, _ := strconv.Atoi(aprvData["pcaAmt"])
			tot, _ := strconv.Atoi(aprvData["totAmt"])

			tmp["rNum"] = row
			tmp["trDt"] = aprvData["trDt"]
			tmp["aprvAmt"] = aprv
			tmp["cashAmt"] = cash
			tmp["pcaAmt"] = pca
			tmp["totAmt"] = tot
			tmp["diffColor"] = aprvData["diffColor"]
			tmp["dayColor"] = aprvData["dayColor"]

			dayList = append(dayList, tmp)
		}
		monthData := make(map[string]interface{})
		monthData["trMonth"] = sumData["trMonth"]
		monthData["dayList"] = dayList
		monthList = append(monthList, monthData)
	}
	data["monthList"] = monthList

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 매출캘린더 카드사별 매입내역 리스트
func GetAprvDailyList(c echo.Context) error {

	dprintf(4, c, "call GetAprvDailyList\n")

	params := cls.GetParamJsonMap(c)
	aprvDailyList, err := cls.GetSelectDataUsingJson(salesql.SelectAprvDailyList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	cashDailyList, err := cls.GetSelectData(salesql.SelectCashDailyList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 합계건수, 금액, 매입상태 생성
	var aprvList []map[string]interface{}
	var aprvCnt, aprvAmt, pcaCnt, pcaAmt, diffAmt, totFee, vatAmt, outpExptAmt int
	var diff1SumCnt, diff2SumCnt, diff3SumCnt, diff4SumCnt int
	for idx, aprvData := range aprvDailyList {
		data := make(map[string]interface{})
		data["rNum"] = idx + 3
		data["cardCd"] = aprvData["cardCd"]
		data["cardNm"] = aprvData["cardNm"]

		tmp, _ := strconv.Atoi(aprvData["aprvCnt"])
		aprvCnt = aprvCnt + tmp
		data["aprvCnt"] = tmp
		tmp, _ = strconv.Atoi(aprvData["aprvAmt"])
		aprvAmt = aprvAmt + tmp
		data["aprvAmt"] = tmp
		tmp, _ = strconv.Atoi(aprvData["pcaCnt"])
		pcaCnt = pcaCnt + tmp
		data["pcaCnt"] = tmp
		tmp, _ = strconv.Atoi(aprvData["pcaAmt"])
		pcaAmt = pcaAmt + tmp
		data["pcaAmt"] = tmp
		tmp, _ = strconv.Atoi(aprvData["diffAmt"])
		diffAmt = diffAmt + tmp
		data["diffAmt"] = tmp
		tmp, _ = strconv.Atoi(aprvData["totFee"])
		totFee = totFee + tmp
		data["totFee"] = tmp
		tmp, _ = strconv.Atoi(aprvData["vatAmt"])
		vatAmt = vatAmt + tmp
		data["vatAmt"] = tmp
		tmp, _ = strconv.Atoi(aprvData["payAmt"])
		outpExptAmt = outpExptAmt + tmp
		data["outpExptAmt"] = tmp

		var diff1Cnt, diff2Cnt, diff3Cnt, diff4Cnt int
		if aprvData["diffInf"][:1] == "1" { // 1:정상매입
			diff1Cnt, _ = strconv.Atoi(aprvData["diffInf"][2:])
			diff1SumCnt = diff1SumCnt + diff1Cnt
		} else if aprvData["diffInf"][:1] == "2" { // 2:매입제외
			idx := strings.Index(aprvData["diffInf"], "!")
			diff2Cnt, _ = strconv.Atoi(aprvData["diffInf"][2:idx])
			diff2SumCnt = diff2SumCnt + diff2Cnt

			diff1Cnt, _ = strconv.Atoi(aprvData["diffInf"][idx+3:])
			diff1SumCnt = diff1SumCnt + diff1Cnt
		} else if aprvData["diffInf"][:1] == "3" { // 3:매입대기
			diff3Cnt, _ = strconv.Atoi(aprvData["diffInf"][2:])
			diff3SumCnt = diff3SumCnt + diff3Cnt
		} else if aprvData["diffInf"][:1] == "4" { // 4:추가매입
			idx := strings.Index(aprvData["diffInf"], "!")
			diff4Cnt, _ = strconv.Atoi(aprvData["diffInf"][2:idx])
			diff4SumCnt = diff4SumCnt + diff4Cnt

			diff1Cnt, _ = strconv.Atoi(aprvData["diffInf"][idx+3:])
			diff1SumCnt = diff1SumCnt + diff1Cnt
		} else {
			dprintf(4, c, "[ERROR] diffInf value, aprvData=%v\n", aprvData)
		}

		data["diffInf0"] = 0
		data["diffInf1"] = diff1Cnt
		data["diffInf2"] = diff2Cnt
		data["diffInf3"] = diff3Cnt
		data["diffInf4"] = diff4Cnt

		aprvList = append(aprvList, data)
	}

	cashCnt, _ := strconv.Atoi(cashDailyList[0]["aprvCnt"])
	//aprvCnt = aprvCnt
	//pcaCnt = pcaCnt
	cashAmt, _ := strconv.Atoi(cashDailyList[0]["aprvAmt"])
	//aprvAmt = aprvAmt
	//pcaAmt = pcaAmt
	//outpExptAmt = outpExptAmt
	outpExptAmt = outpExptAmt + cashAmt

	var diffCash1Cnt int
	if cashDailyList[0]["diffInf"][:1] == "1" { // 1:정상매입
		diffCash1Cnt, _ = strconv.Atoi(cashDailyList[0]["diffInf"][2:])
	}

	var aprvAllList []map[string]interface{}

	// 합계
	sum := make(map[string]interface{})
	sum["rNum"] = 1
	sum["cardCd"] = "99"
	sum["cardNm"] = "합계"
	sum["aprvCnt"] = aprvCnt
	sum["aprvAmt"] = aprvAmt
	sum["pcaCnt"] = pcaCnt
	sum["pcaAmt"] = pcaAmt
	sum["diffAmt"] = diffAmt
	sum["totFee"] = totFee
	sum["vatAmt"] = vatAmt
	sum["outpExptAmt"] = outpExptAmt
	sum["diffInf0"] = 0
	sum["diffInf1"] = diff1SumCnt
	sum["diffInf2"] = diff2SumCnt
	sum["diffInf3"] = diff3SumCnt
	sum["diffInf4"] = diff4SumCnt
	aprvAllList = append(aprvAllList, sum)

	// 현금영수증
	cash := make(map[string]interface{})
	cash["rNum"] = 2
	cash["cardCd"] = "00"
	cash["cardNm"] = "현금영수증"
	cash["aprvCnt"] = cashCnt
	cash["aprvAmt"] = cashAmt
	cash["pcaCnt"] = cashCnt
	cash["pcaAmt"] = cashAmt
	cash["diffAmt"] = 0
	cash["totFee"] = 0
	cash["vatAmt"] = 0
	cash["outpExptAmt"] = cashAmt
	cash["diffInf0"] = 0
	cash["diffInf1"] = diffCash1Cnt
	cash["diffInf2"] = 0
	cash["diffInf3"] = 0
	cash["diffInf4"] = 0
	aprvAllList = append(aprvAllList, cash)

	// 신용카드
	aprvAllList = append(aprvAllList, aprvList...)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = aprvAllList

	return c.JSON(http.StatusOK, m)
}

// 매출캘린더 특정 카드사 매입내역 리스트
func GetAprvDetailList(c echo.Context) error {

	dprintf(4, c, "call GetAprvDetailList\n")

	params := cls.GetParamJsonMap(c)

	var err error
	var aprvDetail []map[string]string
	if params["cardCd"] == "00" { // 현금영수증
		aprvDetail, err = cls.GetSelectDataUsingJson(salesql.SelectCashDetailList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}
	} else { // 체크/신용카드
		aprvDetail, err = cls.GetSelectDataUsingJson(salesql.SelectAprvDetailList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}
	}

	var aprvDetailList []map[string]interface{}
	for idx, detail := range aprvDetail {
		data := make(map[string]interface{})
		data["rNum"] = idx + 1
		data["trDt"] = detail["trDt"]
		data["trTm"] = detail["trTm"]
		data["diffNm"] = detail["diffNm"]
		data["aprvNo"] = detail["aprvNo"]
		data["cardNo"] = detail["cardNo"]
		data["stsCd"] = detail["stsCd"]
		data["instTrm"] = detail["instTrm"]
		data["cardKndNm"] = detail["cardKndNm"]
		data["outpExptDt"] = detail["outpExptDt"]

		aprvAmt, _ := strconv.Atoi(detail["aprvAmt"])
		pcaAmt, _ := strconv.Atoi(detail["pcaAmt"])
		vatAmt, _ := strconv.Atoi(detail["vatAmt"])
		totFee, _ := strconv.Atoi(detail["totFee"])
		payAmt, _ := strconv.Atoi(detail["payAmt"])

		data["aprvAmt"] = aprvAmt
		data["pcaAmt"] = pcaAmt
		data["vatAmt"] = vatAmt
		data["totFee"] = totFee
		data["payAmt"] = payAmt

		aprvDetailList = append(aprvDetailList, data)
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = aprvDetailList

	return c.JSON(http.StatusOK, m)
}

// 입금캘린더
func GetPayCalendar(c echo.Context) error {

	dprintf(4, c, "call GetPayCalendar\n")

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})

	payKey := fmt.Sprintf("%spayCalender", params["bizNum"])
	rst, aprvCalender := RedisGet(payKey)
	if rst > 0 {
		data := make(map[string]interface{})
		if err := json.Unmarshal([]byte(aprvCalender), &data); err == nil {
			m["resultCode"] = "00"
			m["resultMsg"] = "응답 성공"
			m["resultData"] = data

			return c.JSON(http.StatusOK, m)
		}
	}

	if len(PartnerUrl) > 0 {
		go func() {
			resp, _ := http.Get(fmt.Sprintf("%s?bizNum=%s", PartnerUrl, params["bizNum"]))
			resp.Body.Close()
		}()
	}

	t := time.Now()
	params["startDt"] = fmt.Sprintf("%s01", t.AddDate(0, -8, 0).Format("200601"))
	//params["startDt"] = fmt.Sprintf("%s01", t.AddDate(0, -7, 0).Format("200601"))
	params["endDt"] = t.AddDate(0, 0, -1).Format("20060102")
	//params["endDt"] = "20210430"

	paySumList, err := cls.GetSelectDataUsingJson(salesql.SelectPayCalendarSumList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	delete(params, "startDt")
	delete(params, "endDt")

	data := make(map[string]interface{})

	var summary []map[string]interface{}
	for _, sumData := range paySumList {
		tmp := make(map[string]interface{})
		expt, _ := strconv.Atoi(sumData["outpExptAmt"])
		realIn, _ := strconv.Atoi(sumData["realInAmt"])
		diff, _ := strconv.Atoi(sumData["diffAmt"])

		tmp["trMonth"] = sumData["trMonth"][4:]
		tmp["outpExptAmt"] = expt
		tmp["realInAmt"] = realIn
		tmp["diffAmt"] = diff
		tmp["diffColor"] = sumData["diffColor"]

		summary = append(summary, tmp)
	}
	data["summary"] = summary

	var monthList []map[string]interface{}
	for idx, sumData := range paySumList {
		if idx == len(paySumList)-2 {
			break
		}

		dprintf(4, c, "trMonth=%s\n", sumData["trMonth"])
		// 날짜변경을 위해 Time 값으로 변경
		timeTrDt, err := time.Parse("20060102", fmt.Sprintf("%s01", sumData["trMonth"]))
		timeFirst, timeLast := cls.GetFirstAndLastOfMonth(timeTrDt)
		firstDay := cls.GetFirstOfWeek(timeFirst).Format("20060102")
		lastDay := cls.GetEndOfWeek(timeLast).Format("20060102")

		params["startDt"] = firstDay
		params["endDt"] = lastDay
		payList, err := cls.GetSelectDataUsingJson(salesql.SelectPayCalendarList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}
		delete(params, "startDt")
		delete(params, "endDt")

		var dayList []map[string]interface{}
		for _, payData := range payList {
			tmp := make(map[string]interface{})
			row, _ := strconv.Atoi(payData["rNum"])
			expt, _ := strconv.Atoi(payData["outpExptAmt"])
			realIn, _ := strconv.Atoi(payData["realInAmt"])
			diff, _ := strconv.Atoi(payData["diffAmt"])

			tmp["rNum"] = row
			tmp["trDt"] = payData["trDt"]
			tmp["outpExptAmt"] = expt
			tmp["realInAmt"] = realIn
			tmp["diffAmt"] = diff
			tmp["diffColor"] = payData["diffColor"]
			tmp["dayColor"] = payData["dayColor"]

			dayList = append(dayList, tmp)
		}
		monthData := make(map[string]interface{})
		monthData["trMonth"] = sumData["trMonth"]
		monthData["dayList"] = dayList
		monthList = append(monthList, monthData)
	}
	data["monthList"] = monthList

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 입금캘린더 카드사별 매입내역 리스트
func GetPayDailyList(c echo.Context) error {

	dprintf(4, c, "call GetPayDailyList\n")

	params := cls.GetParamJsonMap(c)
	payDailyList, err := cls.GetSelectDataUsingJson(salesql.SelectPayDailyList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	// 합계건수, 금액 생성
	var payList []map[string]interface{}
	var pcaCnt, pcaAmt, totFee, vatAmt, outpExptAmt, realInAmt, diffAmt int
	for idx, payData := range payDailyList {
		data := make(map[string]interface{})
		data["rNum"] = idx + 2
		data["cardCd"] = payData["cardCd"]
		data["cardNm"] = payData["cardNm"]

		tmp, _ := strconv.Atoi(payData["pcaCnt"])
		pcaCnt = pcaCnt + tmp
		data["pcaCnt"] = tmp
		tmp, _ = strconv.Atoi(payData["pcaAmt"])
		pcaAmt = pcaAmt + tmp
		data["pcaAmt"] = tmp
		tmp, _ = strconv.Atoi(payData["totFee"])
		totFee = totFee + tmp
		data["totFee"] = tmp
		tmp, _ = strconv.Atoi(payData["vatAmt"])
		vatAmt = vatAmt + tmp
		data["vatAmt"] = tmp
		tmp, _ = strconv.Atoi(payData["outpExptAmt"])
		outpExptAmt = outpExptAmt + tmp
		data["outpExptAmt"] = tmp
		tmp, _ = strconv.Atoi(payData["realInAmt"])
		realInAmt = realInAmt + tmp
		data["realInAmt"] = tmp
		tmp, _ = strconv.Atoi(payData["diffAmt"])
		diffAmt = diffAmt + tmp
		data["diffAmt"] = tmp

		data["diffNm"] = payData["diffNm"]
		data["diffColor"] = payData["diffColor"]

		payList = append(payList, data)
	}

	var payAllList []map[string]interface{}

	// 합계
	sum := make(map[string]interface{})
	sum["rNum"] = 1
	sum["cardCd"] = "99"
	sum["cardNm"] = "합계"
	sum["pcaCnt"] = pcaCnt
	sum["pcaAmt"] = pcaAmt
	sum["totFee"] = totFee
	sum["vatAmt"] = vatAmt
	sum["outpExptAmt"] = outpExptAmt
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
	payAllList = append(payAllList, sum)

	// 신용카드
	payAllList = append(payAllList, payList...)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = payAllList

	return c.JSON(http.StatusOK, m)
}

// 입금캘린더 특정 카드사 매입내역 리스트
func GetPayDetailList(c echo.Context) error {

	dprintf(4, c, "call GetPayDetailList\n")

	params := cls.GetParamJsonMap(c)
	payDetailSum, err := cls.GetSelectDataUsingJson(salesql.SelectPayDetailSum, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	payDetail, err := cls.GetSelectDataUsingJson(salesql.SelectPayDetailList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	data := make(map[string]interface{})
	mPayDetailSum := make(map[string]interface{})
	mPayDetailSum["cardCd"] = payDetailSum[0]["cardCd"]
	mPayDetailSum["cardNm"] = payDetailSum[0]["cardNm"]
	mPayDetailSum["diffAmt"], _ = strconv.Atoi(payDetailSum[0]["diffAmt"])
	mPayDetailSum["diffColor"] = payDetailSum[0]["diffColor"]
	mPayDetailSum["diffNm"] = payDetailSum[0]["diffNm"]
	mPayDetailSum["outpExptAmt"], _ = strconv.Atoi(payDetailSum[0]["outpExptAmt"])
	mPayDetailSum["outpExptDt"] = payDetailSum[0]["outpExptDt"]
	mPayDetailSum["pcaAmt"], _ = strconv.Atoi(payDetailSum[0]["pcaAmt"])
	mPayDetailSum["pcaCnt"], _ = strconv.Atoi(payDetailSum[0]["pcaCnt"])
	mPayDetailSum["realInAmt"], _ = strconv.Atoi(payDetailSum[0]["realInAmt"])
	mPayDetailSum["totFee"], _ = strconv.Atoi(payDetailSum[0]["totFee"])
	mPayDetailSum["vatAmt"], _ = strconv.Atoi(payDetailSum[0]["vatAmt"])

	data["summary"] = mPayDetailSum

	var payDetailList []map[string]interface{}
	for idx, detail := range payDetail {
		pay := make(map[string]interface{})
		pay["rNum"] = idx + 1
		pay["trDt"] = detail["trDt"]
		pay["aprvNo"] = detail["aprvNo"]
		pay["stsCd"] = detail["stsCd"]
		pay["instTrm"] = detail["instTrm"]
		pay["cardNo"] = detail["cardNo"]
		pay["cardKndNm"] = detail["cardKndNm"]

		pcaAmt, _ := strconv.Atoi(detail["pcaAmt"])
		outpExptAmt, _ := strconv.Atoi(detail["outpExptAmt"])

		pay["pcaAmt"] = pcaAmt
		pay["outpExptAmt"] = outpExptAmt

		payDetailList = append(payDetailList, pay)
	}
	data["list"] = payDetailList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}
