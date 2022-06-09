package billing

var SelectBillingInfo string = `SELECT PAYER_NAME
							,PAY_CARDNUM
							,PAY_CARDNAME
							,NEXT_PAY_DAY
							,END_DATE
							,PAY_YN
							,B.ITEM_NAME
							,B.PRICE
							,B.DC_PRICE
							,DATE_FORMAT(A.REG_DATE, '%Y-%m-%d')  AS REG_DATE
							,A.B_ID
							FROM e_billing AS A
							INNER JOIN e_billing_item AS B ON A.ITEM_CODE = B.ITEM_CODE
							WHERE 
							STORE_ID='#{storeId}'
							`

var SelectUserInfo string = `SELECT USER_NM
											,HP_NO
											,IFNULL(EMAIL,'') AS EMAIL
											,USER_ID
									FROM priv_user_Info
									WHERE 
									USER_ID='#{userId}'
									`
var SelectPayOidSeq string = `SELECT CONCAT('P',IFNULL(LPAD(MAX(SUBSTRING(PCD_PAY_OID, -10)) + 1, 10, 0), '0000000001')) as payOid
								FROM e_billing_payment
							`
var SelectBillingFreeCheck string = `SELECT count(*) as bCnt
							FROM e_billing
							WHERE 
							STORE_ID='#{storeId}' 
							AND FREE_USE_YN='Y'
							`
var SelectBillingItemInfo string = `SELECT ITEM_NAME
									,PRICE
									,DC_PRICE
									,ITEM_DESC
									,PRICE - DC_PRICE as ITEM_PRICE_DC
									FROM e_billing_item
									WHERE 
									ITEM_CODE='#{itemCode}'
									`

var InsertBillingKey string = `INSERT INTO e_billing
		(
			USER_ID
			,STORE_ID
			if #{PCD_PAYER_ID} != '' then ,PAYER_ID	 
			if #{PCD_PAY_CARDNUM} != '' then ,PAY_CARDNUM	
			if #{PCD_PAY_CARDNAME} != '' then ,PAY_CARDNAME	 
			if #{PCD_PAYER_NAME} != '' then ,PAYER_NAME	
			, NEXT_PAY_DAY	
			, REG_DATE
			, ITEM_CODE
			, START_DATE
			, END_DATE
			, PAY_YN
		)
		VALUES
		(
			'#{userId}'
			,'#{storeId}'
			,'#{PCD_PAYER_ID}'	
			,'#{PCD_PAY_CARDNUM}' 	
			,'#{PCD_PAY_CARDNAME}'	 
			,'#{PCD_PAYER_NAME}'	 
			, DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')	 
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, '#{itemCode}'
			, DATE_FORMAT( now(), '%Y-%m-%d')
			, DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
			, '#{payYn}'
		)
		`
var UpdateBillingKey string = `UPDATE e_billing SET 
									 PAYER_ID= '#{PCD_PAYER_ID}'	
									,PAY_CARDNUM= '#{PCD_PAY_CARDNUM}' 	
									,PAY_CARDNAME= '#{PCD_PAY_CARDNAME}'	 
									,PAYER_NAME= '#{PCD_PAYER_NAME}'
									,NEXT_PAY_DAY = CASE WHEN END_DATE >= SYSDATE() THEN  DATE_FORMAT( date_add(NEXT_PAY_DAY, interval '#{useMonth}' month) , '%Y-%m-%d')
														ELSE  DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,MOD_DATE =DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									,ITEM_CODE = '#{itemCode}'
									,START_DATE = CASE WHEN END_DATE >= SYSDATE() THEN  START_DATE
														ELSE   DATE_FORMAT( date_add(START_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,END_DATE = CASE WHEN END_DATE >= SYSDATE() THEN  DATE_FORMAT( date_add(END_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														ELSE  DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,PAY_YN = '#{payYn}'
								WHERE
									B_ID ='#{bId}'
								`

var UpdateBillingCancel string = ` UPDATE e_billing SET
											PAY_YN =  'N',	
											NEXT_PAY_DAY='',
											MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										WHERE 
										STORE_ID = '#{storeId}'
									`
var UpdateBillingCancelReload string = ` UPDATE e_billing SET
									PAY_YN =  'Y',	
									NEXT_PAY_DAY=END_DATE,
									MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
								STORE_ID = '#{storeId}'
								`

