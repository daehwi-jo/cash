package datas

// 가입일 조회
var SelectRegistDate string = `
							SELECT 
								LEFT(reg_dt, 8) AS regDt,
								REST_ID as restId,
								comp_nm as compNm
							FROM cc_comp_inf 
							WHERE 
								biz_num = '#{bizNum}'
							`

var SelectRegistDateRestId string = `
							SELECT 
								LEFT(reg_dt, 8) AS regDt,
								REST_ID as restId,
								comp_nm as compNm
							FROM cc_comp_inf 
							WHERE 
								rest_id = '#{restId}'
							`

// 카드 조회
var SelectMonthCard string = `
							SELECT card_cd, COUNT(*) AS cnt, card_nm 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num ='#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{endDt}' 
							GROUP BY card_cd ORDER BY COUNT(*) desc
`

var SelectYesterdayAmt string = `
							SELECT tot_amt as amt 
							FROM cc_aprv_sum
							WHERE 
								biz_num ='#{bizNum}' 
							AND 
								bs_dt = '#{bsDt}' 
`

// 카드 조회
var SelectWeekCard string = `
							SELECT card_cd, COUNT(*) AS cnt, card_nm 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num ='#{bizNum}' 
							AND bs_dt
							BETWEEN 
								'#{startDt}' 
							AND 
								'#{endDt}'   
							GROUP BY card_cd ORDER BY COUNT(*) desc
`

// 달아요 사용액
var SelectMonthDarayoAmt string = `
							SELECT ifnull(SUM(total_amt),0) AS tot_amt 
							FROM dar_order_info 
							WHERE 
								rest_id = '#{restId}'
							AND 
								LEFT(order_date,6) = '#{endDt}'
							AND order_stat = '20'
`

// 달아요 사용액
var SelectWeekDarayoAmt string = `
							SELECT ifnull(SUM(total_amt),0) AS tot_amt 
							FROM dar_order_info 
							WHERE 
								rest_id = '#{restId}'
							AND order_date
							BETWEEN 
								'#{startDt}' 
							AND 
								'#{endDt}'
							AND order_stat = '20'
`

// 달아요 사용액
var SelectYesterdayDarayoAmt string = `
							SELECT ifnull(SUM(total_amt),0) AS tot_amt 
							FROM dar_order_info 
							WHERE 
								rest_id = '#{restId}'
							AND order_date = '#{bsDt}'
							AND order_stat = '20'
`

// 지난달 입금 정산
var SelectMonthPayAmt string = `
							SELECT 
										tr_month AS trMonth, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt,
										SUM(z.real_in_amt) - SUM(Z.outp_expt_amt) AS diffAmt
									FROM (
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											IFNULL(SUM(pay_amt),0) AS outp_expt_amt,
											0 AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pca_dtl ON dt = outp_expt_dt AND biz_num = '#{bizNum}'
										WHERE 
											left(dt,6) = '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
										
										UNION ALL
										
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											0 AS outp_expt_amt,
											IFNULL(SUM(real_pay_amt),0) AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pay_dtl ON dt = pay_dt AND biz_num = '#{bizNum}'
										WHERE 
											left(dt,6) = '#{endDt}' 
										GROUP BY SUBSTR(dt, 1, 6)
									) z
`

// 지난주 입금 정산
var SelectWeekPayAmt string = `
							SELECT 
										tr_month AS trMonth, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt,
										SUM(z.real_in_amt) - SUM(Z.outp_expt_amt) AS diffAmt
									FROM (
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											IFNULL(SUM(pay_amt),0) AS outp_expt_amt,
											0 AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pca_dtl ON dt = outp_expt_dt AND biz_num = '#{bizNum}'
										WHERE 
											dt
											BETWEEN 
												'#{startDt}' 
											AND 
												'#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
										
										UNION ALL
										
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											0 AS outp_expt_amt,
											IFNULL(SUM(real_pay_amt),0) AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pay_dtl ON dt = pay_dt AND biz_num = '#{bizNum}'
										WHERE 
											dt
											BETWEEN 
												'#{startDt}' 
											AND 
												'#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									) z
`

// 어제 입금 정산
var SelectYesterdayPayAmt string = `
							SELECT 
										tr_month AS trMonth, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt,
										SUM(z.real_in_amt) - SUM(Z.outp_expt_amt) AS diffAmt
									FROM (
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											IFNULL(SUM(pay_amt),0) AS outp_expt_amt,
											0 AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pca_dtl ON dt = outp_expt_dt 
											AND 
												biz_num = '#{bizNum}'
										WHERE 
											dt = '#{bsDt}'
										GROUP BY SUBSTR(dt, 1, 6)
										
										UNION ALL
										
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											0 AS outp_expt_amt,
											IFNULL(SUM(real_pay_amt),0) AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pay_dtl ON dt = pay_dt 
											AND 
												biz_num = '#{bizNum}'
										WHERE 
											dt = '#{bsDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									) z
`

// 금결원 가입여부
var SelectKFTCEnroll string = `
							SELECT kftc_sts_cd AS kftcStsCd
							FROM
								cc_kftc_comp_inf
							WHERE
								biz_num = '#{bizNum}'
							`

// 배달업체 정보
var SelectDeliveryId string = `
							SELECT baemin_id, yogiyo_id, naver_id 
							FROM 
								b_store 
							WHERE 
								biz_num = '#{bizNum}'
							`

var SelectWebViewMonth string = `
							SELECT title, content, position 
							FROM 
								b_web_view 
							WHERE 
								DATE = '#{bsDt}'
							AND TYPE = '2'
                           `

var SelectWebViewMonthDefault string = `
							SELECT title, content, position 
							FROM 
								b_web_view 
							WHERE 
								DATE = '202109'
							AND TYPE = '2'
                           `

var SelectWebViewWeek string = `
							SELECT title, content, position 
							FROM 
								b_web_view 
							WHERE 
								DATE = '#{bsDt}'
							AND TYPE = '1'
                           `

var SelectWebViewWeekDefault string = `
							SELECT title, content, position 
							FROM 
								b_web_view 
							WHERE 
								DATE = '202109'
							AND TYPE = '1'
                           `

var SelectBaeminReview2 string = `
							SELECT round(sum(rating)/COUNT(*),2) AS rating_avg, COUNT(*) AS cnt, LEFT(DATE,6) AS date  
							FROM 
								a_baemin_review 
							WHERE 
								baemin_id ='#{baeminId}' 
							AND LEFT(DATE,6) 
							BETWEEN 
								'#{startDt}' 
							AND 
								'#{endDt}' 
							GROUP BY LEFT(DATE,6) ORDER BY LEFT(DATE,6) desc
                            `

var SelectBaeminReview string = `
							SELECT right(a.date,2) as date, ifnull(rating_avg,0) AS rating_avg, ifnull(b.cnt,0) AS cnt FROM(							
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH),'%Y%m') AS date
							UNION
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 2 MONTH),'%Y%m') AS date
							UNION
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 3 MONTH),'%Y%m') AS date
							UNION
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 4 MONTH),'%Y%m') AS date
							) a left JOIN (SELECT round(sum(rating)/COUNT(*),2) AS rating_avg, COUNT(*) AS cnt, LEFT(DATE,6) AS date  
								FROM 
									a_baemin_review 
								WHERE 
									baemin_id ='#{baeminId}' 
								AND LEFT(DATE,6) 
								BETWEEN 
									'#{startDt}' 
								AND 
									'#{endDt}' 
								GROUP BY LEFT(DATE,6) ORDER BY LEFT(DATE,6) desc) b ON a.date = b.date
                            `

var SelectNaverReview2 string = `
							SELECT round(sum(rating)/COUNT(*),2) AS rating_avg, COUNT(*) AS cnt, LEFT(REPLACE(created,'.',''),6) AS date  
							FROM 
								a_naver_review 
							WHERE 
								naver_id ='#{naverId}'
							AND LEFT(REPLACE(created,'.',''),6) 
							BETWEEN 
								'#{startDt}' 
							AND 
								'#{endDt}' 
							GROUP BY LEFT(REPLACE(created,'.',''),6) 
							ORDER BY LEFT(REPLACE(created,'.',''),6) desc
                            `

var SelectNaverReview string = `
							SELECT right(a.date,2) as date, ifnull(rating_avg,0) AS rating_avg, ifnull(b.cnt,0) AS cnt FROM(							
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH),'%Y%m') AS date
							UNION
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 2 MONTH),'%Y%m') AS date
							UNION
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 3 MONTH),'%Y%m') AS date
							UNION
							SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 4 MONTH),'%Y%m') AS date
							) a left JOIN (SELECT round(sum(rating)/COUNT(*),2) AS rating_avg, COUNT(*) AS cnt, LEFT(REPLACE(created,'.',''),6) AS date  
								FROM 
									a_naver_review 
								WHERE 
									naver_id ='#{naverId}'
								AND LEFT(REPLACE(created,'.',''),6) 
								BETWEEN 
									'#{startDt}' 
								AND 
									'#{endDt}' 
								GROUP BY LEFT(REPLACE(created,'.',''),6) 
								ORDER BY LEFT(REPLACE(created,'.',''),6) desc) b ON a.date = b.date
                            `

var SelectYogiyoReview2 string = `
							SELECT round(sum(rating)/COUNT(*),2) AS rating_avg, COUNT(*) AS cnt, LEFT(REPLACE(time,'-',''),6) AS date  
							FROM 
								a_yogiyo_review 
							WHERE 
								yogiyo_id ='#{yogiyoId}' 
							AND LEFT(REPLACE(time,'-',''),6) 
							BETWEEN 
								'#{startDt}' 
							AND 
								'#{endDt}' 
							GROUP BY LEFT(REPLACE(time,'-',''),6) 
							ORDER BY LEFT(REPLACE(time,'-',''),6) desc
                            `

