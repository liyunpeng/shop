package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/cli"
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
	fmt.Println("cookieName=", cookieName)

	//iris.WithCharset("UTF-8")

	rsp1 := new(models.User)
	err1 := cli.Call("IndexLinks", 10, rsp1)
	if err1 != nil {
		fmt.Println("err =",err1 )
	}else{
		fmt.Println("客户端调用微服务的结果 =", rsp1.Name )
	}
	//rsp := new(models.Order)
	//err := cli.Call("GetOrderByUser", "aa", rsp)
	//var s []Order
	//rsp :=  result.Orders
	//rsp :=  new([]models.Order) //iris.WithCharset("UTF-8"))
	//result.Orders = models.OrderFindByUser("aa")
	//err := cli.Call("GetOrderByUser", "aa", rsp)
	//go func() {
		//rsp := make([]models.Order, 1)


	//} ()

	//go func (){
	//	rsp := new(models.Order)
	//	err := cli.Call("GetOrderById", 1, rsp)
	//	fmt.Println("rsp.name=", rsp.Username)
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
		fmt.Println("用户已经登录")
		result := new( Result)
		orderItems := new( models.OrderItems)
		err := cli.Call("GetOrderByUser", sessionUserName, orderItems)

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
				"UserId": c.Session.GetString(util.SessionUserName),
			},
		}
	} else {
		fmt.Println("用户没有登录")
		return mvc.View{
			Name: "order.html",
			Data: iris.Map{
				"OrderCount": "10002",
				"UserId": c.Session.GetString(util.SessionUserName),
			},
			//Data: map[string] interface{}{
			//	"Result": result,
			//},
		}

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