var InsertBillingPayment string = `INSERT INTO e_billing_payment
										(
										 PCD_PAY_OID
										 if #{bId} != '' then , B_ID 
										 if #{storeId} != '' then , STORE_ID
										 if #{PCD_PAYER_NO} != '' then , PCD_PAYER_NO
										 if #{PCD_PAY_GOODS} != '' then , PCD_PAY_GOODS
										 if #{PCD_PAY_TOTAL} != '' then , PCD_PAY_TOTAL
										 if #{PCD_PAY_ISTAX} != '' then , PCD_PAY_ISTAX
										 if #{PCD_PAY_TAXTOTAL} != '' then , PCD_PAY_TAXTOTAL
										 if #{PCD_PAY_YEAR} != '' then , PCD_PAY_YEAR
										 if #{PCD_PAY_MONTH} != '' then , PCD_PAY_MONTH
 										 if #{PCD_PAYER_NAME} != '' then , PCD_PAYER_NAME
										 if #{PCD_PAY_CARDNAME} != '' then , PCD_PAY_CARDNAME
										 if #{PCD_PAY_CARDNUM} != '' then , PCD_PAY_CARDNUM
										 if #{PCD_PAY_CARDTRADENUM} != '' then , PCD_PAY_CARDTRADENUM
										 if #{PCD_PAY_CARDAUTHNO} != '' then , PCD_PAY_CARDAUTHNO
										 if #{PCD_PAY_CARDRECEIPT} != '' then , PCD_PAY_CARDRECEIPT
										 if #{PCD_PAY_TIME} != '' then , PCD_PAY_TIME
										  , REG_DATE
										 if #{PCD_PAY_CODE} != '' then , PCD_PAY_CODE
										 if #{PCD_PAY_MSG} != '' then , PCD_PAY_MSG
                                          , PAY_TYPE
										)
										VALUES
										(
											'#{PCD_PAY_OID}'
											, '#{bId}'
											, '#{storeId}'
											, '#{PCD_PAYER_NO}'
											, '#{PCD_PAY_GOODS}'
											, '#{PCD_PAY_TOTAL}'
											, '#{PCD_PAY_ISTAX}'
											, '#{PCD_PAY_TAXTOTAL}'
											, '#{PCD_PAY_YEAR}'
											, '#{PCD_PAY_MONTH}'
											, '#{PCD_PAYER_NAME}'
											, '#{PCD_PAY_CARDNAME}'
											, '#{PCD_PAY_CARDNUM}'
											, '#{PCD_PAY_CARDTRADENUM}'
											, '#{PCD_PAY_CARDAUTHNO}'
											, '#{PCD_PAY_CARDRECEIPT}'
											, '#{PCD_PAY_TIME}'
											, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
											, '#{PCD_PAY_CODE}'
											, '#{PCD_PAY_MSG}'
											, '#{payType}'
										)
									`

var InsertBillingPayment_ETC string = `INSERT INTO e_billing_payment
									(
									 PCD_PAY_OID
									, B_ID
									, PAY_STAT
									if #{PCD_PAYER_NO} != '' then , PCD_PAYER_NO
									if #{PCD_PAY_GOODS} != '' then , PCD_PAY_GOODS
									if #{PCD_PAY_TOTAL} != '' then , PCD_PAY_TOTAL
									if #{PCD_PAY_ISTAX} != '' then , PCD_PAY_ISTAX
									if #{PCD_PAY_TAXTOTAL} != '' then , PCD_PAY_TAXTOTAL
									if #{PCD_PAY_YEAR} != '' then , PCD_PAY_YEAR
									if #{PCD_PAY_MONTH} != '' then , PCD_PAY_MONTH
									if #{PCD_PAY_CARDNAME} != '' then , PCD_PAY_CARDNAME
									if #{PCD_PAY_CARDNUM} != '' then , PCD_PAY_CARDNUM
									if #{PCD_PAY_CARDTRADENUM} != '' then , PCD_PAY_CARDTRADENUM
									if #{PCD_PAY_CARDAUTHNO} != '' then , PCD_PAY_CARDAUTHNO
									if #{PCD_PAY_CARDRECEIPT} != '' then , PCD_PAY_CARDRECEIPT
									if #{PCD_PAY_TIME} != '' then , PCD_PAY_TIME
									, REG_DATE
									)
									VALUES
									(
										'#{PCD_PAY_OID}'
										, '#{bId}'
										, '#{payStat}'
										, '#{PCD_PAYER_NO}'
										, '#{PCD_PAY_GOODS}'
										, '#{PCD_PAY_TOTAL}'
										, '#{PCD_PAY_ISTAX}'
										, '#{PCD_PAY_TAXTOTAL}'
										, '#{PCD_PAY_YEAR}'
										, '#{PCD_PAY_MONTH}'
										, '#{PCD_PAY_CARDNAME}'
										, '#{PCD_PAY_CARDNUM}'
										, '#{PCD_PAY_CARDTRADENUM}'
										, '#{PCD_PAY_CARDAUTHNO}'
										, '#{PCD_PAY_CARDRECEIPT}'
										, '#{PCD_PAY_TIME}'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									)
								`


var UpdateBillingPayment string = `UPDATE e_billing_payment SET 
									PAY_STAT= '#{payStat}'
									WHERE 
									PCD_PAY_OID = '#{PCD_PAY_OID}'
									`

var SelectStoreUserInfo string = `SELECT C.LOGIN_ID
											,C.LOGIN_PW
											,A.REST_NM
									FROM priv_rest_info AS A
									INNER JOIN priv_rest_user_info AS B ON A.REST_ID = B.REST_ID
									INNER JOIN priv_user_info AS C ON B.USER_ID = C.USER_ID
									WHERE 
									A.REST_ID = '#{storeId}'
									AND C.USER_ID= '#{userId}'
									`
var InsertSysInfo string = `INSERT INTO sys_user_info
									(USER_ID
									, USER_NM
									, USER_PASS
									, MSG_LANG_CD
									, CONN_ALLOW_YN
									, LAST_CONN_DATE
									, RETRY_CNT
									, INIT_PASS_YN
									, USER_MENU_AUTHOR_YN
									, AUTHOR_CD
									, ADDED_BY
									, ADDED_DATE
									)
									VALUES 
									(
									 '#{loginId}'
									, '#{userNm}'
									, '#{userPass}'
									, 'ko_KR'
									, 'Y'
									, NOW()
									, 0
									, 'N'
									, 'N'
									, '#{authorCd}'
									, 'system'
									, NOW()
									)
									`
var UpdateSysInfo string = `UPDATE sys_user_info SET 
									CONN_ALLOW_YN= 'Y'
									WHERE 
									USER_ID = '#{loginId}'
									`