var SelectYogiyoReview string = `
                                SELECT right(a.date,2) as date, ifnull(rating_avg,0) AS rating_avg, ifnull(b.cnt,0) AS cnt FROM(							
								SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH),'%Y%m') AS date
								UNION
								SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 2 MONTH),'%Y%m') AS date
								UNION
								SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 3 MONTH),'%Y%m') AS date
								UNION
								SELECT DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 4 MONTH),'%Y%m') AS date
								) a left JOIN (SELECT round(sum(rating)/COUNT(*),2) AS rating_avg, COUNT(*) AS cnt, LEFT(REPLACE(time,'-',''),6) AS date  
									FROM 
										a_yogiyo_review 
									WHERE 
										yogiyo_id ='#{yogiyoId}' 
									AND LEFT(REPLACE(time,'-',''),6) 
									BETWEEN 
										'#{startDt}' 
									AND 
										'#{endDt}' 
									GROUP BY LEFT(REPLACE(time,'-',''),6) 
									ORDER BY LEFT(REPLACE(time,'-',''),6) desc) b ON a.date = b.date
                                `

// 주간 매출 분석 (요일별 날자 입력, 지난 7일)
/*
var SelectWeekCash string = `select round(aa.total) AS total
							, round(aa.maxAmt) AS maxAmt
							, (select DAYOFWEEK(b.bs_dt) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt between '20201227' and '20210102' and (b.tot_amt) = aa.maxAmt) as maxWeek
							,round(aa.avgAmt) AS  avgAmt
							,round(aa.minAmt) AS minAmt
							,(select DAYOFWEEK(b.bs_dt) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt between '20201227' and '20210102' and (b.tot_amt) = aa.minAmt) as minWeek
							from (
							select sum(a.tot_amt) as total,
							max(a.tot_amt) as maxAmt,
							avg(a.tot_amt)as avgAmt,
							min(a.tot_amt)as minAmt
							from cc_aprv_sum a
							where a.biz_num = '1231692222' and a.bs_dt between '20201227' and '20210102' ) aa
							`

var SelectWeekCash2 string = `SELECT ROUND(aa.total) AS total
							, ROUND(aa.maxAmt) AS maxAmt
							, (SELECT DAYOFWEEK(b.bs_dt) FROM cc_aprv_sum b WHERE
								b.biz_num = '#{bizNum}'
								AND b.bs_dt BETWEEN '#{startDt}'
								AND '#{endDt}' AND (b.tot_amt) = aa.maxAmt) AS maxWeek
							,ROUND(aa.avgAmt) AS  avgAmt
							,ROUND(aa.minAmt) AS minAmt
							,(SELECT DAYOFWEEK(b.bs_dt) FROM cc_aprv_sum b WHERE
								b.biz_num = '#{bizNum}'
								AND b.bs_dt BETWEEN '#{startDt}'
								AND '#{endDt}' AND (b.tot_amt) = aa.minAmt) AS minWeek
							FROM (
							SELECT SUM(a.tot_amt) AS total,
							MAX(a.tot_amt) AS maxAmt,
							AVG(a.tot_amt) AS avgAmt,
							MIN(a.tot_amt) AS minAmt
							FROM cc_aprv_sum a WHERE
								a.biz_num = '#{bizNum}'
								AND a.bs_dt BETWEEN '#{startDt}'
								AND '#{endDt}'
							) aa
							`
*/

// 주간 건수 분석
/*
var SelectWeekCnt string =
	`
							SELECT bs_dt, CASE DAYOFWEEK(bs_dt)
							WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS dayName,
							LEFT(tr_tm,2) AS tm, LEFT(tr_tm,2)+1 AS tr_tm2, SUM(1) AS cnt
							FROM cc_aprv_dtl
							WHERE
								biz_num = '#{bizNum}'
							AND bs_dt
							BETWEEN '#{startDt}'
							AND '#{endDt}'
							GROUP BY bs_dt, tm ORDER BY cnt DESC limit 3;
    `
*/
// 주간 건수 분석 - 3시간 단위
var SelectWeekCnt string = `
							SELECT CASE z.day_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name,
							SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.cnt, 0)) AS t0003,
							SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.cnt, 0)) AS t0306,
							SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.cnt, 0) )AS t0609,
							SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.cnt, 0)) AS t0912,
							SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.cnt, 0)) AS t1215,
							SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.cnt, 0)) AS t1518,
							SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.cnt, 0)) AS t1821,
							SUM(IF(z.tr_hr >= '21', z.cnt, 0)) AS t2124,
							SUM(z.cnt) AS tot_cnt
							FROM (
								SELECT DAYOFWEEK(bs_dt) AS day_index, LEFT(tr_tm,2) AS tr_hr, COUNT(*) cnt FROM cc_aprv_dtl
								WHERE
									biz_num = '#{bizNum}'
								AND bs_dt
									BETWEEN '#{startDt}'
									AND '#{endDt}'
								GROUP BY DAYOFWEEK(bs_dt), LEFT(tr_tm,2)
							) z GROUP BY z.day_index ORDER BY tot_cnt DESC LIMIT 3
`

// 주간 매출 분석 (요일별 날자 입력, 지난 7일)
var SelectWeekCash string = `
							SELECT if(left(z.minDt,1)='0', SUBSTR(z.minDt,2,5), z.minDt) AS minDt
								, if(left(z.maxDt,1)='0', SUBSTR(z.maxDt,2,5), z.maxDt) AS maxDt
								,ROUND(z.total) AS total
								,ROUND(z.avgAmt) AS avgAmt
								,z.maxAmt
								,z.tot_cnt
								,z.avg_cnt
								,(SELECT GROUP_CONCAT(
										CONCAT(CASE DAYOFWEEK(b.bs_dt)
										WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일'
										WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END, ' ' , if(left(DATE_FORMAT(STR_TO_DATE(b.bs_dt, '%Y%m%d'), '%m/%d'),1)='0', SUBSTR(DATE_FORMAT(STR_TO_DATE(b.bs_dt, '%Y%m%d'), '%m/%d'),2,5), DATE_FORMAT(STR_TO_DATE(b.bs_dt, '%Y%m%d'), '%m/%d'))) SEPARATOR '|')
									FROM cc_aprv_sum b 
									WHERE 
										b.biz_num = '#{bizNum}' 
										AND b.bs_dt 
											BETWEEN '#{startDt}' 
											AND '#{endDt}' 
											AND b.tot_amt = z.maxAmt
								) AS maxWeek
								,z.minAmt
								,(select GROUP_CONCAT(
									CONCAT(CASE DAYOFWEEK(b.bs_dt)
									WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일'
									WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END, ' ' , if(left(DATE_FORMAT(STR_TO_DATE(b.bs_dt, '%Y%m%d'), '%m/%d'),1)='0', SUBSTR(DATE_FORMAT(STR_TO_DATE(b.bs_dt, '%Y%m%d'), '%m/%d'),2,5), DATE_FORMAT(STR_TO_DATE(b.bs_dt, '%Y%m%d'), '%m/%d'))) SEPARATOR '|')
								FROM cc_aprv_sum b 
								WHERE 
									b.biz_num = '#{bizNum}' 
									AND b.bs_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										AND b.tot_amt = z.minAmt
								) AS minWeek
							FROM (
								SELECT
									DATE_FORMAT(STR_TO_DATE(MIN(a.bs_dt), '%Y%m%d'), '%m/%d') AS minDt,
									DATE_FORMAT(STR_TO_DATE(MAX(a.bs_dt), '%Y%m%d'), '%m/%d') AS maxDt,
									SUM(a.tot_amt) AS total,
									SUM(a.tot_cnt) AS tot_cnt,
									ROUND(SUM(a.tot_amt)/SUM(a.tot_cnt),0) AS avg_cnt,
									MAX(a.tot_amt) AS maxAmt,
									(SUM(a.tot_amt)/COUNT(a.bs_dt)) AS avgAmt,
									MIN(a.tot_amt) AS minAmt
								FROM cc_aprv_sum a
								WHERE 
									a.biz_num = '#{bizNum}' 
									AND a.bs_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}'
							) z
							`

// 지난 주 요일 별 매출 항목
var SelectLastWeekAprv string = `
							SELECT DAYOFWEEK(bs_dt) as day_index, DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') as bs_dt, tot_cnt, tot_amt, 
							CASE DAYOFWEEK(bs_dt) WHEN 1 THEN '일' WHEN 2 THEN '월' WHEN 3 THEN '화' WHEN 4 THEN '수' WHEN 5 THEN '목' WHEN 6 THEN '금' ELSE '토' END AS day_name 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND bs_dt 
								BETWEEN '#{lStartDt}' 
								AND '#{lEndDt}' 
							ORDER BY tot_amt desc
`

// 지난 주 요일 별 매출 항목
var SelectLastWeek string = `
							SELECT count(tr_tm) as tot_cnt, SUM(aprv_amt) AS tot_amt, 
							SUM(if(z.tr_tm <=3,1,0)) AS cnt03, SUM(if(z.tr_tm <=3,z.aprv_amt,0)) AS amt03, 
							SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,1,0)) AS cnt36, SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,z.aprv_amt,0)) AS amt36, 
							SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,1,0)) AS cnt69, SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,z.aprv_amt,0)) AS amt69, 
							SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,1,0)) AS cnt912, SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,z.aprv_amt,0)) AS amt912, 
							SUM(if(z.tr_tm >12 AND z.tr_tm <=15,1,0)) AS cnt1215, SUM(if(z.tr_tm >12 AND z.tr_tm <=15,z.aprv_amt,0)) AS amt1215,
							SUM(if(z.tr_tm >15 AND z.tr_tm <=18,1,0)) AS cnt1518, SUM(if(z.tr_tm >15 AND z.tr_tm <=18,z.aprv_amt,0)) AS amt1518,
							SUM(if(z.tr_tm >18 AND z.tr_tm <=21,1,0)) AS cnt1821, SUM(if(z.tr_tm >18 AND z.tr_tm <=21,z.aprv_amt,0)) AS amt1821,
							SUM(if(z.tr_tm >21 AND z.tr_tm <=24,1,0)) AS cnt2124, SUM(if(z.tr_tm >21 AND z.tr_tm <=24,z.aprv_amt,0)) AS amt2124 FROM (
							SELECT bs_dt, LEFT(tr_tm, 2) AS tr_tm, aprv_amt 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								bs_dt = '#{weekDay}'
							) z
`

