package rpc

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	pb "shop/rpc/proto"
)

// 定义服务端实现约定的接口
type UserInfoService struct {
	file *os.File
}

var u = UserInfoService{}

func (u *UserInfoService)bufioWriteFile( s string){
	var err error
	content := []byte(s)
	newWriter := bufio.NewWriterSize(u.file, 1024)
	if _, err = newWriter.Write(content); err != nil {
		fmt.Println(err)
	}
	if err = newWriter.Flush(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("write file successful")
}
// 实现服务端需要首先的接口
func (s *UserInfoService) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	name := req.Name
	// 在数据库查用户信息
	if name == "zhangsan" {
		resp = &pb.UserResponse{
			Id:   1,
			Name: name,
			Age:  22,
			Hobby: []string{"Sing", "run", "basketball"},
		}
	}else{
		resp = &pb.UserResponse{
			Id:   1,
			Name: name,
			Age:  22,
			Hobby: []string{"1", "2", "3"},
		}
	}

	u.bufioWriteFile(name)
	err = nil
	return
}
func (u *UserInfoService) init(){
	var err error
	u.file, err = os.OpenFile("./a.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
	   fmt.Println(err)
	}
}
func (u *UserInfoService) destroy(){
	u.file.Close()
}

func Server() {
	// 1. 监听
	addr := "127.0.0.1:8989"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("监听异常：%s\n", err)
	}
	fmt.Printf("开始监听：%s\n", addr)
	// 2.实例化gRPC
	s := grpc.NewServer()

	u.init()
	defer u.destroy()
	// 3.在gRPC上注册微服务
	// 第二个参数类型需要接口类型的变量
	pb.RegisterUserInfoServiceServer(s, &u)
	// 4.启动gRPC服务
	s.Serve(lis)
}
