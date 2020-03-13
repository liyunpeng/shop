package controllers

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"shop/services"
	"strings"
)

type EtcdManangerController struct {
	Ctx iris.Context

	Service services.EtcdService
}

var v = mvc.View{
	Name: "conf_manager.html",
}

func (e *EtcdManangerController) Get() mvc.Result {
	return v
}

func (e *EtcdManangerController) Post() mvc.Response {
	k := e.Ctx.FormValue("key")
	v := e.Ctx.FormValue("value")

	e.Service.PutKV(k, v)
	return mvc.Response{
		Text: "ok",
	}
}

func (e *EtcdManangerController) GetKv() string {
	k := e.Ctx.FormValue("k")

	resp := e.Service.Get(k)

	var v strings.Builder

	for _, ev := range resp.Kvs {
		v.WriteString(string(ev.Value))
		fmt.Printf("etcd key = %s , etcd value = %s", ev.Key, ev.Value)
	}

	return v.String()
}

//func  (e *EtcdManangerController)PostAdd() mvc.Result{
//	f := e.Ctx.FormValue("data")
//	e.Ctx
//
//	return mvc.Response{
//		//如果不是nil，则会显示此错误
//		Err: err,
//	}
//}
