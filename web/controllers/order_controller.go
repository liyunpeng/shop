package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/util"

	//"github.com/kataras/iris/v12/sessions"
)

type OrderController struct {
	Ctx     iris.Context
	Session *sessions.Session
}

var orderStaticView = mvc.View{
	Name: "order.html",
	Data: iris.Map{"Title": "User Registration"},
}

func (c *OrderController) Get() mvc.Result {

	cookieName := c.Ctx.GetCookie(util.COOKEI_NAME)
	fmt.Println("cookieName=", cookieName)


	if len(c.Session.GetString(util.SessionUserName)) > 0 {
		fmt.Println("用户已经登录")
		return mvc.View{
			Name: "order.html",
			Data: iris.Map{
				"OrderCount": "10",
				"UserId": c.Session.GetString(util.SessionUserName),
			},
		}
	} else {
		fmt.Println("用户没有登录")
		return orderStaticView
	}
}

func (c *OrderController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(util.SessionUserName, 0)
	return userID
}

func (c *OrderController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *OrderController) logout() {
	c.Session.Destroy()
}
