package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"shop/models"
	"shop/validates"
	"shop/web/controllers"
	"shop/web/middleware"
	"strings"
)

func RegisterApi(app *iris.Application){
	api := app.Party("/api", middleware.CorsAuth()).AllowMethods(iris.MethodOptions)

	api.Post("/login", controllers.UserLogin)
	api.PartyFunc("/user", func(party iris.Party){
		casbinMiddleware := middleware.New(models.Enforcer)                  //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
		party.Use(middleware.JwtHandler().Serve) //登录验证
		party.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
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

func isPermRoute(s *router.Route) bool {
	exceptRouteName := []string{"OPTIONS", "GET", "POST", "HEAD", "PUT", "PATCH"}
	for _, er := range exceptRouteName {
		if strings.Contains(s.Name, er) {
			return true
		}
	}
	return false
}

func GetRoutes(api *iris.Application) []*validates.PermissionRequest {
	rs := api.APIBuilder.GetRoutes()
	var rrs []*validates.PermissionRequest
	for _, s := range rs {
		if !isPermRoute(s) {
			path := strings.Replace(s.Path, ":id", "*", 1)
			rr := &validates.PermissionRequest{Name: path, DisplayName: s.Name, Description: s.Name, Act: s.Method}
			rrs = append(rrs, rr)
		}
	}
	return rrs
}

