package users

var SelectUserLoginCheck string = `SELECT A.USER_ID as userId
										,ATLOGIN_YN as atloginYn
										,A.USE_YN
										,ifnull(B.REG_ID,'null') AS regId
										,A.USER_TY
										,ifnull(A.USER_BIRTH,'') AS USER_BIRTH
										,A.LOGIN_ID
										,A.USER_NM
										,A.HP_NO
								    FROM priv_user_info AS A
									LEFT OUTER JOIN sys_reg_info AS B ON A.USER_ID = B.USER_ID
									WHERE 
									LOGIN_ID ='#{loginId}'
									AND LOGIN_PW='#{password}'
									ORDER BY A.USER_TY DESC
									`


var InserPushData string = `INSERT INTO SYS_REG_INFO
							(
							USER_ID
							, REG_ID
							, OS_TY
							, REG_DATE
							, LOGIN_YN
							)
							VALUES
							(
							'#{userId}'
							, '#{regId}'
							, '#{osTy}'
							, DATE_FORMAT(now(), '%Y%m%d%H%i%s')
							, '#{loginYn}'
							)
							`
var UpdatePushData string = `UPDATE SYS_REG_INFO
							SET
								REG_ID = '#{regId}',
								OS_TY = '#{osTy}',
								LOGIN_YN = '#{loginYn}',
								REG_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
							USER_ID = '#{userId}'
							`



var SelectKakaoUserLoginCheck string = `SELECT A.USER_ID as userId
										,ATLOGIN_YN as atloginYn
										,A.USE_YN
										,ifnull(B.REG_ID,'null') AS regId
										,A.USER_TY
										,ifnull(A.USER_BIRTH,'') AS USER_BIRTH
										,A.LOGIN_ID
										,A.USER_NM
										,A.HP_NO
								    FROM priv_user_info AS A
									LEFT OUTER JOIN sys_reg_info AS B ON A.USER_ID = B.USER_ID
									WHERE 
									KAKAO_KEY = '#{socialToken}'
									`

var SelectNaverUserLoginCheck string = `SELECT A.USER_ID as userId
										,ATLOGIN_YN as atloginYn
										,A.USE_YN
										,ifnull(B.REG_ID,'null') AS regId
										,A.USER_TY
										,ifnull(A.USER_BIRTH,'') AS USER_BIRTH
										,A.LOGIN_ID
										,A.USER_NM
										,A.HP_NO
								    FROM priv_user_info AS A
									LEFT OUTER JOIN sys_reg_info AS B ON A.USER_ID = B.USER_ID
									WHERE 
									NAVER_KEY ='#{socialToken}'
									`

var SelectAppleUserLoginCheck string = `SELECT A.USER_ID as userId
										,ATLOGIN_YN as atloginYn
										,A.USE_YN
										,ifnull(B.REG_ID,'null') AS regId
										,A.USER_TY
										,ifnull(A.USER_BIRTH,'') AS USER_BIRTH
										,A.LOGIN_ID
										,A.USER_NM
										,A.HP_NO
								    FROM priv_user_info AS A
									LEFT OUTER JOIN sys_reg_info AS B ON A.USER_ID = B.USER_ID
									WHERE 
									APPLE_KEY ='#{socialToken}'
									`

var SelectEmailDupCheck string = `SELECT count(*) as emailCnt
								FROM priv_user_info
								WHERE 
								LOGIN_ID ='#{email}'`

var SelectLoginIdDupCheck string = `SELECT count(*) as loginIdCnt
								FROM priv_user_info
								WHERE 
								LOGIN_ID ='#{bizNum}'`




var SelectBizNumDupCheck string = `SELECT count(*) as bizCnt
								FROM priv_rest_info
								WHERE 
								BUSID ='#{bizNum}'`


var SelectCreatUserSeq string = `SELECT CONCAT('U',IFNULL(LPAD(MAX(SUBSTRING(USER_ID, -10)) + 1, 10, 0), '0000000001')) as newUserId
								 FROM priv_user_info
								 `