// 어제부터 지난 7일 요일별 매출 항목
var SelectYeasterWeekAprv string = `
							SELECT DAYOFWEEK(bs_dt) as day_index, DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') as bs_dt, ifnull(tot_cnt,0) AS tot_cnt, ifnull(tot_amt,0) AS tot_amt, 
							CASE DAYOFWEEK(bs_dt) WHEN 1 THEN '일' WHEN 2 THEN '월' WHEN 3 THEN '화' WHEN 4 THEN '수' WHEN 5 THEN '목' WHEN 6 THEN '금' ELSE '토' END AS day_name 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND bs_dt 
								BETWEEN '#{yStartDt}' 
								AND '#{yEndDt}' 
							ORDER BY bs_dt
`

// 단골 고객 비율 3개월간 분석
var SelectThreeMonthAprv string = `
							SELECT LEFT(z.bs_dt,6) AS bs_dt, COUNT(distinct(z.bs_dt)) dt_cnt, SUM(if(z.cnt = 1, z.cnt, 0)) AS visit1, SUM(if(z.cnt >= 2 AND z.cnt <=3, z.cnt, 0)) AS visit23, SUM(if(z.cnt > 3, z.cnt, 0)) AS visit4, SUM(z.tot_amt) AS tot_amt, SUM(z.cnt) AS tot_cnt, ROUND(SUM(z.tot_amt)/SUM(z.cnt),0) AS avg_amt, ROUND(SUM(z.cnt)/count(DISTINCT(z.bs_dt)),0) AS avg_cnt 
							FROM (
								SELECT bs_dt, card_no, COUNT(*) AS cnt, SUM(aprv_amt) AS tot_amt 
								FROM cc_aprv_dtl 
								WHERE 
									biz_num = '#{bizNum}' 
								AND 
									LEFT(bs_dt,6) = '#{mStartDt1}'  
								GROUP BY card_no
							) z
							UNION ALL
							SELECT LEFT(z.bs_dt,6) AS bs_dt, COUNT(distinct(z.bs_dt)) dt_cnt, SUM(if(z.cnt = 1, z.cnt, 0)) AS visit1, SUM(if(z.cnt >= 2 AND z.cnt <=3, z.cnt, 0)) AS visit23, SUM(if(z.cnt > 3, z.cnt, 0)) AS visit4, SUM(z.tot_amt) AS tot_amt, SUM(z.cnt) AS tot_cnt, ROUND(SUM(z.tot_amt)/SUM(z.cnt),0) AS avg_amt, ROUND(SUM(z.cnt)/count(DISTINCT(z.bs_dt)),0) AS avg_cnt 
							FROM (
								SELECT bs_dt, card_no, COUNT(*) AS cnt, SUM(aprv_amt) AS tot_amt 
								FROM cc_aprv_dtl 
								WHERE 
									biz_num = '#{bizNum}' 
								AND 
									LEFT(bs_dt,6) = '#{mStartDt2}' 
								GROUP BY card_no
							) z
							UNION ALL
							SELECT LEFT(z.bs_dt,6) AS bs_dt, COUNT(distinct(z.bs_dt)) dt_cnt, SUM(if(z.cnt = 1, z.cnt, 0)) AS visit1, SUM(if(z.cnt >= 2 AND z.cnt <=3, z.cnt, 0)) AS visit23, SUM(if(z.cnt > 3, z.cnt, 0)) AS visit4, SUM(z.tot_amt) AS tot_amt, SUM(z.cnt) AS tot_cnt, ROUND(SUM(z.tot_amt)/SUM(z.cnt),0) AS avg_amt, ROUND(SUM(z.cnt)/count(DISTINCT(z.bs_dt)),0) AS avg_cnt 
							FROM (
								SELECT bs_dt, card_no, COUNT(*) AS cnt, SUM(aprv_amt) AS tot_amt 
								FROM cc_aprv_dtl 
								WHERE 
									biz_num = '#{bizNum}' 
								AND 
									LEFT(bs_dt,6) = '#{mStartDt3}' 
								GROUP BY card_no
							) z
`

/*
var SelectThreeMonthVisitAprv string = `
							SELECT DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') AS dt, CASE DAYOFWEEK(bs_dt)
							WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS dayName,
							LEFT(tr_tm,2) AS tm, LEFT(tr_tm,2)+1 AS tr_tm2, SUM(1) AS cnt
							FROM cc_aprv_dtl
							WHERE
								biz_num = '#{bizNum}'
							AND LEFT(bs_dt,6)
								BETWEEN '#{mStartDt3}'
								AND '#{mStartDt1}'
							GROUP BY dt, tm ORDER BY cnt DESC LIMIT 3;
`
var SelectThreeMonthVisitAprv string = `
							SELECT CASE z.day_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name , z.tm, z.tm2, z.cnt, SUM(z.cnt) AS tot_cnt
							FROM (
								SELECT DAYOFWEEK(bs_dt) AS day_index, LEFT(tr_tm,2) AS tm ,LEFT(tr_tm,2)+1 AS tm2, COUNT(*) cnt FROM cc_aprv_dtl
								WHERE
									biz_num = '#{bizNum}'
								AND LEFT(bs_dt,6)
									BETWEEN '#{mStartDt3}'
									AND '#{mStartDt1}'
								GROUP BY DAYOFWEEK(bs_dt), LEFT(tr_tm,2) ORDER BY cnt desc
							) z GROUP BY day_name ORDER BY tot_cnt DESC LIMIT 3
`
*/
// 3개월 평균 요일별 고객 분포
var SelectThreeMonthVisitAprv string = `
							SELECT CASE z.day_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, 
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.cnt, 0)) AS t0003,
									SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.cnt, 0)) AS t0306,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.cnt, 0) )AS t0609,
									SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.cnt, 0)) AS t0912,
									SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.cnt, 0)) AS t1215,
									SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.cnt, 0)) AS t1518,
									SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.cnt, 0)) AS t1821,
									SUM(IF(z.tr_hr >= '21', z.cnt, 0)) AS t2124,
									SUM(z.cnt) AS tot_cnt
							FROM (
							SELECT DAYOFWEEK(bs_dt) AS day_index, LEFT(tr_tm,2) AS tr_hr, COUNT(*) cnt FROM cc_aprv_dtl 
							WHERE
								biz_num = '#{bizNum}'
							AND LEFT(bs_dt,6) 
								BETWEEN '#{mStartDt3}' 
								AND '#{mStartDt1}' 
							GROUP BY DAYOFWEEK(bs_dt), LEFT(tr_tm,2)
							) z GROUP BY z.day_index ORDER BY tot_cnt DESC LIMIT 3
`
var SelectThreeMonthVisitAprv2 string = `
							SELECT SUM(tot_amt) AS total_amt, SUM(tot_cnt) AS tot_cnt, ROUND(SUM(tot_amt)/SUM(tot_cnt),0) AS tot_avg
							FROM cc_aprv_sum
							WHERE
								biz_num = '#{bizNum}'
							AND LEFT(bs_dt,6)
								BETWEEN '#{mStartDt3}' 
								AND '#{mStartDt1}'
							GROUP BY biz_num;
`

var SelectLastWeekAprvTip string = `
							SELECT 'WD' as day_type, count(tr_tm) as tot_cnt, SUM(aprv_amt) AS tot_amt, 
							SUM(if(z.tr_tm <=3,1,0)) AS cnt03, SUM(if(z.tr_tm <=3,z.aprv_amt,0)) AS amt03, 
							SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,1,0)) AS cnt36, SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,z.aprv_amt,0)) AS amt36, 
							SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,1,0)) AS cnt69, SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,z.aprv_amt,0)) AS amt69, 
							SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,1,0)) AS cnt912, SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,z.aprv_amt,0)) AS amt912, 
							SUM(if(z.tr_tm >12 AND z.tr_tm <=15,1,0)) AS cnt1215, SUM(if(z.tr_tm >12 AND z.tr_tm <=15,z.aprv_amt,0)) AS amt1215,
							SUM(if(z.tr_tm >15 AND z.tr_tm <=18,1,0)) AS cnt1518, SUM(if(z.tr_tm >15 AND z.tr_tm <=18,z.aprv_amt,0)) AS amt1518,
							SUM(if(z.tr_tm >18 AND z.tr_tm <=21,1,0)) AS cnt1821, SUM(if(z.tr_tm >18 AND z.tr_tm <=21,z.aprv_amt,0)) AS amt1821,
							SUM(if(z.tr_tm >21 AND z.tr_tm <=24,1,0)) AS cnt2124, SUM(if(z.tr_tm >21 AND z.tr_tm <=24,z.aprv_amt,0)) AS amt2124 
							FROM (
								SELECT CASE DAYOFWEEK(bs_dt) WHEN 1 THEN 'HD' WHEN 2 THEN 'WD' WHEN 3 THEN 'WD' WHEN 4 THEN 'WD' WHEN 5 THEN 'WD' WHEN 6 THEN 'WD' ELSE 'HD' END AS day_type, LEFT(tr_tm, 2) AS tr_tm, aprv_amt 
								FROM cc_aprv_dtl 
								WHERE 
									biz_num = '#{bizNum}' 
								AND bs_dt 
									BETWEEN '#{lStartDt}' 
									AND '#{lEndDt}' 
							) z WHERE z.day_type = 'WD'
							UNION ALL
							SELECT 'HD' as day_type, count(tr_tm) as tot_cnt, SUM(aprv_amt) AS tot_amt, 
							SUM(if(z.tr_tm <=3,1,0)) AS cnt03, SUM(if(z.tr_tm <=3,z.aprv_amt,0)) AS amt03, 
							SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,1,0)) AS cnt36, SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,z.aprv_amt,0)) AS amt36, 
							SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,1,0)) AS cnt69, SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,z.aprv_amt,0)) AS amt69, 
							SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,1,0)) AS cnt912, SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,z.aprv_amt,0)) AS amt912, 
							SUM(if(z.tr_tm >12 AND z.tr_tm <=15,1,0)) AS cnt1215, SUM(if(z.tr_tm >12 AND z.tr_tm <=15,z.aprv_amt,0)) AS amt1215,
							SUM(if(z.tr_tm >15 AND z.tr_tm <=18,1,0)) AS cnt1518, SUM(if(z.tr_tm >15 AND z.tr_tm <=18,z.aprv_amt,0)) AS amt1518,
							SUM(if(z.tr_tm >18 AND z.tr_tm <=21,1,0)) AS cnt1821, SUM(if(z.tr_tm >18 AND z.tr_tm <=21,z.aprv_amt,0)) AS amt1821,
							SUM(if(z.tr_tm >21 AND z.tr_tm <=24,1,0)) AS cnt2124, SUM(if(z.tr_tm >21 AND z.tr_tm <=24,z.aprv_amt,0)) AS amt2124 
							FROM (
								SELECT CASE DAYOFWEEK(bs_dt) WHEN 1 THEN 'HD' WHEN 2 THEN 'WD' WHEN 3 THEN 'WD' WHEN 4 THEN 'WD' WHEN 5 THEN 'WD' WHEN 6 THEN 'WD' ELSE 'HD' END AS day_type, LEFT(tr_tm, 2) AS tr_tm, aprv_amt 
								FROM cc_aprv_dtl 
								WHERE 
									biz_num = '#{bizNum}' 
								AND bs_dt 
									BETWEEN '#{lStartDt}' 
									AND '#{lEndDt}' 
							) z WHERE z.day_type = 'HD'
							UNION ALL
							SELECT 'ALL' as day_type, count(tr_tm) as tot_cnt, SUM(aprv_amt) AS tot_amt, 
							SUM(if(z.tr_tm <=3,1,0)) AS cnt03, SUM(if(z.tr_tm <=3,z.aprv_amt,0)) AS amt03, 
							SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,1,0)) AS cnt36, SUM(if(z.tr_tm > 3 AND z.tr_tm <= 6,z.aprv_amt,0)) AS amt36, 
							SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,1,0)) AS cnt69, SUM(if(z.tr_tm > 6 AND z.tr_tm <= 9,z.aprv_amt,0)) AS amt69, 
							SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,1,0)) AS cnt912, SUM(if(z.tr_tm > 9 AND z.tr_tm <= 12,z.aprv_amt,0)) AS amt912, 
							SUM(if(z.tr_tm >12 AND z.tr_tm <=15,1,0)) AS cnt1215, SUM(if(z.tr_tm >12 AND z.tr_tm <=15,z.aprv_amt,0)) AS amt1215,
							SUM(if(z.tr_tm >15 AND z.tr_tm <=18,1,0)) AS cnt1518, SUM(if(z.tr_tm >15 AND z.tr_tm <=18,z.aprv_amt,0)) AS amt1518,
							SUM(if(z.tr_tm >18 AND z.tr_tm <=21,1,0)) AS cnt1821, SUM(if(z.tr_tm >18 AND z.tr_tm <=21,z.aprv_amt,0)) AS amt1821,
							SUM(if(z.tr_tm >21 AND z.tr_tm <=24,1,0)) AS cnt2124, SUM(if(z.tr_tm >21 AND z.tr_tm <=24,z.aprv_amt,0)) AS amt2124 
							FROM (
								SELECT bs_dt, LEFT(tr_tm, 2) AS tr_tm, aprv_amt 
								FROM cc_aprv_dtl 
								WHERE 
									biz_num = '#{bizNum}' 
								AND bs_dt 
									BETWEEN '#{lStartDt}' 
									AND '#{lEndDt}' 
							) z
`

