
package reviews



var SelectCompInfo string = `
							SELECT a.rest_id
								, a.biz_num
								, a.comp_nm
								, b.lat
								, b.lng
								, b.addr 
							FROM cc_comp_inf a 
							left JOIN priv_rest_info b ON a.rest_id = b.rest_id  
							WHERE 
							a.comp_sts_cd=1 
							AND length(b.lat) > 0 
							AND b.lng IS NOT null
							AND a.rest_id = '#{restId}'
							`


var SelectStoreInfo string = `
							SELECT rest_id
									, biz_num
									, comp_nm
									, IFNULL(baemin_id,"") as baeminId
									, IFNULL(yogiyo_id,"") as yogiyoId
									, IFNULL(naver_id,"") as naverId
									, IFNULL(COUPANG_ID,"") as coupangId 
							FROM b_store
							WHERE
								rest_id = '#{restId}'
               				 `

var SelectStoreInfo_v2 string = `
							SELECT rest_id
									, biz_num
									, comp_nm
									, IFNULL(baemin_id,"") as baeminId
									, IFNULL(yogiyo_id,"") as yogiyoId
									, IFNULL(naver_id,"") as naverId
									, IFNULL(COUPANG_ID,"") as coupangId 
							FROM b_store
							WHERE 
								 biz_num = '#{bizNum}'
               				 `


var SelectBaeminInfo string = `
							SELECT comp_nm AS name , category_name_kr AS category, logo_url AS logo, address AS addr 
							FROM a_baemin 
							WHERE 
								baemin_id = '#{baeminId}'
`

var SelectNaverInfo string = `
							SELECT comp_nm AS NAME, ctg AS category, jibun_address AS addr, road_address AS road_addr 
							FROM a_naver 
							WHERE 
								naver_id = '#{naverId}'
`

var SelectYogiyoInfo string = `
							SELECT comp_nm AS NAME, categories AS category, thumbnail_url AS logo, address AS addr
							FROM a_yogiyo 
							WHERE 
								yogiyo_id = '#{yogiyoId}'
`



var SelectCoupangInfo string = `
							SELECT  NAME AS NAME
									,categories AS category
									,BRANDLOGO_PATH AS logo
									,concat(ADDRESS,' ',ADDRESS_DETAIL) AS addr
							FROM a_coupang
							WHERE
								COUPANG_ID = '#{coupangId}'
`


var SelectBlackFilter string = `SELECT ifnull(RATING,'') as rating
								,ifnull(KEYWORD,'') as keyword
								FROM b_store_etc
								WHERE 
								REST_ID='#{restId}' 
               				 `




var UpdateStoreBlackFilter string = `UPDATE b_store_etc SET 
										keyword =trim('#{keyword}')
										,rating= '#{rating}'
									WHERE
									REST_ID='#{restId}'
									`
var InsertStoreBlackFilter string = `INSERT INTO b_store_etc
									(
									REST_ID
									, RATING
									, KEYWORD
								    , BIZ_NUM
									)
									VALUES (
									'#{restId}'
									,trim('#{keyword}')
									,'#{rating}'
									,'#{bizNum}'
									)
							`







var SelectReviewRatingContentCnt string = `
							select count(*)  as TOTAL_COUNT
							from (
								(SELECT "b" as type, contents AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, MENU_NAME AS menuNM
								FROM a_baemin_review 
								WHERE 
									baemin_id ='#{baeminId}' 
								AND date 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								UNION ALL
								(SELECT "n" as type, body AS content, date_format(created,'%Y.%m.%d') as date, rating 
								, '' AS menuNM
								FROM a_naver_review 
								WHERE 
									naver_id = '#{naverId}' 
								AND date_format(created,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								UNION ALL
								(SELECT "cp" as type, review_text AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, order_menus AS menuNM
								FROM a_coupang_review 
								WHERE 
									coupang_id = '#{coupangId}' 
								AND date_format(date,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								UNION ALL
								(SELECT "y" as type, COMMENT AS content, date_format(time,'%Y.%m.%d') as date, rating 
								, MENU_SUMMARY AS menuNM
								FROM a_yogiyo_review 
								WHERE 
									yogiyo_id = '#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								) AS A
                             `

