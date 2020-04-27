package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/logger"
	"shop/models"
	"shop/util"

	//"github.com/kataras/iris/v12/sessions"
)


type BuyController struct {
	Ctx iris.Context
	Session *sessions.Session
}

type BuyResult struct {
	Id int	`json:"id"`
	Title string	`json:"title"`
	Item int  `json:"item"`
	Goods *models.Goods
	User *models.User
}

func (c *BuyController) GetBy(goodsID int64) mvc.Result {
	//goods := models.GoodsFindAll()
	//result := GoodsResult{
	//	Item: 2,
	//	Goods: goods,
	//}
	result := new(BuyResult)
	result.Goods = 	models.GoodsFindById(goodsID)
	result.User = models.UserFindById(util.GetCurrentUserID(c.Session))
	logger.Info.Println("用户", result.User.Username, "name=", result.User.Name,  "进入购买页面，商品ID=", result.Goods.ID)

	var buyStaticView = mvc.View{
		Name: "buy.html",
		Data: iris.Map{
			"Result": result,
			"Title": "User Registration",
		},
	}
	return buyStaticView
}