// 월 평균고객 분석
/*
var SelectWeekPersonVisit string = `
								select round(sum(cust_cnt),0) AS visitTotal,
								(select round(sum(a.cust_cnt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range = '01' and a.reg_no = '1231692222' and total_yearmonth = '202012') as visit1,
								(select round(sum(a.cust_cnt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range in ('02', '03') and a.reg_no = '1231692222' and total_yearmonth = '202012') as  visit23,
								(select round(sum(a.cust_cnt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range > '04' and a.reg_no = '1231692222' and total_yearmonth = '202012') as visit4
								from kftc_b_da_regular_cust_cc
								where reg_no = '1231692222' and total_yearmonth = '202012'
								`
*/
// 월 평균고객 분석
var SelectWeekPersonVisit string = `
									SELECT
										SUM(y.visit_person) AS visitTotal,
										SUM(IF(y.visit_cnt = 1, y.visit_person, 0)) AS visit1,
										SUM(IF(y.visit_cnt >= 2 AND y.visit_cnt <= 3, y.visit_person, 0)) AS visit23,
										SUM(IF(y.visit_cnt >= 4, y.visit_person, 0)) AS visit4
									FROM ( 
										SELECT z.visit_cnt, COUNT(*) AS visit_person
										FROM
										(
										SELECT COUNT(*) AS visit_cnt
										FROM cc_aprv_dtl WHERE 
											biz_num = '#{bizNum}'
											and bs_dt between '#{startDt}' 
											and '#{endDt}' 
											GROUP BY card_no
										) z
										GROUP BY z.visit_cnt
									) y
									`

/*
var SelectWeekPersonVisit1 string = `
									SELECT
										SUM(y.visit_person) AS visitTotal,
										SUM(IF(y.visit_cnt = 1, y.visit_person, 0)) AS visit1,
										SUM(IF(y.visit_cnt >= 2 AND y.visit_cnt <= 3, y.visit_person, 0)) AS visit23,
										SUM(IF(y.visit_cnt >= 4, y.visit_person, 0)) AS visit4
									FROM (
										SELECT z.visit_cnt, COUNT(*) AS visit_person
										FROM
										(
										SELECT COUNT(*) AS visit_cnt
										FROM cc_aprv_dtl
											WHERE biz_num = '1231692222'
												and bs_dt between '20201101' and '20210131'
												GROUP BY card_no
										) z
										GROUP BY z.visit_cnt
									) y
									`
*/

// 객단가 예측
var SelectAverageRevenuePerUser string = `
										SELECT ROUND(AVG(tot_amt/tot_cnt), 0) AS arpu
										FROM cc_aprv_sum
										WHERE 
											biz_num = '#{bizNum}'
											AND bs_dt BETWEEN '#{startDt}'
											AND '#{endDt}'
										`

// 지난주 취소 분석
var SelectLastCancleList string = `
									SELECT tr_dt, tr_tm, aprv_no 
									FROM cc_aprv_dtl 
									WHERE 
										biz_num = '#{bizNum}'
									AND aprv_clss = 1
									AND bs_dt 
									BETWEEN 
										'#{startDt}'
									AND 
										'#{endDt}'
							
`

// 어제 취소 분석
var SelectYesterdayCancleList string = `
									SELECT tr_dt, tr_tm, aprv_no 
									FROM cc_aprv_dtl 
									WHERE 
										biz_num = '#{bizNum}'
									AND aprv_clss = 1
									AND bs_dt = '#{bsDt}'
`

// 12주간 매출 비교
var Select12WeekAmt string = `
									SELECT SUM(real_amt) AS amt, 
										'' AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN '#{0startDt}' 
											AND '#{0endDt}' 
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{1startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{1startDt}')
										- 
											WEEKOFYEAR('#{1day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{1startDt}' 
											AND '#{1endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{2startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{2startDt}')
										- 
											WEEKOFYEAR('#{2day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{2startDt}' 
											AND '#{2endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{3startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{3startDt}')
										- 
											WEEKOFYEAR('#{3day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{3startDt}' 
											AND '#{3endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{4startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{4startDt}')
										- 
											WEEKOFYEAR('#{4day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{4startDt}' 
											AND '#{4endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{5startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{5startDt}')
										- 
											WEEKOFYEAR('#{5day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{5startDt}' 
											AND '#{5endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{6startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{6startDt}')
										- 
											WEEKOFYEAR('#{6day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{6startDt}' 
											AND '#{6endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{7startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{7startDt}')
										- 
											WEEKOFYEAR('#{7day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{7startDt}' 
											AND '#{7endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{8startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{8startDt}')
										- 
											WEEKOFYEAR('#{8day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{8startDt}' 
											AND '#{8endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{9startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{9startDt}')
										- 
											WEEKOFYEAR('#{9day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{9startDt}' 
											AND '#{9endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{10startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{10startDt}')
										- 
											WEEKOFYEAR('#{10day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{10startDt}' 
											AND '#{10endDt}'  
										GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
										CONCAT(
											DATE_FORMAT('#{11startDt}',
										'%m'),
										'월 ',
											WEEKOFYEAR('#{11startDt}')
										- 
											WEEKOFYEAR('#{11day1}')
										+ 1, '주차') AS dayIndex 
										FROM cc_day_sale_sum 
										WHERE 
											biz_num='#{bizNum}' 
										AND tr_dt 
											BETWEEN	'#{11startDt}' 
											AND '#{11endDt}'  
										GROUP BY biz_num
`

var Select12MonthAmt string = `
									SELECT SUM(real_amt) AS amt,'지난달' AS dayIndex 
									FROM cc_day_sale_sum 
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{0bsDt}'  
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{1bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{1bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{1bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{2bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{2bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{2bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{3bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{3bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{3bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{4bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{4bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{4bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{5bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{5bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{5bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{6bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{6bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{6bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{7bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{7bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{7bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{8bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{8bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{8bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{9bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{9bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{9bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{10bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{10bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{10bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT SUM(real_amt) AS amt,
									CONCAT(
										DATE_FORMAT('#{11bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{11bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_day_sale_sum
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										left(tr_dt,6) = '#{11bsDt}'   
									GROUP BY biz_num
`

var Select6MonthAmtCard string = `
									SELECT tot_amt AS amt,'지난달' AS dayIndex 
									FROM cc_aprv_sum_month 
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{0bsDt}'  
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{1bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{1bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{1bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{2bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{2bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{2bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{3bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{3bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{3bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{4bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{4bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{4bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{5bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{5bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{5bsDt}'   
									GROUP BY biz_num
`

var Select12MonthAmtCard string = `

									SELECT tot_amt AS amt,'지난달' AS dayIndex
									FROM cc_aprv_sum_month 
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{0bsDt}'  
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{1bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{1bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{1bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{2bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{2bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{2bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{3bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{3bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{3bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{4bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{4bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{4bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{5bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{5bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{5bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{6bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{6bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{6bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{7bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{7bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{7bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{8bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{8bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{8bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{9bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{9bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{9bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{10bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{10bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{10bsDt}'   
									GROUP BY biz_num
									UNION ALL
									SELECT tot_amt AS amt,
									CONCAT(
										DATE_FORMAT('#{11bsDt1}',
									'%Y'),'년 ',
										DATE_FORMAT('#{11bsDt1}',
									'%c'),'월') AS dayIndex 
									FROM cc_aprv_sum_month
									WHERE 
										biz_num = '#{bizNum}' 
									AND 
										bs_dt = '#{11bsDt}'   
									GROUP BY biz_num
`