var SelectReviewRatingContent string = `
							select *
							from (
								(SELECT "b" as type, contents AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, MENU_NAME AS menuNM
								, member_name AS user
								, (SELECT COUNT(member_no) FROM a_baemin_review AS aa WHERE a.member_no = aa.member_no) AS visitCount
								, MEMBER_NO AS memberNo
								FROM a_baemin_review  as a
								WHERE 
									baemin_id ='#{baeminId}' 
								AND date 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								UNION ALL
								(SELECT "n" as type, body AS content, date_format(created,'%Y.%m.%d') as date, rating 
								, '' AS menuNM
								, AUTHOR_NICKNAME as user
								,(SELECT COUNT(author_id) FROM a_naver_review AS aa WHERE a.author_id = aa.author_id) AS visitCount
								, '' AS memberNo
								FROM a_naver_review  as a
								WHERE 
									naver_id = '#{naverId}' 
								AND date_format(created,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								UNION ALL
								(SELECT "cp" as type, review_text AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, order_menus AS menuNM
								, writer as user
								, 1 AS visitCount
								, '' AS memberNo
								FROM a_coupang_review 
								WHERE 
									coupang_id = '#{coupangId}' 
								AND date_format(date,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
								UNION ALL
								(SELECT "y" as type, COMMENT AS content, date_format(time,'%Y.%m.%d') as date, rating 
								, MENU_SUMMARY AS menuNM
								, phone as user
								, 1 AS visitCount
								, '' AS memberNo
								FROM a_yogiyo_review 
								WHERE 
									yogiyo_id = '#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND rating = '#{rating}'
								)
									) AS A
								ORDER BY date DESC
                             `



var SelectCustomReviewRatingContentCnt string = `
							select count(*)  as TOTAL_COUNT
							from (
								(SELECT "b" as type, contents AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, MENU_NAME AS menuNM
								FROM a_baemin_review 
								WHERE 
									baemin_id ='#{baeminId}' 
								AND date 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating}) 
									OR contents REGEXP '#{keyword}')
								)
								UNION ALL
								(SELECT "n" as type, body AS content, date_format(created,'%Y.%m.%d') as date, rating 
								, '' AS menuNM
								FROM a_naver_review 
								WHERE 
									naver_id = '#{naverId}' 
								AND date_format(created,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating}) 
									OR body REGEXP '#{keyword}')
								)
								UNION ALL
								(SELECT "cp" as type, review_text AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, order_menus AS menuNM
								FROM a_coupang_review 
								WHERE 
									coupang_id = '#{coupangId}' 
								AND date_format(date,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating}) 
									OR review_text REGEXP '#{keyword}')
								)
								UNION ALL
								(SELECT "y" as type, COMMENT AS content, date_format(time,'%Y.%m.%d') as date, rating 
								, MENU_SUMMARY AS menuNM
								FROM a_yogiyo_review 
								WHERE 
									yogiyo_id = '#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating}) 
									OR COMMENT REGEXP '#{keyword}')
								)
								) AS A
                             `

var SelectCustomReviewRatingContent string = `
							select *
							from (
								(SELECT "b" as type, contents AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, MENU_NAME AS menuNM
								, member_name AS user
								, (SELECT COUNT(member_no) FROM a_baemin_review AS aa WHERE a.member_no = aa.member_no) AS visitCount
								, MEMBER_NO AS memberNo
								FROM a_baemin_review  as a
								WHERE 
									baemin_id ='#{baeminId}' 
								AND date 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating}) 
									OR contents REGEXP '#{keyword}')
								)
								UNION ALL
								(SELECT "n" as type, body AS content, date_format(created,'%Y.%m.%d') as date, rating 
								, '' AS menuNM
								, AUTHOR_NICKNAME as user
								,(SELECT COUNT(author_id) FROM a_naver_review AS aa WHERE a.author_id = aa.author_id) AS visitCount
								, '' AS memberNo
								FROM a_naver_review  as a
								WHERE 
									naver_id = '#{naverId}' 
								AND date_format(created,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating}) 
									OR body REGEXP '#{keyword}')
								)
								UNION ALL
								(SELECT "cp" as type, review_text AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, order_menus AS menuNM
								, writer as user
								, 1 AS visitCount
								, '' AS memberNo
								FROM a_coupang_review 
								WHERE 
									coupang_id = '#{coupangId}' 
								AND date_format(date,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating})  
									OR review_text REGEXP '#{keyword}')
								
								)
								UNION ALL
								(SELECT "y" as type, COMMENT AS content, date_format(time,'%Y.%m.%d') as date, rating 
								, MENU_SUMMARY AS menuNM
								, phone as user
								, 1 AS visitCount
								, '' AS memberNo
								FROM a_yogiyo_review 
								WHERE 
									yogiyo_id = '#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND (rating in (#{rating})  
									OR COMMENT REGEXP '#{keyword}')
								)
									) AS A
								ORDER BY date DESC
                             `



