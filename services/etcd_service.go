package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	//"github.com/astaxie/beego/logs"
	//producerClient "github.com/coreos/etcd/clientv3"
	client "go.etcd.io/etcd/clientv3"
)

var (
	ConfChan  = make(chan string, 10)
	waitGroup sync.WaitGroup
)

type EtcdService interface {
	PutKV(key string, value string)
	Get(key string) (*client.GetResponse)
}

type etcdService struct {
	EtcdClient *client.Client

	EtcdKV client.KV
}

func NewEtcdService(addrs []string, timeout time.Duration) *etcdService {
	etcdClient, _ := client.New(client.Config{
		Endpoints:   addrs,
		DialTimeout: timeout,
	})
	kv := client.NewKV(etcdClient)

	e := &etcdService{
		EtcdClient: etcdClient,
		EtcdKV:     kv,
	}

	return e
}

func (e *etcdService) PutKV(key string, value string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	_, err := e.EtcdKV.Put(ctx, key, value) //withPrevKV()是为了获取操作前已经有的key-value

	if err != nil {
		panic(err)
	}

	cancel()

	//fmt.Printf("kvs1: %v", putResp.PrevKv)
}


func (e *etcdService) Get(key string) (*client.GetResponse) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	getResp, err := e.EtcdKV.Get(ctx, key) //withPrefix()是未了获取该key为前缀的所有key-value
	if err != nil {
		panic(err)
	}
	//fmt.Printf("kvs2:  %v", getResp.Kvs)
	cancel()
	return getResp
}

func (e *etcdService) EtcdWatch(keys []string) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	var watchChans []client.WatchChan
	for _, key := range keys {
		rch := e.EtcdClient.Watch(context.Background(), key)
		fmt.Println("添加要watch的key，key的值=", key)
		watchChans = append(watchChans, rch)
	}

	for {
		for _, watchC := range watchChans {
			select {
			case wresp := <-watchC:
				for _, ev := range wresp.Events {
					ConfChan <- string(ev.Kv.Value)
					fmt.Printf("etcd服务watch到新的键值对： etcd key = %s , etcd value = %s \n", ev.Kv.Key, ev.Kv.Value)
				}
			default:
			}
		}
		time.Sleep(time.Second)
	}
}

//GetEtcdConfChan is func get etcd conf add to chan
func (e *etcdService) GetEtcdConfChan() chan string {
	return ConfChan
}
