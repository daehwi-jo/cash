package controller

import (
	"cashApi/src/controller/cls"
)

func SetErrResult(code, msg string) ResponseHeader {
	var resCode ResponseHeader
	resCode.ResultCode = code
	resCode.ResultMsg = msg

	cls.Fprint(code, msg)

	return resCode
}

func SetResult(code string, data []byte) Response {
	var resCode Response
	resCode.ResultCode = code
	resCode.ResultMsg = "성공"
	resCode.ResultData = data
	return resCode
}

/*
func getPage(rowNum, sPage, lPage, cPage int) Page {
	var pageInfo Page
	pageInfo.TotalBlock = 1
	pageInfo.StartPage = sPage
	pageInfo.LastPage = lPage
	pageInfo.CurrentBlock = 1
	pageInfo.TotalPage = 1
	pageInfo.NextPage = 0
	pageInfo.PrevPage = 0
	pageInfo.TotalCount = rowNum
	pageInfo.CurrentPage = cPage

	var pageInfoList PageList
	pageInfoList.PageNoText = cPage
	pageInfoList.PageNo = cPage
	pageInfoList.ClassName = "on"

	pageInfo.PagingList = append(pageInfo.PagingList, pageInfoList)

	return pageInfo
}

*/

// common
type ResponseHeader struct {
	ResultCode string `json:"resultCode"` // result code
	ResultMsg  string `json:"resultMsg"`  // result msg

}

// common
type Response struct {
	ResultCode string      `json:"resultCode"` // result
	ResultMsg  string      `json:"resultMsg"`  // result code
	ResultData interface{} `json:"resultData"` // result data
}

type PageList struct {
	PageNoText int    `json:"pageNoText"`
	PageNo     int    `json:"pageNo"`
	ClassName  string `json:"className"`
}

type UseHistory struct {
	UserId    string `json:"UserId"`    // user id
	UserNm    string `json:"UserNm"`    // user name
	RestNm    string `json:"RestNm"`    // 가맹점 이름
	ItemNm    string `json:"ItemNm"`    // 구매 아이템 이름
	OrderDate string `json:"OrderDate"` // 주문날자시간
	UseAmt    string `json:"UseAmt"`    // 주문 금액
}
