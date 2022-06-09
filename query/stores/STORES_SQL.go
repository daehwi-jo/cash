package stores

var SelectBizNumCheck = `SELECT COUNT(*) as bizCnt
						FROM priv_rest_info
						WHERE 
							BUSID='#{bizNum}'
						`
var SelectStoreSeq string = `SELECT CONCAT('S',IFNULL(LPAD(MAX(SUBSTRING(REST_ID, -10)) + 1, 10, 0), '0000000001')) as storeSeq
							FROM priv_rest_info
							`

var InsertStore string = `INSERT INTO priv_rest_info
							(
							  REST_ID
							, REST_NM
							, BUSID
							, CATEGORY
							, BUETY
							, ADDR
							, ADDR2
							, LAT
							, LNG
							, AUTH_STAT
							, USE_YN
							, CEO_BIRTHDAY
							, CEO_TTI
							, CEO_NM
							, TEL
							, EMAIL
							, H_CODE
							, REG_DATE
							)
							VALUES (
							'#{storeId}'
							,'#{storeNm}'
							,'#{bizNum}'
							,'#{category}'
							,'#{kind}'
							,'#{addr}'
							,'#{addr2}'
							,'#{lat}'
							,'#{lng}'
							,'1'
							,'Y'
							,'#{ceoBirthday}'
							,'#{ceoTti}'
							,'#{ceoName}'
							,'#{storeTel}'
							,''
							,'#{hCode}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')						
							)
							`

var InsertStoreUser string = `INSERT INTO priv_rest_user_info
								(REST_ID
								, USER_ID
								, REST_AUTH
								, PRINT_CON_YN
								, DAYSUM_YN
								, MONSUM_YN
								, PAYHIST_YN
								, GRPHIST_YN
								, PREPAID_YN
								, UNPAID_YN
								, ORDER_YN
								, MENY_YN
								, AGRM_YN
								, EVENT_YN
								, PUSH_YN
								, USE_YN
								, REG_DATE
								)
								VALUES (
										'#{storeId}'
										,'#{userId}'
										,'0'
										,'N'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,'Y'
										,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								)
								`
var UpdateStoreInfo string = `UPDATE priv_rest_info SET  MOD_DATE = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s') 
									,CEO_NM='#{ceoNm}'
									,CEO_BIRTHDAY ='#{ceoBirthday}'
									,CEO_TTI ='#{ceoTti}'
									,TEL ='#{storeTel}'
									,EMAIL ='#{storeEmail}'
									,BANK_CD ='#{bankCd}'
									,ACCOUNT_NO ='#{accountNo}'
									,ACCOUNT_NM = '#{accountNm}'
									,ACCOUNT_CERT_YN='Y'
								WHERE 
								REST_ID ='#{storeId}'
								`

var UpdateStore string = `UPDATE priv_rest_info SET MOD_DATE = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s') 
								,REST_NM ='#{restNm}'
								,CATEGORY ='#{category}'
								,BUETY ='#{buety}'
								,ADDR ='#{addr}'
								,ADDR2 ='#{addr2}'
								,TEL ='#{storeTel}'
								,EMAIL ='#{email}'
								,BANK_CD ='#{bankCd}'
								,ACCOUNT_NO ='#{accountNo}'
								,ACCOUNT_NM ='#{accountNm}'
								,ACCOUNT_CERT_YN ='Y'
								,CEO_BIRTHDAY ='#{ceoBirthday}'
								,CEO_TTI ='#{ceoTti}'
								,CEO_NM ='#{ceoNm}'
								,LAT ='#{lat}'
								,LNG ='#{lng}'
								,H_CODE ='#{hCode}'
								WHERE 
								REST_ID='#{storeId}'
								`



var UpdateStoreCComp string = `UPDATE cc_comp_inf SET MOD_DT = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s') 
								,COMP_NM ='#{restNm}'
								WHERE 
								REST_ID='#{storeId}'
								`
