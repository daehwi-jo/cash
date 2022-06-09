package sales

// 일별 카드 승인 합계리스트
var SelectDayCardSum string = `
							SELECT
								bs_dt AS trDt,
								tot_cnt AS totCnt,
								tot_amt AS totAmt,
								aprv_cnt AS aprvCnt,
								aprv_amt AS aprvAmt,
								can_cnt AS canCnt,
								can_amt AS canAmt
							FROM
								cc_aprv_sum
							WHERE
								biz_num = '#{bizNum}'
								AND bs_dt BETWEEN '#{startDt}'
								AND '#{endDt}'
							ORDER BY bs_dt DESC
							`

// 요청일 카드 승인 상세리스트
var SelectDayCardDetail string = `
								SELECT X.* 
								FROM 
								(
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum, 
										T.*
									FROM (
										SELECT
											biz_num AS bizNum,
											bs_dt AS bsDt,
											seq_no AS seqNo,
											tr_dt AS trDt,
											tr_tm AS trTm,
											aprv_no AS aprvNo,
											card_nm AS cardNm,
											card_no AS cardNo,
											sts_cd AS stsCd,
											(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'OK_CAN_DIV' AND cd_key = sts_cd) AS stsNm,
											aprv_amt AS aprvAmt,
											CASE WHEN inst_trm > 1 THEN CONCAT(inst_trm , '개월') ELSE '일시불' END AS instTrm
										FROM
											cc_aprv_dtl
										WHERE
											biz_num = '#{bizNum}'
											AND BS_DT = '#{bsDt}'
										ORDER BY seq_no
									) T, (SELECT @RN:=0) AS R
								) X 
								WHERE 
									X.rNum >= '#{startNum}'
									AND X.rNum < '#{endNum}'
							`

// 요청일 카드 승인 상세리스트 건수
var SelectDayCardDetailCount string = `
								SELECT COUNT(z.seqNo) AS totCnt
								FROM (
									SELECT
										biz_num AS bizNum,
										bs_dt AS bsDt,
										seq_no AS seqNo,
										tr_dt AS trDt,
										tr_tm AS trTm,
										aprv_no AS aprvNo,
										card_nm AS cardNm,
										card_no AS cardNo,
										sts_cd AS stsCd,
										(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'OK_CAN_DIV' AND cd_key = sts_cd) AS stsNm,
										aprv_amt AS aprvAmt,
										CASE WHEN inst_trm > 1 THEN CONCAT(inst_trm , '개월') ELSE '일시불' END AS instTrm
									FROM
										cc_aprv_dtl
									WHERE
										biz_num = '#{bizNum}'
										AND BS_DT = '#{bsDt}'
									ORDER BY seq_no
								) z
							`

// 일별 현금영수증 승인 합계리스트
var SelectDayCashSum string = `
							SELECT
								tr_dt AS trDt,
								aprv_cnt AS aprvCnt,
								aprv_amt AS aprvAmt,
								can_cnt AS canCnt,
								can_amt AS canAmt,
								aprv_cnt + can_cnt AS totCnt,
								aprv_amt + can_amt AS totAmt
							FROM (
								SELECT
									tr_dt,
									CASE sts_cd WHEN '1' THEN COUNT(tot_amt) ELSE 0 END as aprv_cnt,
									CASE sts_cd WHEN '1' THEN SUM(tot_amt) ELSE 0 END as aprv_amt,
									CASE sts_cd WHEN '3' THEN COUNT(tot_amt) ELSE 0 END as can_cnt,
									CASE sts_cd WHEN '3' THEN SUM(tot_amt) ELSE 0 END as can_amt
								FROM  cc_cash_dtl
								WHERE 
									biz_num = '#{bizNum}'
									AND tr_dt BETWEEN '#{startDt}' 
									AND '#{endDt}'
								GROUP BY tr_dt
								ORDER BY tr_dt DESC
							) Z
							`

// 요청일 현금영수증 승인 상세리스트
var SelectDayCashDetail string = `
								SELECT X.* 
								FROM 
								(
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum, 
										T.*
									FROM (
										SELECT
											seq_no AS seqNo,
											tr_dt AS trDt,
											tr_tm AS trTm,
											isu_cd AS isuCd,
											isu_nm AS isuNm,
											aprv_no AS aprvNo,
											sts_cd AS stsCd,
											(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'OK_CAN_DIV' AND cd_key = sts_cd) AS stsNm,
											base_amt AS baseAmt,
											vat_amt AS vatAmt,
											svc_amt AS svcAmt,
											tot_amt AS totAmt
										FROM cc_cash_dtl
										WHERE
											BIZ_NUM = '#{bizNum}'
											AND BS_DT = '#{bsDt}'
											ORDER BY seq_no
									) T, (SELECT @RN:=0) AS R
								) X 
								WHERE 
									X.rNum >= '#{startNum}'
									AND X.rNum < '#{endNum}'
								`

