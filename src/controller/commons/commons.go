package commons

import (
	"cashApi/src/controller"
	"encoding/json"
	"net/http"

	commonsql "cashApi/query/commons"
	pushs "cashApi/src/controller/API/PUSH"
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

// 카테고리
func GetCategoryList(c echo.Context) error {

	dprintf(4, c, "call GetCategoryList\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectType(commonsql.SelectCategoryList, params, c)
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

func GetCodeList(c echo.Context) error {

	dprintf(4, c, "call GetCodeList\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectType(commonsql.SelectCodeist, params, c)
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


type VerisonData struct {
	VersionCode     string `json:"versionCode"`     //
	IsRequireUpdate bool   `json:"isRequireUpdate"` //
}

//앱 최신 버전 호출
func GetVersionsLatest(c echo.Context) error {

	dprintf(4, c, "call GetVersionsLatest\n")

	params := cls.GetParamJsonMap(c)
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(commonsql.SelectVersion, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "query parameter fail"))
	}
	// 쿼리 실행 후 JSON 형태로 결과 받기
	dbData, err := cls.QueryJsonColumn(selectQuery, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "DB fail"))
	}
	var resultData VerisonData
	err = json.Unmarshal(dbData, &resultData)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "Unmarshal fail"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)

}



//신규 fcm v1
func SendCommonPush(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	title := params["title"]
	content := params["content"]
	pushType := params["pushType"]
	menu := params["menu"]
	userType := params["userType"]

	//유저 타입 과 파라미터
	//0 : 사용자 	   userId
	//1 : 가맹점(사장님) restId
	//2 : 가맹점(사장님) bizNum
	//3 : 사용자(장부장) grpId
	pushQuery :=""
	if userType =="0" {
		pushQuery =commonsql.SelectPushUser
	}else if userType =="1"{
		pushQuery =commonsql.SelectPushRest
	}else if userType =="2"{
		pushQuery =commonsql.SelectPushBizNum
	}else if userType =="3"{
		pushQuery =commonsql.SelectPushGrp
	}else{
		dprintf(1, c, "call SendCommonPush : 잘못된 유저 타입입니다. \n")
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "잘못된 유저 타입입니다."))
	}


	resultData, err := cls.GetSelectData(pushQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		dprintf(1, c, "call SendCommonPush : 데이터가 없습니다. \n")
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "데이터가 없습니다."))
	}


	rtoken := resultData[0]["REG_ID"]
	osTy := resultData[0]["OS_TY"]


	//sendToToken(title,content,rtoken)
	//sendToToken(title,content,rtoken)
	//pushs.SendTotal(title,content,rtoken,pushType,menu,osTy)

	pushs.SendTotal(title,content,rtoken,pushType,menu,"",osTy)


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"


	return c.JSON(http.StatusOK, m)
}


func SendCommonPushMsg(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	menu := params["menu"]
	pushType := params["pushType"]
	userType := params["userType"]


	msgInfo, err := cls.GetSelectData(commonsql.SelectPushMsgInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if msgInfo == nil {
		dprintf(1, c, "call SendCommonPush : 데이터가 없습니다. \n")
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "데이터가 없습니다."))
	}



	title := msgInfo[0]["TITLE"]
	content := msgInfo[0]["MSG"]
	//유저 타입 과 파라미터
	//0 : 사용자 	   userId
	//1 : 가맹점(사장님) restId
	//2 : 가맹점(사장님) bizNum
	//3 : 사용자(장부장) grpId
	pushQuery :=""
	if userType =="0" {
		pushQuery =commonsql.SelectPushUser
	}else if userType =="1"{
		pushQuery =commonsql.SelectPushRest
	}else if userType =="2"{
		pushQuery =commonsql.SelectPushBizNum
	}else if userType =="3"{
		pushQuery =commonsql.SelectPushGrp
	}else{
		dprintf(1, c, "call SendCommonPush : 잘못된 유저 타입입니다. \n")
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "잘못된 유저 타입입니다."))
	}


	resultData, err := cls.GetSelectData(pushQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		dprintf(1, c, "call SendCommonPush : 데이터가 없습니다. \n")
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "데이터가 없습니다."))
	}


	rtoken := resultData[0]["REG_ID"]
	osTy := resultData[0]["OS_TY"]


	//sendToToken(title,content,rtoken)
	//pushs.SendTotal(title,content,rtoken,pushType,menu,osTy)
	pushs.SendTotal(title,content,rtoken,pushType,menu,"",osTy)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"


	return c.JSON(http.StatusOK, m)
}


