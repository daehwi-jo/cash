package homes

var SelectSuccessStateReport string = `
									select a.user_id, c.rest_id, a.user_nm, a.join_date, a.login_id, a.hp_no, a.kakao_key, c.rest_nm
									from priv_user_info a
									left outer join priv_rest_user_info b on a.user_id = b.user_id and b.rest_auth = '0'
									left outer join priv_rest_info c on b.rest_id = c.rest_id 
									where a.user_ty = '1' and a.hp_no not in ('01033841751', '01093939899', '01033841753', '01023847287','01032961793','01027709407') 
									and a.user_id in (
										select user_id from priv_rest_user_info 
										where rest_auth = '0'
										and rest_id in (
											select rest_id from priv_rest_info
											where busid is not null and busid not like '123%'
											and use_yn = 'Y' and ceo_nm is not null and ceo_nm != '' and reg_date >= '2019'
											and rest_id in (
												select rest_id from cc_comp_inf
												where LN_FIRST_YN = 'Y'))
											)
`

var SelectLoginFailReport string = `
									SELECT rest_id FROM cc_comp_inf WHERE biz_num IN(
										SELECT biz_num FROM cc_sync_inf 
										WHERE 
											bs_dt = '#{bsDt}' 
										AND err_cd = '0005'
									)
`

var SelectFailStateReport string = `
									select a.user_id,c.rest_id, a.user_nm, a.join_date, a.login_id, a.hp_no, a.kakao_key, c.rest_nm
									from priv_user_info a
									left outer join priv_rest_user_info b on a.user_id = b.user_id and b.rest_auth = '0'
									left outer join priv_rest_info c on b.rest_id = c.rest_id
									where a.user_ty = '1' and a.hp_no not in ('01033841751', '01093939899', '01033841753', '01023847287','01032961793','01027709407') 
									and a.user_id in (
										select user_id from priv_rest_user_info 
										where rest_auth = '0'
										and rest_id in (
											select rest_id from priv_rest_info
											where busid is not null and busid not like '123%'
											and use_yn = 'Y' and ceo_nm is not null and ceo_nm != '' and reg_date >= '2019'
											and rest_id not in (
												select rest_id from cc_comp_inf
												where LN_FIRST_YN = 'Y')
											)
										)
`

var SelectIntroMsg string = `SELECT CODE_NM as introMsg
								FROM b_code
								WHERE CATEGORY_ID='M001'
								order by RAND()
								LIMIT 1
								`

var SelectStoreService string = `SELECT IFNULL(B.B_ID,'N') AS billingYn
										,A.ADDR
										,A.REST_NM
										,AA.REST_AUTH
										,CASE  WHEN LN_JOIN_STS_CD ='1' THEN  IFNULL(C.LN_ID,'N') ELSE 'N' END AS cardSalesYn
										,CASE  WHEN HOMETAX_JOIN_STS_CD ='1' THEN  IFNULL(C.HOMETAX_ID,'N') ELSE 'N' END AS homeTaxYn
								FROM priv_rest_info AS A
								INNER JOIN priv_rest_user_info AS AA ON A.REST_ID = AA.REST_ID  
								LEFT OUTER JOIN e_billing AS B ON A.REST_ID = B.STORE_ID  AND B.END_DATE >= SYSDATE()
								LEFT OUTER JOIN cc_comp_inf AS C ON A.REST_ID = C.REST_ID
								WHERE 
								A.REST_ID= '#{storeId}'
								AND AA.USER_ID= '#{userId}'
								`
var SelectStoreBillingInfo string = `SELECT DATE_FORMAT(A.END_DATE,'%Y년 %m월 %d일') AS END_DATE
									,B.ITEM_NAME
									,B.ITEM_DESC
									,B.CHANNEL_DIV
									,DATE_FORMAT(A.NEXT_PAY_DAY,'%Y년 %m월 %d일') AS NEXT_PAY_DAY
									,A.PAY_YN
									FROM e_billing AS A
									INNER JOIN e_billing_item AS B ON A.ITEM_CODE = B.ITEM_CODE
									WHERE 
									STORE_ID= '#{storeId}'
									AND END_DATE >= SYSDATE()
								`
