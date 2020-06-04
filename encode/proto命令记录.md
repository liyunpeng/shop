### 安装protoc命令
先要下载protoc命令，
下载后有： 
D:\programe\protoc-3.9.0-win64\bin\protoc.exe
将其bin目录路径添加到path环境变量中。

### 用protoc生成go格式文件和micro格式文件 
在本文件所在目录下运行：
```
$ protoc -I ./source/ --go_out=./generate --micro_out=./generate ./source/user.proto
```
参数说明
* -I 指定了proto原始定义文件所在目录
* --go_out=:./generate 
表示用protoc-gen-go生成go格式的文件， 放在./generate目录下
* --micro_out=./generate 
表示用micro-gen-go生成micro格式的文件， 放在./generate目录下
* ./source/user.proto  表示原生proto定义文件
上面命令在generate目录下生成文件如下：
1.序列化文件：*pb.go
2.go-micro微服务的客户端和服务端代码文件： *pb.micro.go
生成文件的包名就是proto文件名，这里为文件名为user.proto, 所以包名为user

#### 解决 WARNING: Missing 'go_package' option in "user.proto",
encode$ protoc -I ./source/ --go_out=./generate --micro_out=./generate ./source/user.proto
2020/05/22 20:33:27 WARNING: Missing 'go_package' option in "user.proto", please specify:
        option go_package = ".;user";
解决办法： 根据提示， 在user.proto中添加：
option go_package = ".;user";

#### 解决micro生成文件import错误问题

编译报错：

```gotemplate
# shop/service
service/micro_service.go:59:45: cannot use service.Server() (type "github.com/micro/go-micro/v2/server".Server) as type "github.com/micro/go-micro/server".Server in argument to user.RegisterUserHandler:
	"github.com/micro/go-micro/v2/server".Server does not implement "github.com/micro/go-micro/server".Server (wrong type for Handle method)
		have Handle("github.com/micro/go-micro/v2/server".Handler) error
		want Handle("github.com/micro/go-micro/server".Handler) error
```
原因：
业务程序用的是micro/v2, 而生成代码还用micro

解决办法：
将微服务·micro生成文件user.pb.micro.go: 
import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)
改为
import (
	context "context"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)