var SelectCompInfo string = `SELECT COMP_NM
								FROM cc_comp_inf
								WHERE
								REST_ID='#{storeId}'
								`


var InsertStoreFees string = `INSERT INTO DAR_REST_FEES
							(
							REST_ID
							, REST_NM
							, PAYMETHOD
							, REST_FEES
							, START_DATE
							, MEMO
							)
							VALUES
							(
							 '#{storeId}'
							, '#{storeNm}'
							, '#{payMethod}'
							, '#{restFees}'
							, DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s') 
							, '최초 기본수수료율 적용'
							)

`
var SelectCash string = `SELECT
							REST_ID AS restId
							, BIZ_NUM AS bizNum
							, COMP_NM AS compNm
							, COMP_STS_CD AS compStsCd
							, SVC_OPEN_DT AS svcOpenDt
							, LN_FIRST_YN AS lnFirstYn
							, LN_JOIN_TY AS lnJoinTy
							, LN_ID AS lnId
							, LN_PSW AS lnPsw
							, LN_JOIN_STS_CD AS lnJoinStsCd
							, HOMETAX_OPEN_DT AS hometaxOpenDt
							, HOMETAX_FIRST_YN AS hometaxFirstYn
							, HOMETAX_JOIN_TY AS hometaxJoinTy
							, HOMETAX_ID AS hometaxId
							, HOMETAX_PSW AS hometaxPsw
							, HOMETAX_JOIN_STS_CD AS hometaxJoinStsCd
							, FAIL_RSN AS failRsn
							, REG_DT AS regDt
							, SER_ID AS serId
							, PUSH_DT AS pushDt
						FROM
							cc_comp_inf
						WHERE
							REST_ID='#{storeId}'
						`

var InsertCash string = `INSERT INTO cc_comp_inf (
							REST_ID,
							BIZ_NUM,
							COMP_NM,
							COMP_STS_CD,
							if #{lnOpenDt} != '' then SVC_OPEN_DT, 
							LN_FIRST_YN, 
							if #{lnJoinTy} != '' then LN_JOIN_TY, 
							if #{lnId} != '' then LN_ID, 
							if #{lnPsw} != '' then LN_PSW, 
							if #{lnJoinStsCd} != '' then LN_JOIN_STS_CD, 
							if #{hometaxOpenDt} != '' then HOMETAX_OPEN_DT, 
							HOMETAX_FIRST_YN, 
							if #{hometaxJoinTy} != '' then HOMETAX_JOIN_TY, 
							if #{hometaxId} != '' then HOMETAX_ID, 
							if #{hometaxPsw} != '' then HOMETAX_PSW, 
							if #{hometaxJoinStsCd} != '' then HOMETAX_JOIN_STS_CD, 
							if #{failRsn} != '' then FAIL_RSN, 
							REG_DT, 
							SER_ID, 
							PUSH_DT
						)
						VALUES (
							'#{storeId}'
							,'#{bizNum}'
							,'#{compNm}'
							,'1'
							,'#{lnOpenDt}'
							,'N'
							,'#{lnJoinTy}'
							,'#{lnId}'
							,'#{lnPsw}'
							,'#{lnJoinStsCd}'
							,'#{hometaxOpenDt}'
							,'N'
							,'#{hometaxJoinTy}'
							,'#{hometaxId}'
							,'#{hometaxPsw}'
							,'#{hometaxJoinStsCd}'
							,'#{failRsn}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							,'CASH_00'
							,''
						)
						`

