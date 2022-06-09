var sessionStorage_transfer = function(event) {
    if(!event) { event = window.event; } // ie suq
    if(!event.newValue) return;          // do nothing if no value to work with
    if (event.key == 'getSessionStorage') {
        // another tab asked for the sessionStorage -> send it
        localStorage.setItem('sessionStorage', JSON.stringify(sessionStorage));
        // the other tab should now have it, so we're done with it.
        localStorage.removeItem('sessionStorage'); // <- could do short timeout as well.
    } else if (event.key == 'sessionStorage' && !sessionStorage.length) {
        // another tab sent data <- get it
        var data = JSON.parse(event.newValue);
        for (var key in data) {
            sessionStorage.setItem(key, data[key]);
        }
    }
};

// listen for changes to localStorage
if(window.addEventListener) {
    window.addEventListener("storage", sessionStorage_transfer, false);
} else {
    window.attachEvent("onstorage", sessionStorage_transfer);
};


// Ask other tabs for session storage (this is ONLY to trigger event)
if (!sessionStorage.length) {
    localStorage.setItem('getSessionStorage', 'foobar');
    localStorage.removeItem('getSessionStorage', 'foobar');
};



function  darago() {
    //XMLHttpRequest 객체 생성
    let xhr = new XMLHttpRequest();
    const lid = document.getElementById("lid").value;
    const lpwd = document.getElementById("lpwd").value;
    let spanErr =document.getElementById("span_err")
    spanErr.style.display="none";

    if (lid == ""){
        spanErr.style.display="block";
        spanErr.innerText="* 아이디를 입력해주세요.";
        return;
    }

    if (lpwd == ""){
        spanErr.style.display="block";
        spanErr.innerText="* 비밀번호를 입력해주세요.";
        return;
    }

    const e_lpwd =  SHA256(lpwd);
    //요청을 보낼 방식, 주소, 비동기여부 설정
    xhr.open('POST', '/parthner/a/loginOk', true);
    //HTTP 요청 헤더 설정
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.responseType='json';
    //요청 전송
    xhr.send("loginId="+lid+"&password="+e_lpwd+"&osTy=AT");
    //통신후 작업
    xhr.onload = function () {
        //통신 성공
        if (xhr.status == 200) {
          //  console.log(xhr.response);
           // console.log(xhr.response);

            if (xhr.response.resultCode == "00"){
                sessionStorage.setItem('Token', xhr.getResponseHeader("Token"));
                sessionStorage.setItem('bizNum', xhr.response.resultData.bizNum);
                sessionStorage.setItem('loginId', xhr.response.resultData.loginId);
                sessionStorage.setItem('storeId', xhr.response.resultData.storeId);
                sessionStorage.setItem('storeNm', xhr.response.resultData.storeNm);
                sessionStorage.setItem('userBirth', xhr.response.resultData.userBirth);
                sessionStorage.setItem('userNm', xhr.response.resultData.userNm);
                sessionStorage.setItem('userTel', xhr.response.resultData.userTel);
                sessionStorage.setItem('userId', xhr.response.resultData.userId);
                location.href="/parthner/p/guideCertify";
            }else if(xhr.response.resultCode == "01"){
                sessionStorage.setItem('userNm', xhr.response.resultData.userNm);
                sessionStorage.setItem('userTel', xhr.response.resultData.hpNo);
                sessionStorage.setItem('userId', xhr.response.resultData.userId);
                sessionStorage.setItem('loginId', xhr.response.resultData.loginId);
                location.href="/parthner/p/joinStep1";
            }
            else{
                let msg = xhr.response.resultMsg;
                if (msg == "가맹점 사용자가 아닙니다."){
                    alert(msg);
                    return;
                }else{
                    spanErr.style.display="block";
                    spanErr.innerText="* " + msg;
                    return;
                }
            }

        } else {
            //통신 실패

        }
    }
};