// 요청일 현금영수증 승인 상세리스트 건수
var SelectDayCashDetailCount string = `
									SELECT COUNT(z.seqNo) AS totCnt
									FROM (
										SELECT
											seq_no AS seqNo,
											tr_dt AS trDt,
											tr_tm AS trTm,
											isu_cd AS isuCd,
											isu_nm AS isuNm,
											aprv_no AS aprvNo,
											sts_cd AS stsCd,
											(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'OK_CAN_DIV' AND cd_key = sts_cd) AS stsNm,
											base_amt AS baseAmt,
											vat_amt AS vatAmt,
											svc_amt AS svcAmt,
											tot_amt AS totAmt
										FROM cc_cash_dtl
										WHERE
											BIZ_NUM = '#{bizNum}'
											AND BS_DT = '#{bsDt}'
											ORDER BY seq_no
									) z
									`

// 일별 입금내역 합계리스트
var SelectDayPaySum string = `
							SELECT
								a.bs_dt AS trDt,
								a.pca_cnt AS pcaCnt,
								a.pca_amt AS pcaAmt,
								a.real_pay_amt AS realPayAmt,
								(SELECT IFNULL(SUM(b.pay_amt),0) FROM cc_pca_dtl b WHERE a.biz_num = b.biz_num AND a.bs_dt = b.outp_expt_dt) - a.real_pay_amt AS delayAmt
							FROM
								cc_pay_lst a
							WHERE
								a.biz_num = '#{bizNum}'
								AND a.bs_dt BETWEEN '#{startDt}' 
								AND '#{endDt}'
								ORDER BY bs_dt DESC
							`

// 요청일 입금내역 상세리스트
var SelectDayPayDetail string = `
								SELECT X.* 
								FROM 
								(
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum, 
										T.*
									FROM (
										SELECT
											bs_dt AS bsDt,
											seq_no AS seqNo,
											pay_dt AS payDt,
											(SELECT a.card_cd FROM cc_card_comp_inf a WHERE a.card_nm = b.card_nm) AS cardCd,
											card_nm AS cardNm,
											TRIM(mer_no) AS merNo,
											pca_cnt AS pcaCnt,
											pca_amt AS pcaAmt,
											rsv_amt AS delayAmt,
											vat_amt AS vatAmt,
											etc_amt AS etcAmt,
											real_pay_amt AS realPayAmt
										FROM
											cc_pay_dtl b
										WHERE
											biz_num = '#{bizNum}'
											AND bs_dt = '#{bsDt}'
											ORDER BY seq_no
									) T, (SELECT @RN:=0) AS R
								) X 
								WHERE 
									X.rNum >= '#{startNum}'
									AND X.rNum < '#{endNum}'
								`

// 요청일 입금내역 상세리스트 건수
var SelectDayPayDetailCount string = `
									SELECT COUNT(z.seqNo) AS totCnt
									FROM (
										SELECT
											bs_dt AS bsDt,
											seq_no AS seqNo,
											pay_dt AS payDt,
											(SELECT a.card_cd FROM cc_card_comp_inf a WHERE a.card_nm = b.card_nm) AS cardCd,
											card_nm AS cardNm,
											TRIM(mer_no) AS merNo,
											pca_cnt AS pcaCnt,
											pca_amt AS pcaAmt,
											rsv_amt AS delayAmt,
											vat_amt AS vatAmt,
											etc_amt AS etcAmt,
											real_pay_amt AS realPayAmt
										FROM
											cc_pay_dtl b
										WHERE
											biz_num = '#{bizNum}'
											AND bs_dt = '#{bsDt}'
											ORDER BY seq_no
									) z
									`

