// file: main.go

package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/web/routes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"

	"shop/config"
	"shop/models"
	"shop/services"
	"shop/web/controllers"
)

var Conf *config.Config

func init() {
	var _path string

	flag.StringVar(&_path, "c", "./config.yaml", "default config path")
	Conf = &config.Config{}

	content, err := ioutil.ReadFile(_path)
	if err == nil {
		err = yaml.Unmarshal(content, Conf)
		fmt.Println("Conf=", Conf)
	}
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	models.Register(Conf)

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

	sessManager := sessions.New(sessions.Config{
		Cookie:  "sessioncookiename",
		Expires: 24 * time.Hour,
	})

	index := mvc.New(app.Party("/index"))
	index.Handle(new(controllers.IndexController))

	self := mvc.New(app.Party("/self"))
	self.Register(
		sessManager.Start,
	)

	self.Handle(new(controllers.SelfController))

	shopCar := mvc.New(app.Party("/shopcar"))
	shopCar.Handle(new(controllers.ShopCarController))

	assort := mvc.New(app.Party("/assort"))
	assort.Handle(new(controllers.AssortController))

	order := mvc.New(app.Party("/order"))
	order.Handle(new(controllers.OrderController))

	user := mvc.New(app.Party("/user"))
	userService := services.NewUserService()
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserGController))

	routes.RegisterApi(app)

	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr(":8082"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