var UpdateCash string = `UPDATE cc_comp_inf SET MOD_DT = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s') 
							,COMP_STS_CD = '#{compStsCd}'
							,SVC_OPEN_DT ='#{lnOpenDt}'
							,LN_FIRST_YN ='#{lnFirstYn}'
							,LN_JOIN_TY ='#{lnJoinTy}'
							,LN_ID ='#{lnId}'
							,LN_PSW ='#{lnPsw}'
							,LN_JOIN_STS_CD ='#{lnJoinStsCd}'
							,HOMETAX_OPEN_DT ='#{hometaxOpenDt}'
							,HOMETAX_FIRST_YN ='#{hometaxFirstYn}'
							,HOMETAX_JOIN_TY ='#{hometaxJoinTy}'
							,HOMETAX_ID ='#{hometaxId}'
							,HOMETAX_PSW ='#{hometaxPsw}'
							,HOMETAX_JOIN_STS_CD ='#{hometaxJoinStsCd}'
							,FAIL_RSN ='#{failRsn}'
							,PUSH_DT ='#{pushDt}'
						WHERE 
							REST_ID='#{storeId}'
						`

var SelectStoreInfo string = `SELECT A.REST_NM
								,A.BUSID
								,IFNULL(A.CATEGORY,'') AS CATEGORY
								,IFNULL(B.CATEGORY_NM,'') AS  CATEGORY_NM
								, CASE A.BUETY WHEN '00' THEN '2104' 
								  		  WHEN '01' THEN '2109' 
								  		  WHEN '02' THEN '2107' 
								  		  WHEN '03' THEN '2110' 
								  		  WHEN '04' THEN '2004' 
								  		  WHEN '05' THEN '2002' 
								  		  WHEN '06' THEN '2111' 
								  		  WHEN 'CA' THEN '2111' 
								  		  WHEN '07' THEN '2199' 
								  		  WHEN '08' THEN '8205' 
								  		  WHEN '09' THEN '7103' 
										  WHEN NULL THEN '2199' 
										  WHEN '' THEN '2199' 
										  ELSE A.BUETY
								 END AS KIND
								,CASE A.BUETY WHEN '00' THEN '한식' 
									  		  WHEN '01' THEN '중식' 
									  		  WHEN '02' THEN '일식' 
									  		  WHEN '03' THEN '양식' 
									  		  WHEN '04' THEN '카페' 
									  		  WHEN '05' THEN '분식' 
									  		  WHEN '06' THEN '부페' 
									  		  WHEN 'CA' THEN '부페' 
									  		  WHEN '07' THEN '기타' 
									  		  WHEN '08' THEN '유통' 
									  		  WHEN '09' THEN '뷰티' 
											  WHEN NULL THEN '기타' 
											  WHEN '' THEN '기타' 
											  ELSE C.CODE_NM
									 END AS KIND_NM
								,A.ADDR
								,A.ADDR2
								,A.TEL
								,A.EMAIL
								,A.BANK_CD
								,A.ACCOUNT_NO
								,A.ACCOUNT_NM
								,A.CEO_BIRTHDAY
								,A.CEO_NM
								,A.LAT 
								,A.LNG 
								,A.H_CODE
								
						FROM priv_rest_info AS a
						LEFT OUTER JOIN b_category AS b ON a.category = b.CATEGORY_ID AND A.USE_YN='Y'
						LEFT OUTER JOIN b_code AS C ON A.BUETY = C.CODE_ID  AND A.USE_YN='Y'
						WHERE 
							REST_ID ='#{storeId}'
						`

var SelectStoreServiceInfo string = `SELECT IFNULL(CEO_NM,'') AS CEO_NM
										,IFNULL(CEO_BIRTHDAY,'') AS CEO_BIRTHDAY
										,IFNULL(BANK_CD,'') AS BANK_CD
										,IFNULL(ACCOUNT_NO,'') AS ACCOUNT_NO
										,IFNULL(ACCOUNT_NM,'') AS ACCOUNT_NM
										FROM priv_rest_info
										WHERE 
											rest_id='#{storeId}'
									`

var SelectStoreCategories string = `SELECT CODE_ID
									,CODE_NM
									,USE_YN
									FROM DAR_CATEGORY_INFO
									WHERE 
									REST_ID='#{storeId}'
									`

var SelectStoreCategoriesSeq string = `SELECT FN_GET_SEQUENCE('CATEGORY_ID') AS CATEGORY_ID`