// 매출캘린더 월별 합계리스트
var SelectAprvCalendarSumList string = `
									SELECT
										tr_month AS trMonth,
										SUM(z.aprv_amt) AS aprvAmt,
										SUM(z.cash_amt) AS cashAmt,
										SUM(z.pca_amt) AS pcaAmt,
										SUM(z.tot_amt) AS totAmt
									FROM (
										SELECT
											SUBSTR(a.TOTAL_DATE, 1, 6) AS tr_month,
                                            IFNULL(sum(b.TOT_AMT),0) AS aprv_amt,
                                            0 AS cash_amt,
                                            0 AS pca_amt,
                                            IFNULL(sum(b.TOT_AMT),0) AS tot_amt
                                        FROM
                                            sys_week_date a
                                            LEFT JOIN cc_aprv_sum b ON a.TOTAL_DATE = b.BS_DT 
												AND b.biz_num = '#{bizNum}'
										WHERE
											a.TOTAL_DATE BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY SUBSTR(a.TOTAL_DATE, 1, 6)
									
										UNION ALL
									
										SELECT
											SUBSTR(dt, 1, 6) AS tr_month,
											0 AS aprv_amt,
											0 AS cash_amt,
											IFNULL(SUM(pca_amt),0) AS pca_amt,
											0 AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_pca_dtl ON dt = org_tr_dt 
												AND biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									
										UNION ALL
									
										SELECT
											SUBSTR(DT, 1, 6) AS tr_month,
											0 AS aprv_amt,
											IFNULL(SUM(tot_amt),0) AS cash_amt,
											0 AS pca_amt,
											IFNULL(SUM(tot_amt),0) AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_cash_dtl ON dt = tr_dt 
												AND biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									) z
									GROUP BY z.tr_month
									ORDER BY z.tr_month DESC
									`
/*
var SelectAprvCalendarSumList string = `
									SELECT
										tr_month AS trMonth,
										SUM(z.aprv_amt) AS aprvAmt,
										SUM(z.cash_amt) AS cashAmt,
										SUM(z.pca_amt) AS pcaAmt,
										SUM(z.tot_amt) AS totAmt
									FROM (
										SELECT
											SUBSTR(DT, 1, 6) AS tr_month,
											IFNULL(SUM(aprv_amt),0) AS aprv_amt,
											0 AS cash_amt,
											0 AS pca_amt,
											IFNULL(SUM(aprv_amt),0) AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_aprv_dtl ON dt = tr_dt
												AND biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}'
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)

										UNION ALL

										SELECT
											SUBSTR(dt, 1, 6) AS tr_month,
											0 AS aprv_amt,
											0 AS cash_amt,
											IFNULL(SUM(pca_amt),0) AS pca_amt,
											0 AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_pca_dtl ON dt = org_tr_dt
												AND biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}'
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)

										UNION ALL

										SELECT
											SUBSTR(DT, 1, 6) AS tr_month,
											0 AS aprv_amt,
											IFNULL(SUM(tot_amt),0) AS cash_amt,
											0 AS pca_amt,
											IFNULL(SUM(tot_amt),0) AS tot_amt
										FROM
											cc_date_info
											LEFT JOIN cc_cash_dtl ON dt = tr_dt
												AND biz_num = '#{bizNum}'
										WHERE
											dt BETWEEN '#{startDt}'
											AND '#{endDt}'
										GROUP BY SUBSTR(dt, 1, 6)
									) z
									GROUP BY z.tr_month
									ORDER BY z.tr_month DESC
									`
*/

/*
diffColor
매입이 적으면 붉은색 code 1
매입이 많으면 청색 code 2
매입 대기 초록색 code 3
정상 code 0
 */
// 매출캘린더 리스트
var SelectAprvCalendarList string = `
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum, 
										tr_dt AS trDt, 
										SUM(z.aprv_amt) AS aprvAmt, 
										SUM(z.cash_amt) AS cashAmt, 
										SUM(z.pca_amt) AS pcaAmt, 
										SUM(z.tot_amt) AS totAmt,
										CASE WHEN SUM(z.pca_amt) = 0 THEN '3' WHEN SUM(z.aprv_amt) > SUM(z.pca_amt) THEN '1' WHEN SUM(z.aprv_amt) < SUM(z.pca_amt) THEN '2' ELSE '0' END diffColor,
										SUM(day_color) AS dayColor
									FROM (
										SELECT 
											a.dt AS tr_dt, 
											IFNULL(SUM(b.tot_amt),0) AS aprv_amt, 
											0 AS cash_amt, 
											0 AS pca_amt, 
											IFNULL(SUM(b.tot_amt),0) AS tot_amt,
											0 as day_color
										FROM 
											tb_date_info a 
											LEFT JOIN cc_aprv_sum b 
												ON a.dt = b.bs_dt AND b.biz_num = '#{bizNum}'
										WHERE 
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt

										UNION ALL

										SELECT total_date as tr_dt,
										0 as aprv_amt,
										0 as cash_amt,
										0 as pca_amt,
										0 as tot_amt,
										CASE WHEN DAY = 7 then 2
										WHEN datekind = 'W' then 1
										ELSE 3 END day_color 
										FROM 
											sys_week_date 
										WHERE 
											total_date 
										BETWEEN '#{startDt}' 
										AND '#{endDt}'

										UNION ALL

										SELECT 
											a.dt AS tr_dt, 
											0 AS aprv_amt, 
											0 AS cash_amt, 
											IFNULL(SUM(b.pca_amt),0) AS pca_amt, 
											0 AS tot_amt,
											0 as day_color
										FROM 
											tb_date_info a 
											LEFT JOIN cc_pca_dtl b 
												ON a.dt = b.org_tr_dt AND b.biz_num = '#{bizNum}'
										WHERE  
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt
										
										UNION ALL
										
										SELECT 
											a.dt AS tr_dt, 
											0 AS aprv_amt, 
											IFNULL(SUM(c.tot_amt),0) AS cash_amt, 
											0 AS pca_amt, 
											IFNULL(SUM(c.tot_amt),0) AS tot_amt,
											0 as day_color
										FROM 
											tb_date_info a 
											LEFT JOIN cc_cash_dtl c 
												ON a.dt = c.tr_dt AND c.biz_num = '#{bizNum}'
										WHERE  
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt
										) z INNER JOIN (SELECT @RN := 0) R
										GROUP BY z.tr_dt
										ORDER BY z.tr_dt
									`