function SHA256(s){

    var chrsz   = 8;
    var hexcase = 0;

    function safe_add (x, y) {
        var lsw = (x & 0xFFFF) + (y & 0xFFFF);
        var msw = (x >> 16) + (y >> 16) + (lsw >> 16);
        return (msw << 16) | (lsw & 0xFFFF);
    }

    function S (X, n) { return ( X >>> n ) | (X << (32 - n)); }
    function R (X, n) { return ( X >>> n ); }
    function Ch(x, y, z) { return ((x & y) ^ ((~x) & z)); }
    function Maj(x, y, z) { return ((x & y) ^ (x & z) ^ (y & z)); }
    function Sigma0256(x) { return (S(x, 2) ^ S(x, 13) ^ S(x, 22)); }
    function Sigma1256(x) { return (S(x, 6) ^ S(x, 11) ^ S(x, 25)); }
    function Gamma0256(x) { return (S(x, 7) ^ S(x, 18) ^ R(x, 3)); }
    function Gamma1256(x) { return (S(x, 17) ^ S(x, 19) ^ R(x, 10)); }

    function core_sha256 (m, l) {

        var K = new Array(0x428A2F98, 0x71374491, 0xB5C0FBCF, 0xE9B5DBA5, 0x3956C25B, 0x59F111F1,
            0x923F82A4, 0xAB1C5ED5, 0xD807AA98, 0x12835B01, 0x243185BE, 0x550C7DC3,
            0x72BE5D74, 0x80DEB1FE, 0x9BDC06A7, 0xC19BF174, 0xE49B69C1, 0xEFBE4786,
            0xFC19DC6, 0x240CA1CC, 0x2DE92C6F, 0x4A7484AA, 0x5CB0A9DC, 0x76F988DA,
            0x983E5152, 0xA831C66D, 0xB00327C8, 0xBF597FC7, 0xC6E00BF3, 0xD5A79147,
            0x6CA6351, 0x14292967, 0x27B70A85, 0x2E1B2138, 0x4D2C6DFC, 0x53380D13,
            0x650A7354, 0x766A0ABB, 0x81C2C92E, 0x92722C85, 0xA2BFE8A1, 0xA81A664B,
            0xC24B8B70, 0xC76C51A3, 0xD192E819, 0xD6990624, 0xF40E3585, 0x106AA070,
            0x19A4C116, 0x1E376C08, 0x2748774C, 0x34B0BCB5, 0x391C0CB3, 0x4ED8AA4A,
            0x5B9CCA4F, 0x682E6FF3, 0x748F82EE, 0x78A5636F, 0x84C87814, 0x8CC70208,
            0x90BEFFFA, 0xA4506CEB, 0xBEF9A3F7, 0xC67178F2);

        var HASH = new Array(0x6A09E667, 0xBB67AE85, 0x3C6EF372, 0xA54FF53A, 0x510E527F,
            0x9B05688C, 0x1F83D9AB, 0x5BE0CD19);

        var W = new Array(64);
        var a, b, c, d, e, f, g, h, i, j;
        var T1, T2;

        m[l >> 5] |= 0x80 << (24 - l % 32);
        m[((l + 64 >> 9) << 4) + 15] = l;

        for ( var i = 0; i<m.length; i+=16 ) {
            a = HASH[0];
            b = HASH[1];
            c = HASH[2];
            d = HASH[3];
            e = HASH[4];
            f = HASH[5];
            g = HASH[6];
            h = HASH[7];

            for ( var j = 0; j<64; j++) {
                if (j < 16) W[j] = m[j + i];
                else W[j] = safe_add(safe_add(safe_add(Gamma1256(W[j - 2]), W[j - 7]), Gamma0256(W[j - 15])), W[j - 16]);

                T1 = safe_add(safe_add(safe_add(safe_add(h, Sigma1256(e)), Ch(e, f, g)), K[j]), W[j]);
                T2 = safe_add(Sigma0256(a), Maj(a, b, c));

                h = g;
                g = f;
                f = e;
                e = safe_add(d, T1);
                d = c;
                c = b;
                b = a;
                a = safe_add(T1, T2);
            }

            HASH[0] = safe_add(a, HASH[0]);
            HASH[1] = safe_add(b, HASH[1]);
            HASH[2] = safe_add(c, HASH[2]);
            HASH[3] = safe_add(d, HASH[3]);
            HASH[4] = safe_add(e, HASH[4]);
            HASH[5] = safe_add(f, HASH[5]);
            HASH[6] = safe_add(g, HASH[6]);
            HASH[7] = safe_add(h, HASH[7]);
        }
        return HASH;
    }

    function str2binb (str) {
        var bin = Array();
        var mask = (1 << chrsz) - 1;
        for(var i = 0; i < str.length * chrsz; i += chrsz) {
            bin[i>>5] |= (str.charCodeAt(i / chrsz) & mask) << (24 - i%32);
        }
        return bin;
    }

    function Utf8Encode(string) {
        string = string.replace(/\r\n/g,"\n");
        var utftext = "";

        for (var n = 0; n < string.length; n++) {

            var c = string.charCodeAt(n);

            if (c < 128) {
                utftext += String.fromCharCode(c);
            }
            else if((c > 127) && (c < 2048)) {
                utftext += String.fromCharCode((c >> 6) | 192);
                utftext += String.fromCharCode((c & 63) | 128);
            }
            else {
                utftext += String.fromCharCode((c >> 12) | 224);
                utftext += String.fromCharCode(((c >> 6) & 63) | 128);
                utftext += String.fromCharCode((c & 63) | 128);
            }

        }

        return utftext;
    }

    function binb2hex (binarray) {
        var hex_tab = hexcase ? "0123456789ABCDEF" : "0123456789abcdef";
        var str = "";
        for(var i = 0; i < binarray.length * 4; i++) {
            str += hex_tab.charAt((binarray[i>>2] >> ((3 - i%4)*8+4)) & 0xF) +
                hex_tab.charAt((binarray[i>>2] >> ((3 - i%4)*8  )) & 0xF);
        }
        return str;
    }

    s = Utf8Encode(s);
    return binb2hex(core_sha256(str2binb(s), s.length * chrsz));

}



