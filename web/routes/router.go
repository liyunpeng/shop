package routes

import (
	"github.com/kataras/iris/v12"
	"shop/web/controllers"
)

func RegisterApi(app *iris.Application){
	api := app.Party("/api")

	api.PartyFunc("/user", func(party iris.Party){
		party.Get("/",  controllers.ApiUserGetAll)
		party.Get("/{id:uint}",  controllers.ApiUserGetById)
		party.Post("/",  controllers.ApiUserPost)
	})
}
