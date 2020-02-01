// file: main.go

package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"

	"github.com/liyunpeng/shop/datasource"
	"github.com/liyunpeng/shop/repositories"
	"github.com/liyunpeng/shop/services"
	"github.com/liyunpeng/shop/web/controllers"
	"github.com/liyunpeng/shop/web/middleware"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

var Conf *Config
var _path string

func init() {
	flag.StringVar(&_path, "c", "./config.yaml", "default config path")
	InitConfig()
}

// 从配置文件中加载配置
func InitConfig() error {
	Conf = &Config{}

	content, err := ioutil.ReadFile(_path)
	if err == nil {
		err = yaml.Unmarshal(content, Conf)
		fmt.Println("Conf=", Conf)
	}
	return err
}

// 总的配置
type Config struct {
	Server ServerConf   `yaml:"server"`
	Mysql  MysqlConf `yaml:"mysql"`
	Redis  RedisConf `yaml:"redis"`
}

// 服务的配置
type ServerConf struct {
	Port int `yaml:"port"`
	List []string `yaml:"list,flow"`
}

type RedisConf struct {
	Enable bool `yaml:"enable"`
}

type MysqlConf struct {
	MaxIdle int `yaml:"maxidle"`
	MaxOpen int `yaml:"maxopen"`
	User string `yaml:"user"`
	Host string `yaml:"host"`
	Password string `yaml:"password"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}


func getDb() (db *gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()

	conf := Conf.Mysql

	//TODO？ 全局db没用
	/*
		链接localhost数据库， 用户名root, 密码root
	*/
	db, err := gorm.Open(
		"mysql", "root:root@/gotest?charset=utf8&parseTime=True&loc=Local")
	//"mysql", "root:password@tcp(192.168.0.220:31111)/gotest?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		db.DB().SetMaxIdleConns(conf.MaxIdle)
		db.DB().SetMaxOpenConns(conf.MaxOpen)
		db.DB().SetConnMaxLifetime(time.Duration(30) * time.Minute)
		err = db.DB().Ping()
		fmt.Println("成功连接数据库 db=", db, "err=", err)
	} else {
		fmt.Println("没有连接到数据库 err= ", err)
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
		//sessManager.Start,
	)
	user.Handle(new(controllers.UserController))

	home := mvc.New(app.Party("/index"))
	home.Handle(new(controllers.HomeController))

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

	userg := mvc.New(app.Party("/userg"))
	userGService := services.NewUserGService(getDb())
	userg.Register(
		userGService,
		sessManager.Start,
	)
	userg.Handle(new(controllers.UserGController))

	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr(":8080"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}
