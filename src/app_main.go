package cashApi

import (
	// login 및 기본

	apis "cashApi/src/controller/API/ETC"
	kakaos "cashApi/src/controller/API/KAKAO"
	pushs "cashApi/src/controller/API/PUSH"
	smss "cashApi/src/controller/API/SMS"
	billings "cashApi/src/controller/billing"
	books "cashApi/src/controller/books"
	commons "cashApi/src/controller/commons"
	datas "cashApi/src/controller/datas"
	delions "cashApi/src/controller/delions"
	"cashApi/src/controller/homes"
	parthners "cashApi/src/controller/parthner"
	"cashApi/src/controller/review/tracking"
	"cashApi/src/controller/reviewMng"
	sales "cashApi/src/controller/sales"
	stores "cashApi/src/controller/stores"
	users "cashApi/src/controller/users"
	views "cashApi/src/controller/views"
	"github.com/labstack/echo/v4/middleware"
	//공통 api
	"cashApi/src/controller/cls"

	"github.com/labstack/echo/v4"
)

var lprintf func(int, string, ...interface{}) = cls.Lprintf

func SvcSetting(e *echo.Echo, fname string) *echo.Echo {
	lprintf(4, "[INFO] sql start \n")

	// 캐쉬컴바인 api
	cash := e.Group("/api/cash")
	cash.HEAD("/health", users.BaseUrl) // 헬스체크
	cash.GET("/health", users.BaseUrl)  // 헬스체크
	cash.PUT("/health", users.BaseUrl)  // 헬스체크

	// login
	cash.POST("/login", users.LoginDarayo)                      // 로그인
	cash.POST("/logOut", users.LoginOut)                        // 로그인
	cash.POST("/socialLogin", users.LoginDarayoSocial)          // 소셜 로그인
	cash.POST("/smsRequest", smss.SendSmsConfirm)               // 문자 인증 요청 // sms api 개발 필요
	cash.POST("/smsCheck", smss.ConfirmCheck)                   // 문자 인증 확인
	cash.GET("/emailCheck", users.GetEmailDupCheck)             // 이메일 중복 체크 (아이디 체크)
	cash.GET("/socialTokenCheck", users.GetSocialTokenDupCheck) // 소셜 로그인 토근 중복확인

	cash.PUT("/pushInfo", users.PushInfo) // 푸쉬 정보 업데이트

	//2021.07.14 추가 개발 시작
	cash.GET("/bizNumCheck", users.GetBizNumDupCheck)           // 사업자 번호 중복 체크 (아이디 체크)
	cash.POST("/bizNumJoin", users.SetUserBizNumJoin)           // 회원 가입 - 사업자 번호로 가입
	cash.GET("/cardsalesLoginCheck", users.CardsalesLoginCheck) // 여신협회 로그인 조회
	cash.GET("/hometaxLoginCheck", users.HometaxLoginCheck)     // 홈텍스 로그인 조회
	//2021.07.14 추가 개발 끝

	cash.POST("/join", users.SetUserJoin)                  // 회원 가입 - 사용자 정보 입력
	cash.POST("/joinStep2", users.SetUserJoinStep2)        // 회원 가입 - 가맹점 정보 입력
	cash.POST("/storeInfoUpdate", users.SetStorInfoUpdate) // 회원 가입 - 대표 && 계좌정보 업데이트
	cash.POST("/socialJoin", users.SetUserJoin)            // 소셜 회원 가입
	cash.POST("/cashJoin", users.SetCashJoin)              // 여신 or 홈텍스 가입
	cash.POST("/cashUpdate", users.SetCashModify)          // 여신 or 홈텍스 정보 변경

	// 공통
	common := cash.Group("/commons")
	common.GET("/category/:grpCode", commons.GetCategoryList) //공통코드
	common.GET("/code/:categoryId", commons.GetCodeList)      //코드
	common.GET("/versions/latest", commons.GetVersionsLatest) //앱 최신 버전 호출
	common.GET("/push", commons.SendCommonPush)               //푸쉬 전송
	common.GET("/pushMsg", commons.SendCommonPushMsg)         //푸쉬 전송

	// 홈화면
	home := cash.Group("/homes")
	home.Use(middleware.JWT([]byte("darago")))
	home.GET("/:storeId", homes.GetHomeData)          //홈화면 데이터
	home.GET("/getSalesList", homes.GetSalesInfo)     // 매출정보
	home.GET("/getPayDailyList", homes.GetPayDayList) // 카드사별 조회

	//설정
	home.PUT("/:storeId/storeInfo", users.SetStoreInfo)               // 설정 - 매장 정보 관리 업데이트
	home.PUT("/:storeId/ceoInfo", users.SetCeoInfo)                   // 설정 - 대표자 정보 관리 업데이트
	home.GET("/:storeId/storeInfo", users.GetStoreInfo)               // 설정 - 매장 정보 관리
	home.GET("/:storeId/storeServiceInfo", users.GetStoreServiceInfo) // 가입 대행 서비스 정보
	home.PUT("/:userId/userInfo", users.SetUserInfo)                  // 내정보 변경
	home.GET("/:userId/setupInfo", users.GetSetupInfo)                // 설정 - 메인
	home.PUT("/:userId/setupInfo", users.SetSetupInfo)                // 설정 - 메인 업데이트
	home.GET("/board", homes.GetBoardList)                            // 공지사항 && 이벤트 리스트
	home.GET("/:storeId/storeChargeList", users.GetStoreChargeList)   // 설정 - 충전 정보
	home.PUT("/:storeId/storeCharge", users.SetStoreCharge)           // 설정 - 충전 정보 수정
	home.PUT("/:storeId/storeChargeYn", users.SetStoreChargeYn)       // 설정 - 선불 충전 사용 여부

	// 장부
	book := cash.Group("/books")
	book.Use(middleware.JWT([]byte("darago")))
	book.GET("/:storeId", books.GetLinkBookList)                    // 연결된 장부 목록
	book.GET("/:storeId/order", books.GetBooksOrders)               // 장부별 주문
	book.GET("/:storeId/:orderNo", books.GetOrderInfo)              // 주문 상세
	book.GET("/:storeId/payment", books.GetBooksPayments)           // 결제 조회
	book.GET("/:storeId/charge", books.GetStoreChargeAmt)           // 장부 충전금액 조회
	book.POST("/:storeId/charging", books.SetStoreCharging)         // 매장 충전
	book.DELETE("/:storeId/charging", books.SetStoreChargingCancel) // 매장 충전 취소
	book.GET("/:storeId/account", books.GetUnPaidList)              // 매장 정산 대상 조회
	book.POST("/:storeId/account", books.SetPaidOk)                 // 매장 정산
	book.DELETE("/:storeId/account", books.SetPaidCancel)           // 매장 정산 취소
	book.DELETE("/orderCancel/:orderNo", books.SetOrderCancel)      // 주문 취소

	// 매장관리
	store := cash.Group("/stores")
	store.Use(middleware.JWT([]byte("darago")))
	store.GET("/:storeId/categories", stores.GetStoreCategories)                   // 카테고리 목록 조회
	store.POST("/:storeId/categories", stores.SetStoreInsertCategories)            // 카테고리 추가
	store.PUT("/:storeId/categories/:categoryId", stores.SetStoreUpdateCategories) // 카테고리 수정
	store.GET("/:storeId/menus", stores.GetStoreMenuList)                          // 메뉴  조회
	store.POST("/:storeId/menus", stores.SetStoreInsertMenu)                       // 메뉴  추가
	store.PUT("/:storeId/menus/:itemNo", stores.SetStoreUpdateMenu)                // 메뉴  수정
	//store.DELETE("/menus/:itemNo", stores.SetStoreDeleteMenu) // 메뉴  제거
	store.GET("/:storeId/services", stores.GetStoreServiceList) // 서비스 목록 조회
	store.PUT("/:storeId/services", stores.SetStoreService)     // 서비스 수정
	store.GET("/:storeId/info", stores.GetStoreInfo)            // 매장 관리

	//store.GET("/:storeId", books.GetLinkBookList) // 매장 상세보기
	//store.GET("/:storeId", books.GetLinkBookList) // 영업관리
	//store.GET("/:storeId", books.GetLinkBookList) // 영업관리 수정

	//기타 api
	etcs := e.Group("/api/etc")
	etcs.GET("/bizNumCheck", apis.BizNumCheck)       // 사업자 등록증 번호 조회
	etcs.GET("/acctNameSearch", apis.AcctNameSearch) // 계좌실명 조회
	etcs.GET("/holiday", apis.SetHoliday)            // 휴일 정보 업데이트
	etcs.POST("/simpleOrhomesder", apis.SimpleOrder) // 태블릿 주문
	etcs.POST("/coupon", apis.CouponUse)             // 쿠폰사용
	etcs.GET("/cardsalesLogin", apis.CardsalesLogin) // 여신협회 로그인 조회
	etcs.GET("/hometaxLogin", apis.HometaxLogin)     // 홈텍스 로그인 조회

	etcs.GET("/kakoalim", kakaos.SendKakaoAlim)    // 홈텍스 로그인 조회
	etcs.GET("/testkakoalim", kakaos.SendAlimRest) // 홈텍스 로그인 조회

	etcs.GET("/kakao/lastMonthReport", kakaos.LastMonthReport3) // 지난달 분석 보고서
	etcs.GET("/kakao/lastWeekReport", kakaos.LastWeekReport)    // 지난주 분석 보고서
	etcs.GET("/kakao/yesterdayReport", kakaos.YesterdayReport)  // 어제 분석 보고서
	etcs.GET("/kakao/stateReport", kakaos.StateReport)          // 가입상태 안내 보고서
	etcs.GET("/kakao/welcome", kakaos.WelcomeMessage)           // 회원가입 메시지

	//푸쉬 api
	push := e.Group("/api/push")
	push.GET("/sendPush", pushs.SendPush) // 푸쉬 보내기

	// 사용자
	user := cash.Group("/users")
	user.GET("/:storeId/cashInfo", users.GetCashInfo) // 캐쉬 인증 화면 데이터

	//sms api
	sms := e.Group("/api/sms")
	sms.GET("/sendSms", smss.SendSms)            // sms 보내기
	sms.POST("/smsRequest", smss.SendSmsConfirm) // 문자 인증 요청 // sms api 개발 필요
	sms.POST("/smsCheck", smss.ConfirmCheck)     // 문자 인증 확인

	//구독
	billing := e.Group("/api/billing")

	billing.GET("/:storeId", billings.GetBillingInfo) // 구독 정보 상세 조회
	//billing.GET("/:storeId/reg", billings.BillingCheck) // 구독 신청
	//billing.PUT("/:storeId", billings.GetBillingInfo)          // 구독 변경
	//billing.DELETE("/:storeId", billings.SetBillingCancel)           // 구독 취소
	//billing.POST("/billingAuth", billings.BillingAuth)               // 구독 승인
	billing.POST("/billingResult", billings.ViewBillingResult) // 구독 결과 - 카드 등록만
	//billing.POST("/billingPayResult", billings.ViewBillingPayResult) // 구독 결과 - 금액 결제와 카드 등록

	//빌링 관련 페이지
	billingPay := e.Group("/api/pay")
	billingPay.GET("/b_reg", billings.BillingRegView)
	billingPay.POST("/billingAuth", billings.BillingAuth)               // 구독 승인
	billingPay.DELETE("/:storeId", billings.SetBillingCancel)           // 구독 취소
	billingPay.POST("/billingRegResult", billings.ViewBillingPayResult) // 구독 결과 - 금액 결제와 카드 등록
	billingPay.GET("/billingSend", billings.CallPayment)                // 결제 요청
	billingPay.DELETE("/billingCancel", billings.CallPaymentCancel)     // 결제 취소

	// 매출정산
	sale := cash.Group("/sales")
	sale.Use(middleware.JWT([]byte("darago")))
	sale.GET("/getCardSum", sales.GetCardSum)               // 매출정산 - 카드승인 합계
	sale.GET("/getCardList", sales.GetCardList)             // 매출정산 - 카드승인 상세
	sale.GET("/getCashSum", sales.GetCashSum)               // 매출정산 - 현금영수증승인 합계
	sale.GET("/getCashList", sales.GetCashList)             // 매출정산 - 현금영수증승인 상세
	sale.GET("/getPaySum", sales.GetPaySum)                 // 매출정산 - 입금정보 합계
	sale.GET("/getPayList", sales.GetPayList)               // 매출정산 - 입금정보 상세
	sale.GET("/getAprvCalendar", sales.GetAprvCalendar)     // 매출정산 - 매출캘린더
	sale.GET("/getAprvDailyList", sales.GetAprvDailyList)   // 매출정산 - 매출캘린더 카드사별 매입내역 리스트
	sale.GET("/getAprvDetailList", sales.GetAprvDetailList) // 매출정산 - 매출캘린더 특정 카드사 매입내역 리스트
	sale.GET("/getPayCalendar", sales.GetPayCalendar)       // 매출정산 - 입금캘린더
	sale.GET("/getPayDailyList", sales.GetPayDailyList)     // 매출정산 - 입금캘린더 카드사별 매입내역 리스트
	sale.GET("/getPayDetailList", sales.GetPayDetailList)   // 매출정산 - 입금캘린더 특정 카드사 매입내역 리스트

	// 그래프 화면
	view := e.Group("/page")
	view.GET("/week", views.GetWeekView)   // 주간 매출분석 , 고객님 방문 분석, 요일별 매출 분석
	view.GET("/month", views.GetMonthView) // 주간 매출분석 , 고객님 방문 분석, 요일별 매출 분석

	view.GET("/pos_week", views.GetWeekViewPos)   // 주간 매출분석 , 고객님 방문 분석, 요일별 매출 분석 -- 포스용
	view.GET("/pos_month", views.GetMonthViewPos) // 주간 매출분석 , 고객님 방문 분석, 요일별 매출 분석 -- 포스용

	// 달아요 파트너
	parthner := e.Group("/parthner/p")

	//parthner.Use(middleware.JWT([]byte("darago")))
	parthner.GET("/login", parthners.ParthnerLogin)         // 로그인
	parthner.GET("/join", parthners.ParthnerJoin)           // 회원 가입
	parthner.GET("/joinStep1", parthners.ParthnerJoinStep1) // 가맹점 정보 입력
	parthner.GET("/joinStep2", parthners.ParthnerJoinStep2) // 인증하기
	parthner.GET("/guideCertify", parthners.GuideCertify)   // 인증 가이드 화면
	parthner.GET("/guidePartner", parthners.GuidePartner)   // 파트너 가이드 화면
	//cls.SetNotLoginUrl("/parthner/login")

	// 달아요 파트너 api
	parthnerApi := e.Group("/parthner/a")
	parthnerApi.POST("/loginOk", parthners.ParthnerLoginOk) // 로그인
	parthnerApi.POST("/home", parthners.Homedata)           //기본 데이터

	// 그래프   data
	data := e.Group("/data")
	data.GET("/week", datas.GetWeekData)
	//data.GET("/week1", datas.GetWeekData1)
	data.GET("/weekTimeAvg", datas.GetWeekAvgTimData)
	//data.GET("/month", datas.GetMonthData)
	data.GET("/month", datas.GetMonthData1)
	data.GET("/monthCompare", datas.GetMonthCompareData)
	//data.GET("/wordCloud", datas.GetWordCloud)

	data.GET("/pos_week", datas.GetPosWeekData)
	data.GET("/pos_weekTimeAvg", datas.GetPosWeekDataTimeAvg)
	data.GET("/pos_month", datas.GetPosMonthData)

	//딜리온 API
	delion := e.Group("/api/delion")
	delion.POST("/join", delions.SetDelionJoin) // 딜리온 회원가입 대행
	delion.POST("/joinTEST", delions.TEST11)    // 딜리온 회원가입 대행

	// 리뷰 API
	review := e.Group("/review")
	review.GET("/tracking", tracking.GetCustomerInfo)
	review.GET("/test", tracking.CrossTest)

	//리뷰 관리

	review.GET("/test", reviewMng.TestPage)                     //샘플
	review.GET("/reviewSetting", reviewMng.ReviewSetting)       // 리뷰 환경 설정
	review.GET("/reviewList", reviewMng.ReviewList)             // 전체리뷰
	review.GET("/customReviewList", reviewMng.CustomReviewList) // 관심리뷰
	review.GET("/reviewMain", reviewMng.ReviewMain)             // 리뷰분석

	review.POST("/api/getReviewMain", reviewMng.GetReviewMain)   // 리뷰 메인 api - 전체별점 , 월별 차트
	review.POST("/api/getReviewMain2", reviewMng.GetReviewMain2) // 리뷰 메인 api(2) - 리뷰 키워드, 꿀팁 , 컨텐츠, 워드클라우드

	review.POST("/api/getReviewSetting", reviewMng.GetReviewSetting)    // 리뷰 환경 설정값 api
	review.PUT("/api/keyword", reviewMng.SetKeywordSetup)               // 관심 별점, 키워드 업데이트
	review.PUT("/api/base", reviewMng.SetStoreReviewInfo)               // 매장주소 수정
	review.POST("/api/deliveryList", reviewMng.GetDeliveryList)         // 가맹점 배달업체 리스트
	review.PUT("/api/matching", reviewMng.SetCompSetting)               // 가맹점 배달업체 매칭
	review.POST("/api/reviewList", reviewMng.GetReviewList)             // 전체 리뷰 내역
	review.POST("/api/customReviewList", reviewMng.GetCustomReviewList) // 전체 리뷰 내역
	review.POST("/api/writer", reviewMng.GetReviewWriter)               // 리뷰 작성자 분석

	/////////////////////////////////////캐쉬 v2 시작  2021.11.01
	cash_v2 := e.Group("/api/cash/v2")

	cash_v2.GET("/idSearch", users.GetSearchId) // 회원 - 아이디 찾기
	cash_v2.GET("/pwSearch", users.GetSearchPw) // 회원 - 비밀번호  찾기
	cash_v2.PUT("/pwCh", users.SetChPw)         // 회원 - 비밀번호 변경

	order_v2 := cash_v2.Group("/order")
	order_v2.GET("/:restId/pickupList", stores.GetOrderPickupList)     // 포장주문현황
	order_v2.PUT("/:restId/pickupStatus", stores.SetOrderPickupStatus) // 포장주문 상태 업데이트

	cash_v2.POST("/getMonthly", datas.GetMonthly) // 월간 보고서 api

	view2 := e.Group("/page/v2")
	view2.GET("/monthly", views.GetMonthlyView_V2) // 월간 보고서

	cls.SetNotLoginUrl("/")
	lprintf(4, "[INFO] page start \n")

	return e

}
