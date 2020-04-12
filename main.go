// file: main.go

package main

import (
	stdContext "context"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"shop/srv"

	//"github.com/kataras/neffos/stackexchange/redis"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"os"
	"os/signal"
	"shop/rpc"
	"shop/util"
	"syscall"
	//"github.com/kataras/iris/v12/websocket"
	//
	//
	"net/http"
	"shop/handler"

	gf "github.com/snowlyg/gotransformer"
	"golang.org/x/net/websocket"
	"shop/config"
	"shop/models"
	"shop/services"
	"shop/transformer"
	"shop/web/routes"
	"time"

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
	//content, err := ioutil.ReadFile(_path)®
	//if err == nil {
	//	err = yaml.Unmarshal(content, Conf)
	//	fmt.Println("Conf=", Conf)
	//}
}

func getTransformConfiguration( irisConfiguration iris.Configuration) *transformer.Conf {
	app := transformer.App{}
	g := gf.NewTransform(&app, irisConfiguration.Other["App"], time.RFC3339)
	_ = g.Transformer()

	db := transformer.Mysql{}
	g.OutputObj = &db
	g.InsertObj = irisConfiguration.Other["Mysql"]
	_ = g.Transformer()

	mongodb := transformer.Mongodb{}
	g.OutputObj = &mongodb
	g.InsertObj = irisConfiguration.Other["Mongodb"]
	_ = g.Transformer()

	redis := transformer.Redis{}
	g.OutputObj = &redis
	g.InsertObj = irisConfiguration.Other["Redis"]
	_ = g.Transformer()

	sqlite := transformer.Sqlite{}
	g.OutputObj = &sqlite
	g.InsertObj = irisConfiguration.Other["Sqlite"]
	_ = g.Transformer()

	testData := transformer.TestData{}
	g.OutputObj = &testData
	g.InsertObj = irisConfiguration.Other["TestData"]
	_ = g.Transformer()

	kafkaConf := transformer.Kafka{}
	g.OutputObj = &kafkaConf
	g.InsertObj = irisConfiguration.Other["Kafka"]
	_ = g.Transformer()

	etcdConf := transformer.EtcdConf{}
	g.OutputObj = &etcdConf
	g.InsertObj = irisConfiguration.Other["Etcd"]
	_ = g.Transformer()

	grpcConf := transformer.GrpcConf{}
	g.OutputObj = &grpcConf
	g.InsertObj = irisConfiguration.Other["Grpc"]
	_ = g.Transformer()

	consulConf := transformer.ConsulConf{}
	g.OutputObj = &consulConf
	g.InsertObj = irisConfiguration.Other["Consul"]
	_ = g.Transformer()

	cf := &transformer.Conf{
		App:      app,
		Mysql:    db,
		Mongodb:  mongodb,
		Redis:    redis,
		Sqlite:   sqlite,
		TestData: testData,
		Kafka: kafkaConf,
		Etcd: etcdConf,
		Grpc: grpcConf,
		Consul: consulConf,
	}

	return cf
}

func GetEtcdKeys() ([]string) {
	var etcdKeys []string
	//ips, err := getLocalIP()
	var ips []string
	//var err error
	ips = append(ips, "192.168.0.1")
	//if err != nil {
	//	fmt.Println("get local ip error:", err)
	//	//return err
	//}
	for _, ip := range ips {
		//key := fmt.Sprintf("/logagent/%s/logconfig", ip)
		etcdKeys = append(etcdKeys, ip)
	}
	fmt.Println("从etcd服务器获取到的以IP名为键的键值对: ", etcdKeys)
	return etcdKeys
}


// values should match with the client sides as well.
const enableJWT = true
const namespace = "default"