var SelectStoreItemCodeSeq string = `SELECT
										IFNULL(A.CATEGORY_ID,FN_GET_SEQUENCE('CATEGORY_ID')) AS categoryId
										,B.REST_NM AS categoryNm
										,CONCAT(IFNULL(LPAD(MAX(SUBSTRING(A.CODE_ID, -5)) + 1, 5, 0), '00001')) as codeId
										FROM DAR_CATEGORY_INFO AS A
										INNER JOIN PRIV_REST_INFO AS B ON A.REST_ID = B.REST_ID
										WHERE A.REST_ID='#{storeId}'
										AND A.CODE_ID <> '99999'
										`

var InsertStoreCategories string = `INSERT INTO DAR_CATEGORY_INFO
										(
											CATEGORY_ID
											, CATEGORY_NM
											, REST_ID
											, CODE_ID
											, CODE_NM
											, USE_YN
										)
										VALUES
										(
											 '#{categoryId}'
											, '#{categoryNm}'
											, '#{storeId}'
											, '#{codeId}'
											, '#{codeNm}'
											, '#{useYn}'
										)`

var UpdateStroeCategories string = `UPDATE  DAR_CATEGORY_INFO SET  
											CODE_NM = '#{codeNm}'
											,USE_YN = '#{useYn}'
										WHERE 
										REST_ID='#{storeId}'
										AND CODE_ID='#{categoryId}'
										`

var SelectStoreMenuList string = `SELECT  B.CODE_ID AS codeId
										,B.CODE_NM AS codeNm
										,A.ITEM_NO
										,A.ITEM_NM
										,A.ITEM_PRICE
										,'' as eventStartDate 
										,'' as eventEndDate 
										,A.ITEM_PRICE as eventPrice
										,A.USE_YN as useYn
										,A.BEST_YN as bestYn
								FROM dar_sale_item_info AS A
								INNER JOIN DAR_CATEGORY_INFO AS B ON A.ITEM_MENU = B.CODE_ID AND A.REST_ID = B.REST_ID
								WHERE A.REST_ID='#{storeId}'
								`

var SelectStoreMenuSeq string = `SELECT
								IFNULL(LPAD(MAX(ITEM_NO) + 1, 10, 0), '0000000001') as itemNo
								FROM dar_sale_item_info	
								`
var InsertStroeMenu string = ` INSERT INTO dar_sale_item_info(
												ITEM_NO
												, ITEM_NM
												, REST_ID
												, ITEM_STAT
												, ITEM_PRICE
												, REG_DATE
												, ITEM_MENU
												, USE_YN
												, BEST_YN
											)
											VALUES (
												'#{itemNo}'
												,'#{itemNm}'
												,'#{storeId}'
												,'1'
												,#{itemPrice}
												,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
												,'#{codeId}'
												,'#{useYn}'
												,'#{bestYn}'
											)`

var UpdateStroeMenu string = `UPDATE  dar_sale_item_info SET  
											ITEM_NM = '#{itemNm}',
											ITEM_PRICE = #{itemPrice},
											ITEM_MENU = '#{codeId}',
											USE_YN='#{useYn}',
											BEST_YN='#{bestYn}',
											MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										WHERE 
										REST_ID = '#{storeId}'
										AND ITEM_NO='#{itemNo}'
										`

var SelectStoreServiceList string = `SELECT SERVICE_ID
										, SERVICE_NM
										, SERVICE_INFO
										, USE_YN
										FROM priv_rest_service
										WHERE 
										REST_ID = '#{storeId}'
										`

var SelectStoreService string = `SELECT SERVICE_ID
										,SERVICE_NM
										,SERVICE_INFO
								FROM priv_rest_service
								WHERE 
								USE_YN='1'
								AND REST_ID = '#{storeId}'
								`