function  smsReq() {
    //XMLHttpRequest 객체 생성
    let lhpNo = document.getElementById("lhpNo").value;
    let smsErr =document.getElementById("smsErr")
    smsErr.style.display="none";
    document.getElementById("smsCheckYn").value="N";

    const dss ="t";

    if (lhpNo == ""){
        smsErr.style.display="block";
        smsErr.innerText="* 휴대폰 번호를 입력해주세요.";
        return;
    }

    let sendData = {
        telNum: lhpNo,
        stype: 'r'
    };
    var opts = {method: 'POST', body: JSON.stringify(sendData), headers: {"Content-Type": "application/json"}};
    fetch('/api/cash/smsRequest', opts).then(function(response) {
        return response.json();
    }).then(function(res) {
        if (res.resultCode == "00"){
            document.getElementById("smsLi").style.display="block";
            document.getElementById("lhpNo").readOnly=true;
            document.getElementById("lhpNoBtn").setAttribute("disabled","disabled");
        }else if (res.resultCode == "99"){

        }else{
            alert("통신 오류");
            return;
        }
    });

};


function  smsCheck() {

    let lhpNo = document.getElementById("lhpNo").value;
    let smsNum = document.getElementById("smsConfirm").value;
    let smsErr =document.getElementById("smsErr");
    smsErr.style.display="none";

    if (smsNum == ""){
        smsErr.style.display="block";
        smsErr.innerText="* 인증번호를 입력해주세요.";
        return false;
    }

    if (lhpNo == ""){
        smsErr.style.display="block";
        smsErr.innerText="* 휴대폰 번호를 입력해주세요.";
        return;
    }

    //XMLHttpRequest 객체 생성
    let xhr = new XMLHttpRequest();
    //요청을 보낼 방식, 주소, 비동기여부 설정
    xhr.open('POST', '/api/cash/smsCheck', true);
    //HTTP 요청 헤더 설정
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.responseType='json';
    //요청 전송
    xhr.send("telNum="+lhpNo+"&smsNum="+smsNum);
    //통신후 작업
    xhr.onload = function () {
        //통신 성공
        if (xhr.status == 200) {

            //	console.log(xhr.response);
            if (xhr.response.resultCode == "00"){
                document.getElementById("smsConfirm").setAttribute("disabled","disabled");
                document.getElementById("smsConfirm").style="background:#d8d4d4;"
                document.getElementById("smsConfirmBtn").setAttribute("disabled","disabled");
                document.getElementById("smsConfirmBtn").value="인증성공";
                document.getElementById("smsCheckYn").value="Y";

                return;
            }else if(xhr.response.resultCode == "99"){
                smsErr.style.display="block";
                smsErr.innerText="*" + xhr.response.resultMsg;
                return;
            }
        } else {
            alert("통신 오류");
            return;
        }
    }
};


