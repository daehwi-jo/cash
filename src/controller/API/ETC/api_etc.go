package apis

import (
	"bytes"
	apisql "cashApi/query/API"
	storesql "cashApi/query/stores"
	"cashApi/src/controller"
	"cashApi/src/controller/cls"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func BaseUrl(c echo.Context) error {
	return c.JSONP(http.StatusOK, "", "sms")
}


// 사업자 번호 조회
func BizNumCheck(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	bizNum :=params["bizNum"]


	//테스트용
	if strings.Contains(bizNum,"12345678") {
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		return c.JSON(http.StatusOK, m)
	}


	url :="https://teht.hometax.go.kr/wqAction.do?actionId=ATTABZAA001R08&screenId=UTEABAAA13&popupYn=false&realScreenId="
	xmlData :="<map id='ATTABZAA001R08'>" +
		"<pubcUserNo/>" +
		"<mobYn>N</mobYn>" +
		"<inqrTrgtClCd>1</inqrTrgtClCd>" +
		"<txprDscmNo>"+bizNum+"</txprDscmNo>" +
		"<dongCode>05</dongCode" +
		"><psbSearch>Y</psbSearch>" +
		"<map id='userReqInfoVO'/>" +
		"</map>"
	buf :=bytes.NewBufferString(xmlData)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var rMap Map
	xml.Unmarshal(body, &rMap)

	checkResult := rMap.SmpcBmanTrtCntn


	if strings.Contains(strings.Replace(checkResult, " ", "", -1),"등록되어있는사업자등록번호입니다") == false {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "사용할수 없는 사업자 번호 입니다."))
	}


	bizNumCnt, err := cls.GetSelectData(storesql.SelectBizNumCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	bizCnt,_:= strconv.Atoi(bizNumCnt[0]["bizCnt"])
	if bizCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 가입된 사업자 번호입니다."))
	}



	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	return c.JSON(http.StatusOK, m)
}





// 계좌실명조회
func AcctNameSearch(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	bankCode :=params["bankCode"]
	accountNo :=params["accountNo"]
	//buyerAuthNum :=params["buyerAuthNum"]

	url :="https://webtx.tpay.co.kr/api/v1/acct_name_search?"
	api_key := "xG3E5I+uuUvo+3ui/PKAPhxhutmQteOf3UiZ3PYG/zpO6fHsJZdlY28GOAWP09Kp7ArmIQdFlG7elvpTf/AKqQ=="
	mid := "darayo001m"
	bank_code := bankCode
	account := accountNo
	buyer_auth_num := "2222"

	urlParameters := "api_key="+api_key+ "&mid="+mid+ "&bank_code="+bank_code+ "&account="+account+"&buyer_auth_num="+buyer_auth_num;
	resp, err := http.Get(url+urlParameters)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result TpayResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	rAccountName :=result.Account_name
	rcode := result.Result_cd
	if rcode !="000"{
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "올바른 '계좌번호'를 넣어주세요."))
	}

	msg :="예금주명 '"+rAccountName+"'(이)가 맞습니까?"

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["msg"] = msg
	m["account_name"] = rAccountName

	return c.JSON(http.StatusOK, m)
}

