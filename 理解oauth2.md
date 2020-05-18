### 一. 认证服务器的log
$ go run server.go
2020/03/24 15:38:51 Server is running at 9096 port.

##### 1. 收到重定向的log：
```
authorise, r= &{
    GET /authorize?client_id=222222&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz
    HTTP/1.1 1 1
    map[
        Accept:[
            text/html,application/xhtml+xml,application/xml;
            q=0.9,image/webp,image/apng,*/*;
            q=0.8,application/signed-exchange;
            v=b3;
            q=0.9
        ]
        Accept-Encoding:[gzip, deflate, br]
        Accept-Language:[zh-CN,zh;q=0.9,en;q=0.8]
        Connection:[keep-alive]
        Cookie:[
            Goland-45831321=2674bde7-8028-4a67-ad16-7d809ec996b3;
            go_session_id=ZjM5ZGFiMGYtMjMxYy00ZTBhLWE3YjMtY2IxMWMxMzYyYzQ0.87e43b03d9eb17b368c02769993866fd4d7bbe02;
            go_session_id=ZTBlNjg3MzktMjAxNS00NTliLTlhNzAtNGJlYjIzMTEzMmVi.1e38b0664d12cee822cae368b934d9f89c05a79c;
            go_session_id=NGJmMWJmY2ItZDU2NC00ZDBkLWE3YzQtYmMyNzE2ZTVjY2Zk.1fa2fd459e3fa2176eaf428fc1c29eb2fb8788bc
        ]
        Sec-Fetch-Mode:[navigate]
        Sec-Fetch-Site:[none]
        Sec-Fetch-User:[?1]
        Upgrade-Insecure-Requests:[1]
        User-Agent:[
            Mozilla/5.0 (Windows NT 6.1; Win64; x64)
            AppleWebKit/537.36 (KHTML, like Gecko)
            Chrome/79.0.3945.130 Safari/537.36
        ]
    ]
    {} <nil> 0 [] false localhost:9096
    map[
        lient_id:[222222]
        redirect_uri:[http://localhost:9094/oauth2]
        response_type:[code]
        scope:[all]
        state:[xyz]
    ]
    map[] <nil>
    map[] [::1]:52928
    /authorize?client_id=222222&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz
    <nil> <nil> <nil> 0xc0000e4180
}
```

#### 2. 点击login的log：
```
login handler session=  &{{{0 0} 0 0 0 0} 0xc0000da4b0 0xc0000f1290 4bf1bfcb-d564-4d0d-a7c4-bc2716e5ccfd 7200 map[LoggedInUserID:000000 ReturnUri:map[client_id:[222222] redirect_uri:[http://localhost:9094/oauth2] response_type:[code] scope:[all] state:[xyz]]]}
login handler requeset=  &{POST /login HTTP/1.1 1 1 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9] Accept-Encoding:[gzip, deflate, br] Accept-Language:[zh-CN,zh;q=0.9,en;q=0.8] Cache-Control:[max-age=0] Connection:[keep-alive] Content-Length:[44] Content-Type:[application/x-www-form-urlencoded] Cookie:[Goland-45831321=2674bde7-8028-4a67-ad16-7d809ec996b3; go_session_id=NGJmMWJmY2ItZDU2NC00ZDBkLWE3YzQtYmMyNzE2ZTVjY2Zk.1fa2fd459e3fa2176eaf428fc1c29eb2fb8788bc] Origin:[http://localhost:9096] Referer:[http://localhost:9096/login] Sec-Fetch-Mode:[navigate] Sec-Fetch-Site:[same-origin] Sec-Fetch-User:[?1] Upgrade-Insecure-Requests:[1] User-Agent:[Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36]] 0xc0000ee840 <nil> 44 [] false localhost:9096 map[] map[] <nil> map[] [::1]:52950 /login <nil> <nil> <nil> 0xc0000ee880}

/auth authHandler, sesstion= &{{{0 0} 0 0 0 0} 0xc0000f6990 0xc000286390 a29031c8-835d-4236-b3c1-b2c95bca403d 7200
map[LoggedInUserID:000000
ReturnUri:map[client_id:[222222]
redirect_uri:[http://localhost:9094/oauth2]
response_type:[code] scope:[all] state:[xyz]]]}
```

