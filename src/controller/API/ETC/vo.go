package apis

import "encoding/xml"


//<map id='' >
//	<map id='resultMsg' >
//		<detailMsg></detailMsg>
//		<msg></msg>
//		<code></code>
//		<result>S</result>
//	</map>
//	<trtEndCd>Y</trtEndCd>
//	<smpcBmanEnglTrtCntn>The business registration number is registered</smpcBmanEnglTrtCntn>
//	<nrgtTxprYn>N</nrgtTxprYn>
//	<smpcBmanTrtCntn>등록되어 있는 사업자등록번호 입니다. </smpcBmanTrtCntn>
//	<trtCntn>부가가치세 일반과세자 입니다.</trtCntn>
//</map>

type ResultMsgMap struct {
	XMLName  xml.Name `xml:"map"`
	Id string `xml:"id,attr"`
	DetailMsg  string   `xml:"detailMsg"`
	Msg  string   `xml:"msg"`
	Code  string   `xml:"code"`
	Result  string   `xml:"result"`
}

type Map struct {
	XMLName  xml.Name `xml:"map"`
	ResultMsgMap ResultMsgMap  `xml:"map"`
	TrtEndCd string   `xml:"trtEndCd"`
	SmpcBmanEnglTrtCntn  string   `xml:"smpcBmanEnglTrtCntn"`
	NrgtTxprYn  string   `xml:"nrgtTxprYn"`
	SmpcBmanTrtCntn  string   `xml:"smpcBmanTrtCntn"`
	TrtCntn  string   `xml:"trtCntn"`
}


type TpayResult struct {
	Result_cd         string `json:"result_cd"`           //
	Result_msg          	string `json:"result_msg"`          //
	Account_name 		string `json:"account_name"` //
}


type HdResult struct {
	Response         HdResponse `json:"response"`
}

type HdResponse struct {
	Header         HdHeader `json:"header"`
	Body           HdBody   `json:"body"`
}

type HdHeader struct {
	ResultCode 		string   `json:"resultCode"`
	ResultMsg 		string   `json:"resultMsg"`
}

type HdBody struct {
	Items         	HdItems `json:"items"`
	NumOfRows 		int     `json:"numOfRows"`
	PageNo 			int     `json:"pageNo"`
	TotalCount 		int     `json:"totalCount"`
}

type HdItems struct {
	Item         []HdItem `json:"item"`
}

type HdItem struct {
	DateKind     string `json:"dateKind"`
	DateName     string `json:"dateName"`
	IsHoliday 	 string `json:"isHoliday"`
	Locdate 	 int    `json:"locdate"`
	Seq 		 int    `json:"seq"`
}
