// file: main.go

package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/web/routes"
	"time"

	gf "github.com/snowlyg/gotransformer"
	"shop/config"
	"shop/models"
	"shop/services"
	"shop/transformer"
	_ "shop/validates"
	"shop/web/controllers"
)

var Conf *config.Config

func init() {
	//var _path string
	//
	//flag.StringVar(&_path, "c", "./config.yaml", "default config path")
	//Conf = &config.Config{}
	//
	//content, err := ioutil.ReadFile(_path)
	//if err == nil {
	//	err = yaml.Unmarshal(content, Conf)
	//	fmt.Println("Conf=", Conf)
	//}
}

var Sc iris.Configuration

func getSysConf() *transformer.Conf {
	app := transformer.App{}
	g := gf.NewTransform(&app, Sc.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Mysql{}
	g.OutputObj = &db
	g.InsertObj = Sc.Other["Mysql"]
	_ = g.Transformer()

	mongodb := transformer.Mongodb{}
	g.OutputObj = &mongodb
	g.InsertObj = Sc.Other["Mongodb"]
	_ = g.Transformer()

	redis := transformer.Redis{}
	g.OutputObj = &redis
	g.InsertObj = Sc.Other["Redis"]
	_ = g.Transformer()

	sqlite := transformer.Sqlite{}
	g.OutputObj = &sqlite
	g.InsertObj = Sc.Other["Sqlite"]
	_ = g.Transformer()

	testData := transformer.TestData{}
	g.OutputObj = &testData
	g.InsertObj = Sc.Other["TestData"]
	_ = g.Transformer()

	cf := &transformer.Conf{
		App:      app,
		Mysql:    db,
		Mongodb:  mongodb,
		Redis:    redis,
		Sqlite:   sqlite,
		TestData: testData,
	}

	return cf
}

func GetEtcdKeys() ([]string) {
	var etcdKeys []string
	//ips, err := getLocalIP()
	var ips []string
	//var err error
	ips = append(ips, "192.168.0.142")
	//if err != nil {
	//	fmt.Println("get local ip error:", err)
	//	//return err
	//}
	for _, ip := range ips {
		key := fmt.Sprintf("/logagent/%s/logconfig", ip)
		etcdKeys = append(etcdKeys, key)
	}
	fmt.Println("从etcd服务器获取到的以IP名为键的键值对: \n", etcdKeys)
	return etcdKeys
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	Sc = iris.TOML("./config/conf.tml")
	rc := getSysConf()
	models.Register(rc)


	etcdService := services.NewEtcdService(
		[]string{"192.168.0.141:2379"}, 5 * time.Second)
	//[]string{"127.0.0.1:2379"}, 5 * time.Second)

	etcdKeys := GetEtcdKeys()
	fmt.Println("到etcd服务器，按指定的键遍历键值对")
	for _, key := range etcdKeys {
		resp := etcdService.Get(key)
		for _, ev := range resp.Kvs {
			services.ConfChan <- string(ev.Value)
			fmt.Printf("etcdkey = %s \n etcdvalue = %s \n", ev.Key, ev.Value)
		}
	}

	// 启动对etcd的监听服务，有新的键值对会被监听到
	go etcdService.EtcdWatch(etcdKeys)


	tailService := services.NewTailService()
	go tailService.RunServer()

	services.NewKafkaService(
		rc.Kafka.Addr, 3)

	/*
		创建iris应用的
		app.Party得到一个路由对象， party的参数就是一个路径，整个应有都是在这个路径下，
		mvc.new 由这个路由对象， 创建一个mvc的app对象。
		这个app就可以做很多事情，
		注册服务啊，注册控制器

	*/
	etcdManagerApp := mvc.New(app.Party("/etcdmanager"))
	etcdManagerApp.Register(etcdService)
	etcdManagerApp.Handle(new(controllers.EtcdManangerController))



	models.DB.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

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
	apiRoutes := routes.GetRoutes(app)
	models.CreateSystemData(apiRoutes)

	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr(":8082"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		//iris.WithOptimizations,
		iris.WithConfiguration(Sc),
	)
}