// 공휴일 업데이트
func SetHoliday(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	solYear :=params["solYear"]

	pURL :="http://apis.data.go.kr/B090041/openapi/service/SpcdeInfoService/getRestDeInfo?pageNo=1&numOfRows=10000&_type=json&"
	serviceKey := "eV8UL93Tl94H%2B%2Bupnq9QEXFE2n15WFpD%2BdNxuHh3w6cmIvXFh4XgJF45HFkY4HdKnuWl%2FU2Xyqw8w5Mnur8Vpg%3D%3D"

	// 승인
	urlParameters := "solYear="+solYear+ "&ServiceKey="+serviceKey;

	resp, err := http.Get(pURL+urlParameters)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result HdResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//result.Response.Body.Items.Item[0].Locdate
	rIem := result.Response.Body.Items.Item

	for i, _ := range rIem {

		params["locDate"]= strconv.Itoa(rIem[i].Locdate)
		params["dateName"] = rIem[i].DateName
		selectQuery, err := cls.GetQueryJson(apisql.UpdateHolyDay, params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "query parameter fail"))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(selectQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultMsg22"] = string(respBody)

	return c.JSON(http.StatusOK, m)
}



// 테블릿 주문
func SimpleOrder(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	// 간편주문 가능한 장부 수 체크  1개만 있어야함
	grpCntData, err := cls.GetSelectData(apisql.SelectSimpleGrpCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	grpCnt,_:= strconv.Atoi(grpCntData[0]["grpCnt"])
	if grpCnt > 1 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "사용중인 장부가 여러개 입니다."))
	}

	//장부 체크
	grpData, err := cls.GetSelectData(apisql.SelectSimpleGrp, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if grpData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "장부 정보가 없습니다."))
	}


	// 식권 체크
	menuData, err := cls.GetSelectData(apisql.SelectStoreMenuView, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if menuData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "지정 메뉴가 없습니다."))
	}

	grpId := grpData[0]["GRP_ID"]
	userId := grpData[0]["USER_ID"]
	authStat := grpData[0]["AUTH_STAT"]
	userNm := grpData[0]["USER_NM"]
	itemNo := menuData[0]["ITEM_NO"]
	menuPrice := menuData[0]["ITEM_PRICE"]
	totalAmt := menuPrice


	params["grpId"] = grpId
	params["userId"] = userId
 	params["totalAmt"] = totalAmt


 	if authStat !="1" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "장부 사용이 불가능한 사용자입니다."))
	}


	// 동일 주문 체크시간 불러오기
	grpChkData, err := cls.GetSelectData(apisql.SelectGrpOrderCheckData, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	params["checkTime"] = grpChkData[0]["CHECK_TIME"]


	//시간 동일 주문 체크 
	orderChkData, err := cls.GetSelectData(apisql.SelectOrderCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	orderCnt,_:= strconv.Atoi(orderChkData[0]["orderCnt"])
	if orderCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "동일주문(주문내역 확인 필요)"))
	}


	orderData, err := cls.GetSelectData(apisql.CreateOrderSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	orderNo := orderData[0]["orderNo"] + userId




	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -태블릿 주문-간편 주문 (SetStoreCharging)  \n")
			tx.Rollback()
		}
	}()

	params["orderNo"] = orderNo
	params["creditAmt"] = totalAmt
	params["orderTy"] = "5"
	params["payTy"] = "1"
	params["qrOrderType"] = "0"


	//주문 등록
	InsertSimpleOrderQuery, err := cls.SetUpdateParam(apisql.InsertSimpleOrder, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertSimpleOrderQuery parameter fail"))
	}

	_, err = tx.Exec(InsertSimpleOrderQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertSimpleOrderQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertSimpleOrderQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}



	params["orderSeq"] = "0"
	params["itemNo"] = itemNo
	params["orderAmt"] = totalAmt
	params["orderQty"] = "1"
	params["userId"] = userId


	//주문 상세 등록
	InsertSimpleOrderDetailQuery, err := cls.SetUpdateParam(apisql.InsertSimpleOrderDetail, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertSimpleOrderDetailQuery parameter fail"))
	}

	_, err = tx.Exec(InsertSimpleOrderDetailQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertSimpleOrderDetailQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertSimpleOrderDetailQuery, err)
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
	m["USER_NM"] = userNm

	return c.JSON(http.StatusOK, m)
}

func accessMainPage() ([]*http.Cookie, int) {
	apiUrl := "https://hometax.go.kr"
	resource := "/permission.do"

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	data := url.Values{}
	data.Set("screenId", "index_pp")
	urlStr := u.String()

	body := "<map id='postParam'><popupYn>false</popupYn></map>"

	urlStr = "https://hometax.go.kr/permission.do?screenId=index_pp"

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(body)))
	if err != nil {
		lprintf(1, "[FAIL]login: http NewRequest (%s) \n", err.Error())
		return nil, -1
	}

	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "20")
	req.Header.Add("Content-Type", "application/xml; charset=UTF-8")
	req.Header.Add("Origin", "https://www.hometax.go.kr")
	req.Header.Add("Referer", "https://hometax.go.kr/websquare/websquare.wq?w2xPath=/ui/pp/index_pp.xml")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		lprintf(1, "[FAIL]permission: http (%s) \n", err)
		return nil, -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		lprintf(4, "[INFO]resp=(%s) \n", resp)
		if err != nil {
			lprintf(1, "[FAIL]permission: %s \n", err.Error())
		}
		return nil, -1
	}

	cookie := resp.Cookies()
	lprintf(4, "[INFO] permission cookies (%v)", cookie)
	return cookie, 1

}

