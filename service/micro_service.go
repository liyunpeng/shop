package service

import (
	"context"
	"github.com/kataras/golog"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/registry/consul/v2"
	"shop/util"
	"time"
)

var service micro.Service
func StartMicroService(){
	urls := util.GetConsulUrls()
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = urls
	})

	// TODO: 检测consul服务发现是否正常启动

	service = micro.NewService(
		micro.Registry(reg),
		micro.Name(config.Get("srv").String("micro.hrefs.srv")),
		micro.WrapHandler(logWrapper),
	)
	server.Init()
	service.Server().Init(server.Wait(nil))
	micro.RegisterHandler(service.Server(), new(Hrefs))
	service.Run()
	service.Server().Stop()
}

func Stop(){
	service.Server().Stop()
}
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	start := time.Now()
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := fn(ctx, req, rsp)
		//utils.WriteErrorLog(req.Endpoint(), err)

		golog.Infof("%s %s", time.Since(start), req.Endpoint())
		return err
	}
}
