package client

import (
	"context"
	"fmt"
	"shop/custchan"
	"shop/logger"
	"shop/util"
	"time"

	//"github.com/astaxie/beego/logs"
	//producerClient "github.com/coreos/etcd/clientv3"
	client "go.etcd.io/etcd/clientv3"
)

var (

	EtcdClientInsance *etcdClientWrap
)

type EtcdClientWrap interface {
	PutKV(key string, value string)
	Get(key string) (*client.GetResponse)
}

type etcdClientWrap struct {
	EtcdClient *client.Client
	EtcdKV client.KV
}

func NewEtcdClientWrap(addrs []string, timeout time.Duration) *etcdClientWrap {
	etcdClient, err := client.New(client.Config{
		Endpoints:   addrs,
		DialTimeout: timeout,
	})
	if err != nil {
		logger.Info.Println("etcd 连接失败， err=", err)
	}else {
		logger.Info.Println("etcd 连接成功")
	}
	kv := client.NewKV(etcdClient)

	e := &etcdClientWrap{
		EtcdClient: etcdClient,
		EtcdKV:         kv,
	}
	return e
}

func (e *etcdClientWrap) PutKV(key string, value string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := e.EtcdKV.Put(ctx, key, value) //withPrevKV()是为了获取操作前已经有的key-value

	if err != nil {
		panic(err)
	}

	cancel()

	//fmt.Printf("kvs1: %v", putResp.PrevKv)
}


func (e *etcdClientWrap) Get(key string) (*client.GetResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	getResp, err := e.EtcdKV.Get(ctx, key) //withPrefix()是未了获取该key为前缀的所有key-value
	if err != nil {
		logger.Info.Println("etcd get key 出错：", err)
	}
	//fmt.Printf("kvs2:  %v", getResp.Kvs)
	cancel()
	return getResp
}
func (e *etcdClientWrap) EtcdWatch(ctx context.Context, keys []string) {
	defer util.WaitGroup.Done()

	defer util.PrintFuncName()

	var watchChans []client.WatchChan
	for _, key := range keys {
		rch := e.EtcdClient.Watch(context.Background(), key)
		logger.Info.Println("添加要watch的key，key的值=", key)
		watchChans = append(watchChans, rch)
	}

	for {
		for _, watchC := range watchChans {
			select {
			case <-util.ChanStop:
				logger.Info.Println("etcd watch 协程 退出")
				return
			case wresp := <-watchC:
				for _, ev := range wresp.Events {
					custchan.ConfChan <- string(ev.Kv.Value)
					fmt.Printf("etcd服务watch到新的键值对： etcd key = %s , etcd value = %s \n", ev.Kv.Key, ev.Kv.Value)
				}
			case <-ctx.Done():
				util.Logger.Info(" 关闭配置通道，custchan.ConfChan 通道长度=", len(custchan.ConfChan))
				close(custchan.ConfChan)
				util.Logger.Info("EtcdWatch 退出对键值对的监控")
				return
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

//GetEtcdConfChan is func get etcd conf add to chan
func (e *etcdClientWrap) GetEtcdConfChan() chan string {
	return custchan.ConfChan
}

func GetEtcdKeys() ([]string) {
	var etcdKeys []string
	//ips, err := getLocalIP()
	var ips []string
	//var err error
	ips = append(ips, "192.168.0.1")
	//if err != nil {
	//	logger.Info.Println("get local ip error:", err)
	//	//return err
	//}
	for _, ip := range ips {
		//key := fmt.Sprintf("/logagent/%s/logconfig", ip)
		etcdKeys = append(etcdKeys, ip)
	}
	logger.Info.Println("从etcd服务器获取到的以IP名为键的键值对: ", etcdKeys)
	return etcdKeys
}
