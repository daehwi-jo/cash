package delions

import (
	delionsql "cashApi/query/delions"
	storesql "cashApi/query/stores"
	usersql "cashApi/query/users"
	"cashApi/src/controller"
	apis "cashApi/src/controller/API/ETC"
	"cashApi/src/controller/cls"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strings"
)


/* log format */
// 로그 레벨(5~1:INFO, DEBUG, GUIDE, WARN, ERROR), 1인 경우 DB 롤백 필요하며, 에러 테이블에 저장
// darayo printf(로그레벨, 요청 컨텍스트, format, arg) => 무엇을(서비스, 요청), 어떻게(input), 왜(원인,조치)
var dprintf func(int, echo.Context, string, ...interface{}) = cls.Dprintf
var lprintf func(int, string, ...interface{}) = cls.Lprintf

// 딜리온 가입
func SetDelionJoin(c echo.Context) error {

	dprintf(4, c, "call SetDelionJoin\n")


	//파라미터 체크
	checkCode,checkMsg := delionParamCheck(c)
	if checkCode == "99"{
		return c.JSON(http.StatusOK, controller.SetErrResult("99", checkMsg))
	}


	params := cls.GetParamJsonMap(c)


	// 임시 테이블 저장
	insertQuery, err := cls.GetQueryJson(delionsql.InsertDelionTemp, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "query parameter fail"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(insertQuery)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}



	//유저 아이디 생성
	resultData, err := cls.GetSelectData(usersql.SelectCreatUserSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(1)"))
	}
	if resultData == nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(2.1)"))
	}
	userId := resultData[0]["newUserId"]
	params["userId"] = userId

	//가맹점 아이디 생성
	storeSeqData, err := cls.GetSelectData(storesql.SelectStoreSeq, params, c)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(2.2)"))
	}
	storeId := storeSeqData[0]["storeSeq"]
	params["storeId"] = storeId



	// 유저 가입  TRNAN 시작
	tx, err := cls.DBc.Begin()
	if err != nil {
		//return "5100", errors.New("begin error")
	}

	// 오류 처리
	defer func() {
		if err != nil {
			// transaction rollback
			dprintf(4, c, "do rollback -유저 가입(SetDelionJoin)  \n")
			tx.Rollback()
		}
	}()

	// transation exec
	// 파라메터 맵으로 쿼리 변환


	// 유저 생성
	params["userTy"] = "1"
	params["atLoginYn"] = "Y"
	params["pushYn"] = "Y"
	params["loginPw"] = sha256Encoding(params["phone"])
	params["kind"] =""
	params["category"] =""

	//회원 가입
	UserCreateQuery, err := cls.SetUpdateParam(delionsql.InserCreateUserDelion, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "관리자에게 문의 바랍니다(3)"))
	}

	_, err = tx.Exec(UserCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", UserCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", UserCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(4)"))
	}

	//회원가입 동의
	TermsCreateQuery, err := cls.SetUpdateParam(delionsql.InsertTermsUserDelion, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "관리자에게 문의 바랍니다(5)"))
	}

	_, err = tx.Exec(TermsCreateQuery)
	dprintf(4, c, "call set Query (%s)\n", TermsCreateQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", TermsCreateQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(6)"))
	}

	// 가맹점 등록

	params["channelCode"] ="DELION"
	insertStore, err := cls.SetUpdateParam(delionsql.InsertStoreDelion, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(7)"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(insertStore)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(8)"))
	}

	// 매장에 사장님 등록
	InsertStoreUserQuery, err := cls.GetQueryJson(delionsql.InsertStoreUserDelion, params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, controller.SetErrResult("98", "관리자에게 문의 바랍니다(9)"))
	}
	_, err = tx.Exec(InsertStoreUserQuery)
	dprintf(4, c, "call set Query (%s)\n", InsertStoreUserQuery)
	if err != nil {
		dprintf(1, c, "Query(%s) -> error (%s) \n", InsertStoreUserQuery, err)
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(10)"))
	}


	params["lnId"] = "k"+params["biz_number"]
	params["lnPsw"] = params["cardsales_pass"]

	// 캐쉬컴바인 조인
	CashJoinDelion, err := cls.SetUpdateParam(delionsql.InsertCashDelion, params)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(11)"))
	}
	// 쿼리 실행
	_, err = cls.QueryDB(CashJoinDelion)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(12)"))
	}


	apis.BillingInsert(storeId,userId,"I0000000001","3")

	// transaction commit
	err = tx.Commit()
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", "관리자에게 문의 바랍니다(13)"))
	}

	// 유저 가입 TRNAN 종료

	//userData := make(map[string]interface{})
	//userData["userId"] = userId

	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	//m["resultData"] = userData

	return c.JSON(http.StatusOK, m)

}



