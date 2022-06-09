package payments

var InsertPrepaid string = `INSERT INTO dar_prepaid_info
							(
								PREPAID_NO
								, GRP_ID
								, REST_ID
								, JOB_TY
								, PREPAID_AMT
								, REG_DATE
							)
							VALUES
							(
								CONCAT(DATE_FORMAT(NOW(), '%Y%m%d%H%i%s'), '#{bookId}' ) 
								, '#{bookId}'
								, '#{storeId}'
								, '#{jobTy}'
								, '#{prepaidAmt}'
								, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							)	
							`


var InsertPaymentHistory string = `INSERT INTO dar_payment_hist
									(
										HIST_ID
										, REST_ID
										, GRP_ID
										, USER_ID
										, CREDIT_AMT
										, USER_TY
										, SEARCH_TY
										, PAYMENT_TY
										 ,PAY_INFO
										, REG_DATE
										, PAY_CHANNEL
										, ADD_AMT
										, MOID
										if #{accStDay} != '' then ,ACC_ST_DAY
									)
									VALUES
									(
										( SELECT FN_GET_SEQUENCE('PAYMENT_HIST_ID') AS TMP )
										, '#{storeId}'
										, '#{bookId}'
										, '#{userId}'
										, '#{creditAmt}'
										, '#{userTy}'
										, '#{searchTy}'
										, '#{paymentTy}'
										, '#{payInfo}'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										, '#{payChannel}'
										, '#{addAmt}'
										, '#{moid}'
										, '#{accStDay}'
									)`

var SelectStoreCancelCnt string = `SELECT COUNT(*) as CancelCnt
											FROM dar_payment_hist
											WHERE 
											MOID='#{moid}'
											AND PAYMENT_TY IN ('1','4')
											`

var SelectStoreChargeInfo string = `SELECT A.MOID
									,DATE_FORMAT(A.REG_DATE,'%Y.%m.%d %p%h:%i')  AS REG_DATE
									,B.GRP_NM
									,A.PAYMENT_TY
									,A.CREDIT_AMT
									,A.ADD_AMT
									,A.CREDIT_AMT + A.ADD_AMT AS TOTAL_AMT
									,IFNULL(DATE_ADD(C.SETTLMNT_DT, INTERVAL 10 DAY),'') AS expectInDate
									,A.USER_ID
									,A.PAY_INFO
									,A.GRP_ID AS BOOK_ID
									,ACC_ST_DAY
								FROM dar_payment_hist AS A
								INNER JOIN priv_grp_info AS B ON A.GRP_ID = B.GRP_ID AND PAYMENT_TY in ('0','3')
								LEFT OUTER JOIN dar_rest_payment AS C ON A.REST_PAYMENT_ID = C.REST_PAYMENT_ID
								WHERE 
								A.MOID='#{moid}'
								AND A.REST_ID='#{storeId}'
								`