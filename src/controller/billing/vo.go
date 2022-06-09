package billing

type AuthSendData struct {
	Cst_id       string `json:"cst_id"`       //
	CustKey      string `json:"custKey"`      //
	PCD_PAY_WORK string `json:"PCD_PAY_WORK"` //
}

type PaymentAuthData struct {
	Cst_id           string `json:"cst_id"`           //
	CustKey          string `json:"custKey"`          //
	PCD_REGULER_FLAG string `json:"PCD_REGULER_FLAG"` //
	PCD_PAY_TYPE     string `json:"PCD_PAY_TYPE"`     //
}

type AuthResult struct {
	Result  string `json:"result"`  //
	CustKey string `json:"custKey"` //
	Cst_id  string `json:"cst_id"`  //
	AuthKey string `json:"AuthKey"` //
}

type PaymentSendData struct {
	PCD_CST_ID         string `json:"PCD_CST_ID"`         //
	PCD_CUST_KEY       string `json:"PCD_CUST_KEY"`       //
	PCD_AUTH_KEY       string `json:"PCD_AUTH_KEY"`       //
	PCD_PAY_TYPE       string `json:"PCD_PAY_TYPE"`       //
	PCD_PAYER_NO       string `json:"PCD_PAYER_NO"`       //
	PCD_PAYER_ID       string `json:"PCD_PAYER_ID"`       //
	PCD_PAY_GOODS      string `json:"PCD_PAY_GOODS"`      //
	PCD_SIMPLE_FLAG    string `json:"PCD_SIMPLE_FLAG"`    //
	PCD_REGULER_FLAG   string `json:"PCD_REGULER_FLAG"`   //
	PCD_PAY_YEAR       string `json:"PCD_PAY_YEAR"`       //
	PCD_PAY_MONTH      string `json:"PCD_PAY_MONTH"`      //
	PCD_PAY_TOTAL      int    `json:"PCD_PAY_TOTAL"`      //
	PCD_PAY_OID        string `json:"PCD_PAY_OID"`        //
	PCD_PAYER_NAME     string `json:"PCD_PAYER_NAME"`     //
	PCD_PAYER_HP       string `json:"PCD_PAYER_HP"`       //
	PCD_PAYER_EMAIL    string `json:"PCD_PAYER_EMAIL"`    //
	PCD_PAY_ISTAX      string `json:"PCD_PAY_ISTAX"`      //
	PCD_PAY_TAXTOTAL   string `json:"PCD_PAY_TAXTOTAL"`   //
	PCD_REFUND_KEY     string `json:"PCD_REFUND_KEY"`     //
	PCD_PAYCANCEL_FLAG string `json:"PCD_PAYCANCEL_FLAG"` //
	PCD_REFUND_TOTAL   int    `json:"PCD_REFUND_TOTAL"`   //
	PCD_PAY_DATE       string `json:"PCD_PAY_DATE"`       //

}

type PaymentResultData struct {
	PCD_PAY_RST          string `json:"PCD_PAY_RST"`          //
	PCD_PAY_MSG          string `json:"PCD_PAY_MSG"`          //
	PCD_PAY_OID          string `json:"PCD_PAY_OID"`          //
	PCD_PAY_TYPE         string `json:"PCD_PAY_TYPE"`         //
	PCD_PAYER_NO         string `json:"PCD_PAYER_NO"`         //
	PCD_PAYER_ID         string `json:"PCD_PAYER_NO"`         //
	PCD_PAYER_NAME       string `json:"PCD_PAYER_NO"`         //
	PCD_PAYER_HP         string `json:"PCD_PAYER_NO"`         //
	PCD_PAYER_EMAIL      string `json:"PCD_PAYER_NO"`         //
	PCD_PAY_YEAR         string `json:"PCD_PAY_YEAR"`         //
	PCD_PAY_MONTH        string `json:"PCD_PAY_MONTH"`        //
	PCD_PAY_GOODS        string `json:"PCD_PAY_GOODS"`        //
	PCD_PAY_TOTAL        string `json:"PCD_PAY_TOTAL"`        //
	PCD_PAY_TAXTOTAL     int    `json:"PCD_PAY_TAXTOTAL"`     //
	PCD_PAY_ISTAX        string `json:"PCD_PAY_ISTAX"`        //
	PCD_PAY_TIME         string `json:"PCD_PAY_TIME"`         //
	PCD_PAY_CARDNAME     string `json:"PCD_PAY_CARDNAME"`     //
	PCD_PAY_CARDNUM      string `json:"PCD_PAY_CARDNUM"`      //
	PCD_PAY_CARDTRADENUM string `json:"PCD_PAY_CARDTRADENUM"` //
	PCD_PAY_CARDAUTHNO   string `json:"PCD_PAY_CARDAUTHNO"`   //
	PCD_PAY_CARDRECEIPT  string `json:"PCD_PAY_CARDRECEIPT"`  //
	PCD_USER_DEFINE1     string `json:"PCD_USER_DEFINE1"`     //

}