// 매출캘린더 카드사별 매입내역 리스트
// 0:거래없음,1:정상매입,2:매입제외,3:매입대기,4:추가매입,-:확인필요
var SelectAprvDailyList string = `
								SELECT 
									z.card_cd AS cardCd, 
									z.card_nm AS cardNm, 
									z.card_tel AS cardTel, 
									z.aprv_cnt AS aprvCnt, 
									z.aprv_amt AS aprvAmt, 
									z.pca_cnt AS pcaCnt, 
									z.pca_amt AS pcaAmt, 
									(z.aprv_amt-z.pca_amt) AS diffAmt,
									z.tot_fee AS totFee, 
									z.vat_amt AS vatAmt, 
									z.pay_amt AS payAmt,
									CASE 
										WHEN z.aprv_cnt > 0 AND z.pca_cnt = 0
											AND (DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), "%Y%m%d") = '#{bsDt}' OR
												 (DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 4 DAY), "%Y%m%d") < '#{bsDt}') AND DAYOFWEEK('#{bsDt}') in ('1','6','7'))
											THEN CONCAT('3:', z.aprv_cnt)
										WHEN z.aprv_cnt = z.pca_cnt THEN CONCAT('1:', z.pca_cnt)
										WHEN (z.aprv_cnt-z.pca_cnt) > 0 THEN CONCAT('2:', z.aprv_cnt-z.pca_cnt, '!1:', z.pca_cnt)
										WHEN (z.aprv_cnt-z.pca_cnt) < 0 THEN CONCAT('4:', z.pca_cnt-z.aprv_cnt, '!1:', z.aprv_cnt)
										ELSE '-' END diffInf
								FROM
								(
									SELECT 
										x.card_cd AS card_cd, 
										y.card_nm AS card_nm, 
										y.card_tel AS card_tel, 
										SUM(x.aprv_cnt) AS aprv_cnt, 
										SUM(x.aprv_amt) AS aprv_amt, 
										SUM(x.pca_cnt) AS pca_cnt, 
										SUM(x.pca_amt) AS pca_amt, 
										SUM(x.tot_fee) AS tot_fee, 
										SUM(x.vat_amt) AS vat_amt, 
										SUM(x.pay_amt) AS pay_amt
									FROM
									(
										SELECT a.card_cd, 
											IF(a.aprv_amt IS NULL, 0, 1) AS aprv_cnt, 
											IFNULL(a.aprv_amt,0) AS aprv_amt, 
											IF(b.pca_amt IS NULL, 0, 1) AS pca_cnt, 
											IFNULL(b.pca_amt,0) AS pca_amt, 
											IFNULL(b.tot_fee,0) AS tot_fee, 
											IFNULL(b.vat_amt,0) AS vat_amt, 
											IFNULL(b.pay_amt,0) AS pay_amt
										FROM cc_aprv_dtl a 
											LEFT JOIN cc_pca_dtl b 
												ON a.biz_num = b.biz_num AND a.tr_dt = b.org_tr_dt AND a.card_cd = b.card_cd 
												AND a.card_no = b.card_no AND a.aprv_no = b.aprv_no AND a.aprv_amt = b.pca_amt
										WHERE 
											a.biz_num = '#{bizNum}' 
											AND a.tr_dt = '#{bsDt}'
										
										UNION ALL
										
										SELECT a.card_cd, 
											IF(b.aprv_amt IS NULL, 0, 1) AS aprv_cnt, 
											IFNULL(b.aprv_amt,0) AS aprv_amt, 
											IF(a.pca_amt IS NULL, 0, 1) AS pca_cnt,
											IFNULL(a.pca_amt,0) AS pca_amt, 
											IFNULL(a.tot_fee,0) AS tot_fee, 
											IFNULL(a.vat_amt,0) AS vat_amt, 
											IFNULL(a.pay_amt,0) AS pay_amt
										FROM cc_pca_dtl a  
											LEFT JOIN cc_aprv_dtl b 
												ON a.biz_num = b.biz_num AND a.org_tr_dt = b.tr_dt AND a.card_cd = b.card_cd 
												AND a.card_no = b.card_no AND a.aprv_no =b.aprv_no AND b.aprv_amt = a.pca_amt
										WHERE 
											a.biz_num = '#{bizNum}'
											AND a.org_tr_dt = '#{bsDt}'
											AND b.aprv_no IS NULL
									) x,
									cc_card_comp_inf y 
									WHERE 
										x.card_cd = y.card_cd
									GROUP BY 
										x.card_cd, y.card_nm, y.card_tel
								) z
								ORDER BY z.card_cd
								`

