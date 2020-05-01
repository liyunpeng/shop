package service

import (
	"google.golang.org/grpc"
	"net"
	"shop/logger"
	"shop/pprofutil"
	pb "shop/rpc/proto"
	"shop/transformer"
	"shop/util"
	"shop/workerpool"
)

var GrpcSever *grpc.Server

var GrpcWorkerPool *workerpool.WorkerPool

func StartGrpcService(grpcConf transformer.GrpcConf) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()
	pprofutil.StartCpuProf()
	// 1. 监听
	addr := grpcConf.Addr
	listenSocket, err := net.Listen("tcp", grpcConf.Addr)
	if err != nil {
		logger.Info.Printf("监听异常：%GrpcSever\n", err)
	}
	logger.Info.Printf("grpc 服务开始监听的地址和端口：%GrpcSever\n", addr)
	// 2.实例化gRPC
	GrpcSever = grpc.NewServer()

	num := 2
	GrpcWorkerPool = workerpool.NewWorkerPool(num)
	GrpcWorkerPool.Run()

	var u = UserInfoService{}
	u.init()
	defer u.destroy()
	// 3.在gRPC上注册微服务
	// 第二个参数类型需要接口类型的变量
	pb.RegisterUserInfoServiceServer(GrpcSever, &u)
	// 4.启动gRPC服务
	logger.Info.Println("启动gRPC服务")

	GrpcSever.Serve(listenSocket)

}

func StopGrpcService(){
	GrpcSever.Stop()
	pprofutil.StopCpuProf()

	pprofutil.SaveMemProf()
	logger.Info.Println(" grpc 服务结束")
}
