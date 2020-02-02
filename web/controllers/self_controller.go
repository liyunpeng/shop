package controllers
import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"

	//"github.com/kataras/iris/v12/context"

	"github.com/kataras/iris/v12/mvc"

	//"github.com/kataras/iris/v12/sessions"
)

type SelfController struct {
	Ctx iris.Context
	Session *sessions.Session
}

var selfStaticView = mvc.View{
	Name: "self.html",
	// TODO 不理解为什么要去掉这句 Get方法里通过viewdata设置的Title才能有效
	//Data: iris.Map{"Title": "User Registration"},
}

func (c *SelfController) Get() mvc.Result {
	if len(c.Session.GetString("UserID")) > 0 {
		fmt.Println("session.GetString=", c.Session.GetString("UserID"))
		c.Ctx.ViewData("Title", c.Session.GetString("UserID"))

	} else {
		fmt.Println("session is nil")
		c.Ctx.ViewData("Title", "oo")
	}
	return selfStaticView
}