// 캐시컴바인 가입일 조회
var SelectCashRegInfo string = `
						SELECT 
							LEFT(a.reg_date, 8) AS regDt, 
							IFNULL(b.comp_sts_cd, "") AS compStsCd, 
							IFNULL(b.svc_open_dt, "") AS lnOpenDt, 
							IFNULL(b.ln_join_sts_cd, "") AS lnJoinStsDt, 
							IFNULL(b.hometax_open_dt, "") AS hometaxOpenDt,
							IFNULL(b.hometax_join_sts_cd, "") AS hometaxJoinStsDt,
							a.rest_Id as restId
						FROM 
							priv_rest_info a LEFT JOIN cc_comp_inf b ON a.rest_id = b.rest_id
						WHERE 
							TRIM(a.busid) = '#{bizNum}'
						`

var SelectRestUserInfo string =`
						SELECT 
							USER_ID as userId
						FROM 
							priv_rest_user_info
						WHERE
						REST_AUTH = 0
						AND
							REST_ID = '#{restId}'
`

var SelectUserInfo string =`
						SELECT 
							USER_BIRTH as birth,
							HP_NO as hp,
							login_id as id,
							user_nm as name,
							DATE_FORMAT(join_date,'%Y년 %c월 %e일') AS jdate
						FROM 
							priv_user_info
						WHERE 
							USER_ID = '#{userId}'
`

// 매출확인 및 매출 예측
var SelectDaySaleData string = `
							SELECT 
								tr_dt AS trDt,
								real_amt AS realAmt, 
								expect_amt AS expectAmt 
							FROM cc_day_sale_sum 
							WHERE 
								biz_num = '#{bizNum}' 
								AND tr_dt = '#{trDt}'
							`
var SelectDayName string = `
							SELECT ifnull(datename,"") AS datename 
							FROM 
								sys_week_date 
							WHERE 
								total_date='#{trDt}'
                         `

// 오늘 입금 예정 금액
var SelectTodayPay string = `
							SELECT ifnull(SUM(pay_amt),0) AS amt 
							FROM cc_pca_dtl 
							WHERE 
								biz_num='#{bizNum}' 
								AND outp_expt_dt='#{trDt}'
`

// 월 카드 누적 매출
var SelectSaleSum string = `
							SELECT 
								IFNULL(SUM(aprv_amt),0) AS aprvSum
							FROM cc_aprv_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND LEFT(tr_dt,6) = LEFT('#{trDt}',6)
							`

// 월 현금영수증 누적 매출
var SelectCashSum string = `
							SELECT 
								IFNULL(SUM(tot_amt),0) AS aprvSum
							FROM cc_cash_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND LEFT(tr_dt,6) = LEFT('#{trDt}',6)
							`

// 월 누적 입금
var SelectPaySum string = `
							SELECT 
								IFNULL(SUM(REAL_PAY_AMT),0) AS paySum
							FROM cc_pay_dtl
							WHERE
							biz_num = '#{bizNum}'
							AND LEFT(pay_dt,6) = LEFT('#{trDt}',6)
							`

// 카드 승인금액
var SelectCardAmt string = `
							SELECT 
								IFNULL(SUM(aprv_amt), 0) AS aprvAmt
							FROM cc_aprv_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND tr_dt = '#{trDt}'
							`

// 현금영수증 승인금액
var SelectCashAmt string = `
							SELECT
								IFNULL(SUM(tot_amt), 0) AS cashAmt
							FROM cc_cash_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND tr_dt = '#{trDt}'
							`

// 매입 금액
var SelectPcaAmt string = `
							SELECT 
								IFNULL(SUM(pca_amt),0) AS pcaAmt
							FROM cc_pca_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND pca_dt = '#{trDt}'
							`

// 매입 예정 금액
var SelectPrePcaAmt string = `
							SELECT 
								IFNULL(SUM(PAY_AMT),0) AS pcaAmt
							FROM cc_pca_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND OUTP_EXPT_DT = '#{trDt}'
							`

