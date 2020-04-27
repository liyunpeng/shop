package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/client"
	"shop/logger"
	"shop/models"
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

type Result struct {
	Id int	`json:"id"`
	Title string	`json:"title"`
	Orders []*models.Order
	Item int  `json:"item"`
}

//type Order struct {
//	Id int 	`json:"id"`
//	Title string 	`json:"title"`
//}
func (c *OrderController) Get() mvc.Result {

	cookieName := c.Ctx.GetCookie(util.COOKEI_NAME)
	logger.Info.Println("cookieName=", cookieName)

	//iris.WithCharset("UTF-8")

	rsp1 := new(models.User)
	err1 := client.MicroCall("IndexLinks", 10, rsp1)
	if err1 != nil {
		logger.Info.Println("err =",err1 )
	}else{
		logger.Info.Println("客户端调用微服务的结果 =", rsp1.Name )
	}

	client.MicroCallUser()
	//rsp := new(models.Order)
	//err := cli.MicroCall("GetOrderByUser", "aa", rsp)
	//var s []Order
	//rsp :=  result.Orders
	//rsp :=  new([]models.Order) //iris.WithCharset("UTF-8"))
	//result.Orders = models.OrderFindByUser("aa")
	//err := cli.MicroCall("GetOrderByUser", "aa", rsp)
	//go func() {
		//rsp := make([]models.Order, 1)


	//} ()

	//go func (){
	//	rsp := new(models.Order)
	//	err := cli.MicroCall("GetOrderById", 1, rsp)
	//	logger.Info.Println("rsp.name=", rsp.Username)
	//	result.Orders = append(result.Orders, rsp)
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//	result.Id = 1001
	//}()

	//s1 := Order{
	//	Id:1,
	//	Title: "titile1",
	//}
	//s = append(s, s1)
	//s2 := Order{
	//	Id:1,
	//	Title: "titile2",
	//}
	//s = append(s, s2)
	//c.Ctx.ViewData("Result", result)
	//mvc.View{}.Data = s
	sessionUserName := 	c.Session.GetString(util.SessionUserName)
	if len(sessionUserName) > 0 {
		logger.Info.Println("用户已经登录")
		result := new( Result)
		orderItems := new( models.OrderItems)
		err := client.MicroCall("GetOrderByUser", sessionUserName, orderItems)
		result.Orders = orderItems.Items
		if err != nil {
			panic(err)
		}
		result.Id = 1001
		return mvc.View{
			Name: "order.html",
			Data: iris.Map{
				"Result": result,
				"OrderCount": "10",
				"UserId": c.Session.GetString(util.SessionUserID),
			},
		}
	} else {
		logger.Info.Println("用户没有登录")
		return mvc.View{
			Name: "order.html",
			Data: iris.Map{
				"OrderCount": "10002",
				"UserId": c.Session.GetString(util.SessionUserID),
			},
		}
	}
}

