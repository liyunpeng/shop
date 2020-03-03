package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/liyunpeng/shop/web/controllers"
)

func RegisterApi(app *iris.Application){
	api := app.Party("/api")
	api.PartyFunc("/user", func(party iris.Party){
		party.Get("/",  controllers.ApiUserGet)
	})
}
