## API网关提供单一的入口把的所有微服务用API统一起来
微服务架构是把应用解耦成逻辑上的相对隔离的服务，  
API网关则是提供单一的入口把服务的API统一起来。  
通过服务发现，Micro API以http方式，将请求动态路由到具体的后台服务接口。

Micro API是基于go-micro开发，所以它天然具备服务发现、负载均衡、编码及RPC通信的能力。  
因此，Micro API也是go-micro体系中的一个微服务，它自身也是可插拔的。

## micro 命令的安装
通过命令进行安装：
```
$ go get -u github/micro/micro
``` 
在GOPATH/bin 会有micro 可执行文件  
命令行输入micro --version就可以看到它的版本信息。

Micro API使用命名空间来在逻辑上区分后台服务及公开的服务.  
命名空间及HTTP请求路径会用于解析服务名与方法.  
Micro API默认的命名空间是go.micro.api.   
比如
```
GET /foo HTTP/1.1
```
会被路由到这个服务：
```
go.micro.api.foo
```

## 