var SelectLastCancleAprv string = `
									SELECT tr_dt, tr_tm 
									FROM cc_aprv_dtl 
									WHERE 
										biz_num = '#{bizNum}'
									AND 
										aprv_no = '#{aprvNo}'
`

// 지난달 취소 분석
var SelectLastMonthCancleList string = `
									SELECT tr_dt, tr_tm, aprv_no 
									FROM cc_aprv_dtl 
									WHERE 
										biz_num = '#{bizNum}'
									AND aprv_clss = 1
									AND 
										LEFT(bs_dt,6) = '#{bsDt}' 
`

// 지난주 리뷰 분석
var SelectDeliveryInfo string = `
									SELECT baemin_id, yogiyo_id, naver_id, coupang_id
									FROM b_store 
									WHERE 
										biz_num ='#{bizNum}'
`

// 가맹점 rating, keyword
var SelectReivewOption string = `
								SELECT rating, keyword 
								FROM b_store_etc 
								WHERE 
									biz_num = '#{bizNum}'
`

var SelectReviews string = `
							(SELECT contents AS content, member_no, rating 
							FROM a_baemin_review 
							WHERE 
								baemin_id ='#{baeminId}'
							AND date 
								BETWEEN '#{startDt}' 
								AND '#{endDt}'
							)
							UNION ALL
							(SELECT body AS content, '' as member_no, rating
							FROM a_naver_review 
							WHERE 
								naver_id = '#{naverId}'
							AND date_format(created,'%Y%m%d') 
								BETWEEN '#{startDt}' 
								AND '#{endDt}'
							)
							UNION ALL
							(SELECT COMMENT AS content, '' as member_no, rating 
							FROM a_yogiyo_review 
							WHERE 
								yogiyo_id = '#{yogiyoId}'
							AND date_format(TIME,'%Y%m%d') 
								BETWEEN '#{startDt}' 
								AND '#{endDt}'
							)
							UNION ALL
							(SELECT REVIEW_TEXT AS content, '' as member_no, rating 
							FROM a_coupang_review 
							WHERE 
								coupang_id = '#{coupangId}'
							AND date 
								BETWEEN '#{startDt}' 
								AND '#{endDt}'
							)
                             `

var SelectYesterdayReviews string = `
							(SELECT contents AS content, member_no, rating 
							FROM a_baemin_review 
							WHERE 
								baemin_id ='#{baeminId}'
							AND date = '#{bsDt}'
							)
							UNION ALL
							(SELECT body AS content, '' as member_no, rating
							FROM a_naver_review 
							WHERE 
								naver_id = '#{naverId}'
							AND date_format(created,'%Y%m%d') = '#{bsDt}'
							)
							UNION ALL
							(SELECT COMMENT AS content, '' as member_no, rating 
							FROM a_yogiyo_review 
							WHERE 
								yogiyo_id = '#{yogiyoId}'
							AND date_format(TIME,'%Y%m%d') = '#{bsDt}'
							)
							UNION ALL
							(SELECT REVIEW_TEXT AS content, '' as member_no, rating 
							FROM a_coupang_review 
							WHERE 
								coupang_id = '#{coupangId}'
							AND date = '#{bsDt}'
							)
                             `

var SelectMonthReviews string = `
							(SELECT contents AS content, member_no, rating 
							FROM a_baemin_review 
							WHERE 
								baemin_id ='#{baeminId}'
							AND 
								LEFT(date,6) = '#{bsDt}'
							)
							UNION ALL
							(SELECT body AS content, '' as member_no, rating
							FROM a_naver_review 
							WHERE 
								naver_id = '#{naverId}'
							AND
								LEFT(date_format(created,'%Y%m%d'),6) = '#{bsDt}'
							)
							UNION ALL
							(SELECT COMMENT AS content, '' as member_no, rating 
							FROM a_yogiyo_review 
							WHERE 
								yogiyo_id = '#{yogiyoId}'
							AND 
								left(date_format(TIME,'%Y%m%d'),6) = '#{bsDt}'
							)
                             `

var SelectBaeminCustomerCount string = `
							SELECT COUNT(*) AS cnt, SUM(rating)/COUNT(*) AS avg 
							FROM a_baemin_review 
							WHERE 
								member_no ='#{customerId}'
							GROUP BY member_no
`

var SelectBaeminCustomerInfo string = `
							SELECT member_name, SUM(rating) AS total, COUNT(*) AS cnt, round(SUM(rating)/COUNT(*),1) AS avg
							FROM a_baemin_customer
							WHERE
								member_no ='#{customerId}'
							GROUP BY member_no
`

// 현금영수증 객단가 예측
var SelectCashAverageRevenuePerUser string = `
										SELECT ROUND(AVG(tot_amt/tot_cnt), 0) AS arpu
										FROM cc_cash_sum
										WHERE 
											biz_num = '#{bizNum}'
											AND bs_dt BETWEEN '#{startDt}'
											AND '#{endDt}'
										`

// 시간별 매출 분석 (요일별 날자 입력, 지난 7일), 3시간 단위
var SelectWeekAvgTime string = `
								SELECT 
									z.tr_dt AS trDt, 
									DAYOFWEEK(z.tr_dt) AS week,
									CASE DAYOFWEEK(z.TR_DT) WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일'
									WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS weekNm,
									z.week_end AS weekEnd,
									SUM(z.amt) AS weekTotal,
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.amt, 0)) AS t0003,
									SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.amt, 0)) AS t0306,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.amt, 0) )AS t0609,
									SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.amt, 0)) AS t0912,
									SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.amt, 0)) AS t1215,
									SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.amt, 0)) AS t1518,
									SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.amt, 0)) AS t1821,
									SUM(IF(z.tr_hr >= '21', z.amt, 0)) AS t2124
								FROM
								(
									SELECT tr_dt, week_end, LEFT(tr_tm,2) AS tr_hr, SUM(aprv_amt) AS amt
									FROM cc_aprv_dtl
									WHERE
										BIZ_NUM = '#{bizNum}'
										AND tr_dt = '#{weekDay}'
									GROUP BY tr_dt, LEFT(tr_tm,2)
								) z 
								GROUP BY z.tr_dt
								`

// 시간별 매출 분석 (요일별 날자 입력, 지난 7일), 6시간 단위
var SelectWeekAvgTime2 string = `
								SELECT 
									z.tr_dt AS trDt, 
									DAYOFWEEK(z.tr_dt) AS week,
									CASE DAYOFWEEK(z.TR_DT) WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일'
									WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS weekNm,
									z.week_end AS weekEnd,
									SUM(z.amt) AS weekTotal,
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '06', z.amt, 0)) AS t0006,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '11', z.amt, 0)) AS t0611,
									SUM(IF(z.tr_hr >= '11' AND z.tr_hr < '14', z.amt, 0) )AS t1114,
									SUM(IF(z.tr_hr >= '14' AND z.tr_hr < '17', z.amt, 0)) AS t1417,
									SUM(IF(z.tr_hr >= '17' AND z.tr_hr < '24', z.amt, 0)) AS t1724
								FROM
								(
									SELECT tr_dt, week_end, LEFT(tr_tm,2) AS tr_hr, SUM(aprv_amt) AS amt
									FROM cc_aprv_dtl
									WHERE
										BIZ_NUM = '#{bizNum}'
										AND tr_dt = '#{weekDay}'
									GROUP BY tr_dt, LEFT(tr_tm,2)
								) z 
								GROUP BY z.tr_dt
								`

/*
var SelectWeekAvgTime2 string = `SELECT IFNULL(sum(a.aprv_amt),0) as weekTotal,
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '00' and '0359'),0) as 'line0003',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '04' and '0659'),0) as 'line0406',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '07' and '0959'),0) as 'line0709',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '10' and '1259'),0) as 'line1012',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '13' and '1559'),0) AS 'line1315',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '16' and '1859'),0) as 'line1618',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
									WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '19' and '2159'),0) as 'line1921',
								IFNULL((select sum(b.aprv_amt) from cc_aprv_dtl b
								WHERE
									b.biz_num = '#{bizNum}'
									and b.tr_dt = a.tr_dt and b.tr_tm between '22' and '2459'),0) as 'line2224'
								FROM cc_aprv_dtl a
								WHERE
									a.biz_num = '#{bizNum}'
								AND a.tr_dt = '#{weekDay}'
								`
*/

var SelectWeekAvgTime1 string = `
								SELECT 
									z.tr_dt AS trDt, 
									DAYOFWEEK(z.tr_dt) AS week,
									CASE DAYOFWEEK(z.TR_DT) WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일'
									WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS weekNm,
									z.week_end AS weekEnd,
									SUM(z.amt) AS totSum,
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.amt, 0)) AS t0003,
									SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.amt, 0)) AS t0306,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.amt, 0) )AS t0609,
									SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.amt, 0)) AS t0912,
									SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.amt, 0)) AS t1215,
									SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.amt, 0)) AS t1518,
									SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.amt, 0)) AS t1821,
									SUM(IF(z.tr_hr >= '21', z.amt, 0)) AS t2124
								FROM
								(
									SELECT tr_dt, week_end, LEFT(tr_tm,2) AS tr_hr, SUM(aprv_amt) AS amt
									FROM cc_aprv_dtl
									WHERE
										BIZ_NUM = '#{bizNum}'
										AND tr_dt BETWEEN '#{startDt}' 
										AND '#{endDt}'
									GROUP BY tr_dt, LEFT(tr_tm,2)
								) z 
								GROUP BY z.tr_dt
								`

// 요일 분석 팁 -> 평일 주말 비교
var SelectWeekdayAnalystic string = `select ifnull(round(sum(a.aprv_amt)),0) AS total,
ifnull((select sum(b.aprv_amt) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt between '20201227' and '20210102'and week_end = 'HD'),0)  as 'holy',
ifnull((select sum(b.aprv_amt) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt between '20201227' and '20210102'and week_end = 'WD'),0)  as 'work'
from cc_aprv_dtl a
where a.biz_num = '1231692222'
and a.bs_dt between '20201227' and '20210102'
									`

