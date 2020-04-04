package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/cli"
	"shop/models"
	"shop/util"
)

type IndexController struct {
	Ctx iris.Context
	Session *sessions.Session
}

var indexStaticView = mvc.View{
	Name: "index.html",
	Data: iris.Map{"Title": "User Registration"},
}

func (c *IndexController) Get() mvc.Result {
	//cookieName := c.Session.Get(util.COOKEI_NAME)
	cookieName := c.Ctx.GetCookie(util.COOKEI_NAME)
	rsp := new(models.User)
	err := cli.Call("IndexLinks", true, rsp)
	if err != nil {
		fmt.Println("err =",err )
	}else{
		fmt.Println("客户端调用微服务的结果 =", rsp.Name )
	}
	fmt.Println("cookiename =",cookieName)
	fmt.Println("session name=", c.Session.Get(util.SessionUserName))
	return indexStaticView
}