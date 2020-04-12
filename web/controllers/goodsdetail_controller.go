package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/models"
	"shop/util"

	//"github.com/kataras/iris/v12/sessions"
)

type GoodsDetailController struct {
	Ctx     iris.Context
	Session *sessions.Session
}

type GoodsDetailResult struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Orders []*models.Order
	Item   int `json:"item"`
}

//type Order struct {
//	Id int 	`json:"id"`
//	Title string 	`json:"title"`
//}
func (c *GoodsDetailController) Get() mvc.Result {

	//result := new(Result)
	//orderItems := new(models.OrderItems)
	//result.Orders = orderItems.Items
	//result.Id = 1001
	return mvc.View{
		Name: "goods-detail.html",
		Data: iris.Map{
			//"Result":     result,
			"OrderCount": "10",
			//"UserId":     c.Session.GetString(util.SessionUserName),
		},
	}
}

func (c *GoodsDetailController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(util.SessionUserName, 0)
	return userID
}

func (c *GoodsDetailController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *GoodsDetailController) logout() {
	c.Session.Destroy()
}