var SelectCustomReviewRatingContentCnt_Norating string = `
							select count(*)  as TOTAL_COUNT
							from (
								(SELECT "b" as type, contents AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, MENU_NAME AS menuNM
								FROM a_baemin_review 
								WHERE 
									baemin_id ='#{baeminId}' 
								AND date 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND contents REGEXP '#{keyword}'
								)
								UNION ALL
								(SELECT "n" as type, body AS content, date_format(created,'%Y.%m.%d') as date, rating 
								, '' AS menuNM
								FROM a_naver_review 
								WHERE 
									naver_id = '#{naverId}' 
								AND date_format(created,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND body REGEXP '#{keyword}'
								)
								UNION ALL
								(SELECT "cp" as type, review_text AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, order_menus AS menuNM
								FROM a_coupang_review 
								WHERE 
									coupang_id = '#{coupangId}' 
								AND date_format(date,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND review_text REGEXP '#{keyword}'
								)
								UNION ALL
								(SELECT "y" as type, COMMENT AS content, date_format(time,'%Y.%m.%d') as date, rating 
								, MENU_SUMMARY AS menuNM
								FROM a_yogiyo_review 
								WHERE 
									yogiyo_id = '#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND COMMENT REGEXP '#{keyword}'
								)
								) AS A
                             `

var SelectCustomReviewRatingContent_Norating string = `
							select *
							from (
								(SELECT "b" as type, contents AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, MENU_NAME AS menuNM
								, member_name AS user
								, (SELECT COUNT(member_no) FROM a_baemin_review AS aa WHERE a.member_no = aa.member_no) AS visitCount
								, MEMBER_NO AS memberNo
								FROM a_baemin_review  as a
								WHERE 
									baemin_id ='#{baeminId}' 
								AND date 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND contents REGEXP '#{keyword}'
								)
								UNION ALL
								(SELECT "n" as type, body AS content, date_format(created,'%Y.%m.%d') as date, rating 
								, '' AS menuNM
								, AUTHOR_NICKNAME as user
								,(SELECT COUNT(author_id) FROM a_naver_review AS aa WHERE a.author_id = aa.author_id) AS visitCount
								, '' AS memberNo
								FROM a_naver_review  as a
								WHERE 
									naver_id = '#{naverId}' 
								AND date_format(created,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND body REGEXP '#{keyword}'
								)
								UNION ALL
								(SELECT "cp" as type, review_text AS content, date_format(date,'%Y.%m.%d') as date, rating 
								, order_menus AS menuNM
								, writer as user
								, 1 AS visitCount
								, '' AS memberNo
								FROM a_coupang_review 
								WHERE 
									coupang_id = '#{coupangId}' 
								AND date_format(date,'%Y%m%d') 
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND review_text REGEXP '#{keyword}'
								)
								UNION ALL
								(SELECT "y" as type, COMMENT AS content, date_format(time,'%Y.%m.%d') as date, rating 
								, MENU_SUMMARY AS menuNM
								, phone as user
								, 1 AS visitCount
								, '' AS memberNo
								FROM a_yogiyo_review 
								WHERE 
									yogiyo_id = '#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')
									BETWEEN '#{startDt}' 
									AND '#{endDt}'
								AND COMMENT REGEXP '#{keyword}'
								)
									) AS A
								ORDER BY date DESC
                             `



