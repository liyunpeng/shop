package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"shop/models"

	//"github.com/kataras/iris/v12/sessions"
)


type BuyController struct {

	Ctx iris.Context

}
type BuyResult struct {
	Id int	`json:"id"`
	Title string	`json:"title"`
	Item int  `json:"item"`
	Goods []*models.Goods
}

func (c *BuyController) Get() mvc.Result {
	goods := models.GoodsFindAll()
	result := GoodsResult{
		Item: 2,
		Goods: goods,
	}

	var buyStaticView = mvc.View{
		Name: "buy.html",
		Data: iris.Map{
			"Result": result,
			"Title": "User Registration",
		},
	}
	return buyStaticView
}


