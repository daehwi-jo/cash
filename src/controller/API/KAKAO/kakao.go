package KAKAO

import (
	"bytes"
	commonsql "cashApi/query/commons"
	datasql "cashApi/query/datas"
	homesql "cashApi/query/homes"
	salesql "cashApi/query/sales"
	storesql "cashApi/query/stores"
	"cashApi/src/controller/cls"
	dataCon "cashApi/src/controller/datas"
	"encoding/json"
	"fmt"
	textrank "github.com/DavidBelicza/TextRank"
	humanize "github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	//템플릿 코드					//템플릿 명
	CASH_901 = "cash_901" // 가입상태 안내 미인증
	CASH_902 = "cash_902" // 가입상태 안내 인증

	CASH_201 = "cash_201" // 월간분석_성공
	CASH_205 = "cash_205" // 월간분석_성공5
	CASH_101 = "cash_101" // 주간분석_성공
	CASH_104 = "cash_104" // 주간분석_성공4

	CASH_016 = "cash_016" // 어제분석_주말실패
	CASH_015 = "cash_015" // 어제분석_2차실패
	CASH_014 = "cash_014" // 어제분석_1차실패
	CASH_013 = "cash_013" // 어제분석_성공
	CASH_019 = "cash_019" // 어제분석_성공4

	CASH_005 = "cash_005" // 여신인증
	CASH_004 = "cash_004" // 가입4번
	CASH_003 = "cash_003" // 가입3번
	CASH_002 = "cash_002" // 가입2번
	CASH_001 = "cash_001" // 가입1번
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf
var UserCode, DeptCode, YellowIdKey, SmsNo, KakaoEndpoint string

type KakaoAlimSend struct {
	Usercode    string          `json:"usercode"`
	Deptcode    string          `json:"deptcode"`
	YellowidKey string          `json:"yellowid_key"`
	Messages    []KakaoMessages `json:"messages"`
}

type KakaoMessages struct {
	Type         string         `json:"type"`
	MessageID    int            `json:"message_id"`
	To           string         `json:"to"`
	Callphone    string         `json:"callphone"`
	Text         string         `json:"text"`
	From         string         `json:"from"`
	TemplateCode string         `json:"template_code"`
	ReservedTime string         `json:"reserved_time"`
	ReSend       string         `json:"re_send"`
	ReTitle      string         `json:"re_title"`
	ReText       string         `json:"re_text"`
	Buttons      []KakaoButtons `json:"buttons"`
}

type KakaoButtons struct {
	ButtonType string `json:"button_type"`
	ButtonName string `json:"button_name"`
	ButtonURL  string `json:"button_url"`
	ButtonURL2 string `json:"button_url2"`
}

type KakaoResult struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Results []KakaoResultDesc `json:"results"`
}

type KakaoResultDesc struct {
	Result    string `json:"result"`
	MessageID int    `json:"message_id"`
}

type structCntSort struct {
	key string
	val int
}

func KakaoConfig(fname string) {
	tmp, err := cls.GetTokenValue("KAKAOALIMTAKL.USERCODE", fname)
	if err == cls.CONF_OK {
		UserCode = tmp
	}

	tmp, err = cls.GetTokenValue("KAKAOALIMTAKL.DEPTCODE", fname)
	if err == cls.CONF_OK {
		DeptCode = tmp
	}

	tmp, err = cls.GetTokenValue("KAKAOALIMTAKL.YELLOID_KEY", fname)
	if err == cls.CONF_OK {
		YellowIdKey = tmp
	}

	tmp, err = cls.GetTokenValue("SMS.API.CALLBACK", fname)
	if err == cls.CONF_OK {
		SmsNo = tmp
	}

	KakaoEndpoint = "biz/at/v2/json"
}

func addComma(v string) string {

	if len(v) == 0 || v == "0" {
		return "0원"
	}

	amt, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		lprintf(1, "[ERROR] parse int 64 err(%s)\n", err.Error())
		return "0원"
	}

	return fmt.Sprintf("%s원", humanize.Comma(amt))
}

func addComma2(v int) string {
	if v == 0 {
		return "0원"
	}

	return fmt.Sprintf("%s원", humanize.Comma(int64(v)))
}

func WelcomeMessage(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	restId := params["restId"]

	if len(restId) > 0 {
		Coupon_KakaoAlim(restId)
	}

	return c.JSON(http.StatusOK, nil)
}

func StateReport(c echo.Context) error {
	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	var rst int
	code := params["code"]
	restId := params["restId"]
	today := params["today"]

	if today != time.Now().Format("20060102") {
		m["resultCode"] = "98"
		m["resultMsg"] = "kakao report fail"
		return c.JSON(http.StatusOK, m)
	}

	lprintf(4, "[INFO] restId(%s), code(%s)\n", restId, code)

	if restId == "all" {

		var queryTmp string
		if code == CASH_901 {
			queryTmp = homesql.SelectFailStateReport
		} else if code == CASH_902 {
			queryTmp = homesql.SelectSuccessStateReport
		} else {
			return c.JSON(http.StatusOK, nil)
		}

		comps, err := cls.SelectData(queryTmp, params)
		if err != nil {
			lprintf(1, "[ERROR] SelectSuccessStateReport err(%s)\n", err.Error())
			return c.JSON(http.StatusOK, nil)
		}

		// 미 인증 comp에 pwd 틀린사람 추가
		if code == CASH_901 {
			params["bsDt"] = time.Now().AddDate(0, 0, -1).Format("20060102")

			failComps, err := cls.SelectData(homesql.SelectLoginFailReport, params)
			if err != nil {
				lprintf(1, "[ERROR] SelectSuccessStateReport err(%s)\n", err.Error())
				return c.JSON(http.StatusOK, nil)
			}

			for _, comp := range failComps {
				comps = append(comps, comp)
			}
		}

		for _, comp := range comps {
			params["restId"] = comp["rest_id"]
			if code == CASH_901 {
				StateFailReport(params)
			} else {
				StateSuccessReport(params)
			}
		}

	} else {
		switch code {
		case CASH_901:
			rst = StateFailReport(params)
		case CASH_902:
			rst = StateSuccessReport(params)
		default:
			rst = -1
		}
	}

	if rst < 0 {
		m["resultCode"] = "98"
		m["resultMsg"] = "kakao report fail"
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
	}

	return c.JSON(http.StatusOK, m)
}

func YesterdayReport(c echo.Context) error {

	params := cls.GetParamJsonMap(c)
	m := make(map[string]interface{})
	var rst int
	code := params["code"]

	lprintf(4, "[INFO] code(%s)\n", code)

	switch code {
	case CASH_013, CASH_019:
		rst = YesterdaySuccessReport(params)
	case CASH_014:
		rst = YesterdayFailReport(params)
	case CASH_015:
		rst = YesterdayFail2Report(params)
	case CASH_016:
		rst = YesterdayFailHolidayReport(params)
	case CASH_005: // 여신인증
		rst = CareSalesAccess(params)
	default:
		rst = -1
	}

	if rst < 0 {
		m["resultCode"] = "98"
		m["resultMsg"] = "kakao report faigl"
	} else {
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
	}

	return c.JSON(http.StatusOK, m)
}

