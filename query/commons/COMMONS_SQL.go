package commons

var PagingQuery string =`
					LIMIT #{pageSize}
					OFFSET #{offSet}
`

var SelectBoardList string = `SELECT
							BOARD_ID
							,TITLE
							,BOARD_TYPE
							,LINK_URL
							,DATE_FORMAT(REG_DATE, '%Y.%m.%d') AS regDate	
							FROM sys_boards
							WHERE 
							START_DATE <= NOW() AND END_DATE >=NOW()
							AND B_KIND IN ('A','1')
							`

var SelectCategoryList string = `SELECT CATEGORY_ID
									,CATEGORY_NM 
								FROM b_category
								WHERE 
								CATEGORY_GRP_CODE='#{grpCode}'
								AND USE_YN='Y'
							`
var SelectCodeist string = `SELECT CODE_ID
								,CODE_NM
								FROM b_code
								WHERE
								CATEGORY_ID='#{categoryId}'
								AND USE_YN='Y'
							`

// 버전
var SelectVersion string = `SELECT 
							  VERSION AS versionCode 
							, CASE WHEN AUTO_YN ='Y' THEN 'true' ELSE 'false' END AS  isRequireUpdate
							FROM sys_version_info
							WHERE 
							USE_YN = 'Y' 
							AND OS_TY = '#{osTy}'
							AND APP_TY = '#{appTy}'
							ORDER BY VERSION_ID DESC
							LIMIT 1
							`


var InsertPushLog string = `INSERT INTO sys_log_push
						(
						NAME
						, TITLE
						, BODY
						, APP_TY
						, REG_DATE
						, REG_ID
						)
						VALUES
						(
                          '#{name}'
						, '#{title}'
						, '#{body}'
						, '#{appTy}'
						, NOW()
						, '#{regId}'
						)
						`

var SelectPushGrp string = `SELECT
								A.USER_ID
								, A.REG_ID
								, A.OS_TY
								, A.LOGIN_YN
								, C.PUSH_YN
							FROM 
								SYS_REG_INFO A,
								PRIV_USER_INFO C
							WHERE 
								EXISTS (SELECT USER_ID FROM PRIV_GRP_USER_INFO B
										WHERE 
										A.USER_ID = B.USER_ID AND B.GRP_ID = '#{grpId}'
										AND B.GRP_AUTH = '0'
								)
							AND A.LOGIN_YN = 'Y'
							AND A.USER_ID = C.USER_ID
							`


var SelectPushUser string = `SELECT
							A.USER_ID
							, A.REG_ID
							, A.OS_TY
							, A.REG_DATE
							, A.LOGIN_YN
							, B.PUSH_YN
							FROM
							SYS_REG_INFO A,
							PRIV_USER_INFO B
							WHERE
							A.USER_ID = B.USER_ID
							AND B.USE_YN = 'Y'
							AND A.USER_ID = '#{userId}'
							`

var SelectPushBizNum string = `SELECT C.USER_ID
							, C.REG_ID
							, C.OS_TY
							, C.LOGIN_YN
							, B.PUSH_YN
							FROM PRIV_REST_INFO AS A
							INNER JOIN PRIV_REST_USER_INFO AS B ON A.REST_ID = B.REST_ID
							INNER JOIN SYS_REG_INFO AS C ON B.USER_ID = C.USER_ID
							WHERE
							BUSID='#{bizNum}'
							AND B.REST_AUTH=0
							`


var SelectPushRest string = `SELECT C.USER_ID
							, C.REG_ID
							, C.OS_TY
							, C.LOGIN_YN
							, B.PUSH_YN
							FROM PRIV_REST_INFO AS A
							INNER JOIN PRIV_REST_USER_INFO AS B ON A.REST_ID = B.REST_ID
							INNER JOIN SYS_REG_INFO AS C ON B.USER_ID = C.USER_ID
							WHERE
							A.REST_ID='#{restId}'
							AND B.REST_AUTH=0
							`


var SelectPushMsgInfo string = `SELECT 	MSG
									,TITLE
								FROM dar_msg_info
								WHERE 
								MSG_CODE='#{msgCode}'
								LIMIT 1;
							`




