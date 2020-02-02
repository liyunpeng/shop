package controllers
import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ShopCarController struct {
	Ctx iris.Context
}

var shopCarStaticView = mvc.View{
	Name: "shopcar.html",
	Data: iris.Map{"Title": "User Registration"},
}

func (c *ShopCarController) Get() mvc.Result {
	return shopCarStaticView
}
