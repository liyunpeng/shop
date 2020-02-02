package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	//"github.com/kataras/iris/v12/sessions"
)


type AssortController struct {

	Ctx iris.Context

}

var assortStaticView = mvc.View{
	Name: "assort.html",
	Data: iris.Map{"Title": "User Registration"},
}

func (c *AssortController) Get() mvc.Result {
	return assortStaticView
}


