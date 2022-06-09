package stores

import (
	"bytes"
	"cashApi/src/controller"
	"cashApi/src/controller/cls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	storesql "cashApi/query/stores"
	apiPush "cashApi/src/controller/API/PUSH"
)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func GetStoreCategories(c echo.Context) error {

	dprintf(4, c, "call GetStoreCategories\n")

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})

	addQuery := ""
	if params["useYn"] != "" {
		addQuery += `AND USE_YN = '#{useYn}'`
	}
	orderQuery := "ORDER BY USE_YN DESC"

	resultList, err := cls.GetSelectType(storesql.SelectStoreCategories+addQuery+orderQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultList == nil {
		m["resultCode"] = "00"
		m["resultMsg"] = "카테고리 정보가 없습니다."
		m["resultList"] = []string{}

		return c.JSON(http.StatusOK, m)
	}


	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)

}

func SetStoreInsertCategories(c echo.Context) error {

	dprintf(4, c, "call SetStoreCategories\n")

	params := cls.GetParamJsonMap(c)
	cateSeq, err := cls.GetSelectDataRequire(storesql.SelectStoreItemCodeSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if cateSeq == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "카테고리 ID  생성 실패."))
	}
	params["codeId"] = cateSeq[0]["codeId"]
	params["categoryId"] = cateSeq[0]["categoryId"]
	params["categoryNm"] = cateSeq[0]["categoryNm"]

	// 파라메터 맵으로 쿼리 변환
	insertCategoryQuery, err := cls.GetQueryJson(storesql.InsertStoreCategories, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(insertCategoryQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

func SetStoreUpdateCategories(c echo.Context) error {

	dprintf(4, c, "call SetStoreCategories\n")
	//
	params := cls.GetParamJsonMap(c)

	// 파라메터 맵으로 쿼리 변환
	updateCategoryQuery, err := cls.SetUpdateParam(storesql.UpdateStroeCategories, params)
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

func GetStoreMenuList(c echo.Context) error {

	dprintf(4, c, "call GetStoreMenuInfo\n")

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})

	addQuery := ""
	if params["codeId"] != "all" {
		addQuery += `AND B.CODE_ID = '#{codeId}'`
	}
	orderQuery := "ORDER BY A.USE_YN DESC"

	resultList, err := cls.GetSelectTypeRequire(storesql.SelectStoreMenuList+addQuery+orderQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {
		m["resultCode"] = "00"
		m["resultMsg"] = "메뉴 정보가 없습니다."
		m["resultList"] = []string{}

		return c.JSON(http.StatusOK, m)
	}

	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)

}

func SetStoreInsertMenu(c echo.Context) error {

	dprintf(4, c, "call SetStoreInsertMenu\n")
	//
	params := cls.GetParamJsonMap(c)
	menuSeq, err := cls.GetSelectData(storesql.SelectStoreMenuSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if menuSeq == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "메뉴 ID  생성 실패."))
	}
	params["itemNo"] = menuSeq[0]["itemNo"]

	// 파라메터 맵으로 쿼리 변환
	insertMenuQuery, err := cls.GetQueryJson(storesql.InsertStroeMenu, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(insertMenuQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

func SetStoreUpdateMenu(c echo.Context) error {

	dprintf(4, c, "call SetStoreUpdateMenu\n")
	//
	params := cls.GetParamJsonMap(c)

	// 파라메터 맵으로 쿼리 변환

	updateCategoryQuery, err := cls.SetUpdateParam(storesql.UpdateStroeMenu, params)
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

func GetStoreServiceList(c echo.Context) error {

	dprintf(4, c, "call GetStoreServiceList\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectType(storesql.SelectStoreServiceList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultList == nil {

		// 서비스 한개도 없으면 넣어줌
		SetStoreBaseService(c)
		return c.JSON(http.StatusOK, "")
		//return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)

}

func SetStoreBaseService(c echo.Context) error {

	dprintf(4, c, "call SetStoreBaseService\n")
	//
	params := cls.GetParamJsonMap(c)

	params["useYn"] = "N"

	// 파라메터 맵으로 쿼리 변환
	insertBaseServiceQuery, err := cls.GetQueryJson(storesql.InsertStoreBaseService, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(insertBaseServiceQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	resultList, err := cls.GetSelectType(storesql.SelectStoreServiceList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
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

func SetStoreService(c echo.Context) error {

	dprintf(4, c, "call SetStoreService\n")
	//
	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
	// 상세 주문데이터 get
	var service []StoreService
	err2 := json.Unmarshal(bodyBytes, &service)
	if err2 != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err2.Error()))
	}

	c.Request().Body.Close() //  must close
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	//params := cls.GetParamJsonMap(c)

	params := make(map[string]string)
	params["storeId"] = c.Param("storeId")
	strUseYn := ""
	for i, _ := range service {

		params["serviceId"] = strconv.Itoa(service[i].ServiceId)
		params["serviceInfo"] = service[i].ServiceInfo

		if service[i].UseYn == true {
			strUseYn = "1"
		} else {
			strUseYn = "0"
		}
		params["useYn"] = strUseYn

		updateStoreService, err := cls.GetQueryJson(storesql.UpdateStoreService, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(updateStoreService)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

type StoreService struct {
	ServiceId   int    `json:"serviceId"`
	ServiceInfo string `json:"serviceInfo"`
	UseYn       bool   `json:"useYn"`
}

func GetStoreInfo(c echo.Context) error {

	dprintf(4, c, "call GetStoreInfo\n")

	params := cls.GetParamJsonMap(c)
	serviceList, err := cls.GetSelectType(storesql.SelectStoreService, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	storeInfo, err := cls.GetSelectType(storesql.SelectStoreDesc, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	store := make(map[string]interface{})
	store["storeInfo"] = storeInfo[0]
	store["serviceList"] = serviceList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = store

	return c.JSON(http.StatusOK, m)

}



func GetOrderPickupList(c echo.Context) error {

	dprintf(4, c, "call GetOrderPickupList\n")

	params := cls.GetParamJsonMap(c)



	totalCount, err := cls.GetSelectTypeRequire(storesql.SelectOrderPickupCount, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	pickOrder, err := cls.GetSelectData(storesql.SelectOrderPickupList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult(err.Error(), "DB fail"))
	}

	orderList := make([]map[string]interface{}, len(pickOrder))
	for i := range pickOrder {

		params["orderNo"]=pickOrder[i]["ORDER_NO"]
		menuList, err := cls.GetSelectType(storesql.SelectOrderPickupMenuList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		order2 := make(map[string]interface{})
		order2["orderNo"] = pickOrder[i]["ORDER_NO"]
		order2["userNm"] = pickOrder[i]["USER_NM"]
		order2["orderDate"] = pickOrder[i]["ORDER_DATE"]
		order2["orderTime"] = pickOrder[i]["ORDER_TIME"]
		order2["pStatus"] = pickOrder[i]["P_STATUS"]
		order2["pStatusNm"] = pickOrder[i]["P_STATUS_NM"]
		order2["totalQty"] = pickOrder[i]["TOTAL_ORDER_QTY"]
		order2["menuList"] = menuList
		orderList[i] = order2
	}


	pickUpData := make(map[string]interface{})
	pickUpData["totalCount"] = totalCount[0]
	pickUpData["orderList"] = orderList




	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = pickUpData

	return c.JSON(http.StatusOK, m)

}



// 포장주문 상태 업데이트
func SetOrderPickupStatus(c echo.Context) error {

	dprintf(4, c, "call SetOrderPickupStatus\n")

	params := cls.GetParamJsonMap(c)

	updateQuery :=""
	pushMsg := ""



	orderUserInfo, err := cls.GetSelectData(storesql.SelectPickupOrderUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if orderUserInfo ==nil{
		m := make(map[string]interface{})
		m["resultCode"] = "99"
		m["resultMsg"] = "주문정보가 없습니다."
		return c.JSON(http.StatusOK, m)
	}

	userId :=orderUserInfo[0]["USER_ID"]
	restNm :=orderUserInfo[0]["REST_NM"]



	if params["pStatus"]=="32"{
		pushMsg = restNm+ " - 주문하신 메뉴를 준비 중입니다."
		updateQuery=storesql.UpdatePickupStatus32
	}else if params["pStatus"]=="34"{
		pushMsg = restNm+ " - 메뉴가 모두 준비되었습니다.\n매장에서 메뉴를 픽업 해주세요."
		updateQuery=storesql.UpdatePickupStatus34
	}else if params["pStatus"]=="36"{
		updateQuery=storesql.UpdatePickupStatus36
	}



	UpdatePickupStatusQuery, err := cls.GetQueryJson(updateQuery, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(UpdatePickupStatusQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}


	//if len(os.Getenv("SERVER_TYPE")) > 0 {
		// 푸쉬 전송 시작
		if pushMsg !=""{
			go apiPush.SendPush_Msg_V1("주문", pushMsg, "M", "0", userId, "", "order")
		}
		// 푸쉬 전송 완료
	//}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	return c.JSON(http.StatusOK, m)

}
