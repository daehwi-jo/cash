package books

import (
	booksql "cashApi/query/books"
	commons "cashApi/query/commons"
	linksql "cashApi/query/links"
	ordersql "cashApi/query/orders"
	paymentsql "cashApi/query/payments"
	"cashApi/src/controller"
	"fmt"
	"net/http"
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
	return c.JSONP(http.StatusOK, "", "books")
}


// 연결된 장부 목록
func GetLinkBookList(c echo.Context) error {

	dprintf(4, c, "call GetLinkBookList\n")

	params := cls.GetParamJsonMap(c)

	// 페이징 처리
	pageSize,_ := strconv.Atoi(params["pageSize"])
	pageNo,_ := strconv.Atoi(params["pageNo"])
	offset := strconv.Itoa((pageNo-1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset

	LinkBooksCnt, err := cls.GetSelectData(booksql.SelectLinkBooksCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	totalCount,_ :=strconv.Atoi(LinkBooksCnt[0]["bookCnt"])
	totalPage := strconv.Itoa((totalCount/pageSize) + 1)



	bookList, err := cls.GetSelectType(booksql.SelectLinkBookList + commons.PagingQuery , params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if bookList == nil {
		LinkData := make(map[string]interface{})
		LinkData["totalPage"] = totalPage
		LinkData["curruntPage"] = pageNo
		LinkData["totalAmt"] = 0
		LinkData["linkList"] = []string{}
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = LinkData
		return c.JSON(http.StatusOK, m)
	}





	LinkData := make(map[string]interface{})
	LinkData["linkCnt"] = totalCount
	LinkData["linkList"] = bookList
	LinkData["totalPage"] = totalPage
	LinkData["curruntPage"] = pageNo

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = LinkData

	return c.JSON(http.StatusOK, m)

}



// 장부별 주문
func GetBooksOrders(c echo.Context) error {

	dprintf(4, c, "call GetBooksOrders\n")

	params := cls.GetParamJsonMap(c)

	searchDate :=params["searchDate"]
	now :=time.Now()
	inputDate :=now.Format("2006-01-02")

	if searchDate == "1W" {
		inputDate = now.AddDate(0,0, - 7).Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate,"-","",-1)
	}else if searchDate == "1M" {
		inputDate = now.AddDate(0,-1, 0).Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate,"-","",-1)
	}else if searchDate == "3M" {
		inputDate = now.AddDate(0,-3, 0).Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate,"-","",-1)
	}else if searchDate == "1D" {
		inputDate = now.Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate, "-", "", -1)
	}
	unpaidYn :=params["unpaidYn"]

	if unpaidYn =="Y" {
		params["inputPaid"] ="N"
		params["orderStat"] ="20"
	}else{
		params["orderStat"] ="20','21"
	}


	//페이징처리
	pageSize,_ := strconv.Atoi(params["pageSize"])
	pageNo,_ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo-1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset


	orderListCount, err := cls.GetSelectData(booksql.SelectBookOrderCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	totalCount,_ := strconv.Atoi(orderListCount[0]["totalCount"])
	totalPage := strconv.Itoa((totalCount/pageSize) + 1)



	orderTotal, err := cls.GetSelectData(booksql.SelectBookOrderTotal, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}



	orderList, err := cls.GetSelectType(booksql.SelectBookOrderList + commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	if orderList == nil {
		LinkData := make(map[string]interface{})
		LinkData["totalPage"] = totalPage
		LinkData["curruntPage"] = pageNo
		LinkData["totalAmt"] = 0
		LinkData["orderList"] = []string{}
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = LinkData
		return c.JSON(http.StatusOK, m)
	}



	totalAmt,_ :=strconv.Atoi(orderTotal[0]["totalAmt"])
	LinkData := make(map[string]interface{})
	LinkData["totalAmt"] = totalAmt
	LinkData["totalPage"] = totalPage
	LinkData["curruntPage"] = pageNo
	LinkData["orderList"] = orderList


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = LinkData

	return c.JSON(http.StatusOK, m)

}




func GetOrderInfo(c echo.Context) error {

	dprintf(4, c, "call GetOrderInfo\n")

	params := cls.GetParamJsonMap(c)

	orderInfo, err := cls.GetSelectData(ordersql.SelectOrderInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if orderInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "잘못된 주문번호 입니다."))
	}



	menuDetail, err := cls.GetSelectType(ordersql.SelectOrderDetail, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if menuDetail == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "응답 실패"))
	}


	order := make(map[string]interface{})
	order["orderNo"] = orderInfo[0]["ORDER_NO"]
	order["storeNm"] = orderInfo[0]["REST_NM"]
	order["bookNm"] = orderInfo[0]["GRP_NM"]
	order["orderMemo"] = orderInfo[0]["ORDER_COMMENT"]
	totalAmt, _ := strconv.Atoi(orderInfo[0]["TOTAL_AMT"])
	order["totalAmt"] = totalAmt
	order["orderStat"] = orderInfo[0]["ORDER_STAT"]
	order["orderDate"] = orderInfo[0]["ORDER_DATE"]
	order["menu"] = menuDetail
	//order["users"] = orderList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = order

	return c.JSON(http.StatusOK, m)

}



