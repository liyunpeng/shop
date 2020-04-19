package client

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	microClient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/config"
)

var microClientObj microClient.Client
var name string

func init() {
	service := micro.NewService()
	service.Init()
	microClientObj = service.Client()
	name = config.Get("srv").String("micro.hrefs.srv")
}

func MicroCall(method string, req interface{}, rsp interface{}) error {
	request := microClientObj.NewRequest(name, fmt.Sprintf("Hrefs.%s", method), req, microClient.WithContentType("application/json"))

	if err := microClientObj.Call(context.TODO(), request, &rsp); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
