package parthner

import (
	homesql "cashApi/query/homes"
	usersql "cashApi/query/users"
	"cashApi/src/controller"
	"cashApi/src/controller/cls"
	users "cashApi/src/controller/users"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	// login 및 기본

)

var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

func ParthnerLogin(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "parthner/fit_darayo_login.htm", m)
}



func ParthnerLoginOk(c echo.Context) error {

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
		users.LoginAcceesLog(c,params)
		return c.JSON(http.StatusOK, controller.SetErrResult("99", "아이디 또는 비밀번호가 잘못되었습니다."))
	}

	useYn := resultData[0]["USE_YN"]
	userId := resultData[0]["userId"]

	loginId:= resultData[0]["LOGIN_ID"]
	userNm:= resultData[0]["USER_NM"]
	userTel:= resultData[0]["HP_NO"]
	userBirth:= resultData[0]["USER_BIRTH"]




	params["userId"] = userId
	if useYn == "N" {
		// 접속 로그
		params["logInOut"] = "I"
		params["ip"] = ip
		params["succYn"] = "N"
		params["type"] = "id"
		users.LoginAcceesLog(c,params)
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
	users.LoginAcceesLog(c,params)


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





func ParthnerJoin(c echo.Context) error {

	params := cls.GetParamJsonMap(c)



	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "parthner/fit_darayo_join.htm", m)
}





func ParthnerJoinStep1(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "parthner/fit_darayo_join_step1.htm", m)
}


func ParthnerJoinStep2(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "parthner/fit_darayo_join_step2.htm", m)
}



func GuideCertify(c echo.Context) error {

	//cMsg,nextGo,introMsg := homedata(c)

	m := make(map[string]interface{})
	//m["cMSG"] = cMsg
	//m["nextGo"] = nextGo
	//m["introMsg"] = introMsg

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "parthner/fit_darayo_guide_certify.htm", m)
}

func GuidePartner(c echo.Context) error {

	//cMsg,nextGo,introMsg := homedata(c)


	m := make(map[string]interface{})
	//m["cMSG"] = cMsg
	//m["nextGo"] = nextGo
	//m["introMsg"] = introMsg

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "parthner/fit_darayo_guide_partner.htm", m)
}


func Homedata(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	resultData, err := cls.GetSelectData(homesql.SelectIntroMsg, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}



	storeData, err := cls.GetSelectData(homesql.SelectStoreService, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	introMsg := resultData[0]["introMsg"]
	billingYn := storeData[0]["billingYn"]
	cardSalesYn := storeData[0]["cardSalesYn"]
	homeTaxYn := storeData[0]["homeTaxYn"]


	billingMsg :=""
	channelDiv :="0"


	if billingYn != "N" {
		billingYn = "Y"

		billingInfo, err := cls.GetSelectData(homesql.SelectStoreBillingInfo, params, c)
		if err != nil {
			return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
		}

		billingDate := billingInfo[0]["END_DATE"]
		channelDiv = billingInfo[0]["CHANNEL_DIV"]


		billingMsg = billingDate +" 까지 사용 가능합니다."
		if channelDiv == "1" {
			nextPayDay := billingInfo[0]["NEXT_PAY_DAY"]
			payYn:= billingInfo[0]["PAY_YN"]

			if payYn== "Y"{
				billingMsg = "파트너 멤버입니다.\n 다음 결제일 : " +nextPayDay
			}else {
				billingMsg = "파트너 멤버입니다.\n " + billingDate +" 까지 사용 가능합니다."
			}

		}



	}

	if cardSalesYn != "N" {
		cardSalesYn = "Y"
	}

	if homeTaxYn != "N" {
		homeTaxYn = "Y"
	}

	cMsg :=""
	nextGo:="N"
	if cardSalesYn == "N" && homeTaxYn == "N" {
		cMsg = "여신협회 및 홈택스 미인증 상태 입니다."
	}else  if cardSalesYn == "N" && homeTaxYn == "Y" {
		cMsg = "홈택스 인증 상태 입니다."
	}else{
		nextGo ="Y"
	}

	homeData := make(map[string]interface{})
	homeData["cMSG"] = cMsg
	homeData["nextGo"] = nextGo
	homeData["introMsg"] = introMsg
	homeData["billingMsg"] = billingMsg
	homeData["billingYn"] = billingYn
	homeData["channelDiv"] = channelDiv


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = homeData


	return c.JSON(http.StatusOK, m)
}