var SelectStoreDesc string = `SELECT REST_NM
									,OPEN_WEEK
									,OPEN_WEEKEND
									,IFNULL(INTRO_SHORT,'') AS INTRO_SHORT
									,IFNULL(B.INTRO,'') AS INTRO
									,IFNULL(B.DELIVERY_YN ,'N') AS DELIVERY_YN
									,IFNULL(B.DELIVERY_LOCATION,'') AS DELIVERY_LOCATION
									FROM priv_rest_info AS A
									LEFT OUTER JOIN priv_rest_etc AS B ON A.REST_ID = B.REST_ID
									WHERE 
									A.REST_ID= '#{storeId}'
								`

var InsertStoreBaseService string = ` 
									INSERT INTO priv_rest_service(REST_ID, SERVICE_ID, SERVICE_NM, SERVICE_INFO, USE_YN)
									SELECT 
										'#{storeId}'
										, SERVICE_ID
										, SERVICE_NM
										, SERVICE_INFO
										, USE_YN
									FROM priv_rest_service
									WHERE 
									REST_ID='R0000000000'
									`

var UpdateStoreService string = `UPDATE priv_rest_service SET 
										 SERVICE_INFO  = '#{serviceInfo}'
										, USE_YN  = '#{useYn}'
										, MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										WHERE 
											REST_ID='#{storeId}'
											AND SERVICE_ID ='#{serviceId}'
											`

var SelectRestUserInfo string = `SELECT C.USER_NM
								,C.HP_NO
								,C.USER_ID
								,C.LOGIN_ID
								FROM priv_rest_info AS A
								INNER JOIN priv_rest_user_info AS B ON A.REST_ID = B.REST_ID AND B.REST_AUTH='0'
								INNER JOIN priv_user_info AS C ON B.USER_ID = C.USER_ID
								WHERE 
									A.REST_ID='#{storeId}'
											`

var SelectSuccessReportSendCheck string = `
							SELECT result 
							FROM sys_alimtalk_log  
							WHERE 
								template_code = '#{code}' 
							AND 
								left(send_date,8) = '#{today}' 
							AND 
								user_id='#{userId}'
`


var SelectStoreChargeList string = `SELECT SEQ_NO
											,AMT
											,ADD_AMT
											,USE_YN
									FROM dar_prepayment_charge_info
									WHERE 
									REST_ID='#{storeId}'
									order by AMT asc
						`


var SelectChargeBase string = `SELECT SEQ_NO
											,AMT
											,ADD_AMT
									FROM dar_prepayment_charge_info
									WHERE 
									REST_ID='R0000000000'
									AND USE_YN='Y'
						`

var InsertStoreChargeList string = `INSERT INTO dar_prepayment_charge_info
									(
									 SEQ_NO
									,REST_ID
									,AMT
									,ADD_AMT
									,USE_YN
									)
									VALUES (
									'#{storeSeqNo}'
									,'#{storeId}'
									,#{amt}
									,0
									,'N'
									)
								`




var UpdateStoreCharge string = `UPDATE dar_prepayment_charge_info SET ADD_AMT=#{addAmt}
																	,USE_YN= '#{useYn}'
								WHERE 
									REST_ID='#{storeId}'
									AND SEQ_NO='#{seqNo}'
								`


var UpdateStoreChargeYn string = `UPDATE priv_rest_info SET 
															PAYMENT_USE_YN='#{paymentUseYn}'
															, MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
									REST_ID='#{storeId}'
								`

var SelectAlimTalkInfo string = `SELECT KAKAO_WEEK
									, KAKAO_MONTH
									, KAKAO_DAILY
									FROM priv_rest_etc
									WHERE
										REST_ID='#{storeId}'
				           `



var InsertStoreAlimtalk string = `INSERT INTO  priv_rest_etc (
										REST_ID
									   ,KAKAO_WEEK
									   ,KAKAO_MONTH
									   ,KAKAO_DAILY
									)
									VALUES(
									    '#{storeId}'
									    ,'#{kakaoWeek}'
									    ,'#{kakaoMonth}'
									    ,'#{kakaoDaily}'
									 )
											`