function  joinOk() {


    let lid = document.getElementById("lid").value;
    let lname = document.getElementById("lname").value;
    let lbirthDay = document.getElementById("lbirthDay").value;
    let lhpNo = document.getElementById("lhpNo").value;
    let lpw = document.getElementById("lpw").value;
    let smsCheckYn = document.getElementById("smsCheckYn").value;
    let emailCheckYn = document.getElementById("emailCheckYn").value;
    let emailDupCheckYn = document.getElementById("emailDupCheckYn").value;



    document.getElementById("emailErr").innerText="";
    document.getElementById("nameErr").innerText="";
    document.getElementById("smsErr").innerText="";
    document.getElementById("pwdErr").innerText="";
    document.getElementById("birthDayErr").innerText="";



    if (lid == ""){
        document.getElementById("emailErr").style.display="block";
        document.getElementById("emailErr").innerText="* 이메일을 입력해주세요.";
        document.getElementById("lid").focus();
        return;
    }



    if (emailDupCheckYn == "N"){
        document.getElementById("emailErr").style.display="block";
        document.getElementById("emailErr").innerText="* 올바른 이메일 형식이 아닙니다.";
        document.getElementById("lid").focus();
        return;
    }


    if (emailCheckYn == "N"){
        document.getElementById("emailErr").style.display="block";
        document.getElementById("emailErr").innerText="* 올바른 이메일 형식이 아닙니다.";
        document.getElementById("lid").focus();
        return;
    }

    if (lname == ""){
        document.getElementById("nameErr").style.display="block";
        document.getElementById("nameErr").innerText="* 이름을 입력하세요.";
        document.getElementById("lname").focus();
        return;
    }

    if (lbirthDay == ""){
        document.getElementById("birthDayErr").style.display="block";
        document.getElementById("birthDayErr").innerText="* 생년월일을 입력하세요.";
        document.getElementById("lbirthDay").focus();
        return;
    }

    if (lhpNo == ""){
        document.getElementById("smsErr").style.display="block";
        document.getElementById("smsErr").innerText="* 전화번호를 입력해주세요.";
        document.getElementById("lhpNo").focus();
        return;
    }

    if (smsCheckYn == "N"){
        document.getElementById("smsErr").style.display="block";
        document.getElementById("smsErr").innerText="* 전화번호 인증을 완료해주세요.";
        document.getElementById("lhpNo").focus();
        return;
    }

    if (lpw == ""){
        document.getElementById("pwdErr").style.display="block";
        document.getElementById("pwdErr").innerText="*비밀번호를 입력해주세요.";
        document.getElementById("lpw").focus();
        return;
    }

    if(!/^[a-zA-Z0-9]{8,20}$/.test(lpw)){
        document.getElementById("pwdErr").style.display="block";
        document.getElementById("pwdErr").innerText="*숫자와 영문자 조합으로 8~20자리를 사용해야 합니다.";
        document.getElementById("lpw").focus();
        return;
    }

    let sendData = {
        email: lid,
        userNm: lname,
        userBirth: lbirthDay,
        userTel: lhpNo,
        loginPw: SHA256(lpw),
        termsOfService: "Y",
        termsOfPersonal: "Y",
        termsOfPayment: "Y",
        termsOfBenefit: "Y"
    };
    var opts = {method: 'POST', body: JSON.stringify(sendData), headers: {"Content-Type": "application/json"}};
    fetch('/api/cash/join', opts).then(function(response) {
        return response.json();
    }).then(function(res) {

       // console.log(res)
        if (res.resultCode == "00"){
            sessionStorage.setItem('userId', res.resultData.userId);
            location.href="/parthner/p/joinStep1";
        }else if (res.resultCode == "99"){
            alert(res.resultMsg);
            return;
        }else{
            alert("통신 오류");
            return;
        }
    });
}

function chkEmail(str) {
    var regExp = /^[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*@[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*.[a-zA-Z]{2,3}$/i;
    if (regExp.test(str)) return true;
    else return false;
}


function emailChk(obj){
    str = obj.value;
    strLenght = obj.value.length
    let chkYn=chkEmail(str)
    if (chkYn == false && strLenght > 3){
        document.getElementById("emailErr").style.display="block";
        document.getElementById("emailErr").innerText="* 올바른 이메일 형식이 아닙니다.";
        document.getElementById("emailCheckYn").value="N";
    }else{
        document.getElementById("emailErr").innerText="";
        document.getElementById("emailCheckYn").value="Y";
    }
}

function pwdChk(obj){
    str = obj.value;
    strLenght = obj.value.length
    let chkYn=/^[a-zA-Z0-9]{8,20}$/.test(str)
    if (chkYn == false){
        document.getElementById("pwdErr").style.display="block";
        document.getElementById("pwdErr").innerText="*숫자와 영문자 조합으로 8~15자리를 사용해야 합니다.";
    }else{
        document.getElementById("pwdErr").innerText="";
    }
}


function  couponOk() {

    //console.log(hCode)
    let couponNo = document.getElementById("couponNo").value;


    document.getElementById("couponNoErr").innerText="";

    if (couponNo == ""){
        document.getElementById("couponNoErr").style.display="block";
        document.getElementById("couponNoErr").innerText="* 쿠폰번호를 입력해주세요.";
        document.getElementById("couponNo").focus();
        return;
    }
    let sendData = {
        userId: sessionStorage.getItem('userId'),
        storeId: sessionStorage.getItem('storeId'),
        couponNo : couponNo

    };
    var opts = {method: 'POST', body: JSON.stringify(sendData), headers: {"Content-Type": "application/json"}};
    fetch('/api/etc/coupon', opts).then(function(response) {
        return response.json();
    }).then(function(res) {
        //console.log(res)
        if (res.resultCode == "00"){
            alert("쿠폰을 등록하였습니다.")
            location.reload();
        }else if (res.resultCode == "99"){
            document.getElementById("couponNoErr").style.display="block";
            document.getElementById("couponNoErr").innerText="*"+res.resultMsg;
            document.getElementById("couponNo").focus();
            return;
        }else{
            alert("통신 오류");
            return;
        }
    });
}

