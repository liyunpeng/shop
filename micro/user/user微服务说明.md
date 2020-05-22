生成go ， micro格式文件的正确姿势：
```
user$ protoc -I ./proto/user/ --micro_out=./proto/user --go_out=./proto/user ./proto/user/user.proto
```

启动user微服务
```
user$ go run main.go plugin.go database.go 
2020-05-22 21:36:02.821240 I | Transport [http] Listening on [::]:63086
2020-05-22 21:36:02.821336 I | Broker [http] Connected to [::]:63087
2020-05-22 21:36:02.821534 I | Registry [mdns] Registering node: go.micro.srv.user-5ffcb092-3bae-48bc-b80d-43b78f48f6de

```

启动微服务网关api
```
encode$ micro api --namespace=go.micro.srv
2020-05-22 21:38:07.299347 I | [api] Registering API Default Handler at /
2020-05-22 21:38:07.300364 I | [api] HTTP API Listening on [::]:8080
2020-05-22 21:38:07.301007 I | [api] Transport [http] Listening on [::]:63795
2020-05-22 21:38:07.302478 I | [api] Broker [http] Connected to [::]:63796
2020-05-22 21:38:07.302703 I | [api] Registry [mdns] Registering node: go.micro.api-77c76e95-407c-4222-9fd1-e28416b94f7e
::1 - - [22/May/2020:21:39:58 +0800] "POST /user/UserService/Register HTTP/1.1" 200 35 "" "PostmanRuntime/7.24.1"
```

postman 测试：
```gotemplate
post http://localhost:8080/user/UserService/Register
{"user":{"id":1,"name":"test","phone":"123","password":"123"}}
```
响应
```gotemplate
{"code":"200","msg":"注册成功"}
```