#### 3. 点击allow后的log
```
authorise, r= &{
    POST /authorize HTTP/1.1 1 1
    map[
        Accept:[
            text/html,application/xhtml+xml,
            application/xml;
            q=0.9,image/webp,image/apng,*/*;
            q=0.8,application/signed-exchange;
            v=b3;q=0.9
        ]
        Accept-Encoding:[gzip, deflate, br]
        Accept-Language:[zh-CN,zh;q=0.9,en;q=0.8]
        Cache-Control:[max-age=0]
        Connection:[keep-alive]
        Content-Length:[0]
        Content-Type:[application/x-www-form-urlencoded]
        Cookie:[
            Goland-45831321=2674bde7-8028-4a67-ad16-7d809ec996b3;
            go_session_id=NGJmMWJmY2ItZDU2NC00ZDBkLWE3YzQtYmMyNzE2ZTVjY2Zk.1fa2fd459e3fa2176eaf428fc1c29eb2fb8788bc
        ]
        Origin:[http://localhost:9096]
        Referer:[http://localhost:9096/auth]
        Sec-Fetch-Mode:[navigate]
        Sec-Fetch-Site:[same-origin]
        Sec-Fetch-User:[?1]
        Upgrade-Insecure-Requests:[1]
        User-Agent:[
            Mozilla/5.0 (Windows NT 6.1; Win64; x64)
            AppleWebKit/537.36 (KHTML, like Gecko)
            Chrome/79.0.3945.130 Safari/537.36
        ]
    ]

    {}
    <nil> 0 [] false
    localhost:9096
    map[
        client_id:[222222]
        redirect_uri:[http://localhost:9094/oauth2]
        response_type:[code]
        scope:[all]
        state:[xyz]
    ]
    map[] <nil>
    map[]
    [::1]:52964 /authorize
    <nil> <nil> <nil> 0xc0000e5340
}


token 请求处理，
w= &{
    0xc0000f2140 0xc00010eb00 0xc0000e5580 0x5598f0 true false false false
    0xc0000e4200 {
        0xc0000d8540 map[Cache-Control:[no-store]
        Content-Type:[application/json;charset=UTF-8]
        Pragma:[no-cache]] false false
    }
    map[
        Cache-Control:[no-store]
        Content-Type:[application/json;charset=UTF-8]
        Pragma:[no-cache]
    ]
    true 303 -1 200 false false [] 0
    [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
    [0 0 0 0 0 0 0 0 0 0] [0 0 0] 0xc000227260 0
}

r= &{
    POST /token HTTP/1.1 1 1
    map[
        Accept-Encoding:[gzip]
        Authorization:[Basic MjIyMjIyOjIyMjIyMjIy]
        Content-Length:[109]
        Content-Type:[application/x-www-form-urlencoded]
        User-Agent:[Go-http-client/1.1]
    ]
    0xc0000e5580
    <nil> 109 [] false localhost:9096
    map[
        code:[OBEAFRNAPVIPVQHRMYCEWG]   生成的授权码
        grant_type:[authorization_code]
        redirect_uri:[http://localhost:9094/oauth2]
    ]
    map[
        code:[OBEAFRNAPVIPVQHRMYCEWG]
        grant_type:[authorization_code]
        redirect_uri:[http://localhost:9094/oauth2]
    ]
    <nil>
    map[] [::1]:52967 /token <nil> <nil> <nil> 0xc0000e55c0
}

authorise, r= &{
    GET /authorize?client_id=222222&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz
    HTTP/1.1 1 1
    map[
        Accept:[image/webp,image/apng,image/*,*/*;q=0.8]
        Accept-Encoding:[gzip, deflate, br]
        Accept-Language:[zh-CN,zh;q=0.9,en;q=0.8]
        Connection:[keep-alive]
        Cookie:[
            Goland-45831321=2674bde7-8028-4a67-ad16-7d809ec996b3;
            go_session_id=NGJmMWJmY2ItZDU2NC00ZDBkLWE3YzQtYmMyNzE2ZTVjY2Zk.1fa2fd459e3fa2176eaf428fc1c29eb2fb8788bc
        ]
        Referer:[http://localhost:9094/oauth2?code=OBEAFRNAPVIPVQHRMYCEWG&state=xyz]
        Sec-Fetch-Mode:[no-cors] Sec-Fetch-Site:[same-site]
        User-Agent:[
            Mozilla/5.0 (Windows NT 6.1; Win64; x64)
            AppleWebKit/537.36 (KHTML, like Gecko)
            Chrome/79.0.3945.130 Safari/537.36
        ]
    ]
    {}
    <nil> 0 [] false
    localhost:9096
    map[
        client_id:[222222]
        redirect_uri:[http://localhost:9094/oauth2]
        response_type:[code]
        scope:[all]
        state:[xyz]
    ]
    map[]
    <nil> map[]
    [::1]:52964
    /authorize?client_id=222222&redirect_uri=http%3A%2F%2Flocalhost%3A9094%2Foauth2&response_type=code&scope=all&state=xyz
    <nil> <nil> <nil> 0xc0000e5800
}
```

### 二. 客户端的log：
在认证服务器，用户点allow之后
```
$ go run client.go
go: finding golang.org/x/oauth2 latest
2020/03/24 15:38:47 Client is running at 9094 port.
获取授权码 code= OBEAFRNAPVIPVQHRMYCEWG
获取 token = &{eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODUwNDI4OTYsInN1YiI6IjAwMDAwMCJ9.wm0Q8hMxCpycJAagou_zywyFpXUVdLFOtxUmhoaLnrDTdIuFNugebozPzMO7tVNlyoGp-kFZbkbUZSZzBxD4nw Bearer NOTSEZYNUU27WJNBJJ8WHA 2020-03-24 17:41:36.6875 +0800 CST m=+7368.752000001 map[access_token:eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODUwNDI4OTYsInN1YiI6IjAwMDAwMCJ9.wm0Q8hMxCpycJAagou_zywyFpXUVdLFOtxUmhoaLnrDTdIuFNugebozPzMO7tVNlyoGp-kFZbkbUZSZzBxD4nw expires_in:7200 refresh_token:NOTSEZYNUU27WJNBJJ8WHA scope:all token_type:Bearer]}
```
