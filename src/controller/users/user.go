package users

import (
	"bytes"
	commonsql "cashApi/query/commons"
	storesql "cashApi/query/stores"
	usersql "cashApi/query/users"
	"cashApi/src/controller"
	kakao "cashApi/src/controller/API/KAKAO"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"cashApi/src/controller/cls"

	"github.com/labstack/echo/v4"
)

/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
// darayo printf(로그레벨, 요청 컨텍스트, format, arg) => 무엇을(서비스, 요청), 어떻게(input), 왜(원인,조치)
var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func BaseUrl(c echo.Context) error {

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func LoginOut(c echo.Context) error {
	params := cls.GetParamJsonMap(c)

	userId := params["userId"]
	params["regId"] = " "
	params["loginYn"] = "N"

	selectQuery, err := cls.SetUpdateParam(usersql.UpdatePushData, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	cls.ClearLoginSession(c, userId)

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

// 로그인
func LoginDarayo(c echo.Context) error {

	dprintf(4, c, "call LoginDarayo\n")
	ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectDataRequire(usersql.SelectUserLoginCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		// 접속 로그
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = "id"
		LoginAcceesLog(c, params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "아이디 또는 비밀번호가 잘못되었습니다."))
	}

	useYn := resultData[0]["USE_YN"]
	userId := resultData[0]["userId"]

	loginId := resultData[0]["LOGIN_ID"]
	userNm := resultData[0]["USER_NM"]
	userTel := resultData[0]["HP_NO"]
	userBirth := resultData[0]["USER_BIRTH"]

	params["userId"] = userId
	if useYn == "N" {
		// 접속 로그
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = "id"
		LoginAcceesLog(c, params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "로그인 실패."))
	} else if useYn == "S" {

		checkData, err := cls.GetSelectData(usersql.SelectJoinCheck, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		c = cls.SetLoginJWT(c, userId)

		joinStep := "joinStep2"

		bizNum := ""
		idCheck := checkId(checkData[0]["LOGIN_ID"])
		if idCheck == "bizNum" {
			bizNum = checkData[0]["LOGIN_ID"]
		}

		chkData := make(map[string]interface{})
		chkData["userId"] = userId
		chkData["joinStep"] = joinStep
		chkData["userNm"] = checkData[0]["USER_NM"]
		chkData["hpNo"] = checkData[0]["HP_NO"]
		chkData["loginId"] = checkData[0]["LOGIN_ID"]
		chkData["bizNum"] = bizNum

		rm := make(map[string]interface{})
		rm["resultCode"] = "01"
		rm["resultMsg"] = "응답 성공"
		rm["resultData"] = chkData

		return c.JSON(http.StatusOK, rm)
	}

	userTy := resultData[0]["USER_TY"]
	if userTy != "1" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "가맹점 사용자가 아닙니다."))
	}

	// 푸쉬 아이디 등록
	regId := resultData[0]["regId"]
	PushQuery := ""

	params["loginYn"] = "Y"
	if regId == "null" {
		// insert
		PushQuery = usersql.InserPushData
	} else {
		// update
		PushQuery = usersql.UpdatePushData
	}
	InsertPushQuery, err := cls.GetQueryJson(PushQuery, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "push reg query parameter fail"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(InsertPushQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail(push reg)"))
	}

	storeInfo, err := cls.GetSelectDataRequire(usersql.SelectUserStoreInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	//토큰 발행

	c = cls.SetLoginJWT(c, userId)

	// 접속 로그
	params["logInOut"] = "I"
	params["ip"] = ip
	params["succYn"] = "Y"
	params["type"] = "id"
	LoginAcceesLog(c, params)

	userData := make(map[string]interface{})
	userData["userId"] = userId
	userData["loginId"] = loginId
	userData["userNm"] = userNm
	userData["userTel"] = userTel
	userData["userBirth"] = userBirth
	userData["storeId"] = storeInfo[0]["STORE_ID"]
	userData["storeNm"] = storeInfo[0]["STORE_NM"]
	userData["bizNum"] = storeInfo[0]["BIZ_NUM"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

// 푸쉬 정보 업데이트
func PushInfo(c echo.Context) error {

	dprintf(4, c, "call PushInfoUpdate \n")
	params := cls.GetParamJsonMap(c)
	params["loginYn"] = "Y"

	InsertPushQuery, err := cls.GetQueryJson(usersql.UpdatePushData, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "push reg query parameter fail"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(InsertPushQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail(push reg)"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 소셜 로그인
func LoginDarayoSocial(c echo.Context) error {

	dprintf(4, c, "call GetUserInfo\n")
	ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	socialType := c.FormValue("socialType")

	LoginSqlQuery := ""
	if socialType == "kakao" {
		LoginSqlQuery = usersql.SelectKakaoUserLoginCheck
	} else if socialType == "naver" {
		LoginSqlQuery = usersql.SelectNaverUserLoginCheck
	} else if socialType == "apple" {
		LoginSqlQuery = usersql.SelectAppleUserLoginCheck
	}

	if LoginSqlQuery == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "지원하지 않는 로그인 방식 입니다."))
	}
	params := cls.GetParamJsonMap(c)

	socialToken := params["socialToken"]
	resultData, err := cls.GetSelectDataRequire(LoginSqlQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		// 접속 로그
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = socialType
		params["loginId"] = socialToken
		LoginAcceesLog(c, params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "아이디 또는 비밀번호가 잘못되었습니다."))
	}

	useYn := resultData[0]["USE_YN"]
	userId := resultData[0]["userId"]
	userBirth := resultData[0]["USER_BIRTH"]

	params["userId"] = userId
	if useYn == "N" {

		// 접속 로그
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = socialType
		params["loginId"] = socialToken
		LoginAcceesLog(c, params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "로그인 실패."))
	} else if useYn == "S" {

		checkData, err := cls.GetSelectData(usersql.SelectJoinCheck, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		joinStep := "joinStep2"

		chkData := make(map[string]interface{})
		chkData["userId"] = userId
		chkData["joinStep"] = joinStep
		chkData["userNm"] = checkData[0]["USER_NM"]
		chkData["hpNo"] = checkData[0]["HP_NO"]
		chkData["loginId"] = checkData[0]["LOGIN_ID"]

		rm := make(map[string]interface{})
		rm["resultCode"] = "01"
		rm["resultMsg"] = "응답 성공"
		rm["resultData"] = chkData

		return c.JSON(http.StatusOK, rm)
	}

	userTy := resultData[0]["USER_TY"]
	if userTy != "1" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "가맹점 사용자가 아닙니다."))
	}

	// 푸쉬 아이디 등록
	regId := resultData[0]["regId"]
	params["loginYn"] = "Y"
	PushQuery := ""
	if regId == "null" {
		// insert
		PushQuery = usersql.InserPushData
	} else {
		// update
		PushQuery = usersql.UpdatePushData
	}
	InsertPushQuery, err := cls.GetQueryJson(PushQuery, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "push reg query parameter fail"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(InsertPushQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail(push reg)"))
	}

	storeInfo, err := cls.GetSelectDataRequire(usersql.SelectUserStoreInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	c = cls.SetLoginJWT(c, userId)

	// 접속 로그
	//	dprintf(4, c, "call GetUserInfo\n", resultData)
	loginId := resultData[0]["LOGIN_ID"]
	userNm := resultData[0]["USER_NM"]
	userTel := resultData[0]["HP_NO"]

	params["loginId"] = loginId
	params["logInOut"] = "I"
	params["ip"] = ip
	params["succYn"] = "Y"
	params["type"] = socialType
	params["loginId"] = socialToken
	LoginAcceesLog(c, params)

	userData := make(map[string]interface{})
	userData["userId"] = userId
	userData["loginId"] = loginId
	userData["userNm"] = userNm
	userData["userTel"] = userTel
	userData["userBirth"] = userBirth
	userData["storeId"] = storeInfo[0]["STORE_ID"]
	userData["storeNm"] = storeInfo[0]["STORE_NM"]
	userData["bizNum"] = storeInfo[0]["BIZ_NUM"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

// 자동 로그인
func LoginDarayoAuto(c echo.Context) error {

	dprintf(4, c, "call LoginDarayoAuto\n")
	ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectDataRequire(usersql.SelectUserLoginCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "로그인 실패"))
	}
	userId := resultData[0]["userId"]
	loginId := resultData[0]["LOGIN_ID"]
	userNm := resultData[0]["USER_NM"]
	userTel := resultData[0]["HP_NO"]

	c = cls.SetLoginJWT(c, userId)

	// 접속 로그
	params["logInOut"] = "I"
	params["ip"] = ip
	params["succYn"] = "Y"
	params["type"] = "id"
	LoginAcceesLog(c, params)

	//	dprintf(4, c, "call GetUserInfo\n", resultData)

	userData := make(map[string]interface{})
	userData["userId"] = userId
	userData["loginId"] = loginId
	userData["userNm"] = userNm
	userData["userTel"] = userTel

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

//회원정보 조회
func GetUserInfo(c echo.Context) error {

	dprintf(4, c, "call GetUserInfo\n")
	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectData(usersql.SelectUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData[0]

	return c.JSON(http.StatusOK, m)

}

// 설정 - 메인 업데이트
func SetSetupInfo(c echo.Context) error {

	dprintf(4, c, "call SetSetupInfo\n")

	params := cls.GetParamJsonMap(c)
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(usersql.UpdateSetupInfo, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "query parameter fail"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

//설정 - 내정보 변경
func SetUserInfo(c echo.Context) error {

	dprintf(4, c, "call SetUserInfo\n")

	params := cls.GetParamJsonMap(c)
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.SetUpdateParam(usersql.UpdateUserInfo, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "query parameter fail"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

func checkId(str string) string {

	var validEmail, _ = regexp.Compile(
		"^[_a-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$",
	)

	result := "bizNum"
	chkResultEmail := validEmail.MatchString(str)

	if chkResultEmail == true {
		result = "email"
	}

	return result
}

// 이메일 중복 체크 (아이디 체크)
func GetEmailDupCheck(c echo.Context) error {

	dprintf(4, c, "call GetEmailDupCheck\n")
	params := cls.GetParamJsonMap(c)

	var validEmail, _ = regexp.Compile(
		"^[_a-z0-9+-.]+@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$",
	)
	email := params["email"]

	println(email)
	chkResult := validEmail.MatchString(email)

	if chkResult == false {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "잘못된 email 입니다."))
	}

	resultData, err := cls.GetSelectData(usersql.SelectEmailDupCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	emailCnt := resultData[0]["emailCnt"]

	if emailCnt != "0" {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "Email 중복"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

type ResultMsgMap struct {
	XMLName   xml.Name `xml:"map"`
	Id        string   `xml:"id,attr"`
	DetailMsg string   `xml:"detailMsg"`
	Msg       string   `xml:"msg"`
	Code      string   `xml:"code"`
	Result    string   `xml:"result"`
}

type Map struct {
	XMLName             xml.Name     `xml:"map"`
	ResultMsgMap        ResultMsgMap `xml:"map"`
	TrtEndCd            string       `xml:"trtEndCd"`
	SmpcBmanEnglTrtCntn string       `xml:"smpcBmanEnglTrtCntn"`
	NrgtTxprYn          string       `xml:"nrgtTxprYn"`
	SmpcBmanTrtCntn     string       `xml:"smpcBmanTrtCntn"`
	TrtCntn             string       `xml:"trtCntn"`
}

func GetBizNumDupCheck(c echo.Context) error {

	dprintf(4, c, "call GetBizNumDupCheck\n")

	params := cls.GetParamJsonMap(c)

	bizNum := params["bizNum"]

	println("11")

	url := "https://teht.hometax.go.kr/wqAction.do?actionId=ATTABZAA001R08&screenId=UTEABAAA13&popupYn=false&realScreenId="
	xmlData := "<map id='ATTABZAA001R08'>" +
		"<pubcUserNo/>" +
		"<mobYn>N</mobYn>" +
		"<inqrTrgtClCd>1</inqrTrgtClCd>" +
		"<txprDscmNo>" + bizNum + "</txprDscmNo>" +
		"<dongCode>05</dongCode" +
		"><psbSearch>Y</psbSearch>" +
		"<map id='userReqInfoVO'/>" +
		"</map>"
	buf := bytes.NewBufferString(xmlData)
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

	if strings.Contains(strings.Replace(checkResult, " ", "", -1), "등록되어있는사업자등록번호입니다") == false {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "사용할수 없는 사업자 번호 입니다."))
	}

	idChk, err := cls.GetSelectData(usersql.SelectLoginIdDupCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	idCnt, _ := strconv.Atoi(idChk[0]["loginIdCnt"])
	if idCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "이미 가입된 사업자 번호입니다."))
	}

	bizNumCnt, err := cls.GetSelectData(usersql.SelectBizNumDupCheck, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	bizCnt, _ := strconv.Atoi(bizNumCnt[0]["bizCnt"])

	if bizCnt > 0 {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "이미 가입된 사업자 번호입니다."))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	return c.JSON(http.StatusOK, m)
}

//소셜 로그인 토큰 중복체크
func GetSocialTokenDupCheck(c echo.Context) error {

	dprintf(4, c, "call GetSocialTokenDupCheck\n")

	socialType := c.FormValue("socialType")

	SqlQuery := ""
	if socialType == "kakao" {
		SqlQuery = usersql.SelectKakaoTokenDupCheck
	} else if socialType == "naver" {
		SqlQuery = usersql.SelectNaverTokenDupCheck
	} else if socialType == "apple" {
		SqlQuery = usersql.SelectAppleTokenDupCheck
	}

	if SqlQuery == "" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "지원하지 않는 로그인 방식 입니다."))
	}

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectData(SqlQuery, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	emailCnt := resultData[0]["tokenCnt"]

	if emailCnt != "0" {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "토큰 중복"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 사용자 패스워드 재설정
func SetPwdResetChange(c echo.Context) error {

	dprintf(4, c, "call SetNotifications\n")

	params := cls.GetParamJsonMap(c)
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(usersql.UpdateUserPasswordChange, params)
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

// 회원 가입 1단계 - 사업자 번호로 가입
func SetUserBizNumJoin(c echo.Context) error {

	dprintf(4, c, "call SetUserBizNumJoin\n")

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(usersql.SelectCreatUserSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "유저아이디 생성 실패(2)"))
	}
	userId := resultData[0]["newUserId"]
	params["userId"] = userId
	// 유저 가입  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -유저 가입(setUserJoin)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	socialType := params["socialType"]
	socialToken := params["socialToken"]
	loginPw := params["loginPw"]
	loginId := params["loginId"]

	if socialType == "kakao" {
		params["kakaoKey"] = socialToken
		params["kakaoPw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	} else if socialType == "naver" {
		params["naverKey"] = socialToken
		params["naverPw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	} else if socialType == "apple" {
		params["appleKey"] = socialToken
		params["applePw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	}

	termsOfBenefit := params["termsOfBenefit"]
	pushYn := "N"
	if termsOfBenefit == "Y" {
		pushYn = "Y"
	}

	// 유저 생성
	params["userTy"] = "1"
	params["atLoginYn"] = "Y"
	params["pushYn"] = pushYn
	UserCreateQuery, err := cls.SetUpdateParam(usersql.InserCreateUser, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UserCreateQuery parameter fail"))
	}

	_, err = tx.Exec(UserCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", UserCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UserCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	TermsCreateQuery, err := cls.SetUpdateParam(usersql.InsertTermsUser, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "TermsCreateQuery parameter fail"))
	}

	_, err = tx.Exec(TermsCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", TermsCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", TermsCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	c = cls.SetLoginJWT(c, userId)

	userData := make(map[string]interface{})
	userData["userId"] = userId
	userData["bizNum"] = loginId

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

// 회원 가입 1단계
func SetUserJoin(c echo.Context) error {

	dprintf(4, c, "call setUserJoin\n")

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(usersql.SelectCreatUserSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "유저아이디 생성 실패(2)"))
	}
	userId := resultData[0]["newUserId"]
	params["userId"] = userId
	// 유저 가입  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -유저 가입(setUserJoin)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환

	socialType := params["socialType"]
	socialToken := params["socialToken"]
	loginPw := params["loginPw"]
	params["loginId"] = params["email"]

	if socialType == "kakao" {
		params["kakaoKey"] = socialToken
		params["kakaoPw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	} else if socialType == "naver" {
		params["naverKey"] = socialToken
		params["naverPw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"
	} else if socialType == "apple" {
		params["appleKey"] = socialToken
		params["applePw"] = loginPw
		params["loginPw"] = "bcb15f821479b4d5772bd0ca866c00ad5f926e3580720659cc80d39c9d09802a"

		if params["userNm"] == "" {
			params["userNm"] = strings.Replace(userId, "U", "A", -1)
		}
		if params["loginId"] == "" {
			params["loginId"] = strings.Replace(userId, "U", "A", -1)
		}
	}

	termsOfBenefit := params["termsOfBenefit"]
	pushYn := "N"
	if termsOfBenefit == "Y" {
		pushYn = "Y"
	}

	// 유저 생성
	params["userTy"] = "1"
	params["atLoginYn"] = "Y"
	params["pushYn"] = pushYn

	UserCreateQuery, err := cls.SetUpdateParam(usersql.InserCreateUser, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UserCreateQuery parameter fail"))
	}

	_, err = tx.Exec(UserCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", UserCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UserCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	TermsCreateQuery, err := cls.SetUpdateParam(usersql.InsertTermsUser, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "TermsCreateQuery parameter fail"))
	}

	_, err = tx.Exec(TermsCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", TermsCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", TermsCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 유저 가입 TRNAN 종료

	c = cls.SetLoginJWT(c, userId)

	userData := make(map[string]interface{})
	userData["userId"] = userId

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

//회원 가입 - 가맹점 정보 입력
func SetUserJoinStep2(c echo.Context) error {

	dprintf(4, c, "call SetUserJoinStep2\n")

	params := cls.GetParamJsonMap(c)

	//가맹점 아이디 생성
	storeSeqData, err := cls.GetSelectData(storesql.SelectStoreSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	storeId := storeSeqData[0]["storeSeq"]
	params["storeId"] = storeId

	UserInfo, err := cls.GetSelectData(usersql.SelectUserInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	storeTel := UserInfo[0]["userTel"]
	birthday := UserInfo[0]["birthday"]
	email := UserInfo[0]["email"]
	userName := UserInfo[0]["userName"]
	//recomCode := UserInfo[0]["recomCode"]

	params["storeTel"] = storeTel
	params["ceoBirthday"] = birthday
	params["storeEmail"] = email
	params["ceoName"] = userName
	params["ceoTti"] = GetBrithdayTti(birthday)

	//  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -가맹점 생성 (SetUserJoinStep2)  \n")
			tx.Rollback()
		}
	}()

	insertStore, err := cls.GetQueryJson(storesql.InsertStore, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = tx.Exec(insertStore)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 매장에 사장님 등록
	InsertStoreUserQuery, err := cls.GetQueryJson(storesql.InsertStoreUser, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertStoreUserQuery parameter fail"))
	}
	_, err = tx.Exec(InsertStoreUserQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertStoreUserQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertStoreUserQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 가맹점 완료 처리
	params["useYn"] = "Y"
	joinFinishUpdate, err := cls.SetUpdateParam(usersql.UpdateUserInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = tx.Exec(joinFinishUpdate)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 수수료 데이터 입력 -CARD
	params["payMethod"] = "CARD"
	params["restFees"] = "3"
	restFeesInsert1, err := cls.SetUpdateParam(storesql.InsertStoreFees, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = tx.Exec(restFeesInsert1)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// 수수료 데이터 입력 -Bank
	params["payMethod"] = "BANK"
	params["restFees"] = "3"
	restFeesInsert2, err := cls.SetUpdateParam(storesql.InsertStoreFees, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = tx.Exec(restFeesInsert2)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	c = cls.SetLoginJWT(c, params["userId"])
	//useMonth :="1"
	//if strings.ToUpper(recomCode) == "FIT6M"  {
	//	useMonth ="6"
	//}
	//apis.BillingInsert(storeId,params["userId"],"I0000000001",useMonth)

	userData := make(map[string]interface{})
	userData["storeId"] = storeId
	userData["userId"] = params["userId"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

//가맹점 정보  - 대표 && 계좌정보 업데이트
func SetStorInfoUpdate(c echo.Context) error {

	dprintf(4, c, "call SetStorInfoUpdate\n")

	params := cls.GetParamJsonMap(c)

	//  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -가맹점 생성 (SetUserJoinStep2)  \n")
			tx.Rollback()
		}
	}()

	params["ceoTti"] = GetBrithdayTti(params["ceoBirthday"])

	// 파라메터 맵으로 쿼리 변환
	UpdateStoreInfoQuery, err := cls.SetUpdateParam(storesql.UpdateStoreInfo, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = tx.Exec(UpdateStoreInfoQuery)
	if err != nil {
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

// 설정 - 매장 정보 관리 업데이트
func SetStoreInfo(c echo.Context) error {

	dprintf(4, c, "call SetStoreInfo\n")

	params := cls.GetParamJsonMap(c)
	params["ceoTti"] = GetBrithdayTti(params["ceoBirthday"])
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.SetUpdateParam(storesql.UpdateStore, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	selectStoreCCompChk, err := cls.GetSelectData(storesql.SelectCompInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	if selectStoreCCompChk != nil {
		UpdateStoreCCompQuery, err := cls.SetUpdateParam(storesql.UpdateStoreCComp, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(UpdateStoreCCompQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 설정 - 대표자 정보 관리 업데이트
func SetCeoInfo(c echo.Context) error {

	dprintf(4, c, "call SetCeoInfo\n")

	params := cls.GetParamJsonMap(c)

	email := params["email"]

	if email == "" {
		params["email"] = " "
	}
	// 파라메터 맵으로 쿼리 변환
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

// 설정 - 매장 정보 관리
func GetStoreInfo(c echo.Context) error {

	dprintf(4, c, "call GetStoreInfo\n")

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectType(storesql.SelectStoreInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData[0]

	return c.JSON(http.StatusOK, m)

}

// 가입대행 서비스 정보
func GetStoreServiceInfo(c echo.Context) error {

	dprintf(4, c, "call GetStoreServiceInfo\n")

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectType(storesql.SelectStoreServiceInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData[0]

	return c.JSON(http.StatusOK, m)

}

// 설정 - 메인
func GetSetupInfo(c echo.Context) error {

	dprintf(4, c, "call GetSetupInfo\n")

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectType(usersql.SelectUserSetupInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "no data"))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData[0]

	return c.JSON(http.StatusOK, m)

}

// 설정 - 충전 정보
func GetStoreChargeList(c echo.Context) error {

	dprintf(4, c, "call GetStoreChargeList\n")

	params := cls.GetParamJsonMap(c)

	paymentInfo, err := cls.GetSelectData(storesql.SelectStorePaymentInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	paymentUseYn := paymentInfo[0]["paymentUseYn"]

	resultData := make(map[string]interface{})
	resultData["paymentUseYn"] = paymentUseYn

	resultList, err := cls.GetSelectType(storesql.SelectStoreChargeList, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultList == nil {

		baseList, err := cls.GetSelectData(storesql.SelectChargeBase, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		tx, err := cls.DBc.Begin()
		if err != nil {
			//return "5100", errors.New("begin error")
		}

		// 오류 처리
		defer func() {
			if err != nil {
				// transaction rollback
				dprintf(4, c, "do rollback -설정 - 충전 정보(GetStoreChargeList)  \n")
				tx.Rollback()
			}
		}()

		for i := range baseList {

			Bparam := make(map[string]string)

			Bparam["storeSeqNo"] = baseList[i]["SEQ_NO"] + "_" + params["storeId"]
			Bparam["amt"] = baseList[i]["AMT"]
			Bparam["storeId"] = params["storeId"]

			InsertStoreChargeListQuery, err := cls.SetUpdateParam(storesql.InsertStoreChargeList, Bparam)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "InsertStoreChargeListQuery parameter fail"))
			}

			_, err = tx.Exec(InsertStoreChargeListQuery)
			dprintf(4, c, "call set Query (%s)\n", InsertStoreChargeListQuery)
			if err != nil {
				dprintf(1, c, "Query(%s) -> error (%s) \n", InsertStoreChargeListQuery, err)
				return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
			}

		}

		// transaction commit
		err = tx.Commit()
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		resultList, err := cls.GetSelectType(storesql.SelectStoreChargeList, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
		}

		resultData["chargeList"] = resultList

	}

	resultData["chargeList"] = resultList

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = resultData

	return c.JSON(http.StatusOK, m)

}

type StoreCharge struct {
	SeqNo  string `json:"seqNo"`
	AddAmt int    `json:"addAmt"`
	UseYn  string `json:"useYn"`
}

// 설정 - 선불 충전 사용 여부
func SetStoreChargeYn(c echo.Context) error {

	dprintf(4, c, "call SetStoreChargeYn\n")

	params := cls.GetParamJsonMap(c)

	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -선불 충전 사용 여부(SetStoreChargeYn)  \n")
			tx.Rollback()
		}
	}()

	UpdateStoreChargeYnQuery, err := cls.GetQueryJson(storesql.UpdateStoreChargeYn, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "UpdateStoreChargeYnQuery parameter fail"))
	}

	_, err = tx.Exec(UpdateStoreChargeYnQuery)
	dprintf(4, c, "call set Query (%s)\n", UpdateStoreChargeYnQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UpdateStoreChargeYnQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	data := make(map[string]string)
	data["paymentUseYn"] = params["paymentUseYn"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)

}

// 설정 - 충전 정보
func SetStoreCharge(c echo.Context) error {

	dprintf(4, c, "call SetStoreCharge\n")
	//
	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
	// 상세 주문데이터 get
	var charge []StoreCharge
	err2 := json.Unmarshal(bodyBytes, &charge)
	if err2 != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err2.Error()))
	}

	c.Request().Body.Close() //  must close
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	//params := cls.GetParamJsonMap(c)

	params := make(map[string]string)
	params["storeId"] = c.Param("storeId")
	for i, _ := range charge {

		params["seqNo"] = charge[i].SeqNo
		params["addAmt"] = strconv.Itoa(charge[i].AddAmt)
		params["useYn"] = charge[i].UseYn

		UpdateStoreChargeQuery, err := cls.GetQueryJson(storesql.UpdateStoreCharge, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
		// 쿼리 실행
		_, err = cls.QueryDB(UpdateStoreChargeQuery)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}

// 홈텍스 or 여신 가입하기
func SetCashJoin(c echo.Context) error {
	dprintf(4, c, "call SetCashJoin\n")

	params := cls.GetParamJsonMap(c)
	storeInfo, err := cls.GetSelectData(storesql.SelectCash, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}

	var query string
	if storeInfo == nil {
		// Insert
		query, err = cls.SetUpdateParam(storesql.InsertCash, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	} else {
		// Update
		query, err = cls.SetUpdateParam(storesql.UpdateCash, params)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}
	}

	// 쿼리 실행
	_, err = cls.QueryDB(query)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	data := make(map[string]interface{})
	if storeInfo == nil {
		if params["lnJoinStsCd"] == "1" {
			data["lnJoinYn"] = "Y"
		} else {
			data["lnJoinYn"] = "N"
		}
		if params["hometaxJoinStsCd"] == "1" {
			data["hometaxJoinYn"] = "Y"
		} else {
			data["hometaxJoinYn"] = "N"
		}
	} else {
		if params["lnJoinStsCd"] == "1" || storeInfo[0]["lnJoinStsCd"] == "1" {
			data["lnJoinYn"] = "Y"
		} else {
			data["lnJoinYn"] = "N"
		}
		if params["hometaxJoinStsCd"] == "1" || storeInfo[0]["hometaxJoinStsCd"] == "1" {
			data["hometaxJoinYn"] = "Y"
		} else {
			data["hometaxJoinYn"] = "N"
		}
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = data

	return c.JSON(http.StatusOK, m)
}

// 홈텍스 or 여신 정보 수정하기
func SetCashModify(c echo.Context) error {
	dprintf(4, c, "call SetCashJoin\n")

	params := cls.GetParamJsonMap(c)
	updateQuery, err := cls.SetUpdateParam(storesql.UpdateCash, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(updateQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)
}

func LoginAcceesLog(c echo.Context, params map[string]string) {

	dprintf(4, c, "call LoginAcceesLog\n")
	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(usersql.InsertLoginAccess, params)
	if err != nil {
		dprintf(4, c, "LoginAcceesLog query parameter fail\n")
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		dprintf(4, c, "LoginAcceesLog DB fail\n")
	}

}

// 캐쉬 인증화면 정보
func GetCashInfo(c echo.Context) error {

	dprintf(4, c, "call GetCashInfo\n")

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectData(usersql.SelectStoreCashInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("01", "정보가 없습니다."))
	}
	homeTaxErrMsg := ""
	lnErrMsg := ""

	lnAuthFail, _ := strconv.Atoi(resultData[0]["LN_AUTH_FAIL"])
	hometaxAuthFail, _ := strconv.Atoi(resultData[0]["HOMETAX_AUTH_FAIL"])

	if lnAuthFail > 3 {
		lnErrMsg = "4회 이상 오류시 30분 이후 인증이 가능합니다.(현재 " + strconv.Itoa(lnAuthFail) + "회 오류)"
	}
	if hometaxAuthFail > 0 {
		homeTaxErrMsg = "3회 이상 오류시 홈텍스 페이지에서 비밀번호 재설정이 필요 할 수 있습니다.(현재 " + strconv.Itoa(hometaxAuthFail) + "회 오류)"
	}

	lnFailDt, _ := strconv.Atoi(resultData[0]["LN_FAIL_DT"])
	homtaxtFailDt, _ := strconv.Atoi(resultData[0]["HOMTAXT_FAIL_DT"])

	lnRemainTime, _ := strconv.Atoi(resultData[0]["LN_REMAIN_TIME"])
	homtaxtRemainTime, _ := strconv.Atoi(resultData[0]["HOMTAXT_REMAIN_TIME"])

	result := make(map[string]interface{})
	result["chkTime"] = 30

	result["hometaxAuthFail"] = hometaxAuthFail
	result["homeTaxErrMsg"] = homeTaxErrMsg
	result["homtaxtFailDt"] = strconv.Itoa(homtaxtFailDt)
	result["hometaxId"] = resultData[0]["HOMETAX_ID"]
	result["hometaxPsw"] = resultData[0]["HOMETAX_PSW"]
	result["hometaxJoinStsCd"] = resultData[0]["HOMETAX_JOIN_STS_CD"]
	result["homtaxtRemainTime"] = homtaxtRemainTime

	result["lnAuthFail"] = lnAuthFail
	result["lnErrMsg"] = lnErrMsg
	result["lnFailDt"] = strconv.Itoa(lnFailDt)
	result["lnId"] = resultData[0]["LN_ID"]
	result["lnPsw"] = resultData[0]["LN_PSW"]
	result["lnJoinStsCd"] = resultData[0]["LN_JOIN_STS_CD"]
	result["lnRemainTime"] = lnRemainTime

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)

}

func GetBrithdayTti(birthday string) string {

	if len(birthday) < 4 {
		return ""
	}

	year := birthday[:4]
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		lprintf(1, "[ERROR] strconv atoi(%s)\n", err.Error())
		return ""
	}

	yearGap := 1991 - yearInt
	var yearSwitch int

	if yearGap >= 0 {
		if yearGap <= 12 {
			yearSwitch = yearGap
		} else {
			yearSwitch = yearGap % 12
		}

		switch yearSwitch {
		case 0:
			return "07"
		case 1:
			return "06"
		case 2:
			return "05"
		case 3:
			return "04"
		case 4:
			return "03"
		case 5:
			return "02"
		case 6:
			return "01"
		case 7:
			return "00"
		case 8:
			return "11"
		case 9:
			return "10"
		case 10:
			return "09"
		case 11:
			return "08"
		case 12:
			return "07"
		}
	} else {
		yearGap = yearGap * (-1)

		if yearGap <= 12 {
			yearSwitch = yearGap
		} else {
			yearSwitch = yearGap % 12
		}

		switch yearSwitch {
		case 0:
			return "07"
		case 1:
			return "08"
		case 2:
			return "09"
		case 3:
			return "10"
		case 4:
			return "11"
		case 5:
			return "00"
		case 6:
			return "01"
		case 7:
			return "02"
		case 8:
			return "03"
		case 9:
			return "04"
		case 10:
			return "05"
		case 11:
			return "06"
		case 12:
			return "07"
		}
	}

	return ""

}

func CardsalesLoginCheck(c echo.Context) error {
	dprintf(4, c, "call CardsalesLoginCheck\n")
	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	//체크
	cashInfo, err := cls.GetSelectData(usersql.SelectStoreCashInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if cashInfo != nil {
		lnAuthFail, _ := strconv.Atoi(cashInfo[0]["LN_AUTH_FAIL"])
		lnRemainTime, _ := strconv.Atoi(cashInfo[0]["LN_REMAIN_TIME"])

		if lnAuthFail > 3 {
			if lnRemainTime > 0 {
				resultData := make(map[string]interface{})
				resultData["lnAuthFail"] = lnAuthFail
				lnFailDt, _ := strconv.Atoi(cashInfo[0]["LN_FAIL_DT"])
				resultData["lnFailDt"] = lnFailDt
				resultData["lnRemainTime"] = lnRemainTime
				resultData["lnErrMsg"] = "4회 이상 오류시 30분 이후 인증이 가능합니다. (현재 " + strconv.Itoa(lnAuthFail) + "회 오류)"

				m["resultCode"] = "99"
				m["resultMsg"] = "login Fail"
				m["resultData"] = resultData
				return c.JSON(http.StatusOK, m)
			}
		}

	}

	//체크 끝

	apiUrl := "https://www.cardsales.or.kr"
	resource := "/authentication"
	data := url.Values{}

	data.Set("j_username", params["loginId"])
	data.Set("j_password", params["password"])
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

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

		println(err)
		m["resultCode"] = "99"
		m["resultMsg"] = err.Error()
		return c.JSON(http.StatusOK, m)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		m["resultCode"] = "00"
		m["resultMsg"] = "login Sucess"
		//m["resp"] = resp
		AuthSuccessUpdate(c, params, "ln")
	} else {
		m["resultCode"] = "99"
		m["resultMsg"] = "login Fail"
		lnAuthFailCnt, lnFailDt := AuthFailUpdate(c, params, "ln")

		resultData := make(map[string]interface{})
		resultData["lnAuthFailCnt"] = lnAuthFailCnt
		resultData["lnFailDt"] = lnFailDt
		resultData["lnErrMsg"] = "4회 이상 오류시 30분 이후 인증이 가능합니다. (현재 " + strconv.Itoa(lnAuthFailCnt) + "회 오류)"

		m["resultData"] = resultData
		//m["resp"] = resp

	}

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

func HometaxLoginCheck(c echo.Context) error {
	dprintf(4, c, "call HometaxLoginCheck\n")

	m := make(map[string]interface{})
	cookie, rst := accessMainPage()
	if rst < 0 {
		m["resultCode"] = "99"
		m["resultMsg"] = "access fail"
		return c.JSON(http.StatusOK, m)
	}

	params := cls.GetParamJsonMap(c)

	//체크
	cashInfo, err := cls.GetSelectData(usersql.SelectStoreCashInfo, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if cashInfo != nil {
		hometaxAuthFail, _ := strconv.Atoi(cashInfo[0]["HOMETAX_AUTH_FAIL"])
		homtaxtRemainTime, _ := strconv.Atoi(cashInfo[0]["HOMTAXT_REMAIN_TIME"])

		if hometaxAuthFail > 3 {
			if homtaxtRemainTime > 0 {
				resultData := make(map[string]interface{})
				resultData["hometaxAuthFail"] = hometaxAuthFail
				resultData["homtaxtFailDt"] = cashInfo[0]["HOMTAXT_FAIL_DT"]
				resultData["homtaxtRemainTime"] = homtaxtRemainTime
				resultData["homeTaxErrMsg"] = "3회 이상 오류시 홈텍스 페이지에서 비밀번호 재설정이 필요 할 수 있습니다.(현재 " + strconv.Itoa(hometaxAuthFail) + "회 오류)"

				m["resultCode"] = "99"
				m["resultMsg"] = "login Fail"
				m["resultData"] = resultData
				return c.JSON(http.StatusOK, m)
			}
		}
	}

	//체크 끝

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
		lprintf(4, "[FAIL]login fail: check id/passwd (%s) \n", result)
		m["resultCode"] = "99"
		m["resultMsg"] = "login fail"

		homeTaxAuthFailCnt, homeTaxFailDt := AuthFailUpdate(c, params, "hometax")

		resultData := make(map[string]interface{})
		resultData["homeTaxAuthFailCnt"] = homeTaxAuthFailCnt
		resultData["homeTaxFailDt"] = homeTaxFailDt
		resultData["homeTaxErrMsg"] = "3회 이상 오류시 홈텍스 페이지에서 비밀번호 재설정이 필요 할 수 있습니다. (현재 " + strconv.Itoa(homeTaxAuthFailCnt) + "회 오류)"

		m["resultData"] = resultData

		return c.JSON(http.StatusOK, m)
	}
	AuthSuccessUpdate(c, params, "hometax")
	m["resultCode"] = "00"
	m["resultMsg"] = "login success"

	return c.JSON(http.StatusOK, m)
}

// 여신협회 or 홈택스 성공
func AuthSuccessUpdate(c echo.Context, params map[string]string, urlType string) {

	dprintf(4, c, "call AuthSuccessUpdate\n")

	storeInfo, err := cls.GetSelectData(storesql.SelectCash, params, c)
	if err != nil {
		dprintf(4, c, "AuthSuccessUpdate storeInfo select  query parameter fail\n")
	}

	var query string
	if storeInfo == nil {
		// Insert
		query, err = cls.SetUpdateParam(storesql.InsertCash, params)
		if err != nil {
			dprintf(4, c, "AuthSuccessUpdate storeInfo insert  query parameter fail\n")
		}
		_, err = cls.QueryDB(query)
		if err != nil {
			dprintf(4, c, "AuthSuccessUpdate storeInfo insert  query parameter fail\n")
		}
	}

	sqlQuery := ""
	if urlType == "ln" {
		sqlQuery = commonsql.UpdateLnAuthSuccess
		//여신협회 성공시 가입3번 알림톡 전송
		storeId := params["storeId"]
		//kakao.JoinOk_KakaoAlim(storeId)
		kakao.Coupon_KakaoAlim(storeId)
	} else {
		sqlQuery = commonsql.UpdateHomeTaxAuthSuccess
	}

	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(sqlQuery, params)
	if err != nil {
		dprintf(4, c, "UpdateCompAuthFail query parameter fail\n")
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		dprintf(4, c, "UpdateCompAuthFail DB fail\n")
	}

}

// 여신협회 or 홈택스 실패 기록
func AuthFailUpdate(c echo.Context, params map[string]string, urlType string) (int, string) {

	dprintf(4, c, "call AuthFailUpdate\n")

	storeInfo, err := cls.GetSelectData(storesql.SelectCash, params, c)
	if err != nil {
		return 0, ""
	}

	var query string
	if storeInfo == nil {
		// Insert
		query, err = cls.SetUpdateParam(storesql.InsertCash, params)
		if err != nil {
			return 0, ""
		}
		_, err = cls.QueryDB(query)
		if err != nil {
			return 0, ""
		}
	}

	authInfo, err := cls.GetSelectData(commonsql.SelectCompAuthInfo, params, c)
	if err != nil {
		return 0, ""
	}
	lnFailDt := ""
	homeFailDt := ""
	return1 := 0
	return2 := ""
	if urlType == "ln" {
		lnAuthFail, _ := strconv.Atoi(authInfo[0]["LN_AUTH_FAIL"])
		lnFailDt = authInfo[0]["LN_FAIL_DT"]
		lnAuthFailCnt := lnAuthFail + 1
		params["lnAuthFailCnt"] = strconv.Itoa(lnAuthFailCnt)

		return1 = lnAuthFailCnt
		return2 = lnFailDt
	} else {
		homeTaxAuthFail, _ := strconv.Atoi(authInfo[0]["HOMETAX_AUTH_FAIL"])
		homeFailDt = authInfo[0]["HOMTAXT_FAIL_DT"]
		homeTaxAuthFailCnt := homeTaxAuthFail + 1
		params["homeTaxAuthFailCnt"] = strconv.Itoa(homeTaxAuthFailCnt)

		return1 = homeTaxAuthFailCnt
		return2 = homeFailDt
	}

	sqlQuery := ""
	if urlType == "ln" {
		sqlQuery = commonsql.UpdateLnAuthFail
	} else {
		sqlQuery = commonsql.UpdateHomeTaxAuthFail
	}

	// 파라메터 맵으로 쿼리 변환
	selectQuery, err := cls.GetQueryJson(sqlQuery, params)
	if err != nil {
		dprintf(4, c, "UpdateCompAuthFail query parameter fail\n")
	}
	// 쿼리 실행
	_, err = cls.QueryDB(selectQuery)
	if err != nil {
		dprintf(4, c, "UpdateCompAuthFail DB fail\n")
	}

	return return1, return2
}

func GetSearchId(c echo.Context) error {

	dprintf(4, c, "call GetSearchId\n")

	params := cls.GetParamJsonMap(c)
	resultList, err := cls.GetSelectTypeRequire(usersql.SelectUserIdSearch, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultList == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "이름, 전화번호에 해당하는 고객을 찾을 수 없습니다."))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultList"] = resultList

	return c.JSON(http.StatusOK, m)

}

func GetSearchPw(c echo.Context) error {

	dprintf(4, c, "call GetSearchPw\n")

	params := cls.GetParamJsonMap(c)
	resultData, err := cls.GetSelectDataRequire(usersql.SelectUserPwSearch, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "DB fail"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "로그인 아이디, 전화번호에 해당하는 고객을 찾을 수 없습니다."))
	}

	loginType := resultData[0]["LOGIN_TYPE"]

	if loginType != "ID" {
		return c.JSON(http.StatusOK, controller.SetErrResult("99", loginType+" 로그인을 사용하는 계정입니다. 비밀번호 없이 해당 소셜 로그인을 기능을 이용하면 됩니다."))
	}

	userData := make(map[string]interface{})
	userData["userId"] = resultData[0]["USER_ID"]

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}

func SetChPw(c echo.Context) error {

	dprintf(4, c, "call SetChPw\n")

	params := cls.GetParamJsonMap(c)

	UpdatePwChQuery, err := cls.GetQueryJson(usersql.UpdateUserPassWd, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(UpdatePwChQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"

	return c.JSON(http.StatusOK, m)

}
