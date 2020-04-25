package routes

import (
	"github.com/gorilla/securecookie"
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

	api.Post("/login", controllers.ApiUserLogin).Name = "登录"
	api.PartyFunc("/user", func(party iris.Party){
		casbinMiddleware := middleware.New(models.Enforcer)                  //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
		party.Use(middleware.JwtHandler().Serve) //登录验证
		party.Use(casbinMiddleware.ServeHTTP) //权限验证
		party.Get("/",  controllers.ApiUserGetAll).Name = "获取所有用户"
		party.Get("/{id:uint}",  controllers.ApiUserGetById).Name = "获取指定用户"
		party.Post("/",  controllers.ApiUserPost).Name = "创建用户"
		party.Put("/",  controllers.ApiUserUpdate).Name = "修改用户"
		party.Post("/insertOrUpdate",  controllers.ApiUserInsertOrUpdate).Name = "创建或修改用户"
	})

	api.PartyFunc("/etcd", func(party iris.Party){
		party.Get("/{key:string}",  controllers.ApiEtcdGetKV).Name = "获取etcdkeyValue"
		party.Get("/listallkeys",  controllers.ApiEtcdListAllKV).Name = "获取etcdkeyValue"
		//party.Get("/{id:uint}",  controllers.ApiEtcdGetKv).Name = "获取kv"
		//party.Post("/",  controllers.ApiEtcdPost).Name = "创建etcdkv"
		//party.Put("/",  controllers.ApiUserUpdate).Name = "修改用户"
		//party.Post("/insertOrUpdate",  controllers.ApiUserInsertOrUpdate).Name = "创建或修改用户"
	})

	api.PartyFunc("/database", func(party iris.Party){
		party.Post("/create", controllers.ApiDatabaseCreate).Name = "创建初始数据库"
	})
	api.PartyFunc("/mem", func(party iris.Party){
		party.Post("/redisset", controllers.ApiRedisSet).Name = "redis操作"
	})

	// cookie 加密操作
	cookieGet(app)
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


var (
	// AES仅支持16,24或32字节的密钥大小。
	//您需要准确提供该密钥字节大小，或者从您键入的内容中获取密钥。
	hashKey  = []byte("the-big-and-secret-fash-key-here")
	blockKey = []byte("lot-secret-of-characters-big-too")
	sc       = securecookie.New(hashKey, blockKey)
)

func cookieGet( app *iris.Application ) {

	app.Get("/cookies/{name}/{value}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		value := ctx.Params().Get("value")
		//加密值
		ctx.SetCookieKV(name, value, iris.CookieEncode(sc.Encode)) // <--设置一个Cookie
		ctx.Writef("cookie added: %s = %s", name, value)
	})
	app.Get("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		//解密值
		value := ctx.GetCookie(name, iris.CookieDecode(sc.Decode)) // <--检索Cookie
		ctx.WriteString(value)
	})
	app.Delete("/cookies/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.RemoveCookie(name) // <-- 删除Cookie
		ctx.Writef("cookie %s removed", name)
	})
}
