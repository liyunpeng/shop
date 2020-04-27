package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/client"
	"shop/logger"
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
	err := client.MicroCall("IndexLinks", 10, rsp)
	if err != nil {
		logger.Info.Println("err =",err )
	}else{
		logger.Info.Println("客户端调用微服务的结果 =", rsp.Name )
	}
	logger.Info.Println("cookiename =",cookieName)
	logger.Info.Println("session name=", c.Session.Get(util.SessionUserID))

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
