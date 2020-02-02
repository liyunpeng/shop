package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type IndexController struct {
	Ctx iris.Context
}

var indexStaticView = mvc.View{
	Name: "index.html",
	Data: iris.Map{"Title": "User Registration"},
}

func (c *IndexController) Get() mvc.Result {
	return indexStaticView
}