// 매출캘린더 현금영수증 합계내역
var SelectCashDailyList string = `
								SELECT
									x.aprv_cnt AS aprvCnt, 
									x.aprv_amt AS aprvAmt, 
									x.aprv_cnt AS payCnt, 
									x.aprv_amt AS payAmt, 
									0 AS diffAmt, 
									0 AS totFee, 
									0 AS vatAmt, 
									x.aprv_amt AS outpExptAmt,
									CASE x.aprv_cnt WHEN 0 THEN '0:0' ELSE CONCAT('1:', x.aprv_cnt) END diffInf
								FROM (
									SELECT
										COUNT(tot_amt) AS aprv_cnt, 
										IFNULL(SUM(tot_amt),0) AS aprv_amt
									FROM cc_cash_dtl 
									WHERE 
										biz_num = '#{bizNum}' 
										AND tr_dt = '#{bsDt}'
								) x
								`

// 매출캘린더 지정카드사 매입내역 리스트
var SelectAprvDetailList string = `
								SELECT 
									a.tr_dt AS trDt, 
									a.tr_tm AS trTm, 
									CASE WHEN a.pca_yn = 'Y' THEN '정상매입'
										WHEN a.pca_yn = 'N' AND a.sts_cd = '1' AND
											 (DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), "%Y%m%d") != a.tr_dt AND
											  !(DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 4 DAY), "%Y%m%d") < a.tr_dt) AND DAYOFWEEK(a.tr_dt) in ('1','6','7')) THEN '매입제외'
										WHEN a.pca_yn = 'N' AND a.sts_cd IN ('2','3') THEN '매입제외'
										WHEN a.pca_yn = 'N' AND a.sts_cd = '1' THEN '매입대기'
										ELSE '확인필요' END diffNm,
									a.aprv_no AS aprvNo, 
									a.card_no AS cardNo, 
									CASE a.sts_cd WHEN '1' THEN '승인' WHEN '2' THEN '승인(-)' WHEN '3' THEN '취소(-)' END stsCd,
									a.aprv_amt AS aprvAmt, 
									IFNULL(b.pca_amt,0) AS pcaAmt,
									CASE WHEN a.inst_trm > 1 THEN CONCAT(a.inst_trm , '개월') ELSE '일시불' END AS instTrm,
									(SELECT cd_val FROM tb_cd_inf WHERE CD_ID = 'CARD_KND' AND cd_key = a.card_knd) AS cardKndNm,
									IFNULL(b.vat_amt,0) AS vatAmt, 
									IFNULL(b.tot_fee,0) AS totFee, 
									IFNULL(b.pay_amt,0) AS payAmt, 
									IFNULL(b.outp_expt_dt,'') AS outpExptDt
								FROM 
									cc_aprv_dtl a 
									LEFT JOIN cc_pca_dtl b 
										ON a.biz_num = b.biz_num AND a.tr_dt = b.org_tr_dt AND a.card_cd = b.card_cd 
										AND a.card_no = b.card_no AND a.aprv_no = b.aprv_no AND a.aprv_amt = b.pca_amt
								WHERE 
									a.biz_num = '#{bizNum}'
									AND a.tr_dt = '#{bsDt}' 
									AND a.card_cd = '#{cardCd}'
									
								UNION ALL
								
								SELECT 
									a.tr_dt AS trDt, 
									'235959' AS trTm, 
									'추가매입' AS diffNm,
									a.aprv_no AS aprvNo, 
									a.card_no AS cardNo, 
									CASE a.sts_cd WHEN '1' THEN '승인' WHEN '2' THEN '승인(-)' WHEN '3' THEN '취소(-)' END stsCd,  
									a.pca_amt AS aprvAmt, 
									a.pca_amt AS pcaAmt,
									CASE WHEN b.inst_trm > 1 THEN CONCAT(b.inst_trm , '개월') ELSE '일시불' END AS instTrm,
									(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'CARD_KND' AND cd_key = a.card_knd) AS cardKndNm,
									a.vat_amt AS vatAmt, 
									a.TOT_FEE AS totFee, 
									a.PAY_AMT AS payAmt, 
									a.OUTP_EXPT_DT AS outpExptDt
								FROM 
									cc_pca_dtl a 
									LEFT JOIN cc_aprv_dtl b 
										ON a.biz_num = b.biz_num AND a.org_tr_dt = b.tr_dt AND a.card_cd = b.card_cd 
										AND a.card_no = b.card_no AND a.aprv_no = b.aprv_no AND a.pca_amt = b.aprv_amt
								WHERE 
									a.biz_num = '#{bizNum}'
									AND a.org_tr_dt = '#{bsDt}' 
									AND a.card_cd = '#{cardCd}'
									AND b.aprv_no IS NULL 
								`

