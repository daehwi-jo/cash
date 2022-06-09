package orders


var SelectUnpaidListCount string = `SELECT COUNT(*) AS orderCnt
								, SUM(total_amt) AS TOTAL_AMT
								,FN_GET_GRPNAME(A.GRP_ID) AS BOOK_NM
								,(SELECT USER_ID FROM PRIV_GRP_USER_INFO AS AA WHERE AA.GRP_ID=A.GRP_ID AND AA.GRP_AUTH='0') AS USER_ID
								FROM  dar_order_info AS A
								INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID
								WHERE 
								A.REST_ID='#{storeId}'
								AND A.order_ty IN ('1','2','3','5')
								AND A.GRP_ID ='#{bookId}'
								AND PAY_TY='1' 
								AND PAID_YN='N'
								AND order_stat = '20'
								AND DATE_FORMAT(A.ORDER_DATE,'%Y-%m-%d') <='#{accStDay}'
								`

var SelectUnpaidList string = `SELECT DATE_FORMAT(A.ORDER_DATE,'%Y.%m.%d %p%h:%i')  AS ORDER_DATE
								,ORDER_TY
								,B.USER_NM AS orderer
								,CASE WHEN A.ORDER_TY ='1' THEN 'pay'
																			WHEN A.ORDER_TY ='2' THEN 'delivery'
																			WHEN A.ORDER_TY ='3' THEN 'takeout'
																			WHEN A.ORDER_TY ='5' THEN ''
																			END AS ORDER_TY
								,A.TOTAL_AMT
								FROM  dar_order_info AS A
								INNER JOIN priv_user_info AS B ON A.USER_ID = B.USER_ID
								WHERE 
								A.REST_ID='#{storeId}'
								AND A.order_ty IN ('1','2','3','5')
								AND A.GRP_ID ='#{bookId}'
								AND PAY_TY='1' 
								AND PAID_YN ='N'
								AND order_stat = '20'
								AND DATE_FORMAT(A.ORDER_DATE,'%Y-%m-%d') <='#{accStDay}'
								ORDER BY ORDER_DATE DESC
								`

var UpdateOrderPaid string = `UPDATE dar_order_info SET PAID_YN='Y'
												,moid = '#{moid}'
												,PAY_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE
								REST_ID='#{storeId}'
								AND GRP_ID ='#{bookId}'
								AND PAY_TY='1' 
								AND PAID_YN='N'
								AND DATE_FORMAT(ORDER_DATE,'%Y-%m-%d') <='#{accStDay}'
								`


var UpdateOrderPaidCancel string = `UPDATE dar_order_info SET PAID_YN='N'
								,moid = NULL
								,PAY_DATE = NULL
				WHERE
				moid= '#{moid}'
				`


var SelectOrder string = `SELECT A.ORDER_NO
							,A.TOTAL_AMT
							,A.ORDER_STAT
							,A.PAY_TY
							,A.GRP_ID AS BOOK_ID
							,A.REST_ID AS STORE_ID
							,A.USER_ID
							,A.POINT_USE
					FROM dar_order_info AS A
					WHERE 
					A.ORDER_NO = '#{orderNo}'
					`

var UpdateOrderCancel string = ` UPDATE dar_order_info SET ORDER_STAT='21'
												, ORDER_CANCEL_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE
								ORDER_NO = '#{orderNo}'
									`



var SelectOrderInfo string = `SELECT A.ORDER_NO
											,B.REST_NM
											,C.GRP_NM
											,A.TOTAL_AMT
											,A.ORDER_STAT
											,DATE_FORMAT(A.ORDER_DATE,'%Y.%m.%d %p%h:%i') AS ORDER_DATE
											,ifnull(A.ORDER_COMMENT,'') AS ORDER_COMMENT
									FROM dar_order_info AS A
									INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
									INNER JOIN priv_grp_info AS C ON A.GRP_ID = C.GRP_ID
									WHERE 
									A.ORDER_NO = '#{orderNo}'
									`


var SelectOrderDetail string = `SELECT CASE WHEN  A.ITEM_NO='9999999999'  THEN '금액권'  ELSE B.ITEM_NM END  as menuNm
										 ,SUM(ORDER_AMT* ORDER_QTY) as menuPrice
										 ,SUM(ORDER_QTY) as menuQty
								FROM dar_order_detail AS A 
								LEFT OUTER JOIN dar_sale_item_info AS B ON A.ITEM_NO = B.ITEM_NO
								WHERE 
								A.ORDER_NO= '#{orderNo}'
                                GROUP BY A.ITEM_NO 
								`