func HometaxLogin(c echo.Context) error {

	m := make(map[string]interface{})

	cookie, rst := accessMainPage()
	if rst < 0 {
		m["resultCode"] = "99"
		m["resultMsg"] = "access fail"
		return c.JSON(http.StatusOK, m)
	}

	params := cls.GetParamJsonMap(c)

	apiUrl := "https://hometax.go.kr"

	encodeId := b64.StdEncoding.EncodeToString([]byte(params["loginId"]))
	hashPass := fmt.Sprintf("%X", sha256.Sum256([]byte(params["password"])))
	encodePass := b64.StdEncoding.EncodeToString([]byte(hashPass))
	lprintf(4, "[INFO] encode : TE5NzQxRjU3MzY0Q0Y1OENEQzMxOTY0MENCNTc2RUY1M0Y3QzAzMkNBQTJEMTg2QTNFMEQwOTM4NkI3NjlDRA== (%v)", encodePass)

	data := url.Values{}
	data.Set("ssoLoginYn", "Y")
	data.Set("secCardLoginYn", "")
	data.Set("secCardId", "")
	data.Set("cncClCd", "01")
	data.Set("id", encodeId)
	data.Set("pswd", encodePass)
	data.Set("ssoStatus", "")
	data.Set("portalStatus", "")
	data.Set("scrnId", "UTXPPABA01")
	data.Set("userScrnRslnXcCnt", "2048")
	data.Set("userScrnRslnYcCnt", "1152")

	resource := "/pubcLogin.do?domain=hometax.go.kr&mainSys=Yn"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	urlStr = "https://hometax.go.kr/pubcLogin.do?domain=hometax.go.kr&mainSys=Yn"
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		lprintf(1, "[FAIL]login: http NewRequest (%s) \n", err.Error())
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()
		return c.JSON(http.StatusOK, m)
	}

	req.Header.Set("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.hometax.go.kr")
	req.Header.Add("Referer", "https://hometax.go.kr/websquare/websquare.wq?w2xPath=/ui/comm/a/b/UTXPPABA01.xml&w2xHome=/ui/pp/&w2xDocumentRoot=")

	for i := range cookie {
		req.AddCookie(cookie[i])
	}

	lprintf(4, "[INFO]request(%v) \n", req)
	lprintf(4, "[INFO]request body(%v) \n", data.Encode())
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		lprintf(1, "[FAIL]login: http (%s) \n", err)
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		lprintf(4, "[INFO]resp=(%s) \n", resp)
		if err != nil {
			lprintf(1, "[FAIL]login: %s \n", err.Error())
		}
		m["resultCode"] = "99"
		m["resultMsg"] = "login fail"
		return c.JSON(http.StatusOK, m)
	}

	cookie = resp.Cookies()
	result := fmt.Sprintf("%v", cookie)
	if strings.Index(result, "NTS_LOGIN_SYSTEM_CODE_P=TXPP") < 0 {
		// login fail
		lprintf(1, "[FAIL]login fail: check id/passwd (%s) \n", result)
		m["resultCode"] = "99"
		m["resultMsg"] = "login fail"
		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "login success"
	return c.JSON(http.StatusOK, m)
}

func CardsalesLogin(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	apiUrl := "https://www.cardsales.or.kr"
	resource := "/authentication"
	data := url.Values{}
	data.Set("j_username", params["loginId"])
	data.Set("j_password", params["password"])
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	m := make(map[string]interface{})

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()
		return c.JSON(http.StatusOK, m)
	}

	req.Header.Set("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://www.cardsales.or.kr")
	req.Header.Add("Referer", "https://www.cardsales.or.kr/signin")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302{
		m["resultCode"] = "00"
		m["resultMsg"] = "login Sucess"
	}else{
		m["resultCode"] = "99"
		m["resultMsg"] = "login Fail"
	}

	return c.JSON(http.StatusOK, m)
}

// 쿠폰 사용
func CouponUse(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	// 쿠폰 유효성 체크
	couponData, err := cls.GetSelectData(apisql.SelectCouponInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if couponData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "사용가능한 쿠폰이 아닙니다."))
	}

	//couponNo := params["couponNo"]
	useType := couponData[0]["USE_TYPE"]
	useYn := couponData[0]["USE_YN"]

