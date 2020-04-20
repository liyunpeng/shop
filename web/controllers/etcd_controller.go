package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"shop/client"
	"shop/config"
	"shop/util"
	"shop/validates"
	"strconv"
	"strings"
)

type EtcdController struct {
	Ctx     iris.Context
	Service client.EtcdClientWrap
}

var v = mvc.View{
	Name: "etcd/etcd_manager.html",
}

func (e *EtcdController) Get() mvc.Result {
	return v
}

type MonitorFile struct {
	FileName     string `json:"filename"`
	FileSize     int    `json:"filesize"`
	FileKeyWords string `json:"filekeywords"`
}

func ApiEtcdGetKV(ctx iris.Context) {
	keya := ctx.Params().Get("key")
	fmt.Println("api调用 ApiEtcdGetKV ,请求参数为", keya)
	var ss []config.LogConfig
	//resp :=	client.EtcdClientInsance.Get("/logagent/192.168.0.142/logconfig")
	resp := client.EtcdClientInsance.Get(keya)
	if resp == nil {
		ctx.JSON(ApiResource(true, nil, "请求出错"))
		return
	} else if resp.Count < 1 {
		ctx.JSON(ApiResource(true, nil, "请求的 key 不存在"))
		return
	}
	for _, ev := range resp.Kvs {
		//v := string(ev.Value)
		//json.Unmarshal([]byte(v), &logConfArr)
		var logConfArr []config.LogConfig
		json.Unmarshal(ev.Value, &logConfArr)
		ss = append(ss, logConfArr...)
		fmt.Printf("etcdkey = %s \n etcdvalue = %s \n", ev.Key, ev.Value)

	}

	//s := []MonitorFile{
	//	{
	//		FileName:"log1.txt",
	//		FileSize: 1000,
	//		FileKeyWords: "abc",
	//	},
	//	{
	//		FileName:"log2.txt",
	//		FileSize: 2000,
	//		FileKeyWords: "edf",
	//	},
	//}
	ctx.JSON(ApiResource(true, ss, "获取etcdkvcheng"))
}

type EtcdOption struct {
	Label     string `json:"label"`
	Etcdkey   string `json:"etcdkey"`
	EtcdValue string `json:"etcdvalue"`
}

func ApiEtcdListAllKV(ctx iris.Context) {
	fmt.Println(" Apiectcd list all kv ")
	s := []EtcdOption{
		{
			//Label:"/logagent/192.168.0.1/logconfig",
			Label:     "192.168.0.1",
			Etcdkey:   "1000",
			EtcdValue: "abfffc",
		},
		{
			//Label:"/logagent/192.168.0.2/logconfig",
			Label:     "192.168.0.2",
			Etcdkey:   "1000",
			EtcdValue: "affbc",
		},
		{
			//Label:"/logagent/192.168.0.3/logconfig",
			Label:     "192.168.0.3",
			Etcdkey:   "10l00",
			EtcdValue: "abc",
		},
	}

	ctx.JSON(ApiResource(true, s, "获取etcdkvcheng"))
}

func (e *EtcdController) Post() mvc.Response {
	fmt.Println("ApiUserPost is called")
	k := e.Ctx.FormValue("key")
	v := e.Ctx.FormValue("value")

	e.Service.PutKV(k, v)
	return mvc.Response{
		Text: "ok",
	}
}

func (e *EtcdController) PostKV() {
	fmt.Println("ApiUserPostKV is called")
	//var i interface{}
	//i = new()
	aul := new(validates.CreateEtcdKVRequest)
	if err := e.Ctx.ReadJSON(aul); err != nil {
		e.Ctx.StatusCode(iris.StatusOK)
		_, _ = e.Ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	e.Service.PutKV(aul.Key, aul.Value)
	_, _ = e.Ctx.JSON(ApiResource(true, nil, "成功添加etcd键值对"))
}
func (e *EtcdController) GetAll() {
	fmt.Println("api GetAll is called")

	k := "/logagent/192.168.0.142/logconfig"
	resp := e.Service.Get(k)

	var v strings.Builder

	m := make(map[string]interface{}, 100)
	for _, ev := range resp.Kvs {
		v.WriteString(string(ev.Value))
		m[k] = string(ev.Value)
		fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
	}

	_, _ = e.Ctx.JSON(ApiResource(true, v.String(), "成功添加etcd键值对"))
}

//func ApiEtcdPost(ctx context2.Context) {
//	fmt.Println("ApiUserPost is called")
//	aul := new(validates.CreateEtcdKVRequest)
//	if err := ctx.ReadJSON(aul); err != nil {
//		ctx.StatusCode(iris.StatusOK)
//		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
//		return
//	}
//
//	//e.Service.PutKV(aul.Key, aul.Value)
//	_, _ = ctx.JSON(ApiResource(true, nil, "成功添加etcd键值对"))
//}

func (e *EtcdController) GetKv() string {
	util.PrintFuncName()
	k := e.Ctx.FormValue("keya")
	resp := e.Service.Get(k)
	var v strings.Builder
	for _, ev := range resp.Kvs {
		v.WriteString(string(ev.Value))
		fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
	}

	return v.String()
}

//func  (e *EtcdController)PostAdd() mvc.Result{
//	f := e.Ctx.FormValue("data")
//	e.Ctx
//
//	return mvc.Response{
//		//如果不是nil，则会显示此错误
//		Err: err,
//	}
//}

func ApiRedisSet(ctx iris.Context) {
	fmt.Println("redis 操作")
	for i :=0; i< 100; i++ {
		client.RedisSetString( strconv.Itoa(i), "aaaaaaaa")
	}
	ctx.JSON(ApiResource(true, nil, "redis操作"))
	return
}
