package books


var SelectLinkBookList string = `SELECT 		
									 B.GRP_ID
									,B.GRP_NM
									,IFNULL(ICON_LINK,'') AS ICON_LINK
									,IFNULL(COMPANY_NM,'') AS COMPANY_NM
									,(SELECT COUNT(*) FROM priv_grp_user_info AS aa WHERE aa.GRP_ID = a.GRP_ID AND aa.GRP_AUTH='1' AND aa.AUTH_STAT='1') AS userCnt
									,IFNULL(CASE B.GRP_TYPE_CD WHEN '10' THEN '4' ELSE B.GRP_TYPE_CD END,'4') AS BOOK_TYPE
									,E.USER_NM AS managerNm
									,A.PAY_TY
									,IFNULL(A.PREPAID_AMT,0) AS  PREPAID_AMT
									,IFNULL((SELECT SUM(TOTAL_AMT) FROM dar_order_info AS AA 
														WHERE 
															AA.REST_ID = A.REST_ID AND AA.GRP_ID = A.GRP_ID 
															AND  ORDER_sTAT='20' AND PAID_YN='N' AND PAY_TY='1' ),0) AS UNPAID_AMT
								FROM org_agrm_info  AS A
									INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID 
									INNER JOIN priv_grp_user_info AS D ON B.GRP_ID = D.GRP_ID AND D.GRP_AUTH='0'
									INNER JOIN priv_user_info AS E ON D.USER_ID = E.USER_ID
									LEFT OUTER JOIN ( SELECT AA.company_nm
															,BB.BOOK_ID
															,'' AS ICON_LINK
													FROM b_company AS AA
													INNER JOIN b_company_book AS BB ON AA.COMPANY_ID = BB.COMPANY_ID) AS C ON B.GRP_ID = C.BOOK_ID			
								WHERE 
									A.REST_ID ='#{storeId}'
									AND (B.GRP_NM LIKE '%#{search}%' 
                                         OR  E.USER_NM LIKE '%#{search}%' )
									order by A.AUTH_DATE DESC
								`
var SelectLinkBooksCnt string = `SELECT 	COUNT(*) AS bookCnt
								FROM org_agrm_info  AS A
									INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID 
									INNER JOIN priv_grp_user_info AS D ON B.GRP_ID = D.GRP_ID AND D.GRP_AUTH='0'
									INNER JOIN priv_user_info AS E ON D.USER_ID = E.USER_ID
								WHERE 
									A.REST_ID ='#{storeId}'
									AND (B.GRP_NM LIKE '%#{search}%' 
                                         OR  E.USER_NM LIKE '%#{search}%' )
								`


var SelectBookOrderCount string = `SELECT
									count(*) as totalCount
									FROM dar_order_info AS A
									INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID
									INNER JOIN priv_user_info AS D ON A.USER_ID = D.USER_ID
									WHERE
									A.order_ty IN ('1','2','3','5')
									AND A.ORDER_STAT in ('#{orderStat}')
									AND A.REST_ID='#{storeId}'
								    AND A.PAY_TY='1'  AND A.PAID_YN ='#{inputPaid}'
									AND LEFT(A.ORDER_DATE,8) >= '#{searchDate}'
									AND LEFT(A.ORDER_DATE,8) >= '#{startDate}' 
									AND LEFT(A.ORDER_DATE,8) <= '#{endDate}' 
									AND (D.USER_NM LIKE '%#{search}%'
									or B.GRP_NM LIKE '%#{search}%')
									`

var SelectBookOrderList string = `SELECT
									A.ORDER_NO
									,DATE_FORMAT(A.ORDER_DATE,'%Y.%m.%d %p %h:%i')  AS ORDER_DATE
									,IFNULL(company_nm,'') AS COMPANY_NM
									,A.GRP_ID
									,B.GRP_NM
									,A.TOTAL_AMT
									,CASE WHEN A.ORDER_TY ='1' THEN 'pay'
										WHEN A.ORDER_TY ='2' THEN 'delivery'
										WHEN A.ORDER_TY ='3' THEN 'takeout'
										WHEN A.ORDER_TY ='5' THEN 'pay'
									END AS ORDER_TY
									,D.USER_NM AS orderer
									,FN_GET_ORDER_USER_CNT(A.ORDER_NO) AS ordererCount
									,ORDER_STAT
									FROM dar_order_info AS A
									INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID
									INNER JOIN priv_user_info AS D ON A.USER_ID = D.USER_ID
									LEFT OUTER JOIN ( SELECT AA.company_nm
										,BB.BOOK_ID
										FROM b_company AS AA
										INNER JOIN b_company_book AS BB ON AA.COMPANY_ID = BB.COMPANY_ID) AS C ON B.GRP_ID = C.BOOK_ID
										WHERE
									A.order_ty IN ('1','2','3','5')
									AND A.ORDER_STAT in ('#{orderStat}')
									AND A.REST_ID='#{storeId}'
									AND A.PAY_TY='1' AND A.PAID_YN ='#{inputPaid}'
									AND LEFT(A.ORDER_DATE,8) >= '#{searchDate}'
									AND LEFT(A.ORDER_DATE,8) >= '#{startDate}' 
									AND LEFT(A.ORDER_DATE,8) <= '#{endDate}' 
									AND (D.USER_NM LIKE '%#{search}%'
									or B.GRP_NM LIKE '%#{search}%')
									ORDER BY A.ORDER_DATE DESC
									`


