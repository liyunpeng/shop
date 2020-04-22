package config

import (
	"github.com/kataras/iris/v12"
	gf "github.com/snowlyg/gotransformer"
	"shop/transformer"
	"time"
)

var TransformConfiguration *transformer.Conf

func GetTransformConfiguration( irisConfiguration iris.Configuration) *transformer.Conf {
	app := transformer.App{}
	g := gf.NewTransform(&app, irisConfiguration.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Mysql{}
	g.OutputObj = &db
	g.InsertObj = irisConfiguration.Other["Mysql"]
	_ = g.Transformer()

	mongodb := transformer.Mongodb{}
	g.OutputObj = &mongodb
	g.InsertObj = irisConfiguration.Other["Mongodb"]
	_ = g.Transformer()

	redis := transformer.Redis{}
	g.OutputObj = &redis
	g.InsertObj = irisConfiguration.Other["Redis"]
	_ = g.Transformer()

	sqlite := transformer.Sqlite{}
	g.OutputObj = &sqlite
	g.InsertObj = irisConfiguration.Other["Sqlite"]
	_ = g.Transformer()

	testData := transformer.TestData{}
	g.OutputObj = &testData
	g.InsertObj = irisConfiguration.Other["TestData"]
	_ = g.Transformer()

	kafkaConf := transformer.Kafka{}
	g.OutputObj = &kafkaConf
	g.InsertObj = irisConfiguration.Other["Kafka"]
	_ = g.Transformer()

	etcdConf := transformer.EtcdConf{}
	g.OutputObj = &etcdConf
	g.InsertObj = irisConfiguration.Other["Etcd"]
	_ = g.Transformer()

	grpcConf := transformer.GrpcConf{}
	g.OutputObj = &grpcConf
	g.InsertObj = irisConfiguration.Other["Grpc"]
	_ = g.Transformer()

	consulConf := transformer.ConsulConf{}
	g.OutputObj = &consulConf
	g.InsertObj = irisConfiguration.Other["Consul"]
	_ = g.Transformer()

	cf := &transformer.Conf{
		App:      app,
		Mysql:    db,
		Mongodb:  mongodb,
		Redis:    redis,
		Sqlite:   sqlite,
		TestData: testData,
		Kafka: kafkaConf,
		Etcd: etcdConf,
		Grpc: grpcConf,
		Consul: consulConf,
	}

	return cf
}
