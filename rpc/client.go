package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "shop/rpc/proto"
)

func Client( msg string) {
	// 1. 创建与gRPC服务端的连接
	conn, err := grpc.Dial("127.0.0.1:8989", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("连接异常：%s\n", err)
	}
	defer conn.Close()
	// 2. 实例化gRPC客户端
	client := pb.NewUserInfoServiceClient(conn)
	// 3. 组装参数
	req := new(pb.UserRequest)
	req.Name = msg
	// 4. 调用接口
	resp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		fmt.Printf("响应异常：%s\n", err)
	}
	fmt.Printf("响应结果: %v\n", resp)
}