// 입금 금액
var SelectPayAmt string = `
							SELECT
								IFNULL(SUM(REAL_PAY_AMT),0) AS payAmt
							FROM cc_pay_dtl
							WHERE 
								biz_num = '#{bizNum}'
								AND pay_dt = '#{trDt}'
							`

// 단골 비율 예측
var SelectExpectPersonVisit string = `
									SELECT
										SUM(z.visit_person) AS visitTotal,
										SUM(IF(z.visitCnt >= 2, z.visit_person, 0)) AS visit2
									FROM
									(
										SELECT y.visit_cnt AS visitCnt,
											   COUNT(y.visit_cnt) AS visit_person
										FROM
										(
											SELECT COUNT(*) AS visit_cnt
											FROM cc_aprv_dtl
											WHERE 
												biz_num = '#{bizNum}'
												and tr_dt between '#{startDt}' 
												AND '#{endDt}'
												AND DAYOFWEEK(tr_dt) = DAYOFWEEK(DATE_FORMAT(NOW(), '%Y%m%d'))
												GROUP BY CARD_NO
										) y
										GROUP BY y.visit_cnt
									) z
									`
// 결제 건수 예측
var SelectExpectCnt string = `
							SELECT round(COUNT(*)/28,0) as avg_cnt 
							FROM cc_aprv_dtl 
							WHERE
								biz_num = '#{bizNum}'
							and tr_dt 
								between '#{startDt}'
								AND '#{endDt}' 		
`
// 바쁜시간 예측
var SelectExpectBusyTime string = `
									SELECT 
										SUM(z.h0003) AS h0003,
										SUM(z.h0306) AS h0306,
										SUM(z.h0609) AS h0609,
										SUM(z.h0912) AS h0912,
										SUM(z.h1215) AS h1215,
										SUM(z.h1518) AS h1518,
										SUM(z.h1821) AS h1821,
										SUM(z.h2124) AS h2124
									FROM (
										SELECT 
											SUM(IF(x.tr_hr >= '00' AND x.tr_hr < '03', x.amt, 0)) AS h0003,
											SUM(IF(x.tr_hr >= '03' AND x.tr_hr < '06', x.amt, 0)) AS h0306,
											SUM(IF(x.tr_hr >= '06' AND x.tr_hr < '09', x.amt, 0)) AS h0609,
											SUM(IF(x.tr_hr >= '09' AND x.tr_hr < '12', x.amt, 0)) AS h0912,
											SUM(IF(x.tr_hr >= '12' AND x.tr_hr < '15', x.amt, 0)) AS h1215,
											SUM(IF(x.tr_hr >= '15' AND x.tr_hr < '18', x.amt, 0)) AS h1518,
											SUM(IF(x.tr_hr >= '18' AND x.tr_hr < '21', x.amt, 0)) AS h1821,
											SUM(IF(x.tr_hr >= '21', x.amt, 0)) AS h2124
										FROM
										(
											SELECT 
												LEFT(tr_tm,2) AS tr_hr, 
												SUM(aprv_amt) AS amt
											FROM cc_aprv_dtl
											WHERE
												biz_num = '#{bizNum}'
												AND tr_dt BETWEEN '#{startDt}' 
												AND '#{endDt}'
												AND DAYOFWEEK(tr_dt) = DAYOFWEEK(DATE_FORMAT(NOW(), '%Y%m%d'))
											GROUP BY LEFT(tr_tm,2)
										) x
										UNION ALL
										SELECT 
											SUM(IF(y.tr_hr >= '00' AND y.tr_hr < '03', y.amt, 0)) AS h0003,
											SUM(IF(y.tr_hr >= '03' AND y.tr_hr < '06', y.amt, 0)) AS h0306,
											SUM(IF(y.tr_hr >= '06' AND y.tr_hr < '09', y.amt, 0)) AS h0609,
											SUM(IF(y.tr_hr >= '09' AND y.tr_hr < '12', y.amt, 0)) AS h0912,
											SUM(IF(y.tr_hr >= '12' AND y.tr_hr < '15', y.amt, 0)) AS h1215,
											SUM(IF(y.tr_hr >= '15' AND y.tr_hr < '18', y.amt, 0)) AS h1518,
											SUM(IF(y.tr_hr >= '18' AND y.tr_hr < '21', y.amt, 0)) AS h1821,
											SUM(IF(y.tr_hr >= '21', y.amt, 0)) AS h2124
										FROM
										(
											SELECT 
												LEFT(tr_tm,2) AS tr_hr, 
												SUM(tot_amt) AS amt
											FROM cc_cash_dtl
											WHERE
												biz_num = '#{bizNum}'
												AND tr_dt BETWEEN '#{startDt}' 
												AND '#{endDt}'
												AND DAYOFWEEK(tr_dt) = DAYOFWEEK(DATE_FORMAT(NOW(), '%Y%m%d'))
											GROUP BY LEFT(tr_tm,2)
										) y 
									) z
									`

