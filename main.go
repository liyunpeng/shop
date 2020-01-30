// file: main.go

package main

import (
	"fmt"
	"time"

	"shop/datasource"
	"shop/repositories"
	"shop/services"
	"shop/web/controllers"
	"shop/web/middleware"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/jinzhu/gorm"
)

func getDb()(db *gorm.DB){
	defer func(){
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()

	//TODO？ 全局db没用
	/*
		链接localhost数据库， 用户名root, 密码root
	*/
	db, err := gorm.Open(
		"mysql", "root:root@/gotest?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		fmt.Println("open db sucess", db)

	} else {
		fmt.Println("open db error ", err)
		panic("数据库错误")
	}

	return db
}

func main() {
	app := iris.New()
	// You got full debug messages, useful when using MVC and you want to make
	// sure that your code is aligned with the Iris' MVC Architecture.
	app.Logger().SetLevel("debug")

	// Load the template files.
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)

	app.HandleDir("/public", "./web/public")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	// ---- Serve our controllers. ----

	// Prepare our repositories and services.
	db, err := datasource.LoadUsers(datasource.Memory)
	if err != nil {
		app.Logger().Fatalf("error while loading the users: %v", err)
		return
	}
	repo := repositories.NewUserRepository(db)
	userService := services.NewUserService(repo)

	// "/users" based mvc application.
	users := mvc.New(app.Party("/users"))
	// Add the basic authentication(admin:password) middleware
	// for the /users based requests.
	users.Router.Use(middleware.BasicAuth)
	// Bind the "userService" to the UserController's Service (interface) field.
	users.Register(userService)
	users.Handle(new(controllers.UsersController))

	// "/user" based mvc application.
	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})
	user := mvc.New(app.Party("/user"))
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserController))


	home := mvc.New(app.Party("/index"))
	home.Handle(new(controllers.HomeController))

	self := mvc.New(app.Party("/self"))
	self.Handle(new(controllers.SelfController))

	shopCar := mvc.New(app.Party("/shopcar"))
	shopCar.Handle(new(controllers.ShopCarController))

	assort := mvc.New(app.Party("/assort"))
	assort.Handle(new(controllers.AssortController))

	order := mvc.New(app.Party("/order"))
	order.Handle(new(controllers.OrderController))

	userg := mvc.New(app.Party("/userg"))
	// Prepare our repositories and services.
	//db1, err := datasource.GetDb(datasource.MySQL)
	//if err != nil {
	//	app.Logger().Fatalf("error while loading the users: %v", err)
	//	return
	//}
	//repo1 := repositories.NewUserRepository(db1)
	userGService := services.NewUserGService(getDb())
	//
	//login.Router.Use(middleware.BasicAuth)
	userg.Register(userGService)
	userg.Handle(new(controllers.UserGController))


	// http://localhost:8080/noexist
	// and all controller's methods like
	// http://localhost:8080/users/1
	// http://localhost:8080/user/register
	// http://localhost:8080/user/login
	// http://localhost:8080/user/me
	// http://localhost:8080/user/logout
	// basic auth: "admin", "password", see "./middleware/basicauth.go" source file.
	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