// if namespace is empty then simply websocket.Events{...} can be used instead.
//var serverEvents = websocket.Namespaces{
//	namespace: websocket.Events{
//		websocket.OnNamespaceConnected: func(nsConn *websocket.NSConn, msg websocket.Message) error {
//			// with `websocket.GetContext` you can retrieve the Iris' `Context`.
//			ctx := websocket.GetContext(nsConn.Conn)
//
//			log.Printf("[%s] connected to namespace [%s] with IP [%s]",
//				nsConn, msg.Namespace,
//				ctx.RemoteAddr())
//			return nil
//		},
//		websocket.OnNamespaceDisconnect: func(nsConn *websocket.NSConn, msg websocket.Message) error {
//			log.Printf("[%s] disconnected from namespace [%s]", nsConn, msg.Namespace)
//			return nil
//		},
//		"chat": func(nsConn *websocket.NSConn, msg websocket.Message) error {
//			// room.String() returns -> NSConn.String() returns -> Conn.String() returns -> Conn.ID()
//			log.Printf("[%s] sent: %s", nsConn, string(msg.Body))
//
//			// Write message back to the client message owner with:
//			// nsConn.Emit("chat", msg)
//			// Write message to all except this client with:
//			nsConn.Conn.GrpcServer().Broadcast(nsConn, msg)
//			return nil
//		},
//	},
//}
//func websocket1(app *iris.Application) {
//	websocketServer := websocket.New(
//		websocket.DefaultGorillaUpgrader, /* DefaultGobwasUpgrader can be used too. */
//		serverEvents)
//
//	//j := jwt.New(jwt.Config{
//	//	// Extract by the "token" url,
//	//	// so the client should dial with ws://localhost:8080/echo?token=$token
//	//	Extractor: jwt.FromParameter("token"),
//	//
//	//	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
//	//		return []byte("My Secret"), nil
//	//	},
//	//
//	//	// When set, the middleware verifies that tokens are signed
//	//	// with the specific signing algorithm
//	//	// If the signing method is not constant the
//	//	// `Config.ValidationKeyGetter` callback field can be used
//	//	// to implement additional checks
//	//	// Important to avoid security issues described here:
//	//	// https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
//	//	SigningMethod: jwt.SigningMethodHS256,
//	//})
//
//	idGen := func(ctx iris.Context) string {
//		if username := ctx.GetHeader("X-Username"); username != "" {
//			return username
//		}
//
//		return websocket.DefaultIDGenerator(ctx)
//	}
//
//	// serves the endpoint of ws://localhost:8080/echo
//	// with optional custom ID generator.
//	//websocketRoute := app.Get("/echo", websocket.Handler(websocketServer, idGen))
//	app.Get("/echo", websocket.Handler(websocketServer, idGen))
//
//}

func startService(transformConfiguration *transformer.Conf) {
	etcdKeys := GetEtcdKeys()
	services.EtcdServiceInsance = services.NewEtcdService(
		[]string{transformConfiguration.Etcd.Addr}, 5 * time.Second)
	go func() {
		fmt.Println("到etcd服务器，按指定的键遍历键值对")
		for _, key := range etcdKeys {
			resp := services.EtcdServiceInsance.Get(key)
			if resp != nil || resp.Count < 1 {
				continue
			}
			for _, ev := range resp.Kvs {
				services.ConfChan <- string(ev.Value)
				fmt.Printf("etcdkey = %s \n etcdvalue = %s \n", ev.Key, ev.Value)
			}
		}
	}()

	// 启动对etcd的监听服务，有新的键值对会被监听到
	util.WaitGroup.Add(1)
	go services.EtcdServiceInsance.EtcdWatch(etcdKeys)

	tailService := services.NewTailService()
	go tailService.RunServer()

	util.WaitGroup.Add(1)
	go services.StartKafkaProducer(
		transformConfiguration.Kafka.Addr, 1)

	util.WaitGroup.Add(1)
	go services.StartKafkaConsumer(transformConfiguration.Kafka.Addr)

	//util.WaitGroup.Add(1)
	go func() {
		//defer 0.198()
		fmt.Println("启动 websocket 服务")
		http.Handle("/ws", websocket.Handler(handler.WebSocketHandle))
		err := http.ListenAndServe(":88", nil)
		if err != nil {
			fmt.Println(err)
			fmt.Println("websocket 启动异常")
		}else{
			fmt.Println("websocket 监听服务")
		}
	}()

	util.WaitGroup.Add(1)
	go rpc.GrpcServer(transformConfiguration.Grpc)

}

func createTestData(transformConfiguration *transformer.Conf){

	services.EtcdServiceInsance.PutKV("192.168.0.1", `
[
	{
		"topic":"nginx_log",
		"log_path":"/Users/admin1/goworkspace/shop/log1.txt",
		"service":"test_service",
		"send_rate":1000
	},
		
	{
		"topic":"nginx_log1",
		"log_path":"/Users/admin1/goworkspace/shop/log2.txt",
		"service":"test_service1",
		"send_rate":1000
	},

	{
		"topic":"nginx_log",
		"log_path":"D:\\goworkspace\\shop\\log1.txt",
		"service":"test_service1",
		"send_rate":1000
	}
]` )
	services.EtcdServiceInsance.PutKV("192.168.0.2", `
[
	{
		"topic":"nginx_log",
		"log_path":"/Users/admin1/goworkspace/shop/123.txt",
		"service":"test_service",
		"send_rate":2000
	}
]` )
}


