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