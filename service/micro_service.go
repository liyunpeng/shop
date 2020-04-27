package service

import (
	"context"
	"github.com/kataras/golog"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/registry"
	etcdv3 "github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/registry/consul/v2"
	"shop/logger"

	//"github.com/micro/go-plugins/registry/etcdv3"
	custconfig "shop/config"
	protobuf "shop/encode/generate"
	"shop/util"
	"time"
)

var service micro.Service
var reg registry.Registry

func InitMicro(){
	// 外部参数
	useEtcd := true

	if useEtcd == true {
		reg = etcdv3.NewRegistry( func(options *registry.Options) {
			// etcd 地址
			//options.Addrs = []string{"127.0.0.1:2379"}
			options.Addrs = []string{custconfig.TransformConfiguration.Etcd.Addr}
			// etcd 用户名密码,如果设置的话
			//etcdv3.Auth("root","password")(options)
		})
	}else{
		urls := util.GetConsulUrls()
		// TODO: 检测consul服务发现是否正常启动
		reg = consul.NewRegistry(func(options *registry.Options) {
			options.Addrs = urls
		})
	}
}
func StartMicroService(){
	service = micro.NewService(
		micro.Registry(reg),
		micro.Name(config.Get("srv").String("micro.hrefs.srv")),
		micro.WrapHandler(logWrapper),
	)

	server.Init()

	service.Server().Init(server.Wait(nil))

	micro.RegisterHandler(service.Server(), new(Hrefs))

	protobuf.RegisterUserHandler(service.Server(), new(User))


	if err := service.Run(); err != nil {
		logger.Info.Println(err)
	}
}

func Stop(){
	service.Server().Stop()
	logger.Info.Println( "micro 微服务结束")
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


func StartMicroService1() {
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("user"),
	)

	service.Init()

	protobuf.RegisterUserHandler(service.Server(), new(User))

	if err := service.Run(); err != nil {
		logger.Info.Println(err)
	}
}
