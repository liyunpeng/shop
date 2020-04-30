package client

import (
	"shop/custchan"
	"shop/logger"
	"shop/transformer"
	"shop/util"
	"time"
)

func StartClient(transformConfiguration *transformer.Conf) {
	etcdKeys := GetEtcdKeys()
	EtcdClientInsance = NewEtcdClientWrap(
		[]string{transformConfiguration.Etcd.Addr}, 5 * time.Second)
	go func() {
		logger.Info.Println("到etcd服务器，按指定的键遍历键值对")
		for _, key := range etcdKeys {
			resp := EtcdClientInsance.Get(key)
			if resp != nil || resp.Count < 1 {
				continue
			}
			for _, ev := range resp.Kvs {
				custchan.ConfChan <- string(ev.Value)
				logger.Info.Printf("etcdkey = %s \n etcdvalue = %s \n", ev.Key, ev.Value)
			}
		}
	}()

	// 启动对etcd的监听服务，有新的键值对会被监听到
	util.WaitGroup.Add(1)

	go EtcdClientInsance.EtcdWatch( util.Ctx , etcdKeys)


	util.WaitGroup.Add(1)
	go StartKafkaProducer(
		transformConfiguration.Kafka.Addr, 1, true)

	util.WaitGroup.Add(1)
	go StartKafkaConsumer(transformConfiguration.Kafka.Addr)

	//go StartOauth2Client()

	go StartGrpcClient()

	StartRedisClient()

	StartRpcClient()
}