// 매출캘린더 혐금영수증 내역 리스트
var SelectCashDetailList string = `
								SELECT
									tr_dt AS trDt,
									tr_tm AS trTm, 
									'현금거래' AS diffNm,
									aprv_no AS aprvNo,
									isu_mtd AS cardNo,
									CASE sts_cd WHEN '1' THEN '승인' WHEN '2' THEN '승인(-)' WHEN '3' THEN '취소(-)' END stsCd,
									tot_amt AS aprvAmt, 
									tot_amt AS pcaAmt,
									'' AS instTrm,
									'현금' AS cardKndNm,
									vat_amt AS vatAmt,
									0 AS totFee,
									0 AS payAmt,
									'' AS outpExptDt
								FROM 
									cc_cash_dtl
								WHERE 
									biz_num = '#{bizNum}'
									AND tr_dt = '#{bsDt}'
								`

// 입금캘린더 월별 합계리스트
var SelectPayCalendarSumList string = `
									SELECT 
										tr_month AS trMonth, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt,
										SUM(z.real_in_amt) - SUM(Z.outp_expt_amt) AS diffAmt,
										CASE WHEN SUM(z.outp_expt_amt) = SUM(z.real_in_amt) THEN '0' 
											WHEN SUM(z.outp_expt_amt) > SUM(z.real_in_amt) THEN '1' ELSE '2' END diffColor 
									FROM (
										SELECT 
											SUBSTR(dt, 1, 6) AS tr_month, 
											IFNULL(SUM(pay_amt),0) AS outp_expt_amt,
											0 AS real_in_amt
										FROM 
											cc_date_info 
											LEFT JOIN cc_pca_dtl ON dt = outp_expt_dt AND biz_num = '#{bizNum}'
										WHERE 
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}' 
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
											dt BETWEEN '#{startDt}' 
											AND '#{endDt}' 
										GROUP BY SUBSTR(dt, 1, 6)
									) z
									GROUP BY z.tr_month
									ORDER BY z.tr_month DESC
									`

// 입금 캘린더
var SelectPayCalendarList string = `
									SELECT 
										CAST(FORMAT(@RN := @RN + 1, 0) as unsigned) AS rNum,
										tr_dt AS trDt, 
										SUM(z.outp_expt_amt) AS outpExptAmt, 
										SUM(z.real_in_amt) AS realInAmt, 
										SUM(z.real_in_amt) - SUM(z.outp_expt_amt) AS diffAmt,
										CASE WHEN SUM(z.outp_expt_amt) = SUM(z.real_in_amt) THEN '0' 
											WHEN SUM(z.outp_expt_amt) > SUM(z.real_in_amt) THEN '1' ELSE '2' END diffColor,
										SUM(day_color) AS dayColor
									FROM (
										SELECT 
											a.dt AS tr_dt, 
											IF(a.dt >= DATE_FORMAT(NOW(), '%Y%m%d'), 0, IFNULL(SUM(b.pay_amt),0)) AS outp_expt_amt, 
											0 AS real_in_amt,
											0 AS day_color
										FROM 
											cc_date_info a 
											LEFT JOIN cc_pca_dtl b 
												ON a.dt = b.outp_expt_dt AND b.biz_num = '#{bizNum}'
										WHERE 
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt

										UNION ALL

										SELECT total_date as tr_dt,
										0 as outp_expt_amt,
										0 as real_in_amt,
										CASE WHEN DAY = 7 then 2
										WHEN datekind = 'W' then 1
										ELSE 3 END day_color 
										FROM 
											sys_week_date 
										WHERE 
											total_date 
										BETWEEN '#{startDt}' 
										AND '#{endDt}'
										
										UNION ALL
										
										SELECT 
											a.dt AS tr_dt, 
											0 AS outp_expt_amt,
											IFNULL(SUM(c.real_pay_amt),0) AS real_in_amt,
											0 AS day_color
										FROM 
											cc_date_info a 
											LEFT JOIN cc_pay_dtl c 
												ON a.dt = c.pay_dt AND c.biz_num = '#{bizNum}'
										WHERE  
											a.dt BETWEEN '#{startDt}' 
											AND '#{endDt}'
										GROUP BY a.dt
									) z INNER JOIN (SELECT @RN := 0) r
									GROUP BY z.tr_dt
									ORDER BY z.tr_dt
									`

