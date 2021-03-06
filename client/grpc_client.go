package client

import (
	"context"
	"google.golang.org/grpc"
	"shop/logger"
	pb "shop/rpc/proto"

	"sync/atomic"
	"time"
)

var grpcClient pb.UserInfoServiceClient

var flag int32

func StartGrpcClient() {
	// 1. 创建与gRPC服务端的连接
	conn, err := grpc.Dial("127.0.0.1:8989", grpc.WithInsecure())
	if err != nil {
		logger.Info.Printf(" 连接异常：%GrpcSever\n", err)
	}
	//defer conn.Close()
	// 2. 实例化gRPC客户端
	grpcClient = pb.NewUserInfoServiceClient(conn)
}

func GrpcCall(msg string) {
	for {
		if atomic.CompareAndSwapInt32(&flag, 0, 1) {
			// 3. 组装参数
			req := new(pb.UserRequest)
			req.Name = msg
			// 4. 调用接口
			resp, err := grpcClient.GetUserInfo(context.Background(), req)
			if err != nil {
				logger.Info.Printf("响应异常：%GrpcSever\n", err)
			}
			logger.Info.Printf("响应结果: %v\n", resp)
			atomic.StoreInt32(&flag, 0)
		}

		logger.Info.Println("cas: GrpcCall 被他人调用，GrpcCall 还在锁定中 ")

		time.Sleep(time.Second)
	}

}
