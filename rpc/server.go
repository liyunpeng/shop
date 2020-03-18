package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	pb "shop/rpc/proto"
)

// 定义服务端实现约定的接口
type UserInfoService struct {
}

var u = UserInfoService{}

func writefile(s string){
	f, _ := os.Create("a.txt")

	f.Write([]byte(s))

	f.Seek(0, os.SEEK_SET)

	p := make([]byte, 5)

	if _, err := f.Read(p); err != nil {
		log.Fatal("[File]", err)
	}

	f.Close()
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
			//切片字段
			Hobby: []string{"Sing", "run", "basketball"},
		}
	}else{
		resp = &pb.UserResponse{
			Id:   1,
			Name: name,
			Age:  22,
			//切片字段
			Hobby: []string{"1", "2", "3"},
		}
	}

	writefile(name)
	err = nil
	return
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
	// 3.在gRPC上注册微服务
	// 第二个参数类型需要接口类型的变量
	pb.RegisterUserInfoServiceServer(s, &u)
	// 4.启动gRPC服务
	s.Serve(lis)
}
