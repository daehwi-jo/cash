package apis

var UpdateHolyDay string = `UPDATE sys_week_date SET HD ='H', DATEKIND = 'G',
							DATENAME = '#{dateName}'
							WHERE 
							TOTAL_DATE ='#{locDate}'
                            `


// 심플오더 관련 시작
var  SelectSimpleGrpCnt string = `SELECT count(*) as grpCnt
								FROM priv_simple_info AS A
								INNER JOIN PRIV_GRP_USER_INFO AS B ON A.GRP_ID = B.GRP_ID
								INNER JOIN PRIV_USER_INFO AS C ON B.USER_ID= C.USER_ID
								WHERE A.USE_YN='Y' 
								AND A.REST_ID = '#{restId}'
								and  B.USER_ID = '#{userId}'
								and C.HP_NO = '#{hpNo}'
							
                             `

var  SelectSimpleGrp string = `SELECT A.GRP_ID 
									  ,B.USER_ID
									  ,B.AUTH_STAT
                                      ,C.USER_NM
								FROM priv_simple_info AS A
								INNER JOIN PRIV_GRP_USER_INFO AS B ON A.GRP_ID = B.GRP_ID
								INNER JOIN PRIV_USER_INFO AS C ON B.USER_ID= C.USER_ID
								WHERE A.USE_YN='Y' 
								AND A.REST_ID = '#{restId}'
								and  B.USER_ID = '#{userId}'
								and C.HP_NO = '#{hpNo}'
                             `

var  SelectStoreMenuView string = `SELECT
									 A.ITEM_NO
								   , A.ITEM_NM
								   , A.REST_ID
								   , A.ITEM_PRICE
								   , A.REG_DATE
								   , A.MOD_DATE
								   , A.USE_YN
								   , A.ITEM_MENU
								   , A.BEST_YN
								FROM DAR_SALE_ITEM_INFO A
								WHERE use_yn='Y'
								AND A.REST_ID = '#{restId}' 
								AND TICKET_YN='Y'
								LIMIT 1
                             `

var  SelectGrpOrderCheckData string = `SELECT CHECK_TIME	
										FROM PRIV_GRP_INFO
										WHERE 
										GRP_ID = '#{grpId}'
                             		`


var  SelectOrderCheck string = `SELECT	 count(*) as orderCnt
								FROM DAR_ORDER_INFO A
								WHERE 
									  A.REST_ID	= '#{restId}'
								AND   A.USER_ID = '#{userId}'
								AND   A.GRP_ID = '#{grpId}'
								AND   A.ORDER_STAT = '20'
								AND   A.TOTAL_AMT = '#{totalAmt}'
								AND   A.ORDER_DATE > DATE_FORMAT(DATE_ADD(SYSDATE(), INTERVAL -#{checkTime} SECOND), '%Y%m%d%H%i%s')
                             		`
var  CreateOrderSeq string = `SELECT FN_GET_SEQUENCE('ORDER_NO') as orderNo `

var  InsertSimpleOrder string = `INSERT INTO DAR_ORDER_INFO
								(
									ORDER_NO
									, REST_ID
									, USER_ID
									, GRP_ID
									, ORDER_STAT
									, ORDER_DATE
									, TOTAL_AMT
									, CREDIT_AMT
									, PAID_YN
									, CONFIRM_YN
									, RO_NO
									, ORDER_TY
									, PRINT_YN
									, PAY_TY
									, QR_ORDER_TYPE
								)
								VALUES
								(
									'#{orderNo}'
									, '#{restId}'
									, '#{userId}'
									, '#{grpId}'
									, '20'
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									, '#{totalAmt}'
									, '#{creditAmt}'
									, 'N'
									, 'N'
									, 0
									, '#{orderTy}'
									, 'N'
									, '#{payTy}'
									, '#{qrOrderType}'
								)
								`


var  InsertSimpleOrderDetail string = `INSERT INTO DAR_ORDER_DETAIL
										(
											  ORDER_NO
											, ORDER_SEQ
											, ITEM_NO
											, USER_ID
											, ORDER_DATE
											, ORDER_QTY
											, ORDER_AMT
										)
										VALUES
										(
											  '#{orderNo}'
											, '#{orderSeq}'
											, '#{itemNo}'
											, '#{userId}'
											, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
											, '#{orderQty}'
											, '#{orderAmt}'
										)
								`
// 심플오더 관련 끝



var  SelectCouponInfo string = `SELECT COUPON_NAME
											,USE_TYPE
											,COUPON_VAL
											,USE_SERVICE
											,ITEM_CODE
											,USE_YN
										FROM b_coupon AS a
										WHERE 
										START_DATE <= DATE_FORMAT(NOW(), '%Y%m%d')
										AND END_DATE  >= DATE_FORMAT(NOW(), '%Y%m%d')
										AND  COUPON_NO = '#{couponNo}'
                             		`


var  SelectCouponChk string = `SELECT COUNT(*) as cnt
									FROM b_coupon_his
									WHERE 
									COUPON_NO =  '#{couponNo}'
									AND USER_KEY= '#{userKey}'
									AND USER_TYPE= '#{userType}'
                             		`




var SelectBillingChk string = `SELECT END_DATE
							,B_ID
							FROM e_billing AS A
							WHERE 
							STORE_ID='#{storeId}'
							`


var InsertBillingCouponUse string = `INSERT INTO e_billing
		(
			USER_ID
			,STORE_ID
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
			, DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')	 
			, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
			, '#{itemCode}'
			, DATE_FORMAT( now(), '%Y-%m-%d')
			, DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
			, '#{payYn}'
		)
		`

var UpdateBillingCouponUse string = `UPDATE e_billing SET
									MOD_DATE =DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									,NEXT_PAY_DAY = CASE WHEN END_DATE >= SYSDATE() THEN  DATE_FORMAT( date_add(END_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														ELSE  DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,ITEM_CODE = '#{itemCode}'
									,START_DATE = CASE WHEN END_DATE >= SYSDATE() THEN  START_DATE
														ELSE   DATE_FORMAT( date_add(START_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,END_DATE = CASE WHEN END_DATE >= SYSDATE() THEN  DATE_FORMAT( date_add(END_DATE, interval '#{useMonth}' month) , '%Y-%m-%d')
														ELSE  DATE_FORMAT( date_add(now(), interval '#{useMonth}' month) , '%Y-%m-%d')
														END 
									,PAY_YN = '#{payYn}'
								WHERE
									STORE_ID ='#{storeId}'
									AND USER_ID = '#{userId}'
								`


var InserCouponHistory string = `INSERT INTO b_coupon_his
								(COUPON_NO
								, USER_KEY
								, USER_TYPE
								, REG_DATE
								)
								VALUES 
								(
									'#{couponNo}'
									, '#{userKey}'
									, '#{userType}'
									, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								)
							`


var InserBillingPayment string = `INSERT INTO e_billing_payment
										(
										 PCD_PAY_OID
										  ,STORE_ID
										 ,REG_DATE
                                         ,PAY_TYPE
										 ,ETC
                                         ,PAY_STAT
										)
										VALUES
										(
											concat('#{couponNo}'
												,'_','#{storeId}')
											, '#{storeId}'
											, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
											, 'C'
											, '#{etc}'
											, '20'
										)
									`