// 입금캘린더 카드사별 입금내역 리스트
var SelectPayDailyList string = `
								SELECT 
									a.card_cd AS cardCd, 
									a.card_nm AS cardNm, 
									a.pca_cnt AS pcaCnt, 
									a.pca_amt AS pcaAmt,
									a.tot_fee AS totFee, 
									a.vat_amt AS vatAmt,
									a.outp_expt_amt AS outpExptAmt, 
									IFNULL(b.real_in_amt, 0) AS realInAmt,
									IFNULL(b.real_in_amt, 0) - a.outp_expt_amt AS diffAmt,
									CASE WHEN IFNULL(b.real_in_amt, 0) = a.outp_expt_amt THEN '일치' 
										WHEN IFNULL(b.real_in_amt, 0) < a.outp_expt_amt THEN '일부입금' ELSE '초과입금' END diffNm, 
									CASE WHEN IFNULL(b.real_in_amt, 0) = a.outp_expt_amt THEN '0' 
										WHEN IFNULL(b.real_in_amt, 0) < a.outp_expt_amt THEN '1' ELSE '2' END diffColor 
								FROM 
								(
									SELECT 
										card_cd, 
										card_nm, 
										COUNT(outp_expt_dt) AS pca_cnt, 
										SUM(pca_amt) AS pca_amt,
										SUM(tot_fee) AS tot_fee, 
										SUM(vat_amt) AS vat_amt, 
										SUM(pay_amt) AS outp_expt_amt
									FROM 
										cc_pca_dtl
									WHERE
										biz_num = '#{bizNum}'
										AND outp_expt_dt = '#{bsDt}'
									GROUP BY card_cd
								) a
								LEFT JOIN 
								(
									SELECT
										(SELECT y.card_cd FROM cc_card_comp_inf y WHERE y.card_nm = z.card_nm) AS card_cd, 
										IFNULL(SUM(real_pay_amt),0) AS real_in_amt
									FROM
										cc_pay_dtl z
									WHERE  
										biz_num = '#{bizNum}'
										AND pay_dt = '#{bsDt}'
									GROUP BY card_cd
								) b ON a.card_cd = b.card_cd
								ORDER BY a.card_cd
								`

// 입금캘린더 카드사별 입금내역 리스트
var SelectPayDailyListHome string = `
								SELECT 
									a.card_cd AS cardCd, 
									a.card_nm AS cardNm,
									a.mer_no AS merNo,
									a.pca_cnt AS pcaCnt, 
									a.pca_amt AS pcaAmt,
									a.tot_fee AS totFee,
									a.outp_expt_amt AS outpExptAmt
								FROM 
								(
									SELECT 
										card_cd, 
										card_nm,
										mer_no,
										COUNT(outp_expt_dt) AS pca_cnt, 
										SUM(pca_amt) AS pca_amt,
										SUM(tot_fee) AS tot_fee, 
										SUM(pay_amt) AS outp_expt_amt
									FROM 
										cc_pca_dtl
									WHERE
										biz_num = '#{bizNum}'
										AND outp_expt_dt = '#{bsDt}'
									GROUP BY card_cd
								) a
								LEFT JOIN 
								(
									SELECT
										(SELECT y.card_cd FROM cc_card_comp_inf y WHERE y.card_nm = z.card_nm) AS card_cd, 
										IFNULL(SUM(real_pay_amt),0) AS real_in_amt
									FROM
										cc_pay_dtl z
									WHERE  
										biz_num = '#{bizNum}'
										AND pay_dt = '#{bsDt}'
									GROUP BY card_cd
								) b ON a.card_cd = b.card_cd
								ORDER BY a.card_cd
								`