var SelectCompAuthInfo string =`SELECT LN_AUTH_FAIL
										,UNIX_TIMESTAMP(SYSDATE()) AS LN_FAIL_DT
										,HOMETAX_AUTH_FAIL
										,UNIX_TIMESTAMP(SYSDATE()) AS HOMTAXT_FAIL_DT
										,TIMESTAMPDIFF(MINUTE,SYSDATE(),DATE_ADD(SYSDATE(), INTERVAL 30 MINUTE)) AS lnRemainTime
										,TIMESTAMPDIFF(MINUTE,SYSDATE(),DATE_ADD(SYSDATE(), INTERVAL 30 MINUTE)) AS homtaxtRemainTime
								FROM cc_comp_inf
								WHERE
									BIZ_NUM='#{bizNum}'

								`

var UpdateLnAuthFail string = `UPDATE cc_comp_inf SET LN_AUTH_FAIL= #{lnAuthFailCnt}
										,LN_FAIL_DT = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s')
										,LN_ID ='#{loginId}'
										,LN_PSW ='#{password}'
										,LN_JOIN_STS_CD =3
										WHERE
										BIZ_NUM='#{bizNum}'
								`


var UpdateLnAuthSuccess string = `UPDATE cc_comp_inf SET LN_AUTH_FAIL= 0
										,LN_FAIL_DT = ''
										,LN_ID ='#{loginId}'
										,LN_PSW ='#{password}'
										,LN_JOIN_STS_CD =1
										,SVC_OPEN_DT= DATE_FORMAT(SYSDATE(),'%Y%m%d')
										WHERE
										BIZ_NUM='#{bizNum}'
								`


var UpdateHomeTaxAuthFail string = `UPDATE cc_comp_inf SET HOMETAX_AUTH_FAIL = #{homeTaxAuthFailCnt}
										,HOMTAXT_FAIL_DT = DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s')
										,HOMETAX_ID= '#{loginId}'
										,HOMETAX_PSW='#{password}'
										,HOMETAX_JOIN_STS_CD=3
										WHERE
										BIZ_NUM='#{bizNum}'
								`

var UpdateHomeTaxAuthSuccess string = `UPDATE cc_comp_inf SET HOMETAX_JOIN_STS_CD=1
										,HOMETAX_ID= '#{loginId}'
										,HOMETAX_PSW='#{password}'
										,HOMETAX_AUTH_FAIL=0
										,HOMTAXT_FAIL_DT=''
										,HOMETAX_OPEN_DT= DATE_FORMAT(SYSDATE(),'%Y%m%d')
										WHERE
										BIZ_NUM='#{bizNum}'
								`


var InsertAlimTalkLog string = `INSERT INTO sys_alimtalk_log
								(
									 HP_NO
									,USER_ID
									if #{messageId} != '' then ,MESSAGE_ID	 
									, TEMPLATE_CODE
									, TEMPLATE_NAME
									, SEND_DATE
									if #{result} != '' then ,RESULT
									)
									VALUES 
									(
                                      '#{hpNo}'
									,  '#{userId}'
									,  '#{messageId}'
									, '#{templateCode}'
									, '#{templateNm}'
									, DATE_FORMAT(SYSDATE(),'%Y%m%d%H%i%s')
									, '#{result}'
									)
								`


var UpdateAlimTalkLog string = `UPDATE  sys_alimtalk_log SET 
										RESULT='#{result}'
								WHERE 
									MESSAGE_ID ='#{messageId}'
								`

var SelctAlimTalkSeq string = `SELECT MAX(seq)+1 as alimSeq
								FROM sys_alimtalk_log
								`


var SelectSendAlimRest string = `SELECT CASE WHEN USE_YN='S' THEN '1' 
												WHEN USE_YN='Y' AND LN_YN <>'1' THEN '2'
												ELSE '3'  END  AS T_CODE
											,USER_ID
											,USER_NM
											,HP_NO
											,USE_YN
											,REST_AUTH
								FROM (
						
										SELECT A.USER_ID
												,A.USER_NM
												,A.HP_NO
												,A.USE_YN
												,IFNULL((SELECT LN_JOIN_STS_CD FROM cc_comp_inf AS aa WHERE c.rest_id = aa.rest_id),0) AS LN_YN
												,IFNULL(B.REST_AUTH,'2') AS REST_AUTH
										FROM priv_user_info AS a 
										LEFT OUTER JOIN priv_rest_user_info AS b ON a.user_id = b.user_id 
										LEFT OUTER JOIN priv_rest_info AS c ON b.rest_id = c.rest_id
										WHERE A.user_ty='1'
										AND LEFT(A.JOIN_DATE,8) IN ('20210714','20210721')
										) AS ZZ
								WHERE REST_AUTH IN('0','2')
								`