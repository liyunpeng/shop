package service

import (
	"flag"
	"github.com/smallnest/rpcx/server"
	"shop/logger"
)

var (
	addr1 = flag.String("addr", "localhost:8972", "server address")
)

func StartRpcxService() {
	flag.Parse()

	addr := "localhost:8972"
	var s *server.Server
	if false {
		s = &server.Server{}
		s.RegisterName("Arith", new(Arith), "")
	}else{
		s = server.NewServer()
		s.Register(new(Arith),"")

	}
	logger.Info.Println("Rpcx 微服务 启动 ")
	s.Serve("tcp",addr)

	logger.Info.Println("Rpcx 微服务运行结束 ")

	//s.Close()
	//select {}
}