func sha256Encoding(str string) string {

	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	mdstr := hex.EncodeToString(md)
	return mdstr

}


func delionParamCheck(c echo.Context) (string,string) {

	code :="00"
	msg :="00"

	darayokey := c.Request().Header.Get("darayokey")
	biz_name := c.Request().FormValue("biz_name")
	biz_number := c.Request().FormValue("biz_number")
	cardsales_pass := c.Request().FormValue("cardsales_pass")
	phone := c.Request().FormValue("phone")
	birth_date := c.Request().FormValue("birth_date")
	ceo_name := c.Request().FormValue("ceo_name")
	//biz_kind := c.Request().FormValue("biz_kind")
	//email := c.Request().FormValue("email")
	chkKey := sha256Encoding("k"+biz_number);

	println(chkKey)
	
	if chkKey != darayokey{
		code="99"
		msg ="api key가 잘못되었습니다."
		return code,msg
	}
	if biz_name =="" {
		code="99"
		msg ="상호명은 필수입니다."
		return code,msg
	}
	if biz_number ==""  {
		code="99"
		msg ="사업자번호는 필수입니다."
		return code,msg
	}
	if cardsales_pass ==""  {
		code="99"
		msg ="여신협회 비밀번호는 필수입니다. "
		return code,msg
	}
	if phone =="" {
		code="99"
		msg ="전화번호는 필수 입니다"
		return code,msg
	}
	if birth_date ==""  {
		code="99"
		msg ="생년월일은 필수입니다."
		return code,msg
	}
	if ceo_name ==""  {
		code="99"
		msg ="대표자 이름은 필수입니다."
		return code,msg
	}
	return code,msg

}


// 딜리온 가입
func TEST11(c echo.Context) error {

	dprintf(4, c, "call TEST11\n")



	params := cls.GetParamJsonMap(c)
	query:=url2.QueryEscape(params["biz_name"])
	url:="https://dapi.kakao.com/v2/local/search/keyword.json?query="+query
	client := &http.Client{}
	req, err := http.NewRequest("GET", url,nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("authorization", "KakaoAK c2bc6457eaaf39999ed317da290848ea")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var result KakaoResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return c.JSON(http.StatusOK, controller.SetErrResult("98", err.Error()))
	}

	categoryNm :=""
	addr :=""
	biz_name := params["biz_name"]
	tel := params["tel"]

	resultlength := len(result.Documents)

	if resultlength == 1 {
		categoryNm = result.Documents[0].Category_name
		addr = result.Documents[0].Address_name
	}else{

		for i := range result.Documents {
			kstoreNm := result.Documents[i].Place_name
			kstoreTel := strings.Replace(result.Documents[i].Phone, "-","",-1)

			if kstoreNm == biz_name &&  kstoreTel ==  tel {
				categoryNm = result.Documents[i].Category_name
				addr = result.Documents[i].Address_name
			}else {
				categoryNm = result.Documents[0].Category_name
				addr = result.Documents[0].Address_name
			}
		}

	}

	fmt.Println(resultlength)
	fmt.Println(categoryNm)
	fmt.Println(addr)




	// fmt.Println(string(body))


	m := make(map[string]interface{})
	m["resultCode"] = "00"
	m["resultMsg"] = "응답 성공"
	m["resultData"] = result

	return c.JSON(http.StatusOK, m)

}




type KakaoResult struct {
	Documents         []KakaoDocuments `json:"documents"`           //
	Meta          	  KakaoMeta 	   `json:"meta"`          //
}


type KakaoDocuments struct {
	Address_name        string `json:"address_name"`           //
	Category_group_code string `json:"category_group_code"`          //
	Category_group_name string `json:"category_group_name"` //
	Category_name 		string `json:"category_name"` //
	Distance 			string `json:"distance"` //
	Id 					string `json:"id"` //
	Phone 				string `json:"phone"` //
	Place_name 			string `json:"place_name"` //
	Place_url 			string `json:"place_url"` //
	Road_address_name 	string `json:"road_address_name"` //
	X 					string `json:"x"` //
	Y 					string `json:"y"` //
}

type KakaoMeta struct {
	Is_end         		bool `json:"is_end"`           //
	Pageable_count      int `json:"pageable_count"`          //
	Same_name		    Same_name `json:"same_name"`          //
	Total_count 		int `json:"total_count"` //
}


type Same_name struct {
	Keyword         		string `json:"keyword"`           //
	//Region      int `json:"region"`          //
	Selected_region 		string `json:"selected_region"` //
}

