package client

import (
	"github.com/smallnest/rpcx/client"
	"shop/logger"
	"shop/service"

	"context"
)

//var (
//	addr1       = flag.String("addr", "127.0.0.1:8972", "server address")
//)

func StartRpcClient() {
	Peer2Peer()
}
func Peer2Peer() {
	//flag.Parse()
	addr := "127.0.0.1:8972"
	d := client.NewPeer2PeerDiscovery("tcp@" + addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &service.Args{
		A: 10,
		B: 20,
	}

	reply := &service.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		logger.Error.Println ("failed to call: %v", err)
	}

	logger.Debug.Printf("rpcx %d * %d = %d", args.A, args.B, reply.C)
}