// 입금캘린더 카드사별 입금내역 리스트
var SelectPayDailyListHome2 string = `
								SELECT 
									a.card_cd AS cardCd, 
									a.card_nm AS cardNm,
									a.mer_no AS merNo,
									a.pca_cnt AS pcaCnt, 
									a.pca_amt AS pcaAmt,
									a.tot_fee AS totFee,
									a.outp_expt_amt AS outpExptAmt
								FROM 
								(
									SELECT 
										card_cd, 
										card_nm,
										mer_no,
										COUNT(outp_expt_dt) AS pca_cnt, 
										SUM(pca_amt) AS pca_amt,
										SUM(tot_fee) AS tot_fee, 
										SUM(pay_amt) AS outp_expt_amt
									FROM 
										cc_pca_dtl
									WHERE
										biz_num = '#{bizNum}'
										AND outp_expt_dt = '#{bsDt2}'
									GROUP BY card_cd
								) a
								LEFT JOIN 
								(
									SELECT
										(SELECT y.card_cd FROM cc_card_comp_inf y WHERE y.card_nm = z.card_nm) AS card_cd, 
										IFNULL(SUM(real_pay_amt),0) AS real_in_amt
									FROM
										cc_pay_dtl z
									WHERE  
										biz_num = '#{bizNum}'
										AND pay_dt = '#{bsDt2}'
									GROUP BY card_cd
								) b ON a.card_cd = b.card_cd
								ORDER BY a.card_cd
								`

// 입금캘린더 지정카드사 입금내역 합계
var SelectPayDetailSum string = `
								SELECT 
									a.card_cd AS cardCd, 
									a.card_nm AS cardNm, 
									a.pca_cnt AS pcaCnt, 
									a.pca_amt AS pcaAmt,
									a.tot_fee AS totFee, 
									a.vat_amt AS vatAmt,
									a.outp_expt_amt AS outpExptAmt, 
									'#{bsDt}' AS outpExptDt,
									IFNULL(b.real_in_amt, 0) AS realInAmt,
									IFNULL(b.real_in_amt, 0) - a.OUTP_EXPT_AMT AS diffAmt,
									CASE WHEN IFNULL(b.real_in_amt, 0) = a.outp_expt_amt THEN '일치' 
										WHEN IFNULL(b.real_in_amt, 0) < a.outp_expt_amt THEN '일부입금' ELSE '초과입금' END diffNm, 
									CASE WHEN IFNULL(b.real_in_amt, 0) = a.outp_expt_amt THEN '0' 
										WHEN IFNULL(b.real_in_amt, 0) < a.outp_expt_amt THEN '1' ELSE '2' END diffColor 
								FROM 
								(
									SELECT 
										card_cd, 
										card_nm, 
										COUNT(outp_expt_dt) AS pca_cnt, 
										SUM(pca_amt) AS pca_amt,
										SUM(tot_fee) AS tot_fee, 
										SUM(vat_amt) AS vat_amt, 
										SUM(pay_amt) AS outp_expt_amt
									FROM 
										cc_pca_dtl
									WHERE
										biz_num = '#{bizNum}'
										AND outp_expt_dt = '#{bsDt}'
									GROUP BY card_cd
								) a
								LEFT JOIN 
								(
									SELECT
										(SELECT y.card_cd FROM cc_card_comp_inf y WHERE y.card_nm = z.card_nm) AS card_cd, 
										IFNULL(SUM(real_pay_amt),0) AS real_in_amt
									FROM
										cc_pay_dtl z
									WHERE  
										BIZ_NUM = '#{bizNum}'
										AND PAY_DT = '#{bsDt}'
									GROUP BY card_cd
								) b ON a.card_cd = b.card_cd
								WHERE 
									a.card_cd = '#{cardCd}'
								ORDER BY a.card_cd
								`

// 입금캘린더 지정카드사 입금내역 리스트
var SelectPayDetailList string = `
								SELECT 
									a.tr_dt AS trDt, 
									a.pca_amt AS pcaAmt, 
									a.aprv_no as aprvNo,
									(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'OK_CAN_DIV' AND cd_key = a.sts_cd) as stsCd,
									(SELECT cd_val FROM tb_cd_inf WHERE cd_id = 'CARD_KND' AND cd_key = a.card_knd ) AS cardKndNm,
									CASE WHEN INST_TRM > 1 THEN CONCAT(b.INST_TRM , '개월') ELSE '일시불' END instTrm,
									a.card_no AS cardNo, 
									a.pay_amt AS outpExptAmt
								FROM 
									cc_pca_dtl a 
									LEFT JOIN cc_aprv_dtl b 
										ON a.biz_num = b.biz_num 
										AND a.org_tr_dt = b.tr_dt 
										AND a.card_no = b.card_no 
										AND a.aprv_no = b.aprv_no 
										AND a.pca_amt = b.aprv_amt
								WHERE 
									a.BIZ_NUM = '#{bizNum}'
									AND a.OUTP_EXPT_DT = '#{bsDt}'
									AND a.CARD_CD = '#{cardCd}' 
								`
