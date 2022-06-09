package delions


var InserCreateUserDelion string = `INSERT INTO priv_user_info
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
										USER_BIRTH,
										USE_YN,
										CHANNEL_CODE,
										JOIN_DATE
									)
									VALUES
									(
										'#{userId}'
										, '#{ceo_name}'
										, '#{biz_number}'
										, '#{loginPw}'
										, '#{userTy}'
										, '#{email}'
										, '#{phone}'
										, '#{atLoginYn}'
										, 'Y'
										, '#{pushYn}'
										, '#{birth_date}'
										, 'Y'
										, 'DELION'
										, DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
										)`


var InsertTermsUserDelion string = `INSERT INTO b_user_terms
							(USER_ID, 
							TERMS_OF_SERVICE, 
							TERMS_OF_PERSONAL, 
							TERMS_OF_PAYMENT, 
							TERMS_OF_BENEFIT, 
							REG_DATE
							)
							VALUES (
							'#{userId}'
							,'Y'
							,'Y'
							,'Y'
							,'Y'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')						
							)
							`

var InsertStoreDelion string = `INSERT INTO priv_rest_info
							(
							  REST_ID,
							  REST_NM,
							  BUSID,
                              CEO_NM,
							  if #{category} != '' then CATEGORY,
							  if #{kind} != '' then BUETY,
							  if #{address} != '' then ADDR,
							  AUTH_STAT,
							  USE_YN,
							  CEO_BIRTHDAY,
							  TEL,
							  if #{email} != '' then EMAIL,
							  REG_DATE,
							  CHANNEL_CODE
							)
							VALUES (
							'#{storeId}'
							,'#{biz_name}'
							,'#{biz_number}'
							,'#{ceo_name}'
							,'#{category}'
							,'#{kind}'
							,'#{address}'
							,'1'
							,'Y'
							,'#{birth_date}'
							,'#{phone}'
							,'#{email}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')		
							,'#{channelCode}'
							)
							`

var InsertStoreUserDelion string = `INSERT INTO priv_rest_user_info
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


var InsertCashDelion string = `INSERT INTO cc_comp_inf (
							REST_ID,
							BIZ_NUM,
							COMP_NM,
							COMP_STS_CD,
							SVC_OPEN_DT,
							LN_FIRST_YN, 
							LN_JOIN_TY, 
							LN_ID, 
							LN_PSW, 
							LN_JOIN_STS_CD,
							REG_DT, 
							SER_ID, 
							PUSH_DT
						)
						VALUES (
							'#{storeId}'
							,'#{biz_number}'
							,'#{biz_name}'
							,'1'
							,DATE_FORMAT(NOW(), '%Y%m%d')
							,'N'
							,'3'
							,'#{lnId}'
							,'#{lnPsw}'
							,'1'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							,'CASH_00'
							,''
						)
						`




var InsertDelionTemp string = `INSERT INTO temp_delion
							(
							  	biz_name, 
								biz_number, 
								cardsales_pass, 
								ceo_name, 
								phone, 
								birth_date, 
								address, 
								biz_kind, 
								email,
								reg_date
							)
							VALUES (
							'#{biz_name}'
							,'#{biz_number}'
							,'#{cardsales_pass}'
							,'#{ceo_name}'
							,'#{phone}'
							,'#{birth_date}'
							,'#{address}'
							,'#{biz_kind}'
							,'#{email}'
							,DATE_FORMAT(NOW(), '%Y%m%d%H%i%s')
							)
							`