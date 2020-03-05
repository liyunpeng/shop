package routes

import (
	"github.com/kataras/iris/v12"
	"shop/web/controllers"
	"shop/web/middleware"
)

func RegisterApi(app *iris.Application){
	api := app.Party("/api", middleware.CorsAuth()).AllowMethods(iris.MethodOptions)

	api.PartyFunc("/user", func(party iris.Party){
		party.Get("/",  controllers.ApiUserGetAll)
		party.Get("/{id:uint}",  controllers.ApiUserGetById)
		party.Post("/",  controllers.ApiUserPost)
		party.Put("/",  controllers.ApiUserUpdate)
		party.Post("/insertOrUpdate",  controllers.ApiUserInsertOrUpdate)
	})

	api.PartyFunc("/database", func(party iris.Party){
		party.Post("/create", controllers.ApiDatabaseCreate)
	})
}
