package service

import (
	"context"
	protobuf "shop/encode/generate"
	"shop/logger"
)

type User struct{}

func (u *User) Hello(ctx context.Context, req *protobuf.Request, res *protobuf.Response) error {
	res.Msg = "Hello " + req.Name
	logger.Info.Println("微服务的服务提供者的服务方法被调用， ", res.Msg)
	return nil
}
