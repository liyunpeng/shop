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

构建镜像
```
user$ docker build . -t user-srv:latest
Sending build context to Docker daemon     64MB
Step 1/3 : FROM alpine
latest: Pulling from library/alpine
cbdbe7a5bc2a: Pull complete 
Digest: sha256:9a839e63dad54c3a6d1834e29692c8492d93f90c59c978c1ed79109ea4fb9a54
Status: Downloaded newer image for alpine:latest
 ---> f70734b6a266
Step 2/3 : ADD user-srv /user-srv
 ---> c5a155480f63
Step 3/3 : ENTRYPOINT [ "/user-srv" ]
 ---> Running in 75a87650e2ed
Removing intermediate container 75a87650e2ed
 ---> 8b0bd45d8e00
Successfully built 8b0bd45d8e00
Successfully tagged user-srv:latest
```

查看镜像
```
hop$ docker images 
 REPOSITORY                                           TAG                 IMAGE ID            CREATED             SIZE
 user-srv                                             latest              8b0bd45d8e00        46 minutes ago      47.4MB
```


##### exec user process caused "exec format error"
容器启动不了的问题：
user$ sudo docker run -p 8080:8080 -it --name docker1235 user-srv:latest 
standard_init_linux.go:211: exec user process caused "exec format error"