var SelectReivewRating string = `SELECT a.rating 
		,ifnull(bamin_cnt,0) AS bamin_cnt
		,ifnull(yogiyo_cnt,0) AS yogiyo_cnt
		,ifnull(naver_cnt,0) AS naver_cnt
		,ifnull(coupang_cnt,0) AS coupang_cnt
		,(ifnull(bamin_cnt,0)+ifnull(yogiyo_cnt,0)+ifnull(naver_cnt,0)+ifnull(coupang_cnt,0)) AS tot_cnt
		FROM
		(
			SELECT 1 AS rating UNION ALL
			SELECT 2 UNION ALL
			SELECT 3 UNION ALL
			SELECT 4 UNION ALL
			SELECT 5 ) A 
		LEFT OUTER JOIN (	
								SELECT rating, COUNT(*) AS bamin_cnt 
								FROM a_baemin_review 
								WHERE 
								BAEMIN_ID='#{baeminId}'
								GROUP BY rating 
							) AS b ON a.rating = b.rating
		LEFT OUTER JOIN (	
								SELECT rating, COUNT(*) AS yogiyo_cnt  
								FROM a_yogiyo_review 
								WHERE 
								yogiyo_id = '#{yogiyoId}'
								GROUP BY rating 
							) AS c ON a.rating = c.rating
		LEFT OUTER JOIN (	
								SELECT z.rat AS rating, SUM(z.cnt) AS naver_cnt FROM(
									SELECT round(rating) AS rat, COUNT(*) as cnt
									FROM a_naver_review 
									WHERE 
									naver_id = '#{naverId}' 
									GROUP BY rating 
								) z GROUP BY z.rat
							) AS d ON a.rating = d.rating
		LEFT OUTER JOIN (	
								SELECT rating, COUNT(*) AS coupang_cnt  
								FROM a_coupang_review 
								WHERE 
								COUPANG_ID = '#{coupangId}'
								GROUP BY rating
							) AS e ON a.rating = e.rating
		order by a.rating asc
		`

var SelectReivewRatingMonth string = `SELECT DATE_FORMAT(A.TOTAL_dATE,'%m월') AS MONTH
		,IFNULL(bamin_avg,0) AS bamin_avg
		,IFNULL(yogiyo_avg,0) AS yogiyo_avg
		,IFNULL(naver_avg,0) AS naver_avg
		,IFNULL(coupang_avg,0) AS coupang_avg
		,ROUND((IFNULL(bamin_avg,0)+IFNULL(yogiyo_avg,0)+IFNULL(naver_avg,0)+IFNULL(coupang_avg,0))
		/(IFNULL(bamin_cnt,0)+IFNULL(yogiyo_cnt,0)+IFNULL(naver_cnt,0)+IFNULL(coupang_cnt,0)) ,2) AS total_avg
		FROM sys_week_date AS A
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(DATE,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS bamin_avg 
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS bamin_cnt
								FROM a_baemin_review 
								WHERE 
								BAEMIN_ID='#{baeminId}' 
								AND DATE >= '#{startDate}' 
								AND DATE <= '#{endDate}' 
								GROUP BY DATE_FORMAT(DATE,'%m') 
								) AS B ON DATE_FORMAT(A.TOTAL_dATE,'%m') = B.t_month
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(time,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS yogiyo_avg
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS yogiyo_cnt
								FROM a_yogiyo_review 
								WHERE 
								yogiyo_id='#{yogiyoId}' 
								AND date_format(time,'%Y%m%d')   >='#{startDate}' 
								AND date_format(time,'%Y%m%d') <= '#{endDate}' 
								GROUP BY DATE_FORMAT(time,'%m') 
								) AS c ON DATE_FORMAT(A.TOTAL_dATE,'%m') = c.t_month
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(created,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS naver_avg
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS naver_cnt
								FROM a_naver_review 
								WHERE 
								naver_id='#{naverId}' 
								AND date_format(created,'%Y%m%d')   >='#{startDate}' 
								AND date_format(created,'%Y%m%d') <= '#{endDate}' 
								GROUP BY DATE_FORMAT(created,'%m') 
								) AS d ON DATE_FORMAT(A.TOTAL_dATE,'%m') = d.t_month
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(date,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS coupang_avg
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS coupang_cnt
								FROM a_coupang_review 
								WHERE 
								coupang_id='#{coupangId}' 
								AND date_format(date,'%Y%m%d')   >= '#{startDate}' 
								AND date_format(date,'%Y%m%d') <= '#{endDate}' 
								GROUP BY DATE_FORMAT(date,'%m') 
								) AS e ON DATE_FORMAT(A.TOTAL_dATE,'%m') = e.t_month
		WHERE  
				TOTAL_DATE >= '#{startDate}' 
				AND TOTAL_DATE <= '#{endDate}' 
		GROUP BY DATE_FORMAT(A.TOTAL_dATE,'%m월')
		ORDER BY A.TOTAL_DATE ASC
	`

