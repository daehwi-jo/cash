package links

var SelectLinkInfo string = `SELECT 
							A.AGRM_ID AS LINK_ID
							, A.GRP_ID
							, FN_GET_GRPNAME(A.GRP_ID) AS BOOK_NM
							, A.REST_ID
							, FN_GET_RESTNAME(A.REST_ID) AS STORE_NM
							, A.REQ_STAT
							, FN_GET_CODENAME('AGRM_STAT', A.REQ_STAT) AS REQ_STAT_NM
							, A.REQ_TY
							, DATE_FORMAT(A.REQ_DATE, '%Y-%m-%d') AS REQ_DATE	
							, DATE_FORMAT(A.AUTH_DATE, '%Y-%m-%d') AS AUTH_DATE
							, A.PAY_TY
							, IFNULL(A.PREPAID_AMT, 0) AS PREPAID_AMT
							, IFNULL(A.PREPAID_POINT, 0) AS PREPAID_POINT
							FROM org_agrm_info  AS a
							INNER join PRIV_REST_INFO AS B on A.REST_ID = B.REST_ID or (B.FRAN_YN = 'Y' and B.FRAN_ID = A.REST_ID)
							WHERE
							B.REST_ID ='#{storeId}'
							AND A.GRP_ID='#{bookId}'
							`
var UpdateLink string = `UPDATE org_agrm_info SET 
							PREPAID_AMT = '#{prepaidAmt}'
							,PREPAID_POINT = #{prepaidPoint}
							,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
							AGRM_ID ='#{linkId}'
						`