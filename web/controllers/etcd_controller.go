package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"shop/services"
	"shop/validates"
	"strings"
)

type EtcdController struct {
	Ctx iris.Context

	Service services.EtcdService
}

var v = mvc.View{
	Name: "conf_manager.html",
}

func (e *EtcdController) Get() mvc.Result {
	return v
}

type MonitorFile struct {
	FileName string `json:"filename"`
	FileSize int  `json:"filesize"`
	FileKeyWords string `json:"filekeywords"`
}

func ApiEtcdGetKV(ctx iris.Context) {

	fmt.Println(" Apiectcd get kv ")
	s := []MonitorFile{
		{
			FileName:"log1.txt",
			FileSize: 1000,
			FileKeyWords: "abc",
		},
		{
			FileName:"log2.txt",
			FileSize: 2000,
			FileKeyWords: "edf",
		},
	}
	ctx.JSON(ApiResource(true, s,  "获取etcdkvcheng"))
}

//func (e *EtcdController) Post() mvc.Response {
//	fmt.Println("ApiUserPost is called")
//	k := e.Ctx.FormValue("key")
//	v := e.Ctx.FormValue("value")
//
//	e.Service.PutKV(k, v)
//	return mvc.Response{
//		Text: "ok",
//	}
//}
func (e *EtcdController) Post() {
	fmt.Println("ApiUserPost is called")
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
	k := e.Ctx.FormValue("k")

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
