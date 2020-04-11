package controllers
import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ShopCarController struct {
	Ctx iris.Context
}


func (c *ShopCarController) Get() mvc.Result {
	result := new( Result)
	result.Item = 3
	shopCarStaticView := mvc.View{
		Name: "shopcar.html",
		Data: iris.Map{
			"Result": result,
			"Title": "User Registration"},
	}

	return shopCarStaticView
}
