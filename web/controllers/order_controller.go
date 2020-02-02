package controllers
import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	//"github.com/kataras/iris/v12/sessions"
)

type OrderController struct {
	Ctx iris.Context
}

var orderStaticView = mvc.View{
	Name: "order.html",
	Data: iris.Map{"Title": "User Registration"},
}

func (c *OrderController) Get() mvc.Result {
	return orderStaticView
}
