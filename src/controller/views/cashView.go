package views

import (
	datasql "cashApi/query/datas"
	"net/http"

	//homesql "cashApi/query/homes"
	"cashApi/src/controller/cls"
	// login 및 기본

	"github.com/labstack/echo/v4"
)

// 주간
func GetWeekView(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["s_Siteid"] = ""
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "views/week.htm", m)
}

//  월간
func GetMonthView(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["s_Siteid"] = ""
	m["bizNum"] = params["bizNum"]
	//m["restId"] = params["restId"]

	regDt, err := cls.GetSelectData(datasql.SelectRegistDate, params, c)
	if err == nil {
		m["restId"] = regDt[0]["restId"]
	}

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "views/month.htm", m)
}

// 주간
func GetWeekViewPos(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	c.Response().Header().Set("Content-Security-Policy", "frame-ancestors *")
	m := make(map[string]interface{})
	m["s_Siteid"] = ""
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "pos_view/fit_darayo_client.htm", m)
}

//  월간
func GetMonthViewPos(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	c.Response().Header().Set("Content-Security-Policy", "frame-ancestors *")

	m := make(map[string]interface{})
	m["s_Siteid"] = ""
	m["bizNum"] = params["bizNum"]

	//lprintf(4, "[INFO] pageNum : %s, siteId : %s\n", pageNum, siteId)

	return c.Render(http.StatusOK, "pos_view/fit_darayo_monthly.htm", m)
}





//  월간
func GetMonthlyView_V2(c echo.Context) error {

	params := cls.GetParamJsonMap(c)

	m := make(map[string]interface{})
	m["bizNum"] = params["bizNum"]

	return c.Render(http.StatusOK, "views/v2/monthly.htm", m)
}