var SelectStorePaymentInfo string = `SELECT PAYMENT_USE_YN as paymentUseYn
									FROM priv_rest_info
									WHERE
										rest_id='#{storeId}'
											`






var SelectOrderPickupCount string = `SELECT SUM(CASE P_STATUS WHEN '30'  THEN 1 ELSE 0 END) AS WATING
										,SUM(CASE P_STATUS WHEN '32'  THEN 1 ELSE 0 END) AS RECEIPT
										,SUM(CASE P_STATUS WHEN '34'  THEN 1 ELSE 0 END) AS COOK
										,SUM(CASE P_STATUS WHEN '36'  THEN 1 ELSE 0 END) AS PICK
								FROM dar_order_info AS a
								INNER JOIN dar_order_pickup AS B ON A.ORDER_NO = B.ORDER_NO
								WHERE 
								A.REST_ID='#{restId}'
								AND LEFT(A.ORDER_DATE,8)=DATE_FORMAT(NOW(), '%Y%m%d')
								`

var SelectOrderPickupList string = `SELECT A.ORDER_NO
										,C.USER_NM
										,DATE_FORMAT(A.ORDER_DATE,'%p %H:%i') AS ORDER_TIME
										,DATE_FORMAT(A.ORDER_DATE,'%Y년 %m월 %d일 %p %H:%i') AS ORDER_DATE
										,B.P_STATUS
									   ,CASE P_STATUS WHEN '30' THEN '접수대기'
													  WHEN '32' THEN '주문접수'
													  WHEN '34' THEN '준비완료'
													  ELSE P_STATUS END AS P_STATUS_NM
										,(SELECT SUM(ORDER_QTY) FROM dar_order_detail AS AA WHERE A.ORDER_NO = AA.ORDER_NO) AS TOTAL_ORDER_QTY
								FROM dar_order_info  AS A
								INNER JOIN dar_order_pickup AS B ON A.ORDER_NO = B.ORDER_NO
								INNER JOIN priv_user_info AS C ON A.USER_ID=C.USER_ID
								WHERE 
								A.REST_ID='#{restId}'
								AND A.ORDER_STAT='20'
								AND P_STATUS IN ('30','32','34')
								AND LEFT(A.ORDER_DATE,8)=DATE_FORMAT(NOW(), '%Y%m%d')
								ORDER BY A.ORDER_DATE DESC
								`


var SelectOrderPickupMenuList string = `SELECT SUM(A.ORDER_QTY) AS ORDER_QTY 
										,CASE A.ITEM_NO WHEN '9999999999' THEN '금액권' ELSE B.ITEM_NM END AS ITEM_NM
								FROM dar_order_detail AS A
								LEFT OUTER JOIN dar_sale_item_info AS B ON A.ITEM_NO = B.ITEM_NO
								WHERE 
								A.ORDER_NO= '#{orderNo}'
								GROUP BY A.ITEM_NO
								`

var UpdatePickupStatus32 string = `UPDATE dar_order_pickup SET P_STATUS = '32'
														,RECEIPT_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
								    ORDER_NO= '#{orderNo}'
								`

var UpdatePickupStatus34 string = `UPDATE dar_order_pickup SET P_STATUS = '34'
														,COOK_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
									ORDER_NO= '#{orderNo}'
								`

var UpdatePickupStatus36 string = `UPDATE dar_order_pickup SET P_STATUS = '36'
														,PICK_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
								WHERE 
									ORDER_NO= '#{orderNo}'
								`



var SelectPickupOrderUserInfo string = `SELECT 
										 A.ORDER_NO
										,A.USER_ID
										,C.REST_NM
									FROM dar_order_info  AS A
									INNER JOIN dar_order_pickup AS B ON A.ORDER_NO = B.ORDER_NO
									INNER JOIN PRIV_REST_INFO AS C ON A.REST_ID = C.REST_ID
									WHERE
									A.ORDER_NO='#{orderNo}'
									AND A.ORDER_STAT='20'
									AND P_STATUS IN ('30','32','34')
								`

