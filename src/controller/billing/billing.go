package billing

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	//commonsql "shopApi/query/commons"
	"cashApi/src/controller"

	"cashApi/src/controller/cls"

	"github.com/labstack/echo/v4"
)

/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
// darayo printf(로그레벨, 요청 컨텍스트, format, arg) => 무엇을(서비스, 요청), 어떻게(input), 왜(원인,조치)
var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func BillingCheck(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	// 맨처음 구독 신청시 6개월 무료
	// 취소후 신청시 기간이 남아있으면 구독 리로드
	// 재 신청시 가긴이 지나 있으면 카드 등록 및 금액 청구 처리

	BillingInfo, _ := cls.GetSelectData(SelectBillingInfo, params, c)
	if BillingInfo == nil {

		SetBillingReg(c)
	} else {
		payYn := BillingInfo[0]["PAY_YN"]
		strEndDate := BillingInfo[0]["END_DATE"]
		endDate, _ := time.Parse("2006-01-02", strEndDate)
		nowDate := time.Now()

		if payYn == "N" && !nowDate.After(endDate) {
			SetBillingCancelReload(c)
		} else if payYn == "Y" && !nowDate.After(endDate) {
			return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 구독중입니다."))
		} else if payYn == "N" && nowDate.After(endDate) {
			ViewBillingReg(c)
		} else if payYn == "Y" && nowDate.After(endDate) {
			ViewBillingReg(c)
		}
	}

	return c.HTML(http.StatusOK, "")

}

// 구독 등록
func ViewBillingReg(c echo.Context) error {

	dprintf(4, c, "call ViewBillingReg\n")

	goUrl := "/billing/b_reg.htm"

	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	papleJsUrl, _ := cls.GetTokenValue("PAYPLE.PAYPLE_JS_URL", fname)

	BillingHistoryYn, err := cls.GetSelectData(SelectBillingFreeCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	bCnt, _ := strconv.Atoi(BillingHistoryYn[0]["bCnt"])

	if bCnt > 0 {

		itemInfo, err := cls.GetSelectData(SelectBillingItemInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		if itemInfo == nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("99", "잘못된 상품입니다."))
		}

		payOidSeq, err := cls.GetSelectData(SelectPayOidSeq, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		todate := time.Now().Format("2006-01")

		itemPirce, _ := strconv.Atoi(itemInfo[0]["PRICE"])
		itemDcPirce, _ := strconv.Atoi(itemInfo[0]["DC_PRICE"])
		params["itemNm"] = itemInfo[0]["ITEM_NAME"]
		params["totalAmt"] = strconv.Itoa(itemPirce - itemDcPirce)
		params["payOid"] = payOidSeq[0]["payOid"]
		params["payYear"] = todate[:4]
		params["payMonth"] = todate[5:]
		goUrl = "/billing/b_pay_reg.html"
	}

	userInfo, err := cls.GetSelectData(SelectUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if userInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "유저정보가 없습니다."))
	}
	userId := params["userId"]
	storeId := params["storeId"]
	itemCode := params["itemCode"]
	os := params["os"]
	userData := userId + "@" + storeId + "@" + itemCode + "@" + os

	m := make(map[string]interface{})
	m["paple_js_url"] = papleJsUrl
	m["userId"] = userInfo[0]["USER_ID"]
	m["hpNo"] = userInfo[0]["USER_TEL"]
	m["email"] = userInfo[0]["EMAIL"]
	m["userNm"] = userInfo[0]["USER_NM"]
	m["userData"] = userData
	m["itemNm"] = params["itemNm"]
	m["totalAmt"] = params["totalAmt"]
	m["payOid"] = params["payOid"]
	m["payYear"] = params["payYear"]
	m["payMonth"] = "02" //params["payMonth"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, goUrl, m)
}