func registerControllers( app *iris.Application) {
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)

	sessManager := sessions.New(sessions.Config{
		Cookie:  util.COOKEI_NAME,
		Expires: 24 * time.Hour,
		Encode:       secureCookie.Encode,
		Decode:       secureCookie.Decode,
		AllowReclaim: true,
	})
	db := redis.New(redis.Config{
		Network:   "tcp",
		Addr:      "127.0.0.1:6379",
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Password:  "",
		Database:  "",
		Prefix:    "",
		Delim:     "-",
		Driver:    redis.Redigo(), // redis.Radix() can be used instead.
	})
	sessManager.UseDatabase(db)

	etcdApp := mvc.New(app.Party("/etcd"))
	etcdApp.Register(services.EtcdServiceInsance)
	etcdApp.Handle(new(controllers.EtcdController))

	home := mvc.New(app.Party("/"))
	home.Register(
		sessManager.Start,
	)
	home.Handle(new(controllers.HomeController))

	self := mvc.New(app.Party("/self"))
	self.Register(
		sessManager.Start,
	)
	self.Handle(new(controllers.SelfController))

	shopCar := mvc.New(app.Party("/shopcar"))
	shopCar.Handle(new(controllers.ShopCarController))

	goodsDetail := mvc.New(app.Party("/goodsdetail"))
	goodsDetail.Handle(new(controllers.GoodsDetailController))

	assort := mvc.New(app.Party("/assort"))
	assort.Handle(new(controllers.AssortController))

	buy := mvc.New(app.Party("/goodsdetail/buy"))
	buy.Register(
		sessManager.Start,
	)
	buy.Handle(new(controllers.BuyController))

	order := mvc.New(app.Party("/order"))
	order.Register(
		sessManager.Start,
	)
	order.Handle(new(controllers.OrderController))

	user := mvc.New(app.Party("/user"))
	userService := services.NewUserService()
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(controllers.UserGController))
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")


	handler.WebsocketChan = make( chan string, 10)

	irisConfiguration := iris.TOML("./config/conf.tml")
	config.TransformConfiguration = getTransformConfiguration(irisConfiguration)

	models.Register(config.TransformConfiguration)
	models.DB.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
		&models.Order{},
		&models.Goods{},
	)

	startService(config.TransformConfiguration)

	/*
		创建iris应用的
		app.Party得到一个路由对象， party的参数就是一个路径，整个应有都是在这个路径下，
		mvc.new 由这个路由对象， 创建一个mvc的app对象。
		这个app就可以做很多事情，
		注册服务啊，注册控制器
	*/
	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "./web/public")
	registerControllers(app)
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	routes.RegisterApi(app)


	apiRoutes := routes.GetRoutes(app)
	models.CreateSystemData(apiRoutes)

	createTestData(config.TransformConfiguration)
	//websocket1(app)

	//setupWebsocket(app)
	go 	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr(":8082"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		//iris.WithOptimizations,
		iris.WithConfiguration(irisConfiguration),
		iris.WithoutInterruptHandler,
		iris.WithCharset("UTF-8"),
	)
	fmt.Println("启动iris服务完毕")


	cookieGet(app)

	go srv.Start()
	control(app)

	fmt.Println("等待所有routine关闭动作完成")
	util.WaitGroup.Wait()
	fmt.Println("所有routine的关闭动作已全部完成")


}

func control (app *iris.Application) {
	signals := make(chan os.Signal, 1)
	//signal.Notify(signals, os.Interrupt)

	signal.Notify(signals,
		// kill -SIGINT XXXX 或 Ctrl+c
		os.Interrupt,
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		os.Kill,
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)

Loopa:
	//go func(){
	for {
		select {
		case <- signals:
			println("shutdown...")

			close(util.ChanStop)
			close(services.KafkaProducerObj.MsgChan)

			timeout := 5 * time.Second
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
			cancel()
			fmt.Println("关闭iris 服务")
			app.Shutdown(ctx)

			fmt.Println("关闭grpc 服务")
			rpc.GrpcSever.Stop()

			fmt.Println("关闭go-micro 微服务")
			srv.Stop()

			break Loopa

		}
	}
	//}()
}
//func setupWebsocket(app *iris.Application) {
//	// create our echo websocket server
//	ws := websocket.New(websocket.Config{
//		ReadBufferSize:  1024,
//		WriteBufferSize: 1024,
//	})
//	ws.OnConnection(handleConnection)
//
//	// register the server on an endpoint.
//	// see the inline javascript code in the websockets.html,
//	// this endpoint is used to connect to the server.
//	app.Get("/echo", ws.Handler())
//	// serve the javascript built'n client-side library,
//	// see websockets.html script tags, this path is used.
//	app.Any("/iris-ws.js", websocket.ClientHandler())
//}
//
//func handleConnection(c websocket.Connection) {
//	// Read events from browser
//	c.On("chat", func(msg string) {
//		// Print the message to the console, c.Context() is the iris's http context.
//		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
//		// Write message back to the client message owner with:
//		// c.Emit("chat", msg)
//		// Write message to all except this client with:
//		c.To(websocket.Broadcast).Emit("chat", msg)
//	})
//}


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
