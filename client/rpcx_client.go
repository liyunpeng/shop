package client

import (
	"bytes"
	"github.com/rpcx-ecosystem/rpcx-gateway"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/codec"
	"io/ioutil"
	"net/http"
	"shop/logger"
	"shop/rpc"

	"context"
)

//var (
//	addr1       = flag.String("addr", "127.0.0.1:8972", "server address")
//)

func StartRpcClient() {
	Peer2Peer()

	CallwithGateway()
}
func Peer2Peer() {
	//flag.Parse()
	addr := "127.0.0.1:8972"
	d := client.NewPeer2PeerDiscovery("tcp@" + addr, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &rpc.Args{
		A: 10,
		B: 20,
	}

	reply := &rpc.Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		logger.Error.Println ("failed to call: %v", err)
	}

	logger.Debug.Printf("rpcx %d * %d = %d", args.A, args.B, reply.C)
}


func CallwithGateway(){
	cc := &codec.MsgpackCodec{}

	args := &rpc.Args{
		A: 100,
		B: 200,
	}

	data, _ := cc.Encode(args)
	// request
	req,err := http.NewRequest("POST","http://127.0.0.1:8972/", bytes.NewReader(data))
	if err != nil{
		logger.Error.Println("failed to create request: ", err)
		return
	}

	// 设置header
	h := req.Header
	h.Set(gateway.XMessageID,"10000")
	h.Set(gateway.XMessageType,"0")
	h.Set(gateway.XSerializeType,"3")
	h.Set(gateway.XServicePath,"Arith")
	h.Set(gateway.XServiceMethod,"Mul")

	// 发送http请求
	//  http请求===>rpcx请求===>rpcx服务===>返回rpcx结果===>转换为http的response===>输出到client
	res, err := http.DefaultClient.Do(req)
	if err != nil{
		logger.Error.Println ("failed to call: ", err)
	}
	defer res.Body.Close()
	// 获取结果
	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil{
		logger.Error.Println ("failed to read response: ", err)
	}
	// 解码
	reply := &rpc.Reply{}
	err = cc.Decode(replyData, reply)
	if err != nil{
		logger.Error.Println ("failed to decode reply: ", err)
	}
	logger.Info.Printf ("rpcx gateway call %d * %d = %d", args.A, args.B, reply.C)
}