// 카드 등록만
func ViewBillingResult(c echo.Context) error {

	dprintf(4, c, "call ViewBillingResult\n")
	params := cls.GetParamJsonMap(c)

	PCD_PAY_RST := params["PCD_PAY_RST"]
	//PCD_PAY_MSG := params["PCD_PAY_MSG"]
	PCD_USER_DEFINE1 := params["PCD_USER_DEFINE1"]
	sdata := strings.Split(PCD_USER_DEFINE1, "@")

	userId := sdata[0]
	storeId := sdata[1]
	itemCode := sdata[2]
	os := sdata[3]

	params["userId"] = userId
	params["storeId"] = storeId
	params["itemCode"] = itemCode
	params["os"] = os

	if PCD_PAY_RST == "success" {

		billingInfo, err := cls.GetSelectData(SelectBillingInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		strQuery := ""
		if billingInfo == nil {
			strQuery = InsertBillingKey
		} else {
			strQuery = UpdateBillingKey
		}

		params["freeUseYn"] = "Y"
		params["useMonth"] = "1"
		params["payYn"] = "Y"

		selectQuery, err := cls.SetUpdateParam(strQuery, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(selectQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "등록 성공"
	m["os"] = os

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "/billing/b_result.html", m)
}

// 등록과 금액 결제 동시
func ViewBillingPayResult(c echo.Context) error {

//	dprintf(4, c, "call ViewBillingResult\n")
	params := cls.GetParamJsonMap(c)

	PCD_PAY_RST := params["PCD_PAY_RST"]
	PCD_PAY_MSG := params["PCD_PAY_MSG"]
	PCD_USER_DEFINE1 := params["PCD_USER_DEFINE1"]
	PCD_PAY_OID := params["PCD_PAY_OID"]
	sdata := strings.Split(PCD_USER_DEFINE1, "@")

	dprintf(4, c, "call BillingResult PAY_OID : %s , PAY_MSG : %S   \n " ,PCD_PAY_OID ,PCD_PAY_MSG )

	userId := sdata[0]
	storeId := sdata[1]
	itemCode := sdata[2]
//os := sdata[3]

	params["userId"] = userId
	params["storeId"] = storeId
	params["itemCode"] = itemCode
	//params["os"] = os

	//결제 히스토리 등록

	params["payType"] = "P"
	InsertBillingHistory, err := cls.SetUpdateParam(InsertBillingPayment, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(InsertBillingHistory)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if PCD_PAY_RST == "success" {

		billingInfo, err := cls.GetSelectData(SelectBillingInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		strQuery := ""
		if billingInfo == nil {
			strQuery = InsertBillingKey
		} else {
			strQuery = UpdateBillingKey
			params["bId"] = billingInfo[0]["B_ID"]
		}

		params["freeUseYn"] = "Y"
		params["useMonth"] = "1"
		params["payYn"] = "Y"



		//카드 등록
		selectQuery, err := cls.SetUpdateParam(strQuery, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(selectQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		//결제 정보 등록

		params["payStat"] = "20"
		//결제 히스토리 등록
		UpdateBilling, err := cls.SetUpdateParam(UpdateBillingPayment, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(UpdateBilling)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		//파트너 페이지 아이디 등록
		//RegBizMember(storeId,userId)

	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "등록 성공"
	//m["os"] = os

	dprintf(4, c, "BillingFinish PAY_OID : %s\n " ,PCD_PAY_OID )

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "/billing/b_result.htm", m)
}

// 로그인
func GetBillingInfo(c echo.Context) error {

	dprintf(4, c, "call GetBillingInfo\n")
	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectType(SelectBillingInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "구독 정보가 없습니다."))
	}

	//	dprintf(4, c, "call GetUserInfo\n", resultData)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData[0]

	return c.JSON(http.StatusOK, m)

}

func BillingAuth(c echo.Context) error {

	dprintf(4, c, "call BillingAuth\n")
	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	cst_id, _ := cls.GetTokenValue("PAYPLE.PAYPLE_CUST_ID", fname)
	pURL, _ := cls.GetTokenValue("PAYPLE.PAYPLE_URL", fname)
	cust_key, _ := cls.GetTokenValue("PAYPLE.PAYPLE_CUST_KEY", fname)
	REFERE_URL, _ := cls.GetTokenValue("PAYPLE.REFERE_URL", fname)


	//pURL := "https://testcpay.payple.kr/php/auth.php"
	//cst_id := "test"
	//cust_key := "abcd1234567890"
	//pcd_refund_key := "a41ce010ede9fcbfb3be86b24858806596a9db68b79d138b147c3e563e1829a0"
	//REFERE_URL := "http://172.30.1.9:7000/"
	work := params["work"]

	var reqData AuthSendData
	reqData.Cst_id = cst_id
	reqData.CustKey = cust_key
	reqData.PCD_PAY_WORK = work
	pbytes, _ := json.Marshal(reqData)
	buff := bytes.NewBuffer(pbytes)

	req, err := http.NewRequest("POST", pURL, buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("referer", REFERE_URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	str := ""
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str = string(respBody)
	}
	//	dprintf(4, c, "call GetUserInfo\n", resultData)

	return c.HTML(http.StatusOK, str)

}

// 구독 일반 , 이벤트 등록
func SetBillingReg(c echo.Context) error {

	dprintf(4, c, "call SetBillingReg\n")
	//
	params := cls.GetParamJsonMap(c)

	// 파라메터 맵으로 쿼리 변환

	params["freeUseYn"] = "Y"
	params["useMonth"] = "6"
	params["payYn"] = "Y"

	selectQuery, err := cls.SetUpdateParam(InsertBillingKey, params)
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

func SetBillingCancel(c echo.Context) error {

	dprintf(4, c, "call SetBillingCancel\n")
	//
	params := cls.GetParamJsonMap(c)

	// 파라메터 맵으로 쿼리 변환

	updateCategoryQuery, err := cls.SetUpdateParam(UpdateBillingCancel, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(updateCategoryQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

func SetBillingCancelReload(c echo.Context) error {

	dprintf(4, c, "call SetBillingCancelReload\n")
	//
	params := cls.GetParamJsonMap(c)

	// 파라메터 맵으로 쿼리 변환

	updateCategoryQuery, err := cls.SetUpdateParam(UpdateBillingCancelReload, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(updateCategoryQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

//결제 취소
func CallPaymentCancel(c echo.Context) error {

	dprintf(4, c, "call CallPaymentCancel\n")
	//
	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	cst_id, _ := cls.GetTokenValue("PAYPLE.PAYPLE_CUST_ID", fname)
	pURL, _ := cls.GetTokenValue("PAYPLE.PAYPLE_URL", fname)
	cust_key, _ := cls.GetTokenValue("PAYPLE.PAYPLE_CUST_KEY", fname)
	REFERE_URL, _ := cls.GetTokenValue("PAYPLE.REFERE_URL", fname)
	refundKey, _ := cls.GetTokenValue("PAYPLE.PAYPLE_REFUND_KEY", fname)
	cancelUrl, _ := cls.GetTokenValue("PAYPLE.PAY_CANCEL_URL", fname)

	// 파라메터 맵으로 쿼리 변환

	var reqData PaymentAuthData
	reqData.Cst_id = cst_id
	reqData.CustKey = cust_key
	reqData.PCD_REGULER_FLAG = "Y"
	pbytes, _ := json.Marshal(reqData)
	buff := bytes.NewBuffer(pbytes)

	// 승인
	req, err := http.NewRequest("POST", pURL, buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("referer", REFERE_URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//str := ""
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	//	str = string(respBody)
	//}

	var authResult AuthResult
	err = json.Unmarshal(respBody, &authResult)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//println(authResult.AuthKey)
	//println(authResult.Cst_id)

	auth_result := authResult.Result
	PCD_CUST_KEY := authResult.CustKey
	PCD_CST_ID := authResult.Cst_id
	PCD_AUTH_KEY := authResult.AuthKey

	println(auth_result)

	if auth_result != "success" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "결제 승인 오류."))
	}

	totalAmt, _ := strconv.Atoi(params["totalAmt"])

	var payData PaymentSendData

	payData.PCD_CST_ID = PCD_CST_ID
	payData.PCD_CUST_KEY = PCD_CUST_KEY
	payData.PCD_AUTH_KEY = PCD_AUTH_KEY
	payData.PCD_REFUND_KEY = refundKey
	payData.PCD_PAYCANCEL_FLAG = "Y"
	payData.PCD_PAY_OID = params["PCD_PAY_OID"]
	payData.PCD_REGULER_FLAG = "Y"
	payData.PCD_PAY_YEAR = params["PCD_PAY_YEAR"]
	payData.PCD_PAY_MONTH = params["PCD_PAY_MONTH"]
	payData.PCD_PAY_DATE = params["PCD_PAY_DATE"]
	payData.PCD_REFUND_TOTAL = totalAmt
	payData.PCD_PAY_TAXTOTAL = params["PCD_PAY_TAXTOTAL"]

	payReqbytes, _ := json.Marshal(payData)
	payReqbuff := bytes.NewBuffer(payReqbytes)

	//결제 요청
	payReq, err := http.NewRequest("POST", cancelUrl, payReqbuff)
	if err != nil {
		panic(err)
	}
	payReq.Header.Add("Content-Type", "application/json")
	payReq.Header.Add("referer", REFERE_URL)

	client_pay := &http.Client{}
	payResp, err := client_pay.Do(payReq)
	if err != nil {
		panic(err)
	}
	defer payResp.Body.Close()
	//str := ""
	// Response 체크.
	payRespBody, err := ioutil.ReadAll(payResp.Body)

	if err == nil {
		println(string(payRespBody))
	}

	var payResult PaymentResultData
	err = json.Unmarshal(payRespBody, &payResult)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

//결제 요청
func CallPayment(c echo.Context) error {

	dprintf(4, c, "call CallPayment\n")
	//
	params := cls.GetParamJsonMap(c)

	fname := cls.Cls_conf(os.Args)
	cst_id, _ := cls.GetTokenValue("PAYPLE.PAYPLE_CUST_ID", fname)
	pURL, _ := cls.GetTokenValue("PAYPLE.PAYPLE_URL", fname)
	cust_key, _ := cls.GetTokenValue("PAYPLE.PAYPLE_CUST_KEY", fname)
	REFERE_URL, _ := cls.GetTokenValue("PAYPLE.REFERE_URL", fname)
	PAY_REQ_URL, _ := cls.GetTokenValue("PAYPLE.PAY_REQ_URL", fname)

	// 파라메터 맵으로 쿼리 변환

	var reqData PaymentAuthData
	reqData.Cst_id = cst_id
	reqData.CustKey = cust_key
	reqData.PCD_PAY_TYPE = "card"
	reqData.PCD_REGULER_FLAG = "Y"
	pbytes, _ := json.Marshal(reqData)
	buff := bytes.NewBuffer(pbytes)

	// 승인
	req, err := http.NewRequest("POST", pURL, buff)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("referer", REFERE_URL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//str := ""
	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	//if err == nil {
	//	str = string(respBody)
	//}

	var authResult AuthResult
	err = json.Unmarshal(respBody, &authResult)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//println(authResult.AuthKey)
	//println(authResult.Cst_id)

	auth_result := authResult.Result
	PCD_CUST_KEY := authResult.CustKey
	PCD_CST_ID := authResult.Cst_id
	PCD_AUTH_KEY := authResult.AuthKey

	println(auth_result)

	if auth_result != "success" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "결제 승인 오류."))
	}

	todate := time.Now().Format("2006-01")
	totalAmt, _ := strconv.Atoi(params["totalAmt"])

	var payData PaymentSendData

	payData.PCD_CST_ID = PCD_CST_ID
	payData.PCD_CUST_KEY = PCD_CUST_KEY
	payData.PCD_AUTH_KEY = PCD_AUTH_KEY
	payData.PCD_PAY_TYPE = "card"
	payData.PCD_PAYER_ID = params["PAYER_ID"]
	payData.PCD_PAYER_NO = params["userId"][1:]
	payData.PCD_PAYER_HP = params["hpNo"]
	payData.PCD_PAYER_EMAIL = params["email"]
	payData.PCD_PAY_GOODS = "달아요-구독 자동결제"
	payData.PCD_PAY_TOTAL = totalAmt
	payData.PCD_PAY_ISTAX = "N"
	//payData.PCD_PAY_OID = payOidSeq[0]["payOid"]
	payData.PCD_PAY_YEAR = todate[:4]
	payData.PCD_PAY_MONTH = todate[5:]
	payData.PCD_REGULER_FLAG = "Y"

	payReqbytes, _ := json.Marshal(payData)
	payReqbuff := bytes.NewBuffer(payReqbytes)

	//결제 요청
	payReq, err := http.NewRequest("POST", PAY_REQ_URL, payReqbuff)
	if err != nil {
		panic(err)
	}
	payReq.Header.Add("Content-Type", "application/json")
	payReq.Header.Add("referer", REFERE_URL)

	client_pay := &http.Client{}
	payResp, err := client_pay.Do(payReq)
	if err != nil {
		panic(err)
	}
	defer payResp.Body.Close()
	//str := ""
	// Response 체크.
	payRespBody, err := ioutil.ReadAll(payResp.Body)

	if err == nil {
		println(string(payRespBody))
	}

	var payResult PaymentResultData
	err = json.Unmarshal(payRespBody, &payResult)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}


// test
func BillingRegView(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(SelectUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "회원 정보가 없습니다."))
	}


	itemData, err := cls.GetSelectData(SelectBillingItemInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if itemData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "아이템 정보가 없습니다."))
	}

	fname := cls.Cls_conf(os.Args)
	paple_js_url, _ := cls.GetTokenValue("PAYPLE.PAYPLE_JS_URL", fname)


	now := time.Now()

	nowTime :=now.Format("2006-01-02 15:04:05")
	nowDate :=now.Format("2006-01-02")
	sNowDate := strings.Split(nowDate,"-")


	storeId := params["storeId"]
	userId := params["userId"]
	itemCode := params["itemCode"]


	pSet := strings.Replace(strings.Replace(strings.Replace(nowTime,"-","",-1),":","",-1) ," " ,"",-1)
	payOid := "P"+pSet+"_"+ storeId

	userData := userId+"@"+storeId+"@"+itemCode;

	m := make(map[string]interface{})
	m["userId"] = resultData[0]["USER_ID"]
	m["userNm"] = resultData[0]["USER_NM"]
	m["hpNo"] = resultData[0]["HP_NO"]
	m["email"] = resultData[0]["EMAIL"]
	m["paple_js_url"] = paple_js_url

	m["itemNm"] = itemData[0]["ITEM_NAME"]
	m["totalAmt"] = itemData[0]["ITEM_PRICE_DC"]
	m["payOid"] = payOid
	m["payYear"] = sNowDate[0]
	m["payMonth"] = sNowDate[1]
	m["userData"] = userData


	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "billing/b_reg.htm", m)
}




func RegBizMember(storeId, userId string) {


	params := make(map[string]string)
	params["userId"] = userId
	params["storeId"] = storeId
	storeInfo, err := cls.SelectData(SelectStoreUserInfo, params)
	if err != nil {
		return
	}

	strQuery := ""
	if storeInfo == nil {
		strQuery = InsertSysInfo
	} else {
		strQuery = UpdateSysInfo
	}
	strQuery = InsertSysInfo
	params["loginId"] = storeInfo[0]["LOGIN_ID"]
	params["userPass"] = storeInfo[0]["LOGIN_PW"]
	params["userNm"] = storeInfo[0]["REST_NM"]
	params["authorCd"] = "SM"

	updateStoreService, err := cls.GetQueryJson(strQuery, params)
	if err != nil {
		return
	}
	// 쿼리 실행
	_, err = cls.QueryDB(updateStoreService)
	if err != nil {
		return
	}



}


