package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	//"github.com/kataras/iris/v12/sessions"
)


type AssortController struct {

	Ctx iris.Context

}



func (c *AssortController) Get() mvc.Result {
	result := Result{
		Item: 2,
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


