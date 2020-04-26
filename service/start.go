package service

import (
	"shop/transformer"
	"shop/util"
)

func StartService(transformConfiguration *transformer.Conf) {
	go StartTailService()

	go StartWebSocketService()

	util.WaitGroup.Add(1)
	go StartGrpcService(transformConfiguration.Grpc)

	go StartMicroService()

	//go StartOauth2Service()
}