// 결제 조회
func GetBooksPayments(c echo.Context) error {

	dprintf(4, c, "call GetBooksOrders\n")

	params := cls.GetParamJsonMap(c)

	searchDate :=params["searchDate"]
	now :=time.Now()
	inputDate :=now.Format("2006-01-02")

	if searchDate == "1W" {
		inputDate = now.AddDate(0,0, - 7).Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate,"-","",-1)
	}else if searchDate == "1M" {
		inputDate = now.AddDate(0,-1, 0).Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate,"-","",-1)
	}else if searchDate == "3M" {
		inputDate = now.AddDate(0,-3, 0).Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate,"-","",-1)
	}else if searchDate == "1D" {
		inputDate = now.Format("2006-01-02")
		params["searchDate"] = strings.Replace(inputDate, "-", "", -1)
	}


	payChannel :=params["appPayYn"]

	if payChannel =="Y" {
		params["payChannel"] ="02"
	}


	//페이징처리
	pageSize,_ := strconv.Atoi(params["pageSize"])
	pageNo,_ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo-1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset


	PaymemtListCount, err := cls.GetSelectData(booksql.SelectBookPaymemtCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	totalCount,_ := strconv.Atoi(PaymemtListCount[0]["totalCount"])
	totalPage := strconv.Itoa((totalCount/pageSize) + 1)


	orderTotal, err := cls.GetSelectData(booksql.SelectBookPaymemtTotal, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	totalAmt,_ :=strconv.Atoi(orderTotal[0]["totalAmt"])


	resultList, err := cls.GetSelectType(booksql.SelectBookPaymemtList  + commons.PagingQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultList == nil {
		LinkData := make(map[string]interface{})
		LinkData["totalPage"] = totalPage
		LinkData["totalAmt"] = totalAmt
		LinkData["curruntPage"] = pageNo
		LinkData["paymentList"] = []string{}
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = LinkData
		return c.JSON(http.StatusOK, m)
	}


	LinkData := make(map[string]interface{})
	LinkData["totalAmt"] = totalAmt
	LinkData["paymentList"] = resultList
	LinkData["totalPage"] = totalPage
	LinkData["curruntPage"] = pageNo

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = LinkData

	return c.JSON(http.StatusOK, m)

}



// 장부 충전금액 조회
func GetStoreChargeAmt(c echo.Context) error {

	dprintf(4, c, "call GetStoreChargeAmt\n")

	params := cls.GetParamJsonMap(c)
	linkInfo, err := cls.GetSelectData(linksql.SelectLinkInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if linkInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "협약 내용이 없습니다."))
	}
	LinkData := make(map[string]interface{})
	LinkData["prePaidAmt"] = linkInfo[0]["PREPAID_AMT"]


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] =LinkData

	return c.JSON(http.StatusOK, m)

}