var SelectContentList string = `
							SELECT TYPE
									,TITLE
								--	,DESCRIPTION
									,URL
									,IMAGE_URL
									,VIDEO_URL
									,content_id
							from b_contents
							WHERE use_yn='Y'
							AND START_DATE <= date_format(NOW(),'%Y-%m-%d')
							AND END_DATE >= date_format(NOW(),'%Y-%m-%d')
							order by rand() LIMIT 5
`
var SelectTipList string = `
							SELECT   
								TITLE
								,IFNULL(CONTENT,'') AS CONTENT
								,LINK_URL
							from sys_boards
							WHERE use_yn='Y'
							AND BOARD_TYPE='2'
							AND B_KIND='1'
							AND START_DATE <= date_format(NOW(),'%Y-%m-%d')
							AND END_DATE >= date_format(NOW(),'%Y-%m-%d')
							order by rand() LIMIT 1
`



var SelectBillingInfo string = `SELECT A.END_DATE
						,B.ITEM_NAME
						,B.ITEM_DESC
						,DATE_FORMAT(A.START_DATE,'%Y년 %m월 %d일') AS START_DATE
						,DATE_FORMAT(A.END_DATE,'%Y년 %m월 %d일') AS END_DATE
						,A.PAY_YN
						FROM e_billing AS A
						INNER JOIN e_billing_item AS B ON A.ITEM_CODE = B.ITEM_CODE
						WHERE 
						A.STORE_ID = '#{restId}' 
						AND END_DATE >= SYSDATE()
`




var SelectReivewRatingMonth_v2 string = `SELECT DATE_FORMAT(A.TOTAL_dATE,'%m월') AS MONTH
		,IFNULL(ROUND((IFNULL(bamin_avg,0)+IFNULL(yogiyo_avg,0)+IFNULL(naver_avg,0)+IFNULL(coupang_avg,0))
		/(IFNULL(bamin_cnt,0)+IFNULL(yogiyo_cnt,0)+IFNULL(naver_cnt,0)+IFNULL(coupang_cnt,0)) ,2),0) AS total_avg
		,(IFNULL(b_cnt,0)+IFNULL(y_cnt,0)+IFNULL(n_cnt,0)+IFNULL(c_cnt,0)) AS total_cnt
		FROM sys_week_date AS A
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(DATE,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS bamin_avg 
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS bamin_cnt
								,COUNT(*) AS b_cnt
								FROM a_baemin_review 
								WHERE 
								BAEMIN_ID='#{baeminId}' 
								AND left(DATE,6) = '#{bsDt}'
								GROUP BY DATE_FORMAT(DATE,'%m') 
								) AS B ON DATE_FORMAT(A.TOTAL_dATE,'%m') = B.t_month
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(time,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS yogiyo_avg
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS yogiyo_cnt
								,COUNT(*) AS y_cnt
								FROM a_yogiyo_review 
								WHERE 
								yogiyo_id='#{yogiyoId}' 
								AND date_format(time,'%Y%m')  = '#{bsDt}'
								GROUP BY DATE_FORMAT(time,'%m') 
								) AS c ON DATE_FORMAT(A.TOTAL_dATE,'%m') = c.t_month
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(created,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS naver_avg
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS naver_cnt
								,COUNT(*) AS n_cnt
								FROM a_naver_review 
								WHERE 
								naver_id='#{naverId}' 
								AND date_format(created,'%Y%m')  = '#{bsDt}'
								GROUP BY DATE_FORMAT(created,'%m') 
								) AS d ON DATE_FORMAT(A.TOTAL_dATE,'%m') = d.t_month
		LEFT OUTER JOIN (
								SELECT DATE_FORMAT(date,'%m') AS t_month ,round(sum(rating)/COUNT(*),2) AS coupang_avg
								,CASE  WHEN IFNULL(sum(rating),0) > 0 THEN 1 ELSE 0 END AS coupang_cnt
								,COUNT(*) AS c_cnt
								FROM a_coupang_review 
								WHERE 
								coupang_id='#{coupangId}' 
								AND date_format(date,'%Y%m')  = '#{bsDt}'
								GROUP BY DATE_FORMAT(date,'%m') 
								) AS e ON DATE_FORMAT(A.TOTAL_dATE,'%m') = e.t_month
		WHERE  
				left(TOTAL_DATE,6) = '#{bsDt}'
		GROUP BY DATE_FORMAT(A.TOTAL_dATE,'%m월')
		ORDER BY A.TOTAL_DATE ASC
	`