var SelectExpectBusyTimeTomorrow string = `
									SELECT 
										SUM(z.h0003) AS h0003,
										SUM(z.h0306) AS h0306,
										SUM(z.h0609) AS h0609,
										SUM(z.h0912) AS h0912,
										SUM(z.h1215) AS h1215,
										SUM(z.h1518) AS h1518,
										SUM(z.h1821) AS h1821,
										SUM(z.h2124) AS h2124
									FROM (
										SELECT 
											SUM(IF(x.tr_hr >= '00' AND x.tr_hr < '03', x.amt, 0)) AS h0003,
											SUM(IF(x.tr_hr >= '03' AND x.tr_hr < '06', x.amt, 0)) AS h0306,
											SUM(IF(x.tr_hr >= '06' AND x.tr_hr < '09', x.amt, 0)) AS h0609,
											SUM(IF(x.tr_hr >= '09' AND x.tr_hr < '12', x.amt, 0)) AS h0912,
											SUM(IF(x.tr_hr >= '12' AND x.tr_hr < '15', x.amt, 0)) AS h1215,
											SUM(IF(x.tr_hr >= '15' AND x.tr_hr < '18', x.amt, 0)) AS h1518,
											SUM(IF(x.tr_hr >= '18' AND x.tr_hr < '21', x.amt, 0)) AS h1821,
											SUM(IF(x.tr_hr >= '21', x.amt, 0)) AS h2124
										FROM
										(
											SELECT 
												LEFT(tr_tm,2) AS tr_hr, 
												SUM(aprv_amt) AS amt
											FROM cc_aprv_dtl
											WHERE
												biz_num = '#{bizNum}'
												AND tr_dt BETWEEN '#{startDt}' 
												AND '#{endDt}'
												AND DAYOFWEEK(tr_dt) = DAYOFWEEK(DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 1 DAY), '%Y%m%d'))
											GROUP BY LEFT(tr_tm,2)
										) x
										UNION ALL
										SELECT 
											SUM(IF(y.tr_hr >= '00' AND y.tr_hr < '03', y.amt, 0)) AS h0003,
											SUM(IF(y.tr_hr >= '03' AND y.tr_hr < '06', y.amt, 0)) AS h0306,
											SUM(IF(y.tr_hr >= '06' AND y.tr_hr < '09', y.amt, 0)) AS h0609,
											SUM(IF(y.tr_hr >= '09' AND y.tr_hr < '12', y.amt, 0)) AS h0912,
											SUM(IF(y.tr_hr >= '12' AND y.tr_hr < '15', y.amt, 0)) AS h1215,
											SUM(IF(y.tr_hr >= '15' AND y.tr_hr < '18', y.amt, 0)) AS h1518,
											SUM(IF(y.tr_hr >= '18' AND y.tr_hr < '21', y.amt, 0)) AS h1821,
											SUM(IF(y.tr_hr >= '21', y.amt, 0)) AS h2124
										FROM
										(
											SELECT 
												LEFT(tr_tm,2) AS tr_hr, 
												SUM(tot_amt) AS amt
											FROM cc_cash_dtl
											WHERE
												biz_num = '#{bizNum}'
												AND tr_dt BETWEEN '#{startDt}' 
												AND '#{endDt}'
												AND DAYOFWEEK(tr_dt) = DAYOFWEEK(DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 1 DAY), '%Y%m%d'))
											GROUP BY LEFT(tr_tm,2)
										) y 
									) z
									`