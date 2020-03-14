package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	//context2 "github.com/kataras/iris/v12/context"
	"shop/models"
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
	aul := new(validates.CreateEtcdKVRequest)
	if err := e.Ctx.ReadJSON(aul); err != nil {
		e.Ctx.StatusCode(iris.StatusOK)
		_, _ = e.Ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	e.Service.PutKV(aul.Key, aul.Value)
	_, _ = e.Ctx.JSON(ApiResource(true, nil, "成功添加etcd键值对"))
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

func ApiEtcdGetKv(ctx iris.Context) {
	fmt.Println("ApiUserPost is called")
	aul := new(validates.CreateUpdateUserRequest)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}
	u := models.User{}
	u.Username = aul.Username
	u.Password = aul.Password
	u.Phonenumber = aul.Phonenumber
	u.Level = aul.Level
	models.UserInsert(&u)
	//k := e.Ctx.FormValue("key")
	//v := e.Ctx.FormValue("value")
	//
	//e.Service.PutKV(k, v)
//	return mvc.Response{
//		Text: "ok",
//	}
//}

	_, _ = ctx.JSON(ApiResource(true, nil, "成功添加数据行"))
}