var InserCreateUser string = `INSERT INTO priv_user_info
									(
										USER_ID,
										USER_NM,
										LOGIN_ID,
										LOGIN_PW,
										USER_TY,
										if #{email} != '' then EMAIL,
										HP_NO,
										ATLOGIN_YN,
										GEOLOC_YN,
										PUSH_YN,
										if #{kakaoPw} != '' then KAKAO_PW,
										if #{kakaoKey} != '' then KAKAO_KEY,
										if #{applePw} != '' then APPLE_PW,
										if #{appleKey} != '' then APPLE_KEY,
										if #{naverPw} != '' then NAVER_PW,
										if #{naverKey} != '' then NAVER_KEY,
										if #{recomCode} != '' then RECOM_CODE,
										if #{channelCode} != '' then CHANNEL_CODE,
										if #{userBirth} != '' then USER_BIRTH,
										USE_YN,
										JOIN_DATE
									)
									VALUES
									(
										'#{userId}'
										, '#{userNm}'
										, '#{loginId}'
										, '#{loginPw}'
										, '#{userTy}'
										, '#{email}'
										, '#{userTel}'
										, '#{atLoginYn}'
										, 'Y'
										, '#{pushYn}'
										, '#{kakaoPw}'
										, '#{kakaoKey}'
										, '#{applePw}'
										, '#{appleKey}'
										, '#{naverPw}'
										, '#{naverKey}'
										, '#{recomCode}'
										, '#{channelCode}'
										, '#{userBirth}'
										, 'S'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										)`

var InsertTermsUser string = `INSERT INTO b_user_terms
							(USER_ID, 
							TERMS_OF_SERVICE, 
							TERMS_OF_PERSONAL, 
							TERMS_OF_PAYMENT, 
							TERMS_OF_BENEFIT, 
							REG_DATE
							)
							VALUES (
							'#{userId}'
							,'#{termsOfService}'
							,'#{termsOfPersonal}'
							,'#{termsOfPayment}'
							,'#{termsOfBenefit}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')						
							)
							`

var SelectUserInfo string = `SELECT USER_NM AS userName
							, HP_NO as userTel
							, ifnull(USER_BIRTH,'') as birthday
							, EMAIL AS email
							, RECOM_CODE as recomCode
							FROM priv_user_info
							WHERE 
							USER_ID ='#{userId}'  
							`


/////////////////////////////////////////////////////////




var SelectUserSetupInfo string = `SELECT C.PUSH_YN
									,IFNULL(A.RECOM_CODE,'') AS RECOM_CODE
							FROM priv_rest_info AS A
							INNER JOIN priv_rest_user_info  AS B ON  a.rest_id = b.rest_id
							INNER JOIN priv_user_info AS c ON b.user_id = c.user_id
							WHERE 
								C.USER_ID='#{userId}'
							`


var UpdateSetupInfo string = `UPDATE priv_user_info SET 
												 MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')	
                                                , PUSH_YN = '#{pushYn}'
							WHERE 
							USER_ID ='#{userId}'
							`

var UpdateUserInfo string = `UPDATE priv_user_info SET 
												 MOD_DATE= DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')	
                                                , USER_BIRTH = '#{birthday}'
												, USER_NM='#{userName}'
												, USER_TEL='#{userTel}'
												, USE_YN='#{useYn}'
												, LOGIN_PW='#{loginPw}'
							WHERE 
							USER_ID ='#{userId}'
							`




var SelectKakaoTokenDupCheck string = `SELECT count(*) as tokenCnt
								FROM priv_user_info
								WHERE 
								KAKAO_KEY ='#{socialToken}'
								AND USER_TY ='1'
								`

var SelectNaverTokenDupCheck string = `SELECT count(*) as tokenCnt
								FROM priv_user_info
								WHERE 
								NAVER_KEY ='#{socialToken}'
								AND USER_TY ='1'
								`

var SelectAppleTokenDupCheck string = `SELECT count(*) as tokenCnt
								FROM priv_user_info
								WHERE 
								APPLE_KEY ='#{socialToken}'
								AND USER_TY ='1'
								`

var UpdateUserPasswordChange string = `UPDATE priv_user_info SET LOGIN_PW = '#{password}'
													,MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
									 WHERE 
									 USER_ID ='#{userId}'
									 `


var SelectJoinCheck string = `SELECT A.USER_ID
									,A.USER_NM
									,A.HP_NO
									,A.LOGIN_ID
									,IFNULL(B.REST_ID,'NONE') AS STORE_ID
									,IFNULL(C.REST_NM,'NONE') AS REST_NM
									,IFNULL(C.CEO_NM,'NONE') AS CEO_NM
								FROM priv_user_info AS A
								LEFT OUTER JOIN priv_rest_user_info AS b ON a.user_id = b.user_id
								LEFT OUTER JOIN priv_rest_info AS c ON b.rest_id = c.rest_id
								WHERE 
								USER_TY='1'
								AND A.USER_ID='#{userId}'
								`



