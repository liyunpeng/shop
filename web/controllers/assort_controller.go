package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"shop/models"

	//"github.com/kataras/iris/v12/sessions"
)


type AssortController struct {

	Ctx iris.Context

}
type GoodsResult struct {
	Id int	`json:"id"`
	Title string	`json:"title"`
	Item int  `json:"item"`
	Goods []*models.Goods
}

func (c *AssortController) Get() mvc.Result {
	goods := models.GoodsFindAll()
	result := GoodsResult{
		Item: 2,
		Goods: goods,
	}

	var assortStaticView = mvc.View{
		Name: "assort.html",
		Data: iris.Map{
			"Result": result,
			"Title": "User Registration",
		},
	}
	return assortStaticView
}