var SelectWeekdayAnalystic2 string = `
										select ifnull(round(sum(a.aprv_amt)),0) AS total,
										ifnull((select sum(b.aprv_amt) from cc_aprv_sum b 
										WHERE 
											b.biz_num = '#{bizNum}' 
											AND b.bs_dt between '#{startDt}' 
											AND '#{endDt}'
											and week_end = 'HD'),0)  as 'holy',
										ifnull((select sum(b.aprv_amt) from cc_aprv_sum b 
										WHERE 
											b.biz_num = '#{bizNum}' 
											and b.bs_dt between '#{startDt}' 
											and '#{endDt}'
											and week_end = 'WD'),0)  as 'work' from cc_aprv_dtl a 
										WHERE 
											a.biz_num = '#{bizNum}'
											and a.bs_dt between '#{startDt}' 
											and '#{endDt}'
									`

var SelectWeekdayAnalystic1 string = `
									SELECT
										z.week_end AS weekEnd,
										SUM(z.amt) AS totSum,
										SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.amt, 0)) AS t0003,
										SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.amt, 0)) AS t0306,
										SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.amt, 0)) AS t0609,
										SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.amt, 0)) AS t0912,
										SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.amt, 0)) AS t1215,
										SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.amt, 0)) AS t1518,
										SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.amt, 0)) AS t1821,
										SUM(IF(z.tr_hr >= '21', z.amt, 0)) AS t2124
									FROM
									(
										SELECT week_end, LEFT(tr_tm,2) AS tr_hr, SUM(aprv_amt) AS amt
										FROM cc_aprv_dtl
										WHERE
											biz_num = '#{bizNum}'
											AND tr_dt BETWEEN '#{startDt}' AND '#{endDt}'
										GROUP BY week_end, LEFT(tr_tm,2)
									) z
									`

// 주말 바쁜 시간대 도출
var SelectWeeBusyTime string = `SELECT atTime,d_amt
								FROM (
								SELECT DISTINCT
								c.atTime, 
								CASE c.atTime
									WHEN '0~3' THEN ifnull((select sum(b.aprv_amt) 
									FROM cc_aprv_dtl b 
									WHERE b.biz_num = a.biz_num 
										AND b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '00' and '0359'),0) 
									when '4~6' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '04' and '0659'),0)
									when '7~9' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '07' and '0959'),0)
									when '10~12' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '10' and '1259'),0)
									when '13~15' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '13' and '1559'),0)
									when '16~18' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num AND b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '16' and '1859'),0)
									when '19~21' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '19' and '2159'),0)
									when '22~24' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'HD' and b.tr_tm between '22' and '2459'),0)
								END AS d_amt 
								FROM cc_aprv_dtl a
								cross join
								(
								SELECT '0~3' as atTime
								union all SELECT '4~6'
								union all SELECT '7~9'
								union all SELECT '10~12'
								union all SELECT '13~15'
								union all SELECT '16~18'
								union all SELECT '19~21'
								union all SELECT '22~24'
								) c
								WHERE
								 a.biz_num = '#{bizNum}' 
								 AND a.tr_dt 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								 ) AS BB
								ORDER BY d_amt DESC
								LIMIT 1
								`

// 주중 바쁜 시간대 도출
var SelectWeeBusyTimeWork string = `
							SELECT atTime,d_amt
								FROM (
								SELECT DISTINCT
								c.atTime, 
								CASE c.atTime
									WHEN '0~3' THEN ifnull((select sum(b.aprv_amt) 
									FROM cc_aprv_dtl b 
									WHERE b.biz_num = a.biz_num 
										AND b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '00' and '0359'),0) 
									when '4~6' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '04' and '0659'),0)
									when '7~9' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '07' and '0959'),0)
									when '10~12' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '10' and '1259'),0)
									when '13~15' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '13' and '1559'),0)
									when '16~18' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num AND b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '16' and '1859'),0)
									when '19~21' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '19' and '2159'),0)
									when '22~24' then ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt 
										BETWEEN '#{startDt}' 
										AND '#{endDt}' 
										and week_end = 'WD' and b.tr_tm between '22' and '2459'),0)
								END AS d_amt 
								FROM cc_aprv_dtl a
								cross join
								(
								SELECT '0~3' as atTime
								union all SELECT '4~6'
								union all SELECT '7~9'
								union all SELECT '10~12'
								union all SELECT '13~15'
								union all SELECT '16~18'
								union all SELECT '19~21'
								union all SELECT '22~24'
								) c
								WHERE
								 a.biz_num = '#{bizNum}' 
								 AND a.tr_dt 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								 ) AS BB
								ORDER BY d_amt DESC
								LIMIT 1
								`

// 월간분석 - 월 매출 분석
var SelectMonthCash string = `select IFNULL(round(sum(a.tot_amt),0),0) as 'allAmt',
								IFNULL((select round(b.tot_amt, 0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt = '202011'),0) as 'beforeMonthAmt',
								IFNULL((select round(avg(b.tot_amt),0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20201231'),0) as 'avgAmt',
								IFNULL((select round(max(b.tot_amt),0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20201231'),0) as 'maxAmt',
								IFNULL((select right(b.bs_dt,2) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20201231' and round(b.tot_amt,0) = 
										(select round(max(b.tot_amt),0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20201231') limit 1 ) ,0) AS 'maxMonth'
								from cc_aprv_sum_month a 
								where a.biz_num = '1231692222' and a.bs_dt < '20201231'
								union all 
								select IFNULL(round(sum(a.tot_amt),0),0) as 'allAmt',
								IFNULL((select round(b.tot_amt, 0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt = '201911'),0) as 'beforeMonthAmt',
								IFNULL((select round(avg(b.tot_amt),0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20191231'),0) as 'avgAmt',
								IFNULL((select round(max(b.tot_amt),0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20191231'),0) as 'maxAmt',
								IFNULL((select right(b.bs_dt,2) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20191231' and round(b.tot_amt,0) = 
								(select round(max(b.tot_amt),0) from cc_aprv_sum_month b where b.biz_num = a.biz_num and b.bs_dt < '20191231') limit 1 ),0) as 'maxMonth'
								from cc_aprv_sum_month a 
								where a.biz_num = '1231692222' and a.bs_dt < '20191231'
								`

var SelectMonthAprv string = `
							SELECT if(left(right(bs_dt,2),1)=0, RIGHT(bs_dt,1), RIGHT(bs_dt,2)) AS bs_dt, tot_cnt, tot_amt, round(tot_amt/tot_cnt,0) as tot_avg 
							FROM cc_aprv_sum_month 
							WHERE 
								BIZ_NUM = '#{bizNum}' 
							AND bs_dt 
								BETWEEN '#{startDt}' 
								AND '#{endDt}'
`

var SelectMonthDayAprv string = `
							SELECT DAYOFWEEK(bs_dt) as day_index, sum(aprv_amt) AS tot_amt, COUNT(*) AS cnt, ROUND(SUM(aprv_amt)/COUNT(*),0) AS avg_amt 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}' 
							GROUP BY day_index					
`

var SelectMonthDayAprv2 string = `
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=1
							UNION ALL
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=2
							UNION ALL
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=3
							UNION ALL
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=4
							UNION ALL
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=5
							UNION ALL
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=6
							UNION ALL
							SELECT ifnull(DAYOFWEEK(bs_dt),1) as day_index, ifnull(sum(tot_amt),0) AS tot_amt, ifnull(sum(tot_cnt),0) AS cnt, ifnull(ROUND(sum(tot_amt)/sum(tot_cnt),0),0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{lastDt}'
							AND DAYOFWEEK(bs_dt)=7

`

var SelectMonthCash1 string = `
							SELECT
								if(LEFT(MIN(RIGHT(z.month, 2)),1)=0,right(MIN(RIGHT(z.month, 2)),1),MIN(RIGHT(z.month, 2))) as fromDt,
								if(LEFT(MAX(RIGHT(z.month, 2)),1)=0,right(MAX(RIGHT(z.month, 2)),1),MAX(RIGHT(z.month, 2))) as toDt,
								CONCAT(LEFT('#{endDt}', 4), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%m')) AS beforeMonth,
								IFNULL((SELECT b.tot_amt FROM cc_aprv_sum_month b WHERE biz_num = z.biz_num AND bs_dt = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 MONTH), '%Y%m')), 0) AS beforeMonthAmt,
								ROUND(AVG(z.tot_amt),0) AS avgAmt,
								(SELECT LEFT(c.bs_dt,6) FROM cc_aprv_sum_month c WHERE c.tot_amt = MAX(z.tot_amt) AND c.biz_num = z.biz_num) AS maxMonth,
								MAX(z.tot_amt) AS maxAmt,
								SUM(z.tot_amt) AS allAmt
							FROM 
							(
								SELECT 
									b.biz_num, 
									LEFT(a.dt,6) AS 'month', 
									IFNULL(b.tot_amt,0) AS tot_amt
								FROM cc_date_info a
									LEFT JOIN cc_aprv_sum_month b ON LEFT(a.dt,6) = b.bs_dt 
									AND 
										b.biz_num = '#{bizNum}'
								WHERE 
									LEFT(a.dt,6) 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								GROUP BY LEFT(a.dt,6)
							) z
							`

var SelectMonthCash2 string = `
							SELECT
								if(LEFT(MIN(RIGHT(z.month, 2)),1)=0,right(MIN(RIGHT(z.month, 2)),1),MIN(RIGHT(z.month, 2))) as fromDt,
								if(LEFT(MAX(RIGHT(z.month, 2)),1)=0,right(MAX(RIGHT(z.month, 2)),1),MAX(RIGHT(z.month, 2))) as toDt,
								MIN(z.month) AS minMonth,
								MIN(z.tot_amt) AS minMonthAmt,
								ROUND(AVG(z.tot_amt),0) AS avgAmt,
								(SELECT LEFT(c.bs_dt,6) FROM cc_aprv_sum_month c WHERE c.tot_amt = MAX(z.tot_amt) AND c.biz_num = z.biz_num) AS maxMonth,
								MAX(z.tot_amt) AS maxAmt,
								SUM(z.tot_amt) AS allAmt
							FROM 
							(
								SELECT 
									b.biz_num, 
									LEFT(a.dt,6) AS 'month', 
									IFNULL(b.tot_amt,0) AS tot_amt
								FROM cc_date_info a
									LEFT JOIN cc_aprv_sum_month b ON LEFT(a.dt,6) = b.bs_dt 
									AND 
										b.biz_num = '#{bizNum}'
								WHERE 
									LEFT(a.dt,6) 
										BETWEEN '#{startDt}' 
										AND '#{endDt}'
								GROUP BY LEFT(a.dt,6)
							) z

							`

