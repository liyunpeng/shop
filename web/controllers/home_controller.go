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

type HomeController struct {
	Ctx iris.Context
	Session *sessions.Session
}



func (c *HomeController) Get() mvc.Result {
	//cookieName := c.Session.Get(util.COOKEI_NAME)
	cookieName := c.Ctx.GetCookie(util.COOKEI_NAME)
	rsp := new(models.User)
	err := cli.Call("IndexLinks", 10, rsp)
	if err != nil {
		fmt.Println("err =",err )
	}else{
		fmt.Println("客户端调用微服务的结果 =", rsp.Name )
	}
	fmt.Println("cookiename =",cookieName)
	fmt.Println("session name=", c.Session.Get(util.SessionUserName))

	result := Result{
		Item: 1,
	}
	var indexStaticView = mvc.View{
		Name: "home.html",
		Data: iris.Map{
			"Result": result,
			"Title": "User Registration",
		},
	}
	return indexStaticView
}