func StateSuccessReport(params map[string]string) int {

	header := make(map[string]string)

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	var queryTmp string

	if len(params["bizNum"]) > 0 && params["bizNum"] != "all" {
		queryTmp = datasql.SelectRegistDate
	} else {
		queryTmp = datasql.SelectRegistDateRestId
	}

	regDt, err := cls.SelectData(queryTmp, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		return -1
	}
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	header["NAME"] = userInfo[0]["name"]
	header["JOIN_DATE"] = userInfo[0]["jdate"]
	header["LOGIN_ID"] = userInfo[0]["id"]
	header["START_DATE"] = "8월 23일"

	// body content
	contents := AlimTemplateCode(CASH_902, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_902
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_902

	// button
	sendJButton.ButtonName = "앱으로 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "웹으로 자세히 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "http://www.darayocash.com"
	sendJButton.ButtonURL2 = "http://www.darayocash.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_902
	params["templateNm"] = "가입상태 안내_인증"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	if err != nil {
		lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

func StateFailReport(params map[string]string) int {

	header := make(map[string]string)

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	/*
		var queryTmp string

		if len(params["bizNum"]) > 0 && params["bizNum"] != "all"{
			queryTmp = datasql.SelectRegistDate
		}else{
			queryTmp = datasql.SelectRegistDateRestId
		}

		regDt, err := cls.SelectData(queryTmp, params)
		if err != nil {
			lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
			return -1
		}
		params["restId"] = regDt[0]["restId"]
	*/

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]

	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	header["NAME"] = userInfo[0]["name"]
	header["JOIN_DATE"] = userInfo[0]["jdate"]
	header["LOGIN_ID"] = userInfo[0]["id"]
	header["START_DATE"] = "8월 23일"

	// body content
	contents := AlimTemplateCode(CASH_901, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_901
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_901

	// button
	sendJButton.ButtonName = "앱으로 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "웹으로 자세히 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "http://www.darayocash.com"
	sendJButton.ButtonURL2 = "http://www.darayocash.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_901
	params["templateNm"] = "가입상태 안내_미인증"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	if err != nil {
		lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

func CareSalesAccess(params map[string]string) int {

	header := make(map[string]string)
	today := time.Now()
	yDt := today.AddDate(0, 0, -1).Format("20060102")
	params["bsDt"] = yDt
	params["startDt"] = yDt
	params["endDt"] = yDt

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	regDt, err := cls.SelectData(datasql.SelectRegistDate, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		return -1
	}
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	// 어제 달아요 사용액
	yesterDarayoAmt, err := cls.SelectData(datasql.SelectYesterdayDarayoAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		return -1
	}

	// 오늘 예측
	params["startDt"] = today.AddDate(0, 0, -29).Format("20060102") // 어제 기준 4주전
	params["endDt"] = yDt
	params["trDt"] = today.Format("20060102")

	todaySales, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	// 내일 예측
	params["startDt"] = today.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = today.Format("20060102")
	params["trDt"] = today.AddDate(0, 0, 1).Format("20060102")

	todaySales2, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime2, err := cls.SelectData(homesql.SelectExpectBusyTimeTomorrow, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTimeTomorrow err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData2, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData2, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	header["YDAY_DARAYO_AMT"] = addComma(yesterDarayoAmt[0]["tot_amt"])

	if len(todaySales) > 0 {
		header["TDAY_SALES"] = addComma(todaySales[0]["expectAmt"])
	} else {
		header["TDAY_SALES"] = "0원"
	}

	if len(todaySales2) > 0 {
		header["TOMO_SALES"] = addComma(todaySales2[0]["expectAmt"])
	} else {
		header["TOMO_SALES"] = "0원"
	}

	if len(expectBusyTime) > 0 {
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
		header["TDAY_BUSY"] = busyTime
	} else {
		header["TDAY_BUSY"] = ""
	}

	if len(expectBusyTime2) > 0 {
		var maxTime string
		var maxValue int
		for key, value := range expectBusyTime2[0] {
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
		header["TOMO_BUSY"] = busyTime
	} else {
		header["TOMO_BUSY"] = ""
	}

	var arpuAmt, cashArpuAmt int
	var arpuAmt2, cashArpuAmt2 int
	if len(arpuData) > 0 {
		arpuAmt, _ = strconv.Atoi(arpuData[0]["arpu"])
	}
	if len(cashArpuData) > 0 {
		cashArpuAmt, _ = strconv.Atoi(cashArpuData[0]["arpu"])
	}

	header["TDAY_PAY_SET"] = addComma2(arpuAmt + cashArpuAmt)

	if len(arpuData2) > 0 {
		arpuAmt2, _ = strconv.Atoi(arpuData2[0]["arpu"])
	}
	if len(cashArpuData2) > 0 {
		cashArpuAmt2, _ = strconv.Atoi(cashArpuData2[0]["arpu"])
	}

	header["TOMO_PAY_SET"] = addComma2(arpuAmt2 + cashArpuAmt2)

	// body content
	contents := AlimTemplateCode(CASH_005, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_005
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_005

	// button
	sendJButton.ButtonName = "여신협회 인증하기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/RKxq"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/RKxq"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	lprintf(1, "%s", string(respBody))

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_005
	params["templateNm"] = "여신인증"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	if err != nil {
		lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

// 어제 분석 성공
func YesterdaySuccessReport(params map[string]string) int {

	header := make(map[string]string)
	today := time.Now()
	yDt := today.AddDate(0, 0, -1).Format("20060102")
	params["bsDt"] = yDt
	params["startDt"] = yDt
	params["endDt"] = yDt

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	regDt, err := cls.SelectData(datasql.SelectRegistDate, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		return -1
	}
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	yesterdayAmt, err := cls.SelectData(datasql.SelectYesterdayAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	// 어제 입금 정산 결과
	yesterPayAmt, err := cls.SelectData(datasql.SelectYesterdayPayAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthPayAmt err(%s)\n", err.Error())
		return -1
	}

	// 어제 달아요 사용액
	yesterDarayoAmt, err := cls.SelectData(datasql.SelectYesterdayDarayoAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		return -1
	}

	// 어제 취소 분석
	cancleList, err := cls.SelectData(datasql.SelectYesterdayCancleList, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectLastCancleList err(%s)\n", err.Error())
		return -1
	}

	// 배달업체 정보
	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectDeliveryInfo err(%s)\n", err.Error())
		return -1
	}

	if len(deliveryInfo[0]["baemin_id"]) > 0 {
		params["baeminId"] = deliveryInfo[0]["baemin_id"]
	} else {
		params["baeminId"] = "baeminId"
	}

	if len(deliveryInfo[0]["yogiyo_id"]) > 0 {
		params["yogiyoId"] = deliveryInfo[0]["yogiyo_id"]
	} else {
		params["yogiyoId"] = "yogiyoId"
	}

	if len(deliveryInfo[0]["naver_id"]) > 0 {
		params["naverId"] = deliveryInfo[0]["naver_id"]
	} else {
		params["naverId"] = "naverId"
	}

	if len(deliveryInfo[0]["coupang_id"]) > 0 {
		params["coupangId"] = deliveryInfo[0]["coupang_id"]
	} else {
		params["coupangId"] = "coupangId"
	}

	// 어제 리뷰 분석
	reviews, err := cls.GetSelectData2(datasql.SelectYesterdayReviews, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthReviews err(%s)\n", err.Error())
		return -1
	}

	// 카드사 별 입금 예정
	params["bsDt"] = today.Format("20060102")
	payDailyList, err := cls.SelectData(salesql.SelectPayDailyListHome, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectPayDailyListHome err(%s)\n", err.Error())
		return -1
	}

	// 오늘 예측
	params["startDt"] = today.AddDate(0, 0, -29).Format("20060102") // 어제 기준 4주전
	params["endDt"] = yDt
	params["trDt"] = today.Format("20060102")

	todaySales, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	// 내일 예측
	params["startDt"] = today.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = today.Format("20060102")
	params["trDt"] = today.AddDate(0, 0, 1).Format("20060102")

	todaySales2, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime2, err := cls.SelectData(homesql.SelectExpectBusyTimeTomorrow, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTimeTomorrow err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData2, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData2, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	if len(yesterdayAmt) > 0 {
		header["YDAY_SALES"] = addComma(yesterdayAmt[0]["amt"])
	} else {
		header["YDAY_SALES"] = "0원"
	}

	header["YDAY_DEPOS"] = addComma(yesterPayAmt[0]["realInAmt"])
	header["YDAY_MISS"] = addComma(yesterPayAmt[0]["diffAmt"])

	lprintf(4, "[YDAY_MISS] = %s\n", header["YDAY_MISS"])

	header["YDAY_DARAYO_AMT"] = addComma(yesterDarayoAmt[0]["tot_amt"])

	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int
	if len(cancleList) > 0 {
		for _, v := range cancleList {
			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.SelectData(datasql.SelectLastCancleAprv, params)
			if err != nil {
				noCancle++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0 {
				noCancle++ // 미 승인 취소
				continue
			}

			tr, _ := strconv.Atoi(v["tr_tm"][:2])
			otr, _ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10 {
				nightCancle++ // 심야 취소
			} else if v["tr_dt"] != cancleAprv[0]["tr_dt"] {
				dayCancle++ // 일 취소
			} else if tr-otr > 3 {
				timeCancle++ // 시간 취소
			} else {
				okCancle++ // 결제 취소
			}
		}

		header["YDAY_CANC"] = fmt.Sprintf("%d건", okCancle)
		header["YDAY_CANC_ALRT"] = fmt.Sprintf("%d건", timeCancle+dayCancle+nightCancle+noCancle)
	} else {
		header["YDAY_CANC"] = fmt.Sprintf("0건")
		header["YDAY_CANC_ALRT"] = fmt.Sprintf("0건")
	}

	// 가맹점 rating, keyword
	reviewOption, err := cls.GetSelectData2(datasql.SelectReivewOption, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectReivewOption err(%s)\n", err.Error())
		return -1
	}

	var reviewKeywords []string
	//reviewKeywords = append(reviewKeywords, "등록된 키워드가 없습니다")

	var reviewPoint []float64
	//reviewPoint = 1

	if len(reviewOption) > 0 {
		ratings := strings.Split(reviewOption[0]["rating"], ",")
		keywords := strings.Split(reviewOption[0]["keyword"], "|")

		for _, v := range ratings {
			rating, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			reviewPoint = append(reviewPoint, rating)
		}

		for _, v := range keywords {
			reviewKeywords = append(reviewKeywords, v)
		}

	}

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0 {

	Loop1:
		for _, v := range reviews {
			for _, key := range reviewKeywords { // 키워드 포함 리뷰
				if strings.Contains(v["content"], key) {
					keywordReview++
					rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
					if rst == 1 {
						newCustomerTotal += tot
						newCustomer++
					} else if rst == 2 {
						oldCustomerTotal += tot
						oldCustomer++
					}
					continue Loop1
				}
			}

			for _, key := range reviewPoint { // 평점 포함 리뷰
				r, _ := strconv.ParseFloat(v["rating"], 64)
				if r == key {
					lowReview++
					rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
					if rst == 1 {
						newCustomerTotal += tot
						newCustomer++
					} else if rst == 2 {
						oldCustomerTotal += tot
						oldCustomer++
					}
					continue Loop1
				}
			}

			okReview++
		}

		header["YDAY_VIEW"] = fmt.Sprintf("%d건", okReview+lowReview+keywordReview)
		header["YDAY_VIEW_STAR"] = fmt.Sprintf("%d건", lowReview)
		header["YDAY_VIEW_ALRT"] = fmt.Sprintf("%d건", keywordReview)
	} else {
		header["YDAY_VIEW"] = fmt.Sprintf("0건")
		header["YDAY_VIEW_STAR"] = fmt.Sprintf("0건")
		header["YDAY_VIEW_ALRT"] = fmt.Sprintf("0건")
	}

	var keyword string
	for _, key := range reviewKeywords {
		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0 {
		header["YDAY_VIEW_ALRT_DETAIL"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	} else {
		header["YDAY_VIEW_ALRT_DETAIL"] = "등록된 키워드가 없습니다"
	}

	/*
		var reviewKeywords []string
		reviewKeywords = append(reviewKeywords, "맛있어요")

		var reviewPoint float64
		reviewPoint = 1

		var okReview, lowReview, keywordReview int

		if len(reviews) > 0 {

		Loop1:
			for _, v := range reviews {
				for _, key := range reviewKeywords { // 키워드 포함 리뷰
					if strings.Contains(v["content"], key) {
						keywordReview++
						continue Loop1
					}
				}

				r, _ := strconv.ParseFloat(v["rating"], 64)
				if r <= reviewPoint { // 평점 1점 이하 리뷰
					lowReview++
					continue
				}

				okReview++
			}

			header["YDAY_VIEW"] = fmt.Sprintf("%d건", okReview)
			header["YDAY_VIEW_ALRT"] = fmt.Sprintf("%d건", lowReview+keywordReview)

		} else {
			header["YDAY_VIEW"] = fmt.Sprintf("0건")
			header["YDAY_VIEW_ALRT"] = fmt.Sprintf("0건")
		}
	*/

	header["TDAY_DEPOS_KM"] = "0원"
	header["TDAY_DEPOS_SH"] = "0원"
	header["TDAY_DEPOS_BC"] = "0원"

	header["TDAY_DEPOS_LT"] = "0원"

	header["TDAY_DEPOS_HD"] = "0원"
	header["TDAY_DEPOS_SS"] = "0원"
	header["TDAY_DEPOS_NH"] = "0원"

	header["TDAY_DEPOS_HN"] = "0원"

	var tmp, outpExptAmt int
	for _, payData := range payDailyList {
		switch payData["cardNm"] {
		case "KB카드":
			header["TDAY_DEPOS_KM"] = addComma(payData["outpExptAmt"])
		case "신한카드":
			header["TDAY_DEPOS_SH"] = addComma(payData["outpExptAmt"])
		case "비씨카드":
			header["TDAY_DEPOS_BC"] = addComma(payData["outpExptAmt"])
		case "롯데카드":
			header["TDAY_DEPOS_LT"] = addComma(payData["outpExptAmt"])
		case "현대카드":
			header["TDAY_DEPOS_HD"] = addComma(payData["outpExptAmt"])
		case "삼성카드":
			header["TDAY_DEPOS_SS"] = addComma(payData["outpExptAmt"])
		case "농협NH카드":
			header["TDAY_DEPOS_NH"] = addComma(payData["outpExptAmt"])
		case "하나카드":
			header["TDAY_DEPOS_HN"] = addComma(payData["outpExptAmt"])
		}

		tmp, _ = strconv.Atoi(payData["outpExptAmt"])
		outpExptAmt = outpExptAmt + tmp
	}

	header["TDAY_DEPOS"] = addComma2(outpExptAmt)

	if len(todaySales) > 0 {
		header["TDAY_SALES"] = addComma(todaySales[0]["expectAmt"])
	} else {
		header["TDAY_SALES"] = "0원"
	}

	if len(todaySales2) > 0 {
		header["TOMO_SALES"] = addComma(todaySales2[0]["expectAmt"])
	} else {
		header["TOMO_SALES"] = "0원"
	}

	if len(expectBusyTime) > 0 {
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
		header["TDAY_BUSY"] = busyTime
	} else {
		header["TDAY_BUSY"] = ""
	}

	if len(expectBusyTime2) > 0 {
		var maxTime string
		var maxValue int
		for key, value := range expectBusyTime2[0] {
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
		header["TOMO_BUSY"] = busyTime
	} else {
		header["TOMO_BUSY"] = ""
	}

	var arpuAmt, cashArpuAmt int
	var arpuAmt2, cashArpuAmt2 int
	if len(arpuData) > 0 {
		arpuAmt, _ = strconv.Atoi(arpuData[0]["arpu"])
	}
	if len(cashArpuData) > 0 {
		cashArpuAmt, _ = strconv.Atoi(cashArpuData[0]["arpu"])
	}

	header["TDAY_PAY_SET"] = addComma2(arpuAmt + cashArpuAmt)

	if len(arpuData2) > 0 {
		arpuAmt2, _ = strconv.Atoi(arpuData2[0]["arpu"])
	}
	if len(cashArpuData2) > 0 {
		cashArpuAmt2, _ = strconv.Atoi(cashArpuData2[0]["arpu"])
	}

	header["TOMO_PAY_SET"] = addComma2(arpuAmt2 + cashArpuAmt2)

	// body content
	contents := AlimTemplateCode(CASH_019, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_019
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_019

	// button
	sendJButton.ButtonName = "앱에서 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "어제 리뷰 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = fmt.Sprintf("https://cashapi.darayo.com:7788/review/reviewList?restId=%s&startDt=%s&endDt=%s", regDt[0]["restId"], today.AddDate(0, 0, -1).Format("2006-01-02"), today.Format("2006-01-02"))
	sendJButton.ButtonURL2 = fmt.Sprintf("https://cashapi.darayo.com:7788/review/reviewList?restId=%s&startDt=%s&endDt=%s", regDt[0]["restId"], today.AddDate(0, 0, -1).Format("2006-01-02"), today.Format("2006-01-02"))
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	lprintf(4, "%s", string(respBody))

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_019
	params["templateNm"] = "어제분석_성공4"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	if err != nil {
		lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

func YesterdayFailReport(params map[string]string) int {

	header := make(map[string]string)
	today := time.Now()
	yDt := today.AddDate(0, 0, -1).Format("20060102")
	params["bsDt"] = yDt
	params["startDt"] = yDt
	params["endDt"] = yDt

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	regDt, err := cls.SelectData(datasql.SelectRegistDate, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		return -1
	}
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	// 어제 달아요 사용액
	yesterDarayoAmt, err := cls.SelectData(datasql.SelectYesterdayDarayoAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		return -1
	}

	// 카드사 별 입금 예정
	params["bsDt"] = today.Format("20060102")
	payDailyList, err := cls.SelectData(salesql.SelectPayDailyListHome, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectPayDailyListHome err(%s)\n", err.Error())
		return -1
	}

	// 오늘 예측
	params["startDt"] = today.AddDate(0, 0, -29).Format("20060102") // 어제 기준 4주전
	params["endDt"] = yDt
	params["trDt"] = today.Format("20060102")

	todaySales, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	// 내일 예측
	params["startDt"] = today.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = today.Format("20060102")
	params["trDt"] = today.AddDate(0, 0, 1).Format("20060102")

	todaySales2, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime2, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData2, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData2, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]
	header["SEND_TIME"] = "오후 2시"

	header["YDAY_DARAYO_AMT"] = addComma(yesterDarayoAmt[0]["tot_amt"])

	header["TDAY_DEPOS_KM"] = "0원"
	header["TDAY_DEPOS_SH"] = "0원"
	header["TDAY_DEPOS_BC"] = "0원"

	header["TDAY_DEPOS_LT"] = "0원"

	header["TDAY_DEPOS_HD"] = "0원"
	header["TDAY_DEPOS_SS"] = "0원"
	header["TDAY_DEPOS_NH"] = "0원"

	header["TDAY_DEPOS_HN"] = "0원"

	var tmp, outpExptAmt int
	for _, payData := range payDailyList {
		switch payData["cardNm"] {
		case "KB카드":
			header["TDAY_DEPOS_KM"] = addComma(payData["outpExptAmt"])
		case "신한카드":
			header["TDAY_DEPOS_SH"] = addComma(payData["outpExptAmt"])
		case "비씨카드":
			header["TDAY_DEPOS_BC"] = addComma(payData["outpExptAmt"])
		case "롯데카드":
			header["TDAY_DEPOS_LT"] = addComma(payData["outpExptAmt"])
		case "현대카드":
			header["TDAY_DEPOS_HD"] = addComma(payData["outpExptAmt"])
		case "삼성카드":
			header["TDAY_DEPOS_SS"] = addComma(payData["outpExptAmt"])
		case "농협NH카드":
			header["TDAY_DEPOS_NH"] = addComma(payData["outpExptAmt"])
		case "하나카드":
			header["TDAY_DEPOS_HN"] = addComma(payData["outpExptAmt"])
		}

		tmp, _ = strconv.Atoi(payData["outpExptAmt"])
		outpExptAmt = outpExptAmt + tmp
	}

	header["TDAY_DEPOS"] = addComma2(outpExptAmt)

	if len(todaySales) > 0 {
		header["TDAY_SALES"] = addComma(todaySales[0]["expectAmt"])
	} else {
		header["TDAY_SALES"] = "0원"
	}

	if len(todaySales2) > 0 {
		header["TOMO_SALES"] = addComma(todaySales2[0]["expectAmt"])
	} else {
		header["TOMO_SALES"] = "0원"
	}

	if len(expectBusyTime) > 0 {
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
		busyTime := fmt.Sprintf("%s:00~%s:00시", bytes[1:3], bytes[3:])
		header["TDAY_BUSY"] = busyTime
	} else {
		header["TDAY_BUSY"] = ""
	}

	if len(expectBusyTime2) > 0 {
		var maxTime string
		var maxValue int
		for key, value := range expectBusyTime2[0] {
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
		busyTime := fmt.Sprintf("%s:00~%s:00시", bytes[1:3], bytes[3:])
		header["TOMO_BUSY"] = busyTime
	} else {
		header["TOMO_BUSY"] = ""
	}

	var arpuAmt, cashArpuAmt int
	var arpuAmt2, cashArpuAmt2 int
	if len(arpuData) > 0 {
		arpuAmt, _ = strconv.Atoi(arpuData[0]["arpu"])
	}
	if len(cashArpuData) > 0 {
		cashArpuAmt, _ = strconv.Atoi(cashArpuData[0]["arpu"])
	}

	header["TDAY_PAY_SET"] = addComma2(arpuAmt + cashArpuAmt)

	if len(arpuData2) > 0 {
		arpuAmt2, _ = strconv.Atoi(arpuData2[0]["arpu"])
	}
	if len(cashArpuData2) > 0 {
		cashArpuAmt2, _ = strconv.Atoi(cashArpuData2[0]["arpu"])
	}

	header["TOMO_PAY_SET"] = addComma2(arpuAmt2 + cashArpuAmt2)

	// body content
	contents := AlimTemplateCode(CASH_014, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_014
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_014

	// button
	sendJButton.ButtonName = "앱으로 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "웹으로 자세히 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "https://partner.darayo.com"
	sendJButton.ButtonURL2 = "https://partner.darayo.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_014
	params["templateNm"] = "어제분석_1차실패"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

func YesterdayFail2Report(params map[string]string) int {

	header := make(map[string]string)
	today := time.Now()
	yDt := today.AddDate(0, 0, -1).Format("20060102")
	params["bsDt"] = yDt
	params["startDt"] = yDt
	params["endDt"] = yDt

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	regDt, err := cls.SelectData(datasql.SelectRegistDate, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		return -1
	}
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	// 어제 달아요 사용액
	yesterDarayoAmt, err := cls.SelectData(datasql.SelectYesterdayDarayoAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		return -1
	}

	// 카드사 별 입금 예정
	params["bsDt"] = today.Format("20060102")
	payDailyList, err := cls.SelectData(salesql.SelectPayDailyListHome, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectPayDailyListHome err(%s)\n", err.Error())
		return -1
	}

	// 오늘 예측
	params["startDt"] = today.AddDate(0, 0, -29).Format("20060102") // 어제 기준 4주전
	params["endDt"] = yDt
	params["trDt"] = today.Format("20060102")

	todaySales, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	// 내일 예측
	params["startDt"] = today.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = today.Format("20060102")
	params["trDt"] = today.AddDate(0, 0, 1).Format("20060102")

	todaySales2, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime2, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData2, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData2, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	header["YDAY_DARAYO_AMT"] = addComma(yesterDarayoAmt[0]["tot_amt"])

	header["TDAY_DEPOS_KM"] = "0원"
	header["TDAY_DEPOS_SH"] = "0원"
	header["TDAY_DEPOS_BC"] = "0원"

	header["TDAY_DEPOS_LT"] = "0원"

	header["TDAY_DEPOS_HD"] = "0원"
	header["TDAY_DEPOS_SS"] = "0원"
	header["TDAY_DEPOS_NH"] = "0원"

	header["TDAY_DEPOS_HN"] = "0원"

	var tmp, outpExptAmt int
	for _, payData := range payDailyList {
		switch payData["cardNm"] {
		case "KB카드":
			header["TDAY_DEPOS_KM"] = addComma(payData["outpExptAmt"])
		case "신한카드":
			header["TDAY_DEPOS_SH"] = addComma(payData["outpExptAmt"])
		case "비씨카드":
			header["TDAY_DEPOS_BC"] = addComma(payData["outpExptAmt"])
		case "롯데카드":
			header["TDAY_DEPOS_LT"] = addComma(payData["outpExptAmt"])
		case "현대카드":
			header["TDAY_DEPOS_HD"] = addComma(payData["outpExptAmt"])
		case "삼성카드":
			header["TDAY_DEPOS_SS"] = addComma(payData["outpExptAmt"])
		case "농협NH카드":
			header["TDAY_DEPOS_NH"] = addComma(payData["outpExptAmt"])
		case "하나카드":
			header["TDAY_DEPOS_HN"] = addComma(payData["outpExptAmt"])
		}

		tmp, _ = strconv.Atoi(payData["outpExptAmt"])
		outpExptAmt = outpExptAmt + tmp
	}

	header["TDAY_DEPOS"] = addComma2(outpExptAmt)

	if len(todaySales) > 0 {
		header["TDAY_SALES"] = addComma(todaySales[0]["expectAmt"])
	} else {
		header["TDAY_SALES"] = "0원"
	}

	if len(todaySales2) > 0 {
		header["TOMO_SALES"] = addComma(todaySales2[0]["expectAmt"])
	} else {
		header["TOMO_SALES"] = "0원"
	}

	if len(expectBusyTime) > 0 {
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
		busyTime := fmt.Sprintf("%s:00~%s:00시", bytes[1:3], bytes[3:])
		header["TDAY_BUSY"] = busyTime
	} else {
		header["TDAY_BUSY"] = ""
	}

	if len(expectBusyTime2) > 0 {
		var maxTime string
		var maxValue int
		for key, value := range expectBusyTime2[0] {
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
		busyTime := fmt.Sprintf("%s:00~%s:00시", bytes[1:3], bytes[3:])
		header["TOMO_BUSY"] = busyTime
	} else {
		header["TOMO_BUSY"] = ""
	}

	var arpuAmt, cashArpuAmt int
	var arpuAmt2, cashArpuAmt2 int
	if len(arpuData) > 0 {
		arpuAmt, _ = strconv.Atoi(arpuData[0]["arpu"])
	}
	if len(cashArpuData) > 0 {
		cashArpuAmt, _ = strconv.Atoi(cashArpuData[0]["arpu"])
	}

	header["TDAY_PAY_SET"] = addComma2(arpuAmt + cashArpuAmt)

	if len(arpuData2) > 0 {
		arpuAmt2, _ = strconv.Atoi(arpuData2[0]["arpu"])
	}
	if len(cashArpuData2) > 0 {
		cashArpuAmt2, _ = strconv.Atoi(cashArpuData2[0]["arpu"])
	}

	header["TOMO_PAY_SET"] = addComma2(arpuAmt2 + cashArpuAmt2)

	// body content
	contents := AlimTemplateCode(CASH_015, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_015
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_015

	// button
	sendJButton.ButtonName = "앱으로 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "웹으로 자세히 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "https://partner.darayo.com"
	sendJButton.ButtonURL2 = "https://partner.darayo.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_015
	params["templateNm"] = "어제분석_2차실패"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

func YesterdayFailHolidayReport(params map[string]string) int {

	header := make(map[string]string)
	today := time.Now()
	yDt := today.AddDate(0, 0, -1).Format("20060102")
	params["bsDt"] = yDt
	params["startDt"] = yDt
	params["endDt"] = yDt

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		return -1
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	regDt, err := cls.SelectData(datasql.SelectRegistDate, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		return -1
	}
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.SelectData(homesql.SelectRestUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		return -1
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.SelectData(homesql.SelectUserInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		return -1
	}

	// 어제 달아요 사용액
	yesterDarayoAmt, err := cls.SelectData(datasql.SelectYesterdayDarayoAmt, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		return -1
	}

	// 오늘 예측
	params["startDt"] = today.AddDate(0, 0, -29).Format("20060102") // 어제 기준 4주전
	params["endDt"] = yDt
	params["trDt"] = today.Format("20060102")

	todaySales, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	// 내일 예측
	params["startDt"] = today.AddDate(0, 0, -28).Format("20060102")
	params["endDt"] = today.Format("20060102")
	params["trDt"] = today.AddDate(0, 0, 1).Format("20060102")

	todaySales2, err := cls.SelectData(homesql.SelectDaySaleData, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 바쁜 시간
	expectBusyTime2, err := cls.SelectData(homesql.SelectExpectBusyTime, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectExpectBusyTime err(%s)\n", err.Error())
		return -1
	}

	// 객단가 예측
	arpuData2, err := cls.SelectData(datasql.SelectAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}
	cashArpuData2, err := cls.SelectData(datasql.SelectCashAverageRevenuePerUser, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectCashAverageRevenuePerUser err(%s)\n", err.Error())
		return -1
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	header["YDAY_DARAYO_AMT"] = addComma(yesterDarayoAmt[0]["tot_amt"])

	if len(todaySales) > 0 {
		header["TDAY_SALES"] = addComma(todaySales[0]["expectAmt"])
	} else {
		header["TDAY_SALES"] = "0원"
	}

	if len(todaySales2) > 0 {
		header["TOMO_SALES"] = addComma(todaySales2[0]["expectAmt"])
	} else {
		header["TOMO_SALES"] = "0원"
	}

	if len(expectBusyTime) > 0 {
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
		busyTime := fmt.Sprintf("%s:00~%s:00시", bytes[1:3], bytes[3:])
		header["TDAY_BUSY"] = busyTime
	} else {
		header["TDAY_BUSY"] = ""
	}

	if len(expectBusyTime2) > 0 {
		var maxTime string
		var maxValue int
		for key, value := range expectBusyTime2[0] {
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
		busyTime := fmt.Sprintf("%s:00~%s:00시", bytes[1:3], bytes[3:])
		header["TOMO_BUSY"] = busyTime
	} else {
		header["TOMO_BUSY"] = ""
	}

	var arpuAmt, cashArpuAmt int
	var arpuAmt2, cashArpuAmt2 int
	if len(arpuData) > 0 {
		arpuAmt, _ = strconv.Atoi(arpuData[0]["arpu"])
	}
	if len(cashArpuData) > 0 {
		cashArpuAmt, _ = strconv.Atoi(cashArpuData[0]["arpu"])
	}

	header["TDAY_PAY_SET"] = addComma2(arpuAmt + cashArpuAmt)

	if len(arpuData2) > 0 {
		arpuAmt2, _ = strconv.Atoi(arpuData2[0]["arpu"])
	}
	if len(cashArpuData2) > 0 {
		cashArpuAmt2, _ = strconv.Atoi(cashArpuData2[0]["arpu"])
	}

	header["TOMO_PAY_SET"] = addComma2(arpuAmt2 + cashArpuAmt2)

	// body content
	contents := AlimTemplateCode(CASH_016, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_016
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_016

	// button
	sendJButton.ButtonName = "앱으로 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "웹으로 자세히 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "https://partner.darayo.com"
	sendJButton.ButtonURL2 = "https://partner.darayo.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		return -1
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		return -1
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_016
	params["templateNm"] = "어제분석_주말실패"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return -1
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		return -1
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return -1
	}

	return 1
}

func LastWeekReport(c echo.Context) error {

	// bizNum
	params := cls.GetParamJsonMap(c)
	header := make(map[string]string)
	m := make(map[string]interface{})

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	regDt, err := cls.GetSelectData(datasql.SelectRegistDate, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// content data
	dt := time.Now().AddDate(0, 0, -7)
	startWeek := cls.GetFirstOfWeek(dt)
	endWeek := cls.GetEndOfWeek(dt)
	wStartDt := startWeek.Format("20060102")
	wEndDt := endWeek.Format("20060102")

	params["startDt"] = wStartDt
	params["endDt"] = wEndDt
	params["restId"] = regDt[0]["restId"]

	// get user id
	rUserInfo, err := cls.GetSelectData(homesql.SelectRestUserInfo, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.GetSelectData(homesql.SelectUserInfo, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 주간 매출 분석 (요일별 날자 입력, 지난 7일)
	weekSalesData, err := cls.GetSelectData(datasql.SelectWeekCash, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectWeekCash err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난주 카드
	lastWeekCard, err := cls.GetSelectDataUsingJson(datasql.SelectWeekCard, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthCard err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난주 달아요 사용액
	lastWeekDarayoAmt, err := cls.GetSelectDataUsingJson(datasql.SelectWeekDarayoAmt, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난주 입금 정산 결과
	lastWeekPayAmt, err := cls.GetSelectDataUsingJson(datasql.SelectWeekPayAmt, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthPayAmt err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난주 취소 분석
	cancleList, err := cls.GetSelectData(datasql.SelectLastCancleList, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectLastCancleList err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난주 고객님 결제 건수 분석 디테일
	selectWeekCntDetail, err := cls.GetSelectDataUsingJson(datasql.SelectWeekCntDetail, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectWeekCntDetail err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 배달업체 정보
	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectDeliveryInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	if len(deliveryInfo[0]["baemin_id"]) > 0 {
		params["baeminId"] = deliveryInfo[0]["baemin_id"]
	} else {
		params["baeminId"] = "baeminId"
	}

	if len(deliveryInfo[0]["yogiyo_id"]) > 0 {
		params["yogiyoId"] = deliveryInfo[0]["yogiyo_id"]
	} else {
		params["yogiyoId"] = "yogiyoId"
	}

	if len(deliveryInfo[0]["naver_id"]) > 0 {
		params["naverId"] = deliveryInfo[0]["naver_id"]
	} else {
		params["naverId"] = "naverId"
	}

	if len(deliveryInfo[0]["coupang_id"]) > 0 {
		params["coupangId"] = deliveryInfo[0]["coupang_id"]
	} else {
		params["coupangId"] = "coupangId"
	}

	// 지난주 리뷰 분석
	reviews, err := cls.GetSelectData2(datasql.SelectReviews, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthReviews err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	header["WEEK_SALES"] = addComma(weekSalesData[0]["total"])
	header["WEEK_BEST"] = addComma(weekSalesData[0]["maxAmt"])
	header["WEEK_BEST_TIME"] = weekSalesData[0]["maxWeek"]
	header["WEEK_WRST"] = addComma(weekSalesData[0]["minAmt"])
	header["WEEK_WRST_TIME"] = weekSalesData[0]["minWeek"]

	header["WEEK_BCARD"] = lastWeekCard[0]["card_nm"]

	header["WEEK_DARAYO"] = addComma(lastWeekDarayoAmt[0]["tot_amt"])

	header["WEEK_DEPOSIT"] = addComma(lastWeekPayAmt[0]["realInAmt"])
	header["WEEK_MISS"] = addComma(lastWeekPayAmt[0]["diffAmt"])

	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int
	if len(cancleList) > 0 {
		for _, v := range cancleList {
			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0 {
				noCancle++ // 미 승인 취소
				continue
			}

			tr, _ := strconv.Atoi(v["tr_tm"][:2])
			otr, _ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10 {
				nightCancle++ // 심야 취소
			} else if v["tr_dt"] != cancleAprv[0]["tr_dt"] {
				dayCancle++ // 일 취소
			} else if tr-otr > 3 {
				timeCancle++ // 시간 취소
			} else {
				okCancle++ // 결제 취소
			}
		}

		header["WEEK_CANCEL"] = fmt.Sprintf("%d건", okCancle+timeCancle+dayCancle+nightCancle+noCancle)
		header["WEEK_CNIGHT"] = fmt.Sprintf("%d건", nightCancle)
		header["WEEK_CHOUR"] = fmt.Sprintf("%d건", timeCancle)
	} else {
		header["WEEK_CANCEL"] = fmt.Sprintf("0건")
		header["WEEK_CNIGHT"] = fmt.Sprintf("0건")
		header["WEEK_CHOUR"] = fmt.Sprintf("0건")
	}

	// 가맹점 rating, keyword
	reviewOption, err := cls.GetSelectData2(datasql.SelectReivewOption, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectReivewOption err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	var reviewKeywords []string
	//reviewKeywords = append(reviewKeywords, "등록된 키워드가 없습니다")

	var reviewPoint []float64
	//reviewPoint = 1

	if len(reviewOption) > 0 {
		ratings := strings.Split(reviewOption[0]["rating"], ",")
		keywords := strings.Split(reviewOption[0]["keyword"], "|")

		for _, v := range ratings {
			rating, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			reviewPoint = append(reviewPoint, rating)
		}

		for _, v := range keywords {
			reviewKeywords = append(reviewKeywords, v)
		}

	}

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0 {

	Loop1:
		for _, v := range reviews {
			for _, key := range reviewKeywords { // 키워드 포함 리뷰
				if strings.Contains(v["content"], key) {
					keywordReview++
					rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
					if rst == 1 {
						newCustomerTotal += tot
						newCustomer++
					} else if rst == 2 {
						oldCustomerTotal += tot
						oldCustomer++
					}
					continue Loop1
				}
			}

			for _, key := range reviewPoint { // 평점 포함 리뷰
				r, _ := strconv.ParseFloat(v["rating"], 64)
				if r == key {
					lowReview++
					rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
					if rst == 1 {
						newCustomerTotal += tot
						newCustomer++
					} else if rst == 2 {
						oldCustomerTotal += tot
						oldCustomer++
					}
					continue Loop1
				}
			}

			/*
				r, _ := strconv.ParseFloat(v["rating"], 64)
				if r <= reviewPoint { // 평점 1점 이하 리뷰
					lowReview++
					rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
					if rst == 1 {
						newCustomerTotal += tot
						newCustomer++
					} else if rst == 2 {
						oldCustomerTotal += tot
						oldCustomer++
					}
					continue
				}

			*/

			okReview++
		}

		header["WEEK_REVIEW"] = fmt.Sprintf("%d건", okReview+lowReview+keywordReview)
		header["WEEK_RLOW"] = fmt.Sprintf("%d건", lowReview)
		header["WEEK_RKEY"] = fmt.Sprintf("%d건", keywordReview)

		header["NEW_CUST_LIST"] = fmt.Sprintf("%d명", newCustomer)
		header["OLD_CUST_LIST"] = fmt.Sprintf("%d명", oldCustomer)

		if newCustomer > 0 {
			header["NEW_CUST_POINT"] = fmt.Sprintf("%s점", strconv.FormatFloat(newCustomerTotal/float64(newCustomer), 'f', 2, 64))
		} else {
			header["NEW_CUST_POINT"] = fmt.Sprintf("0점")
		}

		if oldCustomer > 0 {
			header["OLD_CUST_POINT"] = fmt.Sprintf("%s점", strconv.FormatFloat(oldCustomerTotal/float64(oldCustomer), 'f', 2, 64))
		} else {
			header["OLD_CUST_POINT"] = fmt.Sprintf("0점")
		}
	} else {
		header["WEEK_REVIEW"] = fmt.Sprintf("0건")
		header["WEEK_RLOW"] = fmt.Sprintf("0건")
		header["WEEK_RKEY"] = fmt.Sprintf("0건")

		header["NEW_CUST_LIST"] = fmt.Sprintf("0명")
		header["OLD_CUST_LIST"] = fmt.Sprintf("0명")

		header["NEW_CUST_POINT"] = fmt.Sprintf("0점")
		header["OLD_CUST_POINT"] = fmt.Sprintf("0점")
	}

	var keyword string
	for _, key := range reviewKeywords {
		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0 {
		header["REVIEW_KEYWORD"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	} else {
		header["REVIEW_KEYWORD"] = "등록된 키워드가 없습니다"
	}

	var weekCntSorts []structCntSort
	weekCntMsg := "%s %s시"

	for i := 0; i < len(selectWeekCntDetail); i++ {
		t0003, _ := strconv.Atoi(selectWeekCntDetail[i]["t0003"])
		t0306, _ := strconv.Atoi(selectWeekCntDetail[i]["t0306"])
		t0609, _ := strconv.Atoi(selectWeekCntDetail[i]["t0609"])
		t0912, _ := strconv.Atoi(selectWeekCntDetail[i]["t0912"])
		t1215, _ := strconv.Atoi(selectWeekCntDetail[i]["t1215"])
		t1518, _ := strconv.Atoi(selectWeekCntDetail[i]["t1518"])
		t1821, _ := strconv.Atoi(selectWeekCntDetail[i]["t1821"])
		t2124, _ := strconv.Atoi(selectWeekCntDetail[i]["t2124"])

		var weekCntSort structCntSort
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "0~3")
		weekCntSort.val = t0003
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "3~6")
		weekCntSort.val = t0306
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "6~9")
		weekCntSort.val = t0609
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "9~12")
		weekCntSort.val = t0912
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "12~15")
		weekCntSort.val = t1215
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "15~18")
		weekCntSort.val = t1518
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "18~21")
		weekCntSort.val = t1821
		weekCntSorts = append(weekCntSorts, weekCntSort)
		weekCntSort.key = fmt.Sprintf(weekCntMsg, selectWeekCntDetail[i]["day_name"], "21~24")
		weekCntSort.val = t2124
		weekCntSorts = append(weekCntSorts, weekCntSort)

		//tTime := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}
		//tm,_ := dataCon.FindBusyTime2(tTime)
		//header[fmt.Sprintf("WEEK_BUSY%d",i+1)] = fmt.Sprintf("%s %s시", selectWeekCntDetail[i]["day_name"], tm)
	}

	// 동일 요일 포함해서 top3 색출
	sort.SliceStable(weekCntSorts, func(i, j int) bool {
		return weekCntSorts[i].val > weekCntSorts[j].val
	})

	for i, v := range weekCntSorts {
		//fmt.Println(fmt.Sprintf("key : %s, val : %d", v.key, v.val))
		header[fmt.Sprintf("WEEK_BUSY%d", i+1)] = v.key
		if i == 2 {
			break
		}
	}

	for i := 1; i < 3; i++ {
		val, exists := header[fmt.Sprintf("WEEK_BUSY%d", i)]
		if !exists {
			header[fmt.Sprintf("WEEK_BUSY%d", i)] = ""
		} else if strings.Contains(val, "0~") {
			header[fmt.Sprintf("WEEK_BUSY%d", i)] = ""
		}
	}

	/*
		for i:=0; i<2; i++{
			if weekCntSorts[i].key[len(weekCntSorts[i].key)-1] == '0'{
				header[fmt.Sprintf("WEEK_BUSY%d",i+1)] = weekCntSorts[i].key
			} else{
				header[fmt.Sprintf("WEEK_BUSY%d",i+1)] = ""
			}
		}
	*/

	// body content
	contents := AlimTemplateCode(CASH_104, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_104
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_104

	// button
	sendJButton.ButtonName = "앱에서 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "지난주 리뷰 보기"
	sendJButton.ButtonType = "WL"

	sendJButton.ButtonURL = fmt.Sprintf("https://cashapi.darayo.com:7788/review/reviewList?restId=%s&startDt=%s&endDt=%s", regDt[0]["restId"], startWeek.Format("2006-01-02"), endWeek.Format("2006-01-02"))
	sendJButton.ButtonURL2 = fmt.Sprintf("https://cashapi.darayo.com:7788/review/reviewList?restId=%s&startDt=%s&endDt=%s", regDt[0]["restId"], startWeek.Format("2006-01-02"), endWeek.Format("2006-01-02"))
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		lprintf(1, "[ERROR] json Marshal err(%s)\n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		lprintf(1, "[ERROR] HttpsJson err(%s)\n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	lprintf(4, "[INFO] %v", result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_104
	params["templateNm"] = "주간분석_성공4"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func LastMonthReport3(c echo.Context) error {

	// bizNum, userId

	params := cls.GetParamJsonMap(c)
	header := make(map[string]string)
	m := make(map[string]interface{})

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	// content data
	// 가입일 조회
	regDt, err := cls.GetSelectData(datasql.SelectRegistDate, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	now := time.Now()
	thisYear := now.AddDate(0, 0, 0).Format("2006")
	// 지난달
	params["endDt"] = now.AddDate(0, 0, -now.Day()).Format("200601")
	params["bsDt"] = params["endDt"]
	params["restId"] = regDt[0]["restId"]

	// 월간분석 - 월 매출 분석(올해)
	if thisYear == regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 수집월 (지난달 부터 시작)
		mon := regDt[0]["regDt"][4:6]
		monInt, _ := strconv.Atoi(mon)
		monInt -= 1
		if monInt < 10 {
			params["startDt"] = fmt.Sprintf("%s0%d", regDt[0]["regDt"][:4], monInt)
		} else {
			params["startDt"] = fmt.Sprintf("%s%d", regDt[0]["regDt"][:4], monInt)
		}
	} else {
		// 작년 가입한 경우 1월, 혹은 프리미엄
		params["startDt"] = fmt.Sprintf("%s01", thisYear)
	}

	// get user id
	rUserInfo, err := cls.GetSelectData(homesql.SelectRestUserInfo, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.GetSelectData(homesql.SelectUserInfo, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	// body content
	contents := AlimTemplateCode(CASH_205, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "ai"
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_205
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	//sendJMessage.ReSend = "R"
	//sendJMessage.ReText = CASH_203

	// button
	sendJButton.ButtonName = "보고서 보러가기"
	sendJButton.ButtonType = "WL"

	// android
	sendJButton.ButtonURL = fmt.Sprintf("https://cashapi.darayo.com:7788/page/v2/monthly?bizNum=%s", params["bizNum"])
	// ios
	sendJButton.ButtonURL2 = fmt.Sprintf("https://cashapi.darayo.com:7788/page/v2/monthly?bizNum=%s", params["bizNum"])
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		lprintf(1, "[ERROR] json Marshal err(%s)\n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		lprintf(1, "[ERROR] HttpsJson err(%s)\n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	lprintf(4, "[INFO] %s", string(respBody))

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	lprintf(4, "[INFO] %v", result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_205
	params["templateNm"] = "월간분석_성공5"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func LastMonthReport(c echo.Context) error {

	// bizNum, userId

	params := cls.GetParamJsonMap(c)
	header := make(map[string]string)
	m := make(map[string]interface{})

	// kakao message index
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params)
	if err != nil {
		lprintf(1, "[ERROR] SelctAlimTalkSeq err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	// content data
	// 가입일 조회
	regDt, err := cls.GetSelectData(datasql.SelectRegistDate, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectRegistDate err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	now := time.Now()
	thisYear := now.AddDate(0, 0, 0).Format("2006")
	// 지난달
	params["endDt"] = now.AddDate(0, 0, -now.Day()).Format("200601")
	params["bsDt"] = params["endDt"]
	params["restId"] = regDt[0]["restId"]

	// 월간분석 - 월 매출 분석(올해)
	if thisYear == regDt[0]["regDt"][:4] {
		// 올해 가입한 경우 수집월 (지난달 부터 시작)
		mon := regDt[0]["regDt"][4:6]
		monInt, _ := strconv.Atoi(mon)
		monInt -= 1
		if monInt < 10 {
			params["startDt"] = fmt.Sprintf("%s0%d", regDt[0]["regDt"][:4], monInt)
		} else {
			params["startDt"] = fmt.Sprintf("%s%d", regDt[0]["regDt"][:4], monInt)
		}
	} else {
		// 작년 가입한 경우 1월, 혹은 프리미엄
		params["startDt"] = fmt.Sprintf("%s01", thisYear)
	}

	// get user id
	rUserInfo, err := cls.GetSelectData(homesql.SelectRestUserInfo, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectRestUserInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	params["userId"] = rUserInfo[0]["userId"]
	userInfo, err := cls.GetSelectData(homesql.SelectUserInfo, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUserInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 월간분석 - 월 매출 분석(올해)
	thisSalesData, err := cls.GetSelectDataUsingJson(datasql.SelectMonthCash1, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthCash1 err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난달 카드
	lastMonthCard, err := cls.GetSelectDataUsingJson(datasql.SelectMonthCard, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthCard err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난달 달아요 사용액
	lastMonthDarayoAmt, err := cls.GetSelectDataUsingJson(datasql.SelectMonthDarayoAmt, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthDarayoAmt err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난달 입금 정산 결과
	lastMonthPayAmt, err := cls.GetSelectDataUsingJson(datasql.SelectMonthPayAmt, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthPayAmt err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난달 취소 분석
	cancleList, err := cls.GetSelectData(datasql.SelectLastMonthCancleList, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectLastCancleList err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 지난달 고객님 결제 건수 분석 디테일
	selectMonthCntDetail, err := cls.GetSelectDataUsingJson(datasql.SelectMonthCntDetail, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthCntDetail err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 배달업체 정보
	deliveryInfo, err := cls.GetSelectData2(datasql.SelectDeliveryInfo, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectDeliveryInfo err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	if len(deliveryInfo[0]["baemin_id"]) > 0 {
		params["baeminId"] = deliveryInfo[0]["baemin_id"]
	} else {
		params["baeminId"] = "baeminId"
	}

	if len(deliveryInfo[0]["yogiyo_id"]) > 0 {
		params["yogiyoId"] = deliveryInfo[0]["yogiyo_id"]
	} else {
		params["yogiyoId"] = "yogiyoId"
	}

	if len(deliveryInfo[0]["naver_id"]) > 0 {
		params["naverId"] = deliveryInfo[0]["naver_id"]
	} else {
		params["naverId"] = "naverId"
	}

	// 지난달 리뷰 분석
	reviews, err := cls.GetSelectData2(datasql.SelectMonthReviews, params)
	if err != nil {
		lprintf(1, "[ERROR] SelectMonthReviews err(%s)\n", err.Error())
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	header["SHOP_NAME"] = regDt[0]["compNm"]

	header["MONTH_SALES"] = addComma(thisSalesData[0]["beforeMonthAmt"])
	header["MONTH_AVRG"] = addComma(thisSalesData[0]["avgAmt"])
	header["MONTH_BEST"] = addComma(thisSalesData[0]["maxAmt"])

	header["MONTH_BCARD"] = lastMonthCard[0]["card_nm"]

	header["MONTH_DARAYO"] = addComma(lastMonthDarayoAmt[0]["tot_amt"])

	header["MONTH_DEPOSIT"] = addComma(lastMonthPayAmt[0]["realInAmt"])
	header["MONTH_MISS_DEPO"] = addComma(lastMonthPayAmt[0]["diffAmt"])

	var okCancle, timeCancle, dayCancle, nightCancle, noCancle int
	if len(cancleList) > 0 {
		for _, v := range cancleList {
			params["aprvNo"] = v["aprv_no"]

			cancleAprv, err := cls.GetSelectData(datasql.SelectLastCancleAprv, params, c)
			if err != nil {
				noCancle++ // 미 승인 취소
				continue
			}

			if len(cancleAprv) == 0 {
				noCancle++ // 미 승인 취소
				continue
			}

			tr, _ := strconv.Atoi(v["tr_tm"][:2])
			otr, _ := strconv.Atoi(cancleAprv[0]["tr_tm"][:2])

			if tr < 10 {
				nightCancle++ // 심야 취소
			} else if v["tr_dt"] != cancleAprv[0]["tr_dt"] {
				dayCancle++ // 일 취소
			} else if tr-otr > 3 {
				timeCancle++ // 시간 취소
			} else {
				okCancle++ // 결제 취소
			}
		}

		header["MONTH_CANCEL"] = fmt.Sprintf("%d건", okCancle+timeCancle+dayCancle+nightCancle+noCancle)
		header["MONTH_CNIGHT"] = fmt.Sprintf("%d건", nightCancle)
		header["MONTH_CHOUR"] = fmt.Sprintf("%d건", timeCancle)
	} else {
		header["MONTH_CANCEL"] = fmt.Sprintf("0건")
		header["MONTH_CNIGHT"] = fmt.Sprintf("0건")
		header["MONTH_CHOUR"] = fmt.Sprintf("0건")
	}

	var reviewKeywords []string
	reviewKeywords = append(reviewKeywords, "맛있어요")

	var reviewPoint float64
	reviewPoint = 1

	var okReview, lowReview, keywordReview int
	var newCustomer, oldCustomer int
	var newCustomerTotal, oldCustomerTotal float64

	if len(reviews) > 0 {

	Loop1:
		for _, v := range reviews {
			for _, key := range reviewKeywords { // 키워드 포함 리뷰
				if strings.Contains(v["content"], key) {
					keywordReview++
					rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
					if rst == 1 {
						newCustomerTotal += tot
						newCustomer++
					} else if rst == 2 {
						oldCustomerTotal += tot
						oldCustomer++
					}
					continue Loop1
				}
			}

			r, _ := strconv.ParseFloat(v["rating"], 64)
			if r <= reviewPoint { // 평점 1점 이하 리뷰
				lowReview++
				rst, tot := dataCon.CheckReivewer(v["member_no"], v["rating"])
				if rst == 1 {
					newCustomerTotal += tot
					newCustomer++
				} else if rst == 2 {
					oldCustomerTotal += tot
					oldCustomer++
				}
				continue
			}

			okReview++
		}

		header["MONTH_REVIEW"] = fmt.Sprintf("%d건", okReview+lowReview+keywordReview)
		header["MONTH_RLOW"] = fmt.Sprintf("%d건", lowReview)
		header["MONTH_RKEY"] = fmt.Sprintf("%d건", keywordReview)

		header["NEW_CUST_LIST"] = fmt.Sprintf("%d명", newCustomer)
		header["OLD_CUST_LIST"] = fmt.Sprintf("%d명", oldCustomer)

		if newCustomer > 0 {
			header["NEW_CUST_POINT"] = fmt.Sprintf("%s점", strconv.FormatFloat(newCustomerTotal/float64(newCustomer), 'f', 2, 64))
		} else {
			header["NEW_CUST_POINT"] = fmt.Sprintf("0점")
		}

		if oldCustomer > 0 {
			header["OLD_CUST_POINT"] = fmt.Sprintf("%s점", strconv.FormatFloat(oldCustomerTotal/float64(oldCustomer), 'f', 2, 64))
		} else {
			header["OLD_CUST_POINT"] = fmt.Sprintf("0점")
		}
	} else {
		header["MONTH_REVIEW"] = fmt.Sprintf("0건")
		header["MONTH_RLOW"] = fmt.Sprintf("0건")
		header["MONTH_RKEY"] = fmt.Sprintf("0건")

		header["NEW_CUST_LIST"] = fmt.Sprintf("0명")
		header["OLD_CUST_LIST"] = fmt.Sprintf("0명")

		header["NEW_CUST_POINT"] = fmt.Sprintf("0점")
		header["OLD_CUST_POINT"] = fmt.Sprintf("0점")
	}

	var keyword string
	for _, key := range reviewKeywords {
		keyword += fmt.Sprintf("'%s',", key)
	}

	if len(keyword) > 0 {
		header["REVIEW_KEYWORD"] = fmt.Sprintf("%s", keyword[:len(keyword)-1])
	} else {
		header["REVIEW_KEYWORD"] = "설정된 키워드 없음"
	}

	// review text ranking
	tr := textrank.NewTextRank()
	rule := textrank.NewDefaultRule()
	language := textrank.NewDefaultLanguage()
	algorithmDef := textrank.NewDefaultAlgorithm()

	for _, v := range reviews {
		tr.Populate(v["content"], language, rule)
	}

	tr.Ranking(algorithmDef)

	//rankedPhrases := textrank.FindPhrases(tr)
	words := textrank.FindSingleWords(tr)
	if len(words) > 0 {
		var word1, word2, word3 string
		var word1cnt, word2cnt, word3cnt int

		for k, v := range words {
			if k == 0 {
				word1 = v.Word
				word1cnt = v.Qty
			} else if k == 1 {
				word2 = v.Word
				word2cnt = v.Qty
			} else if k == 2 {
				word3 = v.Word
				word3cnt = v.Qty
				break
			}
		}

		if word1cnt > 0 {
			header["REVIEW_BWORD1"] = fmt.Sprintf("%s (%d건)", word1, word1cnt)
		} else {
			header["REVIEW_BWORD1"] = "작성된 리뷰가 없습니다."
		}

		if word2cnt > 0 {
			header["REVIEW_BWORD2"] = fmt.Sprintf("%s (%d건)", word2, word2cnt)
		} else {
			header["REVIEW_BWORD2"] = "작성된 리뷰가 없습니다."
		}

		if word3cnt > 0 {
			header["REVIEW_BWORD3"] = fmt.Sprintf("%s (%d건)", word3, word3cnt)
		} else {
			header["REVIEW_BWORD3"] = "작성된 리뷰가 없습니다."
		}

	} else {
		header["REVIEW_BWORD1"] = "작성된 리뷰가 없습니다."
		header["REVIEW_BWORD2"] = "작성된 리뷰가 없습니다."
		header["REVIEW_BWORD3"] = "작성된 리뷰가 없습니다."
	}

	var monthCntSorts []structCntSort

	for i := 0; i < len(selectMonthCntDetail); i++ {
		t0003, _ := strconv.Atoi(selectMonthCntDetail[i]["t0003"])
		t0306, _ := strconv.Atoi(selectMonthCntDetail[i]["t0306"])
		t0609, _ := strconv.Atoi(selectMonthCntDetail[i]["t0609"])
		t0912, _ := strconv.Atoi(selectMonthCntDetail[i]["t0912"])
		t1215, _ := strconv.Atoi(selectMonthCntDetail[i]["t1215"])
		t1518, _ := strconv.Atoi(selectMonthCntDetail[i]["t1518"])
		t1821, _ := strconv.Atoi(selectMonthCntDetail[i]["t1821"])
		t2124, _ := strconv.Atoi(selectMonthCntDetail[i]["t2124"])

		var monthCntSort structCntSort
		monthCntMsg := "%s %s시"

		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "0~3")
		monthCntSort.val = t0003
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "3~6")
		monthCntSort.val = t0306
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "6~9")
		monthCntSort.val = t0609
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "9~12")
		monthCntSort.val = t0912
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "12~15")
		monthCntSort.val = t1215
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "15~18")
		monthCntSort.val = t1518
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "18~21")
		monthCntSort.val = t1821
		monthCntSorts = append(monthCntSorts, monthCntSort)
		monthCntSort.key = fmt.Sprintf(monthCntMsg, selectMonthCntDetail[i]["day_name"], "21~24")
		monthCntSort.val = t2124
		monthCntSorts = append(monthCntSorts, monthCntSort)

		//tTime := []int{t0003, t0306, t0609, t0912, t1215, t1518, t1821, t2124}
		//tm,_ := dataCon.FindBusyTime2(tTime)
		//header[fmt.Sprintf("WEEK_BUSY%d",i+1)] = fmt.Sprintf("%s %s시", selectWeekCntDetail[i]["day_name"], tm)
	}

	// 동일 요일 포함해서 top3 색출
	sort.SliceStable(monthCntSorts, func(i, j int) bool {
		return monthCntSorts[i].val > monthCntSorts[j].val
	})

	for i, v := range monthCntSorts {
		//fmt.Println(fmt.Sprintf("key : %s, val : %d", v.key, v.val))
		header[fmt.Sprintf("MONTH_BUSY%d", i+1)] = v.key
		if i == 2 {
			break
		}
	}

	for i := 1; i < 3; i++ {
		val, exists := header[fmt.Sprintf("MONTH_BUSY%d", i)]
		if !exists {
			header[fmt.Sprintf("MONTH_BUSY%d", i)] = ""
		} else if strings.Contains(val, "0~") {
			header[fmt.Sprintf("MONTH_BUSY%d", i)] = ""
		}
	}

	/*
		for i:=0; i<2; i++{
			header[fmt.Sprintf("MONTH_BUSY%d",i+1)] = monthCntSorts[i].key
		}

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
			tm,_ := dataCon.FindBusyTime2(tTime)

			header[fmt.Sprintf("MONTH_BUSY%d",i+1)] = fmt.Sprintf("%s %s시", selectMonthCntDetail[i]["day_name"], tm)
		}
	*/

	header["MONTH_BUSY_DAY"] = selectMonthCntDetail[0]["day_name"]

	// body content
	contents := AlimTemplateCode(CASH_201, header)

	// struct
	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	//sendJMessage.To = fmt.Sprintf("82%s", userInfo[0]["hp"][1:])
	sendJMessage.Callphone = fmt.Sprintf("82-%s", userInfo[0]["hp"][1:])
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_201
	sendJMessage.MessageID = messageId

	// retry
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_201

	// button
	sendJButton.ButtonName = "앱으로 보기"
	sendJButton.ButtonType = "AL"
	// android
	sendJButton.ButtonURL = "https://darayos.page.link/vn1s"
	// ios
	sendJButton.ButtonURL2 = "https://darayos.page.link/vn1s"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "웹으로 자세히 보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "https://partner.darayo.com"
	sendJButton.ButtonURL2 = "https://partner.darayo.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	// kakao send data
	req, err := json.Marshal(sendJson)
	if err != nil {
		lprintf(1, "[ERROR] json Marshal err(%s)\n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// kakao send
	resp, err := cls.HttpsJson("POST", "rest.surem.com", "443", KakaoEndpoint, req)
	if err != nil {
		lprintf(1, "[ERROR] HttpsJson err(%s)\n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	var result KakaoResult
	json.Unmarshal(respBody, &result)

	lprintf(4, "[INFO] %v", result)

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userInfo[0]["hp"]
	params["templateCode"] = CASH_201
	params["templateNm"] = "월간분석_성공"
	params["result"] = sendResult
	params["userId"] = rUserInfo[0]["userId"]

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		m["resultCode"] = "98"
		m["resultMsg"] = err.Error()

		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func SendKakaoAlim(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	// 전송 시작
	//fname := cls.Cls_conf(os.Args)
	//usercode, _ := cls.GetTokenValue("KAKAOALIMTAKL.USERCODE", fname)
	//deptcode, _ := cls.GetTokenValue("KAKAOALIMTAKL.DEPTCODE", fname)
	//yellowid_key, _ := cls.GetTokenValue("KAKAOALIMTAKL.YELLOID_KEY", fname)

	templateCode := params["templateCode"]

	contents := AlimTemplate(templateCode, "김환")

	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	sendJMessage.Callphone = "82-1020047042"
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = templateCode
	sendJMessage.ReSend = "N"

	sendJButton.ButtonName = "달아요캐시 앱"
	sendJButton.ButtonType = "AL"
	sendJButton.ButtonURL = "https://play.google.com/store/apps/details?id=com.fit.darayos"
	sendJButton.ButtonURL2 = "https://itunes.apple.com/kr/app/%EB%8B%AC%EC%95%84%EC%9A%94-darayo-%EC%82%AC%EC%9E%A5%EB%8B%98-%EB%AA%A8%EB%B0%94%EC%9D%BC%EC%9E%A5%EB%B6%80/id1185364935?mt=8"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "달아요캐시 알아보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "http://www.darayocash.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	pbytes, _ := json.Marshal(sendJson)
	buff := bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://rest.surem.com/%s", KakaoEndpoint), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	str := string(respBody)
	//	println(str)
	//}

	//println(string(pbytes))

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = string(str)

	return c.JSON(http.StatusOK, m)

}

func AlimTemplateCode(tCode string, header map[string]string) string {

	var body string

	switch tCode {
	case CASH_001:
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입은 완료하셨지만\n매장정보(상호, 업종 등) 미입력 상태입니다.\n\n2)앱에서 매장정보를 입력 하시면\n매출, 고객, 리뷰 분석 서비스와\n기업 연결을 제공해 드립니다.\n\n3)이를 활용하면 고객과 매출을 이해하며\n단골과 고정매출 늘어나는 가게가 됩니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_002:
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입은 완료하셨지만\n여신협회 인증정보 미입력 상태입니다.\n\n2) 앱에서 인증정보를 입력해 주시면\n매출, 고객, 리뷰 분석 서비스와\n기업 연결을 제공해 드립니다.\n\n3)이를 활용하면 고객과 매출을 이해하며\n단골과 고정매출 늘어나는 가게가 됩니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_003:
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입이 완료되었으며\n데이터 수접 및 분석 중 입니다.\n\n2) 분석 완료 후 (1~2시간 소요)\n단골이 늘어나고 가게운영 돕는\n매출분석, 고객분석, 리뷰분석 서비스와\n기업연결 서비스를 받으실 수 있습니다.\n\n3) 파트너 회원에게는 매일, 매주, 매월\n알림톡으로 매출, 고객, 리뷰 분석 보고서를\n보내드리며, 보다 편하게 고객과 매출을 이해하고 단골과 고정매출이 늘어나는 가게가 됩니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_005:
		body = "[ 분석 결과 안내 ]\n안녕하세요. \"#{SHOP_NAME}\" 매장의\n어제분석 오늘예측 결과 안내입니다.\n\n[ 어제분석(땀) ]\n여신협회 사이트 로그인 과정에서\n비밀번호 오류로 어제 자료 수집이 불가 합니다.\n앱에 접속하여 여신협회 인증정보를 수정해주세요.\n\n1) 달아요 장부 사용액 [ #{YDAY_DARAYO_AMT} ]\n\n[ 오늘예측(오케이) ]\n1) 예상매출 [ #{TDAY_SALES} ]\n2) 바쁜시간 [ #{TDAY_BUSY} ]\n3) 결제단가 [ #{TDAY_PAY_SET} ]\n\n[ 내일을 준비해요(브이) ]\n1) 예상매출 [ #{TOMO_SALES} ]\n2) 바쁜시간 [ #{TOMO_BUSY} ]\n3) 결제단가 [ #{TOMO_PAY_SET} ]\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_013:
		body = "[분석 결과 안내]\n안녕하세요. 등록하신 \"#{SHOP_NAME}\" 의\n어제분석 오늘예측 결과 안내입니다.\n\n[ 어제분석 (최고)]\n1) 전체 매출 [ #{YDAY_SALES} ]\n2) 카드 입금 [ #{YDAY_DEPOS} ]\n3) 미입금 [ #{YDAY_MISS}]\n4) 달아요 장부 사용액 [ #{YDAY_DARAYO_AMT} ]\n\n► 취소 분석\n1)결제 취소 건수 [ #{YDAY_CANC}]\n2)관심 취소 [ #{YDAY_CANC_ALRT}]\n- 심야취소, 특정 시간 이후 취소 등\n\n► 고객 분석\n1) 리뷰 등록 건 [ #{YDAY_VIEW} ]\n2) 중요 리뷰 [ #{YDAY_VIEW_ALRT} ]\n- 설정한 평점 및 키워드 포함\n\n[ 오늘예측 (오케이) ]\n1) 예상매출 [ #{TDAY_SALES} ]\n2) 바쁜시간 [ #{TDAY_BUSY} ]\n3) 결제단가 [ #{TDAY_PAY_SET} ]\n\n►카드사 입금예정\n- 전체 [ #{TDAY_DEPOS} ]\n- 비씨 [ #{TDAY_DEPOS_BC} ]\n- 삼성 [ #{TDAY_DEPOS_SS} ]\n- 신한 [ #{TDAY_DEPOS_SH} ]\n- 국민 [ #{TDAY_DEPOS_KM} ]\n- 농협 [ #{TDAY_DEPOS_NH} ]\n- 현대 [ #{TDAY_DEPOS_HD} ]\n- 롯데 [ #{TDAY_DEPOS_LT} ]\n- 하나 [ #{TDAY_DEPOS_HN} ]\n\n[ 내일을 준비해요(브이) ]\n1) 예상매출 [ #{TOMO_SALES} ]\n2) 바쁜시간 [ #{TOMO_BUSY} ]\n3) 결제단가 [ #{TOMO_PAY_SET} ]\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_019:
		body = "[ #{SHOP_NAME} 어제 매출 안내 ]\n안녕하세요.\n어제를 분석하고 오늘을 예측한 결과를 안내해 드립니다.\n\n[ 어제분석 (최고) ]\n1) 전체 매출 [ #{YDAY_SALES} ]\n2) 카드 입금 [ #{YDAY_DEPOS} ]\n3) 미입금 [ #{YDAY_MISS} ]\n4) 달아요 장부 사용액 [ #{YDAY_DARAYO_AMT} ]\n\n► 어제 취소\n1) 결제 취소 건수 [ #{YDAY_CANC} ]\n2) 관심 취소 [ #{YDAY_CANC_ALRT} ]\n- 심야 취소, 예전 결제 취소 등\n\n► 어제 리뷰\n1) 전체 등록 리뷰 [ #{YDAY_VIEW} ]\n2) 관심 별점 리뷰 [ #{YDAY_VIEW_STAR} ]\n2) 관심 키워드 포함 리뷰 [ #{YDAY_VIEW_ALRT} ]\n등록 키워드: #{YDAY_VIEW_ALRT_DETAIL}\n\n[ 오늘 예측 (오케이) ]\n1) 예상 매출 [ #{TDAY_SALES} ]\n2) 바쁜 시간 [ #{TDAY_BUSY} ]\n3) 결제 단가 [ #{TDAY_PAY_SET} ]\n\n►카드사 입금예정\n- 전체 [ #{TDAY_DEPOS} ]\n- 비씨 [ #{TDAY_DEPOS_BC} ]\n- 삼성 [ #{TDAY_DEPOS_SS} ]\n- 신한 [ #{TDAY_DEPOS_SH} ]\n- 국민 [ #{TDAY_DEPOS_KM} ]\n- 농협 [ #{TDAY_DEPOS_NH} ]\n- 현대 [ #{TDAY_DEPOS_HD} ]\n- 롯데 [ #{TDAY_DEPOS_LT} ]\n- 하나 [ #{TDAY_DEPOS_HN} ]\n\n우리가게 돕는 오른팔\n달아요캐시에서 자세히 보고\n리뷰 관리도 해보세요"
	case CASH_014:
		body = "[분석 결과 안내]\n안녕하세요. \"#{SHOP_NAME}\" 매장의\n어제분석 오늘예측 결과 안내입니다.\n\n[ 어제분석(땀) ]\n여신협회 사이트 접속 문제로\n어제 자료 수집이 지연되고 있습니다.\n\n달아요에서 계속 체크해서 #{SEND_TIME}에 다시 알려드릴께요.\n\n1) 달아요 장부 사용액 [ #{YDAY_DARAYO_AMT} ]\n\n[ 오늘예측(오케이) ]\n1) 예상매출 [ #{TDAY_SALES} ]\n2) 바쁜시간 [ #{TDAY_BUSY} ]\n3) 결제단가 [ #{TDAY_PAY_SET} ]\n\n► 카드사 입금예정:\n- 전체 [ #{TDAY_DEPOS} ]\n- 비씨 [ #{TDAY_DEPOS_BC} ]\n- 삼성 [ #{TDAY_DEPOS_SS} ]\n- 신한 [ #{TDAY_DEPOS_SH} ]\n- 국민 [ #{TDAY_DEPOS_KM} ]\n- 농협 [ #{TDAY_DEPOS_NH} ]\n- 현대 [ #{TDAY_DEPOS_HD} ]\n- 롯데 [ #{TDAY_DEPOS_LT} ]\n- 하나 [ #{TDAY_DEPOS_HN} ]\n\n[ 내일을 준비해요(브이) ]\n1) 예상매출 [ #{TOMO_SALES} ]\n2) 바쁜시간 [ #{TOMO_BUSY} ]\n3) 결제단가 [ #{TOMO_PAY_SET} ]\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_015:
		body = "[분석 결과 안내]\n안녕하세요. \"#{SHOP_NAME}\" 매장의\n어제분석 오늘예측 결과 안내입니다.\n\n[ 어제분석(땀) ]\n여신협회 사이트 접속 문제로\n어제 자료 수집이 여전히 지연되고 있습니다.\n\n시간이 지난 후 앱에 접속하여 확인해 보세요.\n\n1) 달아요 장부 사용액 [ #{YDAY_DARAYO_AMT} ]\n\n[ 오늘예측(오케이) ]\n1) 예상매출 [ #{TDAY_SALES} ]\n2) 바쁜시간 [ #{TDAY_BUSY} ]\n3) 결제단가 [ #{TDAY_PAY_SET} ]\n\n► 카드사 입금예정:\n- 전체 [ #{TDAY_DEPOS} ]\n- 비씨 [ #{TDAY_DEPOS_BC} ]\n- 삼성 [ #{TDAY_DEPOS_SS} ]\n- 신한 [ #{TDAY_DEPOS_SH} ]\n- 국민 [ #{TDAY_DEPOS_KM} ]\n- 농협 [ #{TDAY_DEPOS_NH} ]\n- 현대 [ #{TDAY_DEPOS_HD} ]\n- 롯데 [ #{TDAY_DEPOS_LT} ]\n- 하나 [ #{TDAY_DEPOS_HN} ]\n\n[ 내일을 준비해요(브이) ]\n1) 예상매출 [ #{TOMO_SALES} ]\n2) 바쁜시간 [ #{TOMO_BUSY} ]\n3) 결제단가 [ #{TOMO_PAY_SET} ]\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_016:
		body = "[분석 결과 안내]\n안녕하세요. \"#{SHOP_NAME}\" 매장의\n어제분석 오늘예측 결과 안내입니다.\n\n[ 어제분석(땀) ]\n여신협회 사이트에 등록된\n주말 매출 자료가 완전하지 않아 데이터가 부족합니다.\n\n평일에 앱에 접속하여 확인해 보세요.\n\n달아요 장부 사용액 [ #{YDAY_DARAYO_AMT} ]\n\n[ 오늘예측(오케이) ]\n1) 예상매출 [ #{TDAY_SALES} ]\n2) 바쁜시간 [ #{TDAY_BUSY} ]\n3) 결제단가 [ #{TDAY_PAY_SET} ]\n\n► 주말에는 카드사 입금이 없습니다.\n\n[ 내일을 준비해요(브이) ]\n1) 예상매출 [ #{TOMO_SALES} ]\n2) 바쁜시간 [ #{TOMO_BUSY} ]\n3) 결제단가 [ #{TOMO_PAY_SET} ]\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	case CASH_101:
		body = "[지난주 분석 보고서]\n안녕하세요. 등록하신 #{SHOP_NAME}의\n지난 주 분석이 완료되었습니다.\n\n[ 지난주 매출 분석(최고) ]\n1) 지난주 매출[ #{WEEK_SALES} ]\n2) 지난주 최고 [#{WEEK_BEST} ]\n3) 지난주 최저 [#{WEEK_WRST} ]\n4) 결제가 많은 카드 [ #{WEEK_BCARD} ]\n5) 달아요 장부 사용액 [ #{WEEK_DARAYO} ]\n\n► 지난주 입금 정산 결과\n1) 카드사 입금 총액 [ #{WEEK_DEPOSIT} ]\n2) 카드사 미입금액 [ #{WEEK_MISS} ]\n\n► 지난주 취소 분석\n1) 지난주 전체 취소 [ #{WEEK_CANCEL} ]\n2) 심야 취소 [ #{WEEK_CNIGHT} ]\n3) 3시간 이후 취소 [ #{WEEK_CHOUR} ]\n\n[ 고객 및 리뷰 분석(오케이) ]\n1) 전체 등록 리뷰 [ #{WEEK_REVIEW} ]\n2) 중요리뷰 (평점) [ #{WEEK_RLOW} ]\n3) 키워드 포함 리뷰 [ #{WEEK_RKEY} ]\n키워드: #{REVIEW_KEYWORD}\n\n► 중요 리뷰 작성자 정보\n1) 신규 작성 고객 [ #{NEW_CUST_LIST} ]\n- 다른 가게 작성 평점 [ #{NEW_CUST_POINT} ]\n2) 단골 고객 [ #{OLD_CUST_LIST} ]\n- 우리 가게 작성 평점 [ #{OLD_CUST_POINT} ]\n\n► 고객들이 많이 방문한 시간\n1) 방문 best 1 [ #{WEEK_BUSY1} ]\n2) 방문 best 2 [ #{WEEK_BUSY2} ]\n3) 방문 best 3 [ #{WEEK_BUSY3} ]\n\n더 자세한 내용은 앱이나 웹에서 확인 하세요.\n\n지난 한 주 동안 수고 많으셨어요.\n사장님 오른팔 달아요 캐시가\n이번 주도 응원할께요!"
	case CASH_104:
		body = "[ #{SHOP_NAME} 지난주 장사 보고서 ]\n안녕하세요.\n지난주 분석을 확인하고\n한 주를 돌아보면서 이번 주를 준비해 보세요.\n\n[ 지난주 매출 (최고) ]\n1) 지난주 매출[ #{WEEK_SALES} ]\n2) 지난주 최고 #{WEEK_BEST_TIME} [ #{WEEK_BEST} ]\n3) 지난주 최저 #{WEEK_WRST_TIME} [ #{WEEK_WRST} ]\n4) 결제가 많은 카드 [ #{WEEK_BCARD} ]\n5) 달아요 장부 사용액 [ #{WEEK_DARAYO} ]\n\n► 지난주 입금 및 취소\n1) 카드사 입금 총액 [ #{WEEK_DEPOSIT} ]\n2) 카드사 미입금액 [ #{WEEK_MISS} ]\n3) 지난주 전체 취소 [ #{WEEK_CANCEL} ]\n4) 심야 취소 [ #{WEEK_CNIGHT} ]\n5) 3시간 이후 취소 [ #{WEEK_CHOUR} ]\n\n► 지난주 리뷰\n1) 전체 등록 리뷰 [ #{WEEK_REVIEW} ]\n2) 관심 별점 리뷰 [ #{WEEK_RLOW} ]\n3) 관심 키워드 포함 리뷰 [ #{WEEK_RKEY} ]\n키워드: #{REVIEW_KEYWORD}\n\n지난 한 주 동안 수고 많으셨어요.\n사장님 오른팔 달아요캐시가\n이번 주도 응원할게요!"
	case CASH_201:
		body = "[지난달 분석 보고서]\n안녕하세요. 등록하신 #{SHOP_NAME}의\n지난 달 분석이 완료되었습니다.\n\n[ 지난달 매출 분석(최고) ]\n1) 지난달 매출[ #{MONTH_SALES} ]\n2) 올해 평균 매출 [ #{MONTH_AVRG} ]\n3) 올해 최고 매출 [ #{MONTH_BEST} ]\n4) 결제가 많은 카드 [ #{MONTH_BCARD} ]\n5) 달아요 장부 사용액 [ #{MONTH_DARAYO} ]\n\n► 지난달 입금 정산 결과\n1) 카드사 입금액 [ #{MONTH_DEPOSIT}]\n2) 카드사 미입금액 [ #{MONTH_MISS_DEPO} ]\n\n► 지난달 취소 분석\n1) 지난달 전체 취소 [ #{MONTH_CANCEL} ]\n2) 심야 취소 [ #{MONTH_CNIGHT} ]\n3) 3시간 이후 취소 [ #{MONTH_CHOUR} ]\n\n[ 고객 및 리뷰 분석(오케이) ]\n1) 전체 등록 리뷰 [ #{MONTH_REVIEW} ]\n2) 중요리뷰 (평점) [ #{MONTH_RLOW} ]\n3) 키워드 포함 리뷰 [ #{MONTH_RKEY} ]\n키워드 #{REVIEW_KEYWORD}\n\n► 중요 리뷰 작성자 정보\n1) 신규 작성 고객 [ #{NEW_CUST_LIST} ]\n- 다른 가게 작성 평점 [ #{NEW_CUST_POINT} ]\n2) 단골 고객 [ #{OLD_CUST_LIST} ]\n- 우리 가게 작성 평점 [ #{OLD_CUST_POINT} ]\n\n► 고객들이 많이 남긴 키워드\n1) [ #{REVIEW_BWORD1} ]\n2) [ #{REVIEW_BWORD2} ]\n3) [ #{REVIEW_BWORD3} ]\n\n► 고객들이 많이 방문한 시간\n1) 방문 best 1 [ #{MONTH_BUSY1} ]\n2) 방문 best 2 [ #{MONTH_BUSY2} ]\n3) 방문 best 3 [#{MONTH_BUSY3} ]\n4) 가장 바쁜 요일 [ #{MONTH_BUSY_DAY} ]\n\n더 자세한 내용은 앱이나 웹에서 확인 하세요."
	case CASH_205:
		body = "[ #{SHOP_NAME} 지난달을 확인해 보세요 ]\n\n안녕하세요. 사장님 지난 한 달 동안 수고 많으셨습니다.\n달아요캐시에서 준비한 지난달 우리 가게 보고서에서 다양한 자료와 컨텐츠들을 확인해\n보시고 이번 달을 준비해 보세요.\n\n그리고 이번 달 대박 나세요!!\n\n사장님 오른팔 달아요캐시"
	case CASH_902:
		body = "[달아요캐시] 회원가입 상태 알림\n\n안녕하세요 #{NAME}님\n#{JOIN_DATE} 가입하신\n달아요에서 현재 서비스 가입상태를 알려드립니다.\n\n달아요는 사장님을 돕기 위해\n가맹점의 어제를 분석하고 오늘을 예측하는\n매출분석 서비스를 새롭게 추가하여\n달아요캐시로 서비스 하고 있습니다.\n\n[가입 상태 및 활용 안내]\n\n1) 회원가입 아이디는 \"#{LOGIN_ID}\" 이며\n여신협회 연결 상태입니다.\n\n장부 및 선불 서비스, 기업 장부(연결)과\n매출 및 고객분석, 리뷰분석 등 가게 운영을 돕는\n분석 서비스를 앱으로 확인할 수 있습니다.\n\n2) 본 메시지를 수신하신 고마운 기존 회원님께\n\n#{START_DATE} 부터\n매출 및 고객, 리뷰 분석 리포트를\n매일/주간/월간 알림톡으로 보내드리는\n파트너 서비스를 평생 무료로 제공합니다.\n\n3) 파트너 멤버는 알림톡 이외에도 웹을 통해\n승인-취소관리, 상세한 주간/월간매출 및 리뷰 보고서,\n리뷰 작성자 분석 서비스를 받을 수 있습니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔(최고)\n달아요 캐시에서 뵙겠습니다."
	case CASH_901:
		body = "[달아요] 회원가입 상태 알림\n\n안녕하세요 #{NAME}님\n#{JOIN_DATE} 가입하신\n달아요 가입 상태를 알려드립니다.\n\n달아요는 사장님을 돕기 위해\n가맹점의 어제를 분석하고 오늘을 예측하는\n매출분석 서비스를 새롭게 추가하여\n달아요캐시로 서비스 하고 있습니다.\n\n[가입 상태 및 활용 안내]\n\n1) 회원가입 아이디는 \"#{LOGIN_ID}\" 이며\n여신협회는 미연결 상태입니다.\n\n장부 및 선불 관리 기능, 기업 장부(연결)은 가능하지만\n새롭게 리뉴얼 된 매출 및 고객분석, 리뷰분석 등\n가게 운영을 돕는 분석 서비스를 앱에서 보실 수 없습니다.\n\n2) 본 메시지를 수신한 기존 회원님께서\n앱 또는 파트너 웹에서 여신협회 정보를 연결 하시면\n\n#{START_DATE}부터\n매출 및 고객, 리뷰 분석 리포트를\n매일/주간/월간 알림톡으로 보내드리는\n파트너 서비스를 평생 무료로 제공합니다.\n\n3) 파트너 멤버는 알림톡 이외에도 웹을 통해\n승인-취소관리, 상세한 주간/월간매출 및 리뷰 보고서,\n리뷰 작성자 분석 서비스를 받을 수 있습니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔(최고)\n달아요 캐시에서 뵙겠습니다."
	case CASH_004:
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입이 완료되었으며\n데이터 수접 및 분석 중 입니다.\n\n2) 분석 완료 후 (1~2시간 소요)\n단골이 늘어나고 가게운영 돕는\n매출분석, 고객분석, 리뷰분석 서비스와\n기업연결 서비스를 받으실 수 있습니다.\n\n3) 파트너 회원에게는 매일, 매주, 매월\n알림톡으로 매출, 고객, 리뷰 분석 보고서를\n보내드리며, 보다 편하게 고객과 매출을 이해하고 단골과 고정매출이 늘어 나도록\n달아요 캐시가 함께합니다.\n\n(별) 사장님에게만 드리는 #{COUPON_DATE} 무료 쿠폰!\n아래 '쿠폰 등록하기' 버튼을 클릭 후 #{LOGIN_ID}으로 접속하고\n발급받은 쿠폰을 등록해 보세요!\n\n#{COUPON} (크크)\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	default:
		return ""
	}

	for k, v := range header {
		body = strings.ReplaceAll(body, fmt.Sprintf("#{%s}", k), v)
	}

	return body
}

func AlimTemplate(templateCode string, data1 string) string {

	body := ""
	if templateCode == "cash_001" {
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입은 완료하셨지만\n매장정보(상호, 업종 등) 미입력 상태입니다.\n\n2)앱에서 매장정보를 입력 하시면\n매출, 고객, 리뷰 분석 서비스와\n기업 연결을 제공해 드립니다.\n\n3)이를 활용하면 고객과 매출을 이해하며\n단골과 고정매출 늘어나는 가게가 됩니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	} else if templateCode == "cash_002" {
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입은 완료하셨지만\n여신협회 인증정보 미입력 상태입니다.\n\n2) 앱에서 인증정보를 입력해 주시면\n매출, 고객, 리뷰 분석 서비스와\n기업 연결을 제공해 드립니다.\n\n3)이를 활용하면 고객과 매출을 이해하며\n단골과 고정매출 늘어나는 가게가 됩니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	} else if templateCode == "cash_003" {
		body = "[달아요캐시] 회원가입 완료\n\n안녕하세요 #{NAME}님\n어제를 분석하고 오늘을 예측으로\n사장님 돕는 오른팔\n\n달아요캐시에 가입 해주셔서\n진심으로 감사드립니다.\n\n[가입 상태 및 활용 안내]\n\n1) 현재 가입이 완료되었으며\n데이터 수접 및 분석 중 입니다.\n\n2) 분석 완료 후 (1~2시간 소요)\n단골이 늘어나고 가게운영 돕는\n매출분석, 고객분석, 리뷰분석 서비스와\n기업연결 서비스를 받으실 수 있습니다.\n\n3) 파트너 회원에게는 매일, 매주, 매월\n알림톡으로 매출, 고객, 리뷰 분석 보고서를\n보내드리며, 보다 편하게 고객과 매출을 이해하고 단골과 고정매출이 늘어나는 가게가 됩니다.\n\n오늘도 대박나는 하루 되시길 바라며\n사장님 돕는 오른팔\n달아요 캐시에서 뵙겠습니다."
	}

	body = strings.Replace(body, "#{NAME}", data1, -1)

	return body
}

func SendAlimRest(c echo.Context) error {

	RestKakaoAlim()

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

func Coupon_KakaoAlim(restId string) string {

	params_em := make(map[string]string)
	params_em["storeId"] = restId
	restInfo, err := cls.SelectData(storesql.SelectRestUserInfo, params_em)
	if err != nil {
		return "fail"
	}
	userNm := restInfo[0]["USER_NM"]
	userTel := restInfo[0]["HP_NO"]
	userId := restInfo[0]["USER_ID"]

	params_em["userId"] = userId

	header := make(map[string]string)
	header["NAME"] = userNm
	header["COUPON_DATE"] = "6개월"
	header["LOGIN_ID"] = restInfo[0]["LOGIN_ID"]
	header["COUPON"] = "cash6m"

	params_em["code"] = CASH_004
	params_em["today"] = time.Now().Format("20060102")

	// 이미 성공 리포트 보냈는지
	sendComp, err := cls.SelectData(storesql.SelectSuccessReportSendCheck, params_em)
	if err != nil {
		lprintf(1, "[ERROR] SelectSuccessReportSendCheck err(%s) \n", err.Error())
		return "fail"
	}

	if len(sendComp) > 0 && sendComp[0]["result"] == "success" {
		return "ok"
	}

	NewUserTel := "82-" + userTel[1:len(userTel)] //"821020047042" //
	//contents := AlimTemplate(CASH_004,userNm)
	contents := AlimTemplateCode(CASH_004, header)

	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params_em)

	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	sendJMessage.Callphone = NewUserTel
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = CASH_004
	sendJMessage.MessageID = messageId

	//재 전송값
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = CASH_004

	//버튼
	sendJButton.ButtonName = "달아요캐시 앱"
	sendJButton.ButtonType = "AL"
	sendJButton.ButtonURL = "https://play.google.com/store/apps/details?id=com.fit.darayos"
	sendJButton.ButtonURL2 = "https://itunes.apple.com/kr/app/%EB%8B%AC%EC%95%84%EC%9A%94-darayo-%EC%82%AC%EC%9E%A5%EB%8B%98-%EB%AA%A8%EB%B0%94%EC%9D%BC%EC%9E%A5%EB%B6%80/id1185364935?mt=8"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "파트너 쿠폰 등록하기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "https://partner.darayo.com/store/storeInfo/membership"
	sendJButton.ButtonURL2 = "https://partner.darayo.com/store/storeInfo/membership"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	pbytes, _ := json.Marshal(sendJson)
	buff := bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://rest.surem.com/%s", KakaoEndpoint), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	str := string(respBody)
	println(str)
	//}

	var result KakaoResult
	//err = json.Unmarshal(respBody, result)

	json.Unmarshal(respBody, &result)
	//fmt.Printf("Results: %v\n", result)
	//println(err)
	//if err != nil {
	//	return
	//}

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params := make(map[string]string)
	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userTel
	params["templateCode"] = CASH_004
	params["templateNm"] = "가입4번"
	params["result"] = sendResult
	params["userId"] = userId

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return ""
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	if err != nil {
		lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		return ""
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return ""
	}

	//println(string(pbytes))

	return "ok"
}

// 가입 완료
func JoinOk_KakaoAlim(restId string) string {

	// 전송 시작
	//fname := cls.Cls_conf(os.Args)
	//usercode, _ := cls.GetTokenValue("KAKAOALIMTAKL.USERCODE", fname)
	//deptcode, _ := cls.GetTokenValue("KAKAOALIMTAKL.DEPTCODE", fname)
	//yellowid_key, _ := cls.GetTokenValue("KAKAOALIMTAKL.YELLOID_KEY", fname)
	//smsNo, _ := cls.GetTokenValue("SMS.API.CALLBACK", fname)

	params_em := make(map[string]string)
	params_em["storeId"] = restId
	restInfo, err := cls.SelectData(storesql.SelectRestUserInfo, params_em)
	if err != nil {
		return "fail"
	}
	userNm := restInfo[0]["USER_NM"]
	userTel := restInfo[0]["HP_NO"]
	userId := restInfo[0]["USER_ID"]

	templateCode := "cash_003"
	NewUserTel := "82-" + userTel[1:len(userTel)] //"821020047042" //
	contents := AlimTemplate(templateCode, userNm)

	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params_em)

	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	sendJMessage.Type = "at"
	sendJMessage.Callphone = NewUserTel
	//sendJMessage.To = NewUserTel
	sendJMessage.Text = contents
	sendJMessage.TemplateCode = templateCode
	sendJMessage.MessageID = messageId

	//재 전송값
	sendJMessage.From = SmsNo
	sendJMessage.ReSend = "R"
	sendJMessage.ReText = templateCode

	//버튼
	sendJButton.ButtonName = "달아요캐시 앱"
	sendJButton.ButtonType = "AL"
	sendJButton.ButtonURL = "https://play.google.com/store/apps/details?id=com.fit.darayos"
	sendJButton.ButtonURL2 = "https://itunes.apple.com/kr/app/%EB%8B%AC%EC%95%84%EC%9A%94-darayo-%EC%82%AC%EC%9E%A5%EB%8B%98-%EB%AA%A8%EB%B0%94%EC%9D%BC%EC%9E%A5%EB%B6%80/id1185364935?mt=8"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "달아요캐시 알아보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "http://www.darayocash.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJson.Messages = append(sendJson.Messages, sendJMessage)

	pbytes, _ := json.Marshal(sendJson)
	buff := bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://rest.surem.com/%s", KakaoEndpoint), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	str := string(respBody)
	println(str)
	//}

	var result KakaoResult
	//err = json.Unmarshal(respBody, result)

	json.Unmarshal(respBody, &result)
	//fmt.Printf("Results: %v\n", result)
	//println(err)
	//if err != nil {
	//	return
	//}

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	params := make(map[string]string)
	params["messageId"] = strconv.Itoa(messageId)
	params["hpNo"] = userTel
	params["templateCode"] = templateCode
	params["templateNm"] = "가입3번"
	params["result"] = sendResult
	params["userId"] = userId

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
	if err != nil {
		return ""
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	if err != nil {
		lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		return ""
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return ""
	}

	//println(string(pbytes))

	return "ok"
}

// 미가입 가맹점 알림톡
func RestKakaoAlim() {

	// 전송 시작
	//fname := cls.Cls_conf(os.Args)
	//usercode, _ := cls.GetTokenValue("KAKAOALIMTAKL.USERCODE", fname)
	//deptcode, _ := cls.GetTokenValue("KAKAOALIMTAKL.DEPTCODE", fname)
	//yellowid_key, _ := cls.GetTokenValue("KAKAOALIMTAKL.YELLOID_KEY", fname)

	//smsNo, _ := cls.GetTokenValue("SMS.API.CALLBACK", fname)

	params_em := make(map[string]string)
	seqAlim, err := cls.SelectData(commonsql.SelctAlimTalkSeq, params_em)

	messageId, _ := strconv.Atoi(seqAlim[0]["alimSeq"])

	var sendJson KakaoAlimSend
	var sendJMessage KakaoMessages
	var sendJButton KakaoButtons

	sendJson.Usercode = UserCode
	sendJson.Deptcode = DeptCode
	sendJson.YellowidKey = YellowIdKey

	//버튼
	sendJButton.ButtonName = "달아요캐시 앱"
	sendJButton.ButtonType = "AL"
	sendJButton.ButtonURL = "https://play.google.com/store/apps/details?id=com.fit.darayos"
	sendJButton.ButtonURL2 = "https://itunes.apple.com/kr/app/%EB%8B%AC%EC%95%84%EC%9A%94-darayo-%EC%82%AC%EC%9E%A5%EB%8B%98-%EB%AA%A8%EB%B0%94%EC%9D%BC%EC%9E%A5%EB%B6%80/id1185364935?mt=8"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	sendJButton.ButtonName = "달아요캐시 알아보기"
	sendJButton.ButtonType = "WL"
	sendJButton.ButtonURL = "http://www.darayocash.com"
	sendJMessage.Buttons = append(sendJMessage.Buttons, sendJButton)

	restList, err := cls.SelectData(commonsql.SelectSendAlimRest, params_em)
	if err != nil {
		lprintf(4, "가맹점 가입 알림톡 전송 대상 없음 \n", err)
		return
	}

	seq := 0
	for i := range restList {

		userNm := restList[i]["USER_NM"]
		userTel := restList[i]["HP_NO"]
		tCode := restList[i]["T_CODE"]
		userId := restList[i]["USER_ID"]

		templateCode := ""

		if tCode == "1" {
			templateCode = "cash_001"
		} else if tCode == "2" {
			templateCode = "cash_002"
		}
		NewUserTel := "82-" + userTel[1:len(userTel)] //"821020047042" //
		contents := AlimTemplate(templateCode, userNm)
		messageId = messageId + seq

		sendJMessage.Type = "at"
		sendJMessage.Callphone = NewUserTel
		sendJMessage.Text = contents
		sendJMessage.TemplateCode = templateCode
		sendJMessage.MessageID = messageId

		//재 전송값
		sendJMessage.From = SmsNo
		sendJMessage.ReSend = "R"
		sendJMessage.ReText = templateCode

		sendJson.Messages = append(sendJson.Messages, sendJMessage)

		params := make(map[string]string)
		params["userId"] = userId
		params["messageId"] = strconv.Itoa(messageId)
		params["hpNo"] = userTel
		params["templateCode"] = templateCode
		params["templateNm"] = "가입" + tCode + "번"

		tx, err := cls.DBc.Begin()
		if err != nil {
			//return "5100", errors.New("begin error")
		}

		// 오류 처리
		defer func() {
			if err != nil {
				// transaction rollback
				lprintf(1, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
				tx.Rollback()
			}
		}()

		// 파라메터 맵으로 쿼리 변환
		insertFilterQuery, err := cls.SetUpdateParam(commonsql.InsertAlimTalkLog, params)
		if err != nil {
			return
		}
		// 쿼리 실행
		_, err = tx.Exec(insertFilterQuery)
		lprintf(4, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
		if err != nil {
			return
		}

		// transaction commit
		err = tx.Commit()
		if err != nil {
			return
		}

		seq = seq + 1

	}

	pbytes, _ := json.Marshal(sendJson)
	buff := bytes.NewBuffer(pbytes)

	//println(string(pbytes))
	//return

	req, err := http.NewRequest("POST", fmt.Sprintf("https://rest.surem.com/%s", KakaoEndpoint), buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	str := string(respBody)
	println(str)
	//}

	var result KakaoResult
	//err = json.Unmarshal(respBody, result)

	json.Unmarshal(respBody, &result)
	//fmt.Printf("Results: %v\n", result)
	//println(err)
	//if err != nil {
	//	return
	//}

	sendResult := result.Results[0].Result
	sendRMessageID := result.Results[0].MessageID

	println(sendRMessageID)

	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			lprintf(4, "do rollback -알림톡 로그 InsertAlimTalkLog) \n")
			tx.Rollback()
		}
	}()

	u_params := make(map[string]string)
	u_params["messageId"] = strconv.Itoa(sendRMessageID)
	u_params["result"] = sendResult

	// 파라메터 맵으로 쿼리 변환
	insertFilterQuery, err := cls.SetUpdateParam(commonsql.UpdateAlimTalkLog, u_params)
	if err != nil {
		return
	}
	// 쿼리 실행
	_, err = tx.Exec(insertFilterQuery)
	lprintf(1, "Query(%s) -> error (%s) \n", insertFilterQuery, err)
	if err != nil {
		return
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return
	}

}