//월평균 요일별 매출 (요일별 날자 입력)
var SelectMonthAvgWeek string = `select round(avg(a.aprv_amt),0) as total,
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201228','20201221','20201214','20201207')) as 'mon',
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201229','20201222','20201215','20201208','20201201')) as 'tue',
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201230','20201223','20201216','20201209','20201202')) as 'wed',
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201231','20201224','20201217','20201210','20201203')) as 'thr',
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201225','20201218','20201211','20201204')) as 'fri',
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201226','20201219','20201212','20201205')) as 'sat',
								(select round(avg(b.aprv_amt),0) from cc_aprv_sum b where b.biz_num = '1231692222' and b.bs_dt in ('20201227','20201220','20201213','20201206')) as 'sun'
								from cc_aprv_sum a
								where a.biz_num = '1231692222' and a.bs_dt between '20201201' and '20201231'
								`

/*
var SelectMonthCntDetail string = `
									SELECT DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') AS dt, CASE DAYOFWEEK(bs_dt)
									WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS dayName,
									LEFT(tr_tm,2) AS tm, LEFT(tr_tm,2)+1 AS tr_tm2, SUM(1) AS cnt
									FROM cc_aprv_dtl
									WHERE
										biz_num = '#{bizNum}'
									AND
										LEFT(bs_dt,6) = '#{bsDt}'
									GROUP BY dt, tm ORDER BY cnt DESC LIMIT 3;
                                  `
*/
var SelectMonthCntDetail string = `
							SELECT CASE z.day_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, 
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.cnt, 0)) AS t0003,
									SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.cnt, 0)) AS t0306,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.cnt, 0) )AS t0609,
									SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.cnt, 0)) AS t0912,
									SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.cnt, 0)) AS t1215,
									SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.cnt, 0)) AS t1518,
									SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.cnt, 0)) AS t1821,
									SUM(IF(z.tr_hr >= '21', z.cnt, 0)) AS t2124,
									SUM(z.cnt) AS tot_cnt
							FROM (
							SELECT DAYOFWEEK(bs_dt) AS day_index, LEFT(tr_tm,2) AS tr_hr, COUNT(*) cnt FROM cc_aprv_dtl 
							WHERE
								biz_num = '#{bizNum}'
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY DAYOFWEEK(bs_dt), LEFT(tr_tm,2)
							) z GROUP BY z.day_index ORDER BY tot_cnt DESC LIMIT 3	
`

var SelectWeekCntDetail string = `
							SELECT CASE z.day_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, 
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.cnt, 0)) AS t0003,
									SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.cnt, 0)) AS t0306,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.cnt, 0) )AS t0609,
									SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.cnt, 0)) AS t0912,
									SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.cnt, 0)) AS t1215,
									SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.cnt, 0)) AS t1518,
									SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.cnt, 0)) AS t1821,
									SUM(IF(z.tr_hr >= '21', z.cnt, 0)) AS t2124,
									SUM(z.cnt) AS tot_cnt
							FROM (
							SELECT DAYOFWEEK(bs_dt) AS day_index, LEFT(tr_tm,2) AS tr_hr, COUNT(*) cnt FROM cc_aprv_dtl 
							WHERE
								biz_num = '#{bizNum}'
							AND bs_dt
							BETWEEN
								'#{startDt}'
							AND
								'#{endDt}'
							GROUP BY DAYOFWEEK(bs_dt), LEFT(tr_tm,2)
							) z GROUP BY z.day_index ORDER BY tot_cnt DESC LIMIT 3	
`

var SelectMonthTmCnt string = `
							SELECT LEFT(tr_tm,2) AS tm, LEFT(tr_tm,2)+1 AS tr_tm2, COUNT(*) AS cnt
							FROM cc_aprv_dtl
							WHERE
								biz_num = '#{bizNum}'
							AND 
								LEFT(bs_dt,6) = '#{bsDt}'
							GROUP BY tm ORDER BY cnt DESC LIMIT 4;
`

var SelectLastMonthAnal string = `
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') as dt, DAYOFWEEK(bs_dt) AS dt_index, COUNT(*) AS cnt, SUM(aprv_amt) AS amt, ROUND(SUM(aprv_amt)/COUNT(*),0) AS avg_amt 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY cnt DESC LIMIT 1) AS a
							UNION ALL
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') as dt, DAYOFWEEK(bs_dt) AS dt_index, COUNT(*) AS cnt, SUM(aprv_amt) AS amt, ROUND(SUM(aprv_amt)/COUNT(*),0) AS avg_amt 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY amt DESC LIMIT 1) AS b
							UNION ALL
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') as dt, DAYOFWEEK(bs_dt) AS dt_index, COUNT(*) AS cnt, SUM(aprv_amt) AS amt, ROUND(SUM(aprv_amt)/COUNT(*),0) AS avg_amt 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY avg_amt DESC LIMIT 1) AS c
							UNION ALL
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(STR_TO_DATE(bs_dt, '%Y%m%d'), '%m/%d') as dt, DAYOFWEEK(bs_dt) AS dt_index, COUNT(*) AS cnt, SUM(aprv_amt) AS amt, ROUND(SUM(aprv_amt)/COUNT(*),0) AS avg_amt 
							FROM cc_aprv_dtl 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY amt LIMIT 1) AS d
`

var SelectLastMonthAnal2 string = `
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(bs_dt, '%c월 %e일') as dt, DAYOFWEEK(bs_dt) AS dt_index, tot_cnt AS cnt, SUM(aprv_amt) AS amt, ROUND(tot_amt/tot_cnt,0) AS avg_amt
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY cnt DESC LIMIT 3) AS a
							UNION ALL
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(bs_dt, '%c월 %e일') as dt, DAYOFWEEK(bs_dt) AS dt_index, tot_cnt AS cnt, SUM(aprv_amt) AS amt, ROUND(tot_amt/tot_cnt,0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY amt DESC LIMIT 3) AS b
							UNION ALL
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(bs_dt, '%c월 %e일') as dt, DAYOFWEEK(bs_dt) AS dt_index, tot_cnt AS cnt, SUM(aprv_amt) AS amt, ROUND(tot_amt/tot_cnt,0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY avg_amt DESC LIMIT 3) AS c
							UNION ALL
							SELECT dt, CASE dt_index WHEN 1 THEN '일요일' WHEN 2 THEN '월요일' WHEN 3 THEN '화요일' WHEN 4 THEN '수요일' WHEN 5 THEN '목요일' WHEN 6 THEN '금요일' ELSE '토요일' END AS day_name, cnt, amt, avg_amt FROM (
							SELECT DATE_FORMAT(bs_dt, '%c월 %e일') as dt, DAYOFWEEK(bs_dt) AS dt_index, tot_cnt AS cnt, SUM(aprv_amt) AS amt, ROUND(tot_amt/tot_cnt,0) AS avg_amt 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND 
								LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY bs_dt ORDER BY amt LIMIT 3) AS d
`

var SelectMonthTotalAmt string = `
							SELECT SUM(tot_amt) AS total_amt, SUM(tot_cnt) AS tot_cnt, ROUND(SUM(tot_amt)/SUM(tot_cnt),0) AS tot_avg
							FROM cc_aprv_sum
							WHERE
								biz_num = '#{bizNum}'
							AND 
								LEFT(bs_dt,6) = '#{bsDt}'
							GROUP BY biz_num
`

var SelectMonthCnt string = `
							SELECT if(LEFT(right(left(bs_dt,6),2),1)=0,right(left(bs_dt,6),1),right(left(bs_dt,6),2)) as bs_dt, SUM(tot_amt) AS total_amt, SUM(tot_cnt) AS tot_cnt, ROUND(SUM(tot_amt)/SUM(tot_cnt),0) AS tot_avg 
							FROM cc_aprv_sum 
							WHERE 
								biz_num = '#{bizNum}' 
							AND LEFT(bs_dt,6) = '#{bsDt}' 
							GROUP BY biz_num
                            `

var SelectMonthAvgWeek1 string = `
								SELECT 
									y.week_cd AS weekCd,
									y.tot_sum AS total,
									ROUND(y.tot_sum/week_cnt, 0) AS totAvg
								FROM (
									SELECT 
										z.week_cd AS week_cd,
										COUNT(z.week_cd) AS week_cnt,
										SUM(z.aprv_amt) AS tot_sum
									FROM 
									(
										SELECT 
										tr_dt, 
										DAYOFWEEK(tr_dt) AS week_cd, 
										SUM(aprv_amt) AS aprv_amt
										FROM cc_aprv_dtl
										WHERE
										biz_num = '#{bizNum}'
										AND LEFT(tr_dt, 6) = '#{bsDt}'
										GROUP BY tr_dt, DAYOFWEEK(tr_dt)
									) z
									GROUP BY z.week_cd
								) y
								GROUP BY y.week_cd
								ORDER BY y.week_cd
								`

// 달아요팁
var SelectMonthDarayoTip string = `select round(avg(a.aprv_amt),0) as darayoTiptotal,
									(select round(avg(b.aprv_amt),0) from cc_aprv_dtl b where b.biz_num = '1231692222' and b.tr_dt between '20201201' and '20201231'and week_end = 'HD') as 'holy',
									(select round(avg(b.aprv_amt),0) from cc_aprv_dtl b where b.biz_num = '1231692222' and b.tr_dt between '20201201' and '20201231'and ifnull(week_end,'0') != 'HD') as 'work'
									from cc_aprv_dtl a
									where a.biz_num = '1231692222'
									and a.tr_dt between '20201201' and '20201231'
								`

