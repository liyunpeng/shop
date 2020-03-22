// file: main.go

package main

import (
	stdContext "context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
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
func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	handler.WebsocketChan = make( chan string, 1000)
	irisConfiguration := iris.TOML("./config/conf.tml")
	transformConfiguration := getTransformConfiguration(irisConfiguration)
	models.Register(transformConfiguration)

	etcdService := services.NewEtcdService(
		[]string{transformConfiguration.Etcd.Addr}, 5 * time.Second)
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

	etcdService.PutKV("/logagent/192.168.0.142/logconfig", `
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
	// 启动对etcd的监听服务，有新的键值对会被监听到
	go etcdService.EtcdWatch(etcdKeys)

	tailService := services.NewTailService()
	go tailService.RunServer()

	go services.StartKafkaProducer(
		transformConfiguration.Kafka.Addr, 1)


	go services.StartKafkaConsumer(transformConfiguration.Kafka.Addr)
	/*
		创建iris应用的
		app.Party得到一个路由对象， party的参数就是一个路径，整个应有都是在这个路径下，
		mvc.new 由这个路由对象， 创建一个mvc的app对象。
		这个app就可以做很多事情，
		注册服务啊，注册控制器

	*/
	etcdApp := mvc.New(app.Party("/etcd"))
	etcdApp.Register(etcdService)
	etcdApp.Handle(new(controllers.EtcdController))

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

	//websocket1(app)

	go func() {
		fmt.Println("启动 websocket 服务")
		http.Handle("/ws", websocket.Handler(handler.Handle))
		err := http.ListenAndServe(":88", nil)
		if err != nil {
			fmt.Println(err)
			fmt.Println("websocket 启动异常")
		}else{
			fmt.Println("websocket 监听服务")
		}
	}()


	go rpc.GrpcServer(transformConfiguration.Grpc)

	//setupWebsocket(app)
	fmt.Println("启动  iris 服务 ")
	go 	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr(":8082"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		//iris.WithOptimizations,
		iris.WithConfiguration(irisConfiguration),
		iris.WithoutInterruptHandler,
	)
	fmt.Println("启动  iris 服务 1 ")


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

					return
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


