package client

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	microClient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
	protobuf "shop/encode/generate" // 应用的目录可以和包名不同
	"shop/logger"
)

var microClientObj microClient.Client
var microServiceName string

var microClientService protobuf.UserService
func init() {
	service := micro.NewService()
	service.Init()
	microClientObj = service.Client()
	microServiceName = config.Get("srv").String("micro.hrefs.srv")


	microClientService = protobuf.NewUserService(
		microServiceName,
		microClientObj)

}

func MicroCall(method string, req interface{}, rsp interface{}) error {
	request := microClientObj.NewRequest(microServiceName, fmt.Sprintf("Hrefs.%s", method), req, microClient.WithContentType("application/json"))

	if err := microClientObj.Call(context.TODO(), request, &rsp); err != nil {
		logger.Info.Println(err)
		return err
	}

	return nil
}

func MicroCallUser(){
	res, err := microClientService.Hello(context.TODO(), &protobuf.Request{Name: "World ^_^"})
	if err != nil {
		logger.Info.Println(err)
	}
	logger.Info.Println("微服务的客户端收到服务端的响应消息=", res.Msg)
}