var SelectMonthDarayoTip1 string = `
								SELECT
									z.week_end AS weekEnd,
									SUM(z.amt) AS total,
									SUM(IF(z.tr_hr >= '00' AND z.tr_hr < '03', z.amt, 0)) AS t0003,
									SUM(IF(z.tr_hr >= '03' AND z.tr_hr < '06', z.amt, 0)) AS t0306,
									SUM(IF(z.tr_hr >= '06' AND z.tr_hr < '09', z.amt, 0) )AS t0609,
									SUM(IF(z.tr_hr >= '09' AND z.tr_hr < '12', z.amt, 0)) AS t0912,
									SUM(IF(z.tr_hr >= '12' AND z.tr_hr < '15', z.amt, 0)) AS t1215,
									SUM(IF(z.tr_hr >= '15' AND z.tr_hr < '18', z.amt, 0)) AS t1518,
									SUM(IF(z.tr_hr >= '18' AND z.tr_hr < '21', z.amt, 0)) AS t1821,
									SUM(IF(z.tr_hr >= '21', z.amt, 0)) AS t2124
								FROM
								(
									SELECT 
										week_end, 
										LEFT(tr_tm,2) AS tr_hr, 
										SUM(aprv_amt) AS amt
									FROM cc_aprv_dtl
									WHERE
										biz_num = '#{bizNum}'
										AND LEFT(tr_dt,6) = '#{bsDt}'
									GROUP BY week_end, LEFT(tr_tm,2)
								) z 
								GROUP BY z.week_end
								`

var SelectMonthBusyTime string = `select ifnull(sum(a.aprv_amt),0) as total,
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '00' and '0359'),0) as '0003',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '04' and '0659'),0) as '0406',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '07' and '0959'),0) as '0709',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '10' and '1259'),0) as '1012',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '13' and '1559'),0) as '1315',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '16' and '1859'),0) as '1618',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '19' and '2159'),0) as '1921',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201201' and '20201231' and week_end = 'HD' and b.tr_tm between '22' and '2459'),0) as '2224'
								from cc_aprv_dtl a
								where a.biz_num = '1231692222' and a.tr_dt between '20201201' and '20201231'
								union
								select ifnull(sum(a.aprv_amt),0) as total,
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '00' and '0359'),0) as '0003',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '04' and '0659'),0) as '0406',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '07' and '0959'),0) as '0709',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '10' and '1259'),0) as '1012',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '13' and '1559'),0) as '1315',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '16' and '1859'),0) as '1618',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '19' and '2159'),0) as '1921',
								ifnull((select sum(b.aprv_amt) from cc_aprv_dtl b where b.biz_num = a.biz_num and b.tr_dt between '20201224' and '20201231' and ifnull(week_end,'0') != 'HD' and b.tr_tm between '22' and '2459'),0) as '2224'
								from cc_aprv_dtl a
								where a.biz_num = '1231692222' and a.tr_dt between '20201201' and '20201231' 
								`

// 방문분석
var SelectMonthVisit string = `select round(sum(app_amt),0) as visitTotal, 
							(select round(sum(a.app_amt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range = '01' and a.reg_no = '1231692222' and total_yearmonth = '202012') as 'visit01',
							(select round(sum(a.app_amt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range in ('02', '03') and a.reg_no = '1231692222' and total_yearmonth = '202012') as 'visit25',
							(select round(sum(a.app_amt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range = '04'  and a.reg_no = '1231692222' and total_yearmonth = '202012') as 'visit69',
							(select round(sum(a.app_amt), 0) from kftc_b_da_regular_cust_cc a where a.visit_cnt_range = '05' and a.reg_no = '1231692222' and total_yearmonth = '202012') as 'visit10'
							from kftc_b_da_regular_cust_cc
							where 
							reg_no = '1231692222' and total_yearmonth = '202012'
							`

var SelectMonthVisit1 string = `
								SELECT
									SUM(y.visit_person) AS visitTotal,
									SUM(y.aprv_amt) AS visitTotalAmt,
									SUM(IF(y.visit_cnt = 1, y.visit_person, 0)) AS visit1,
									SUM(IF(y.visit_cnt = 1, y.aprv_amt, 0)) AS visitAmt1,
									SUM(IF(y.visit_cnt >= 2 AND y.visit_cnt <= 3, y.visit_person, 0)) AS visit23,
									SUM(IF(y.visit_cnt >= 2 AND y.visit_cnt <= 3, y.aprv_amt, 0)) AS visitAmt23,
									SUM(IF(y.visit_cnt >= 4 AND y.visit_cnt <= 9, y.visit_person, 0)) AS visit49,
									SUM(IF(y.visit_cnt >= 4 AND y.visit_cnt <= 9, y.aprv_amt, 0)) AS visitAmt49,
									SUM(IF(y.visit_cnt >= 10, y.visit_person, 0)) AS visit10,
									SUM(IF(y.visit_cnt >= 10, y.aprv_amt, 0)) AS visitAmt10
								FROM (
									SELECT 
										z.visit_cnt, 
										COUNT(*) AS visit_person, 
										SUM(z.aprv_amt) AS aprv_amt
									FROM
									(
									SELECT 
										COUNT(*) AS visit_cnt,
										SUM(aprv_amt) AS aprv_amt
									FROM cc_aprv_dtl
									WHERE 
										biz_num = '#{bizNum}'
										AND LEFT(bs_dt,6) = '#{bsDt}'
									GROUP BY card_no
									ORDER BY COUNT(*)
									) z
									GROUP BY z.visit_cnt
								) y

								`

var SelectMonthAroundStroeVisit string = `SELECT mv_rtl_name, mv_cnt, mv_cd_lv2_nm 
										FROM kftc_b_ag_svr_rltn_day_cc
										WHERE 
										reg_no = '1010202362' 
										AND mv_cnt != 0 
										order by mv_cnt desc                              
										LIMIT 4
										`

var SelectMonthCompareKorea string = ` 
										SELECT reg_no 
										FROM kftc_b_da_retl_mkt_cc
										WHERE LEFT(CD_LV2_NM,2) = '한식' 
										AND 
											total_yearmonth = '#{bsDt}' 
										ORDER BY avg_amt DESC LIMIT 1
										`
var SelectMonthCompareList string = `
										SELECT a.reg_no, a.sido_nm, b.buety, left(a.CD_LV2_NM,2) as cd, a.avg_cnt, a.avg_amt, a.avg_sido_cnt, a.avg_sido_amt, a.avg_gungu_cnt, a.avg_gungu_amt, a.avg_total_cnt, a.avg_total_amt  
										FROM kftc_b_da_retl_mkt_cc a LEFT JOIN priv_rest_info b 
										ON a.REG_NO=b.BUSID 
										WHERE 
											a.total_yearmonth = '#{bsDt}'
										AND
											b.buety = '#{buety}'
`

var SelectPrivRestInfo string = `
										SELECT buety, left(addr,4) as addr FROM priv_rest_info 
										WHERE 
											busid='#{bizNum}'
`

/*
var SelectMonthCompareKorea string = `

										SELECT reg_no
										FROM kftc_b_da_retl_mkt_cc
										WHERE LEFT(CD_LV2_NM,2) = '한식'
										AND
											total_yearmonth = '202101'
										ORDER BY avg_amt DESC LIMIT 1
										`
*/
// '#{bsDt}' -> 202101
var SelectMonthCompareAll string = ` 
										SELECT 
											ROUND(avg_cnt,0) AS my_cnt, ROUND(avg_amt,0) AS my_amt, ROUND(avg_amt/avg_cnt,0) as my_avgprice,
											ROUND(avg_sido_cnt,0) AS sido_cnt, ROUND(avg_sido_amt, 0) AS sido_amt, ROUND(avg_sido_amt/avg_sido_cnt,0) as sido_avgprice,
											ROUND(avg_gungu_cnt,0) AS gungu_cnt, ROUND(avg_gungu_amt, 0) AS gungu_amt, ROUND(avg_gungu_amt/avg_gungu_cnt,0) as gungu_avgprice,
											ROUND(avg_total_cnt,0) AS total_cnt, ROUND(avg_total_amt, 0) AS total_amt, ROUND(avg_total_amt/avg_total_cnt,0) as total_avgprice
										FROM kftc_b_da_retl_mkt_cc	
										WHERE 
											reg_no = '#{bizNum}' 
										AND total_yearmonth = '#{bsDt}'
										`
var SelectCompAprv string = `
										SELECT tot_cnt, tot_amt, ROUND(tot_amt/tot_cnt,0) AS tot_avg 
										FROM cc_aprv_sum_month 
										WHERE 
											biz_num='#{bizNum}' 
										AND 
											bs_dt='#{bsDt}'
`

var SelectLastSaleInfo string = `
								SELECT CARD_CD AS MAX_CARD_CD
										,CARD_NM AS MAX_CARD_NM
										,A.TOT_AMT AS APRV_AMT
										,a.CAN_CNT
										,a.TOT_CNT
								FROM cc_aprv_sum_month AS A
								INNER JOIN cc_aprv_lst AS B ON A.BIZ_NUM = B.BIZ_NUM AND A.BS_DT = left(B.BS_DT,6)
								WHERE 
								A.biz_num= '#{bizNum}' 
								AND A.BS_DT= '#{bsDt}'
								GROUP BY CARD_CD,CARD_NM
								ORDER BY SUM(B.APRV_CNT) DESC
								LIMIT 1
						`

var Select6MonthSales string = `
								SELECT  BS_DT
										,left(TOT_AMT ,length(TOT_AMT) - 3) as TOT_AMT
								FROM cc_aprv_sum_month AS A
								WHERE 
								biz_num= '#{bizNum}' 
								AND BS_DT >= '#{startBsDt}'
								AND BS_DT <= '#{endBsDt}'
								ORDER BY BS_DT ASC
						`
var SelectOrderAmt string = `
								SELECT IFNULL(SUM(TOTAL_AMT),0) AS ORDER_AMT
								FROM PRIV_REST_INFO AS A
								INNER JOIN DAR_ORDER_INFO AS B ON A.REST_ID = B.REST_ID
								WHERE
								A.BUSID='#{bizNum}' 
								AND LEFT(ORDER_DATE,6)='#{bsDt}'
								AND ORDER_STAT='20'
						`