//	couponName := couponData[0]["COUPON_NAME"]
	itemCode := couponData[0]["ITEM_CODE"]
	couponVal := couponData[0]["COUPON_VAL"]

	if useType=="0" && useYn=="Y"{
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 사용된 쿠폰입니다."))
	}


	//사용 히스토리 체크
	params["userKey"] = params["storeId"]
	params["userType"] ="S"

	couponChk, err := cls.GetSelectData(apisql.SelectCouponChk, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	couponCnt,_ := strconv.Atoi(couponChk[0]["cnt"])
	if couponCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 사용된 쿠폰입니다."))
	}

	runType :="I"
	alimTalkInfo, err := cls.GetSelectData(storesql.SelectAlimTalkInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if alimTalkInfo == nil{
		runType="I"
	}



	// TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -파트너 쿠폰 사용 (CouponUse)  \n")
			tx.Rollback()
		}
	}()

	//파트너등록

	params["itemCode"]= itemCode
	params["useMonth"]= couponVal
	params["payYn"]= "N"
	params["etc"]= " "

	billingChk, err := cls.GetSelectData(apisql.SelectBillingChk, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	strQuery := ""
	if billingChk == nil {
		strQuery = apisql.InsertBillingCouponUse
	} else {
		strQuery = apisql.UpdateBillingCouponUse
	}


	// 빌링 등록
	BillingCouponUseQuery, err := cls.GetQueryJson(strQuery, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "BillingCouponUseQuery parameter fail"))
	}
	_, err = tx.Exec(BillingCouponUseQuery)
	dprintf(4, c, "call set Query (%s)\n", BillingCouponUseQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", BillingCouponUseQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 빌링 결제 내역 등록
	BillingPaymentQuery, err := cls.SetUpdateParam(apisql.InserBillingPayment, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertBillingCouponUse parameter fail"))
	}
	_, err = tx.Exec(BillingPaymentQuery)
	dprintf(4, c, "call set Query (%s)\n", BillingPaymentQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", BillingPaymentQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	// 쿠폰 사용 기록 등록
	InserCouponHistoryQuery, err := cls.SetUpdateParam(apisql.InserCouponHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertBillingCouponUse parameter fail"))
	}
	_, err = tx.Exec(InserCouponHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InserCouponHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InserCouponHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	if runType=="I"{
		params["kakaoWeek"]= "Y"
		params["kakaoMonth"]= "Y"
		params["kakaoDaily"]= "Y"
		insertStoreAlimtalkQuery, err := cls.GetQueryJson(storesql.InsertStoreAlimtalk, params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UserCreateQuery parameter fail"))
		}

		_, err = tx.Exec(insertStoreAlimtalkQuery)
		if err != nil {
			dprintf(1, c, "Query(%s) -> error (%s) \n", insertStoreAlimtalkQuery, err)
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
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







func BillingInsert(storeId,userId,itemCode,useMonth string) int {

	params := make(map[string]string)


	// TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollbac
			tx.Rollback()
		}
	}()

	//파트너등록
	params["itemCode"]= itemCode
	params["storeId"]= storeId
	params["useMonth"]= useMonth
	params["userId"]= userId
	params["payYn"]= "N"
	params["couponNo"]= "free"
	params["etc"]= "가입시 무료한달"
	


	// 빌링 등록
	BillingCouponUseQuery, err := cls.GetQueryJson(apisql.InsertBillingCouponUse, params)
	if err != nil {
		return 1
	}
	_, err = tx.Exec(BillingCouponUseQuery)
	lprintf(4, "call set Query (%s)\n" , BillingCouponUseQuery)
	if err != nil {
		lprintf(4, "Query(%s) -> error (%s) \n", BillingCouponUseQuery, err)
		return 1
	}

	// 빌링 결제 내역 등록
	BillingPaymentQuery, err := cls.SetUpdateParam(apisql.InserBillingPayment, params)
	if err != nil {
		lprintf(4, "Query(%s) -> error (%s) \n", BillingPaymentQuery, err)
		return 1
	}
	_, err = tx.Exec(BillingPaymentQuery)
	lprintf(4, "call set Query (%s)\n" , BillingPaymentQuery)
	if err != nil {
		lprintf(4, "Query(%s) -> error (%s) \n", BillingPaymentQuery, err)
		return 1
	}



	// transaction commit
	err = tx.Commit()
	if err != nil {
		return 1
	}

	return 0
}