var SelectBookOrderTotal string = `SELECT
									SUM(A.TOTAL_AMT) as totalAmt
									FROM dar_order_info AS A
									INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID
									INNER JOIN priv_user_info AS D ON A.USER_ID = D.USER_ID
									WHERE
									A.order_ty IN ('1','2','3','5')
									AND A.ORDER_STAT ='20'
									AND A.REST_ID='#{storeId}'
								    AND A.PAY_TY='1'  AND A.PAID_YN ='#{inputPaid}' 
									AND LEFT(A.ORDER_DATE,8) >= '#{searchDate}'
									AND LEFT(A.ORDER_DATE,8) >= '#{startDate}' 
									AND LEFT(A.ORDER_DATE,8) <= '#{endDate}' 
									AND (D.USER_NM LIKE '%#{search}%'
									or B.GRP_NM LIKE '%#{search}%')
									`


var SelectBookPaymemtCount string = `SELECT count(*) as totalCount
								FROM dar_payment_hist AS A
								INNER JOIN priv_grp_info AS B ON A.grp_id = B.grp_id AND A.PAY_CHANNEL<>'05'
								INNER JOIN priv_user_info AS C ON A.USER_ID = C.USER_ID
								WHERE 
								A.REST_ID='#{storeId}'
								AND LEFT(A.REG_DATE,8) >= '#{searchDate}'
								AND LEFT(A.REG_DATE,8) >= '#{startDate}'
								AND LEFT(A.REG_DATE,8) <= '#{endDate}' 
								AND A.PAY_CHANNEL ='#{payChannel}' 
								AND A.PAYMENT_TY IN ('0','3')
								AND (C.USER_NM LIKE '%#{search}%'
									or B.GRP_NM LIKE '%#{search}%')
								ORDER BY A.REG_DATE DESC
								`

var SelectBookPaymemtList string = `SELECT A.MOID
									,DATE_FORMAT(A.REG_DATE,'%Y.%m.%d %p%h:%i')  AS REG_DATE
									,A.grp_id
									,B.grp_nm
									,A.PAYMENT_TY
									,A.CREDIT_AMT
									,DATE_FORMAT(DATE_ADD(A.REG_DATE,INTERVAL 3 DAY),'%Y.%m.%d')   AS INPUT_DATE
									,A.PAY_CHANNEL
									,DATE_FORMAT(A.ACC_ST_DAY,'%Y.%m.%d')  AS ACC_ST_DAY
								FROM dar_payment_hist AS A
								INNER JOIN priv_grp_info AS B ON A.grp_id = B.grp_id AND A.PAY_CHANNEL<>'05'
								INNER JOIN priv_user_info AS C ON A.USER_ID = C.USER_ID
								WHERE 
								A.REST_ID='#{storeId}'
								AND LEFT(A.REG_DATE,8) >= '#{searchDate}'
								AND LEFT(A.REG_DATE,8) >= '#{startDate}'
								AND LEFT(A.REG_DATE,8) <= '#{endDate}' 
								AND A.PAY_CHANNEL ='#{payChannel}' 
								AND (C.USER_NM LIKE '%#{search}%'
									or B.GRP_NM LIKE '%#{search}%')
								ORDER BY A.REG_DATE DESC
								`

var SelectBookPaymemtTotal string = `SELECT 
									SUM(A.CREDIT_AMT) as totalAmt
									FROM dar_payment_hist AS A
									INNER JOIN priv_grp_info AS B ON A.grp_id = B.grp_id AND A.PAY_CHANNEL<>'05'
									INNER JOIN priv_user_info AS C ON A.USER_ID = C.USER_ID
									
									WHERE 
									A.REST_ID='#{storeId}'
									AND LEFT(A.REG_DATE,8) >= '#{searchDate}'
									AND LEFT(A.REG_DATE,8) >= '#{startDate}'
									AND LEFT(A.REG_DATE,8) <= '#{endDate}' 
									AND A.PAY_CHANNEL ='#{payChannel}' 
									AND A.PAYMENT_TY IN ('0','3')
									AND (C.USER_NM LIKE '%#{search}%'
										or B.GRP_NM LIKE '%#{search}%')
									ORDER BY A.REG_DATE DESC
									`

var UpdateBookUserSupportBalance string = `UPDATE priv_grp_user_info SET  
								SUPPORT_BALANCE = SUPPORT_BALANCE + #{orderAmt},
								MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
								GRP_ID = '#{bookId}'
								AND USER_ID ='#{userId}'`


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