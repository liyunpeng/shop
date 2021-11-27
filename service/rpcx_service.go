package service

import (
	"flag"
	"github.com/smallnest/rpcx/server"
	"shop/logger"
)

var (
	addr1 = flag.String("addr", "localhost:8972", "server address")
)


var rpcxServer *server.Server
func StartRpcxService() {
	flag.Parse()

	addr := "localhost:8972"
	if false {
		rpcxServer = &server.Server{}
		rpcxServer.RegisterName("Arith", new(Arith), "")
	}else{
		rpcxServer = server.NewServer()
		rpcxServer.Register(new(Arith),"")
	}
	logger.Info.Println("Rpcx 微服务 启动 ")
	rpcxServer.Serve("tcp",addr)

	logger.Info.Println("Rpcx 微服务运行结束 ")

	//select {}
}


func StopRpcxService() {

	rpcxServer.Close()
	logger.Info.Println("控制Rpcx微服务运行结束 ")
}