// 매장 충전
func SetStoreCharging(c echo.Context) error {

	dprintf(4, c, "call SetStoreCharging\n")

	params := cls.GetParamJsonMap(c)

	linkInfo, err := cls.GetSelectData(linksql.SelectLinkInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if linkInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "협약 내용이 없습니다."))
	}

	creditAmt, _ := strconv.Atoi(params["creditAmt"])
	addAmt, _ := strconv.Atoi(params["addAmt"])
	linkId := linkInfo[0]["LINK_ID"]
	prepaidAmt, _ := strconv.Atoi(linkInfo[0]["PREPAID_AMT"])

	params["payChannel"] = "01"
	params["linkId"] = linkId

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -매장 충전(SetStoreCharging)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 금액 충전
	params["prepaidAmt"] = strconv.Itoa(prepaidAmt + creditAmt + addAmt)
	UpdateLinkQuery, err := cls.SetUpdateParam(linksql.UpdateLink, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateLinkQuery parameter fail"))
	}

	_, err = tx.Exec(UpdateLinkQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateLinkQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateLinkQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 충전금 내역 insert
	params["prepaidAmt"] = strconv.Itoa(creditAmt + addAmt)
	params["jobTy"] = "0"
	InsertPrepaidQuery, err := cls.SetUpdateParam(paymentsql.InsertPrepaid, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPrepaidQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPrepaidQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPrepaidQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPrepaidQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 충전 히스토리 insert
	params["prepaidAmt"] = strconv.Itoa(creditAmt + addAmt)
	params["jobTy"] = "0"
	params["searchTy"] = "1"
	params["paymentTy"] = "0"
	params["userTy"] = "1"

	now := time.Now()
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	diff := now.Sub(then)

	moid := fmt.Sprintf("%d", diff.Milliseconds())
	moid = moid + params["userId"]
	params["moid"] = moid

	InsertPaymentHistoryQuery, err := cls.SetUpdateParam(paymentsql.InsertPaymentHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPaymentHistoryQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPaymentHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPaymentHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPaymentHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 매장 충전 취소
func SetStoreChargingCancel(c echo.Context) error {

	dprintf(4, c, "call SetStoreChargingCancel\n")

	params := cls.GetParamJsonMap(c)

	cancelInfo, err := cls.GetSelectData(paymentsql.SelectStoreCancelCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	cancelCnt, _ := strconv.Atoi(cancelInfo[0]["CancelCnt"])

	if cancelCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 취소된 결제 입니다."))
	}

	chargeInfo, err := cls.GetSelectData(paymentsql.SelectStoreChargeInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if chargeInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "결제 내역이 없습니다."))
	}

	//creditAmt, _ := strconv.Atoi(chargeInfo[0]["CREDIT_AMT"])
	//addAmt, _ := strconv.Atoi(chargeInfo[0]["ADD_AMT"])
	totalAmt, _ := strconv.Atoi(chargeInfo[0]["TOTAL_AMT"])

	params["creditAmt"] = chargeInfo[0]["CREDIT_AMT"]
	params["addAmt"] = chargeInfo[0]["ADD_AMT"]

	params["userId"] = chargeInfo[0]["USER_ID"]
	params["payInfo"] = chargeInfo[0]["PAY_INFO"]
	params["bookId"] = chargeInfo[0]["BOOK_ID"]

	linkInfo, err := cls.GetSelectData(linksql.SelectLinkInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if linkInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "협약 내용이 없습니다."))
	}

	prepaidAmt, _ := strconv.Atoi(linkInfo[0]["PREPAID_AMT"])
	linkId := linkInfo[0]["LINK_ID"]
	params["payChannel"] = "01"
	params["linkId"] = linkId

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -매장 충전 취소 (SetStoreChargingCancel)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 금액 충전
	params["prepaidAmt"] = strconv.Itoa(prepaidAmt - totalAmt)
	UpdateLinkQuery, err := cls.SetUpdateParam(linksql.UpdateLink, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateLinkQuery parameter fail"))
	}

	_, err = tx.Exec(UpdateLinkQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateLinkQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateLinkQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 환불금 내역 insert
	params["prepaidAmt"] = strconv.Itoa(totalAmt)
	params["jobTy"] = "1"
	InsertPrepaidQuery, err := cls.SetUpdateParam(paymentsql.InsertPrepaid, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPrepaidQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPrepaidQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPrepaidQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPrepaidQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 충전 히스토리 insert
	params["prepaidAmt"] = strconv.Itoa(totalAmt)
	params["jobTy"] = "1"
	params["searchTy"] = "1"
	params["paymentTy"] = "1"
	params["userTy"] = "1"

	InsertPaymentHistoryQuery, err := cls.SetUpdateParam(paymentsql.InsertPaymentHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPaymentHistoryQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPaymentHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPaymentHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPaymentHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}



func GetUnPaidList(c echo.Context) error {

	dprintf(4, c, "call GetUnPaidList\n")

	params := cls.GetParamJsonMap(c)

	//페이징처리
	pageSize,_ := strconv.Atoi(params["pageSize"])
	pageNo,_ := strconv.Atoi(params["pageNo"])

	offset := strconv.Itoa((pageNo-1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	params["offSet"] = offset


	resultData, err := cls.GetSelectData(ordersql.SelectUnpaidListCount, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUnpaidListCount err(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	totalCount,_ := strconv.Atoi(resultData[0]["orderCnt"])
	totalPage := strconv.Itoa((totalCount/pageSize) + 1)





	unPaidList, err := cls.GetSelectType(ordersql.SelectUnpaidList  + commons.PagingQuery, params, c)
	if err != nil {
		lprintf(1, "[ERROR] SelectUnpaidList err(%s) \n", err.Error())
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	if unPaidList == nil {
		LinkData := make(map[string]interface{})
		LinkData["totalPage"] = totalPage
		LinkData["curruntPage"] = pageNo
		LinkData["totalAmt"] = 0
		LinkData["linkList"] = []string{}
		LinkData["bookNm"] =""
		m := make(map[string]interface{})
		m["resultCode"] = "00"
		m["resultMsg"] = "응답 성공"
		m["resultData"] = LinkData
		return c.JSON(http.StatusOK, m)
	}

	totalAmt, _ := strconv.Atoi(resultData[0]["TOTAL_AMT"])
	orderCnt, _ := strconv.Atoi(resultData[0]["orderCnt"])

	result := make(map[string]interface{})
	result["totalCnt"] = orderCnt
	result["totalAmt"] = totalAmt
	result["totalPage"] = totalPage
	result["curruntPage"] = pageNo
	result["accountList"] = unPaidList
	result["bookNm"] = resultData[0]["BOOK_NM"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)

}


// 매장 정산 처리
func SetPaidOk(c echo.Context) error {

	dprintf(4, c, "call SetPaidOk\n")

	params := cls.GetParamJsonMap(c)

	totalAmt, _ := strconv.Atoi(params["totalAmt"])

	if totalAmt == 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "데이터가 부족합니다."))
	}

	unPaidInfo, err := cls.GetSelectData(ordersql.SelectUnpaidListCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if unPaidInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "요청하신 정산 기준일에 미처리 정산 내역이 없습니다."))
	}

	rtotalAmt, _ := strconv.Atoi(unPaidInfo[0]["TOTAL_AMT"])
	userId := unPaidInfo[0]["USER_ID"]

	if rtotalAmt != totalAmt {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "정산금액이 요청하신 내역과 다릅니다."))
	}



	params["payChannel"] = "01"
	params["userId"] = userId

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -매장 충전(SetStoreCharging)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 충전 히스토리 insert
	params["creditAmt"] = strconv.Itoa(totalAmt)
	params["addAmt"] = "0"
	params["searchTy"] = "2"
	params["paymentTy"] = "3"
	params["userTy"] = "1"

	now := time.Now()
	then := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	diff := now.Sub(then)

	moid := fmt.Sprintf("%d", diff.Milliseconds())
	moid = moid + params["userId"]
	params["moid"] = moid

	UpdateOrderPaidQuery, err := cls.SetUpdateParam(ordersql.UpdateOrderPaid, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateOrderPaidQuery parameter fail"))
	}

	_, err = tx.Exec(UpdateOrderPaidQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateOrderPaidQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateOrderPaidQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	InsertPaymentHistoryQuery, err := cls.SetUpdateParam(paymentsql.InsertPaymentHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPaymentHistoryQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPaymentHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPaymentHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPaymentHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 매장 정산 취소
func SetPaidCancel(c echo.Context) error {

	dprintf(4, c, "call SetPaidCancel\n")

	params := cls.GetParamJsonMap(c)

	cancelInfo, err := cls.GetSelectData(paymentsql.SelectStoreCancelCnt, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	cancelCnt, _ := strconv.Atoi(cancelInfo[0]["CancelCnt"])

	if cancelCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이미 취소된 결제 입니다."))
	}

	chargeInfo, err := cls.GetSelectData(paymentsql.SelectStoreChargeInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if chargeInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "결제 내역이 없습니다."))
	}

	params["creditAmt"] = chargeInfo[0]["CREDIT_AMT"]
	params["addAmt"] = chargeInfo[0]["ADD_AMT"]
	params["userId"] = chargeInfo[0]["USER_ID"]
	params["payInfo"] = chargeInfo[0]["PAY_INFO"]
	params["bookId"] = chargeInfo[0]["BOOK_ID"]
	params["accStDay"] = chargeInfo[0]["ACC_ST_DAY"]

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -매장 충전(SetStoreCharging)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 충전 히스토리 insert

	params["searchTy"] = "2"
	params["paymentTy"] = "4"
	params["userTy"] = "1"
	params["payChannel"] = "01"

	UpdateOrderPaidQuery, err := cls.SetUpdateParam(ordersql.UpdateOrderPaidCancel, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateOrderPaidCancel parameter fail"))
	}

	_, err = tx.Exec(UpdateOrderPaidQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateOrderPaidQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateOrderPaidQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	InsertPaymentHistoryQuery, err := cls.SetUpdateParam(paymentsql.InsertPaymentHistory, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertPaymentHistoryQuery parameter fail"))
	}

	_, err = tx.Exec(InsertPaymentHistoryQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertPaymentHistoryQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertPaymentHistoryQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}



// 매장 주문취소
func SetOrderCancel(c echo.Context) error {

	dprintf(4, c, "call SetOrderCancel\n")

	params := cls.GetParamJsonMap(c)

	orderInfo, err := cls.GetSelectData(ordersql.SelectOrder, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if orderInfo == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "주문 내용이 없습니다."))
	}

	payTy := orderInfo[0]["PAY_TY"]
	orderStat := orderInfo[0]["ORDER_STAT"]
	totalAmt, _ := strconv.Atoi(orderInfo[0]["TOTAL_AMT"])
	bookId := orderInfo[0]["BOOK_ID"]
	storeId := orderInfo[0]["STORE_ID"]
	userId := orderInfo[0]["USER_ID"]
	pointUse,_ :=  strconv.Atoi(orderInfo[0]["POINT_USE"])


	params["bookId"] = bookId
	params["storeId"] = storeId
	params["userId"] = userId

	if orderStat != "20" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "취소가 불가능한 주문입니다."))
	}

	// 매장 충전  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -주문 취소(SetOrderCancel)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	// 선불일 경우 금액 환불
	if payTy == "0" {

		linkInfo, err := cls.GetSelectData(linksql.SelectLinkInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		if linkInfo == nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("99", "협약이 내용이 없습니다."))
		}
		prepaidAmt, _ := strconv.Atoi(linkInfo[0]["PREPAID_AMT"])
		linkId := linkInfo[0]["LINK_ID"]

		// 포인트 화불
		params["linkId"] = linkId
		params["prepaidAmt"] = strconv.Itoa(prepaidAmt + totalAmt)
		prepaidPoint, _ := strconv.Atoi(linkInfo[0]["PREPAID_POINT"])
		params["prepaidPoint"] = strconv.Itoa(prepaidPoint + pointUse)


		UpdateLinkQuery, err := cls.SetUpdateParam(linksql.UpdateLink, params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", err.Error()))
		}

		_, err = tx.Exec(UpdateLinkQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

	}

	//지원금 환불

	params["orderAmt"]= strconv.Itoa(totalAmt)
	UserSupportBalanceUpdateQuery, err := cls.SetUpdateParam(booksql.UpdateBookUserSupportBalance, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", err.Error()))
	}
	_, err = tx.Exec(UserSupportBalanceUpdateQuery)
	if err != nil {

		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//주문 취소
	OrderCancelQuery, err := cls.SetUpdateParam(ordersql.UpdateOrderCancel, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "OrderCancel parameter fail"))
	}
	_, err = tx.Exec(OrderCancelQuery)
	if err != nil {

		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)




}