var InsertLoginAccess string =`	INSERT INTO SYS_LOG_ACCESS (
											USER_ID
											,ADDED_DATE
											,LOG_IN_OUT
											,IP
											,SUCC_YN 
											,SERVICE 
											,TYPE
								) VALUES (
									'#{loginId}'
									,SYSDATE()
									,'#{logInOut}'
									,'#{ip}'
									,'#{succYn}'
									,'#{osTy}'
									,'#{type}'
								)
								`

var SelectUserStoreInfo string = `SELECT A.REST_ID AS STORE_ID
											,B.REST_NM AS STORE_NM
											,B.BUSID  AS BIZ_NUM
									FROM priv_rest_user_info AS A
									INNER JOIN priv_rest_info AS B ON A.REST_ID = B.REST_ID
									WHERE 
									A.USER_ID ='#{userId}'
									`

var SelectStoreCashInfo string = `SELECT IFNULL(LN_ID,'') AS LN_ID
										,IFNULL(LN_PSW,'') AS LN_PSW
										,LN_JOIN_STS_CD 
										,LN_AUTH_FAIL
										,UNIX_TIMESTAMP(LN_FAIL_DT) AS LN_FAIL_DT
										,IFNULL(TIMESTAMPDIFF(MINUTE,SYSDATE(),DATE_ADD(LN_FAIL_DT, INTERVAL 30 MINUTE)),0) AS LN_REMAIN_TIME
										,IFNULL(HOMETAX_ID,'') AS HOMETAX_ID
										,IFNULL(HOMETAX_PSW,'') AS HOMETAX_PSW
										,HOMETAX_JOIN_STS_CD
										,HOMETAX_AUTH_FAIL
										,UNIX_TIMESTAMP(HOMTAXT_FAIL_DT) AS HOMTAXT_FAIL_DT
										,IFNULL(TIMESTAMPDIFF(MINUTE,SYSDATE(),DATE_ADD(HOMTAXT_FAIL_DT, INTERVAL 30 MINUTE)),0) AS HOMTAXT_REMAIN_TIME
									FROM cc_comp_inf
									WHERE 
										REST_ID ='#{storeId}'
									AND BIZ_NUM='#{bizNum}'
									`






var SelectUserIdSearch string = `SELECT 
									CASE 	WHEN kakao_key IS NOT NULL THEN ''
											WHEN apple_key IS NOT NULL THEN ''
											WHEN naver_key IS NOT NULL THEN ''
									ELSE CONCAT(LEFT(CONCAT(SUBSTRING(LOGIN_ID,1,2),'*********************')
											,CASE INSTR(LOGIN_ID,"@")-1 WHEN  -1 THEN LENGTH(LOGIN_ID)  ELSE INSTR(LOGIN_ID,"@")-1 END) 
											,SUBSTRING(LOGIN_ID,INSTR(LOGIN_ID,"@"),30))
									 END AS LOGIN_ID
								   ,DATE_FORMAT(JOIN_DATE, '%Y.%m.%d') AS JOIN_DATE
								   ,CASE 	WHEN kakao_key IS NOT NULL THEN 'KAKAO'
												WHEN apple_key IS NOT NULL THEN 'APPLE'
												WHEN naver_key IS NOT NULL THEN 'NAVER'
									ELSE 'ID' END AS LOGIN_TYPE
								FROM priv_user_info
								WHERE 
								USER_NM='#{userNm}'
								AND HP_NO='#{userTel}'
								AND USER_TY='1'
								`



var SelectUserPwSearch string = `SELECT USER_ID
									 ,CASE 	WHEN kakao_key IS NOT NULL THEN 'KAKAO'
												WHEN apple_key IS NOT NULL THEN 'APPLE'
												WHEN naver_key IS NOT NULL THEN 'NAVER'
										ELSE 'ID' END AS LOGIN_TYPE
									FROM priv_user_info
									WHERE
									LOGIN_ID = '#{loginId}'
									AND HP_NO='#{userTel}'
									AND USER_TY='1'
								`



var UpdateUserPassWd string = `UPDATE priv_user_info
							SET
								LOGIN_PW = '#{loginPw}',
								MOD_DATE = DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							WHERE 
							USER_ID = '#{userId}'
								`
