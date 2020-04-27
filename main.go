// file: main.go

package main

import (
	stdContext "context"
	"github.com/kataras/iris/v12"
	"golang.org/x/net/context"
	_ "net/http/pprof"
	"shop/client"
	"shop/custchan"
	"shop/logger"
	"shop/service"
	"shop/sys"
	"shop/web/controllers"

	"os"
	"os/signal"
	"shop/util"
	"syscall"
	//"github.com/kataras/iris/v12/websocket"
	//
	//
	"net/http"

	"shop/config"
	"shop/models"
	"shop/web/routes"
	"time"

	_ "shop/validates"
)

//var Conf *config.Config

func init() {
	logger.InitCustLogger()
	go func() {
		err := http.ListenAndServe(":9909", nil)
		if err != nil {
			panic(err)
		}
	}()

	util.Ctx, util.Cancel = context.WithCancel(context.Background())
	//var _path string
	//
	//flag.StringVar(&_path, "c", "./config.yaml", "default config path")
	//Conf = &config.Config{}
	//
	//content, err := ioutil.ReadFile(_path)®
	//if err == nil {
	//	err = yaml.Unmarshal(content, Conf)
	//	logger.Info.Println("Conf=", Conf)
	//}

	err := config.InitConfig("./config/conf.json")
	if err != nil {
		return
	}
	// 分别打印http db rabbitmq配置
	logger.Info.Println("host=", config.HttpConfig.Host)
	logger.Info.Println("port=", config.DBConfig.Port)
	logger.Info.Println("vhost=", config.AmqpConfig.Vhost)

	config.IrisConfiguration = iris.TOML("./config/conf.tml")
	config.TransformConfiguration = config.GetTransformConfiguration(config.IrisConfiguration)

	//util.Logger.SetOutput()


	//if err := pprof.PPCmd("cpu 10s"); err != nil {
	//	panic(err)
	//}
	//
	//if err := pprof.PPCmd("mem"); err != nil {
	//	panic(err)
	//}
}

// values should match with the client sides as well.
//const enableJWT = true
//const namespace = "default"
//
////if namespace is empty then simply websocket.Events{...} can be used instead.
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
//			nsConn.Conn.StartGrpcService().Broadcast(nsConn, msg)
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

	logger.InitCustLogger()
	// TODO
	//args := os.Args
	//if len(args) < 1 || args == nil {
	//	switch args[1] {
	//	case "gen":
	//
	//		break
	//	}
	//}else{
	//
	//}




	//defer logger.Info.Println("主routine完全退出")
	//defer logger.Info.Println("主routine内存分析完毕")
	//defer profile.StartMicroService(profile.MemProfile).Stop()

	//cpuf, err := os.Create("cpu_profile")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(cpuf)
	//defer pprof.StopCPUProfile()
	//
	//defer func() {
	//	memf, err := os.Create("mem_profile")
	//	if err != nil {
	//		log.Fatal("could not create memory profile: ", err)
	//	}
	//	if err := pprof.WriteHeapProfile(memf); err != nil {
	//		log.Fatal("could not write memory profile: ", err)
	//	}
	//	memf.Close()
	//}()

	app := iris.New()

	util.SetLogger(app)
	//outputlog := &client.LoggerOutput{}
	//util.Logger.SetOutput(outputlog)

	//f.Write()
	service.WebsocketChan = make(chan string, 10)

	service.StartService(config.TransformConfiguration)

	client.StartClient(config.TransformConfiguration)

	models.Register(config.TransformConfiguration)
	models.DB.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
		&models.Order{},
		&models.Goods{},
	)

	tmpl := iris.HTML("./web/views", ".html").
		Layout("shared/layout.html")
	tmpl.Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "./web/public")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	/*
		创建iris应用的
		app.Party得到一个路由对象， party的参数就是一个路径，整个应有都是在这个路径下，
		mvc.new 由这个路由对象， 创建一个mvc的app对象。
		这个app就可以做很多事情，
		注册服务啊，注册控制器
	*/
	routes.RegisterApi(app)
	controllers.RegisterControllers(app)

	sys.CreateSystemData(app) //apiRoutes)

	//websocket1(app)

	//setupWebsocket(app)
	go app.Run(
		// Starts the web server at localhost:8080
		iris.Addr(":8082"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		//iris.WithOptimizations,
		iris.WithConfiguration(config.IrisConfiguration),
		iris.WithoutInterruptHandler,
		iris.WithCharset("UTF-8"),
	)
	util.Logger.Debug("启动iris服务完毕")

	stopControl(app)
	logger.Info.Println("等待所有routine关闭动作完成")
	util.WaitGroup.Wait()
	logger.Info.Println("所有routine的关闭动作已全部完成")
}

func stopControl(app *iris.Application) {
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
		case <-signals:
			logger.Info.Println("shutdown...")

			util.Cancel()
			close(util.ChanStop)
			close(custchan.KafkaProducerMsgChan)

			timeout := 5 * time.Second
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
			cancel()
			logger.Info.Println("关闭iris 服务")
			app.Shutdown(ctx)

			logger.Info.Println("关闭grpc 服务")
			service.GrpcSever.Stop()

			logger.Info.Println("关闭go-micro 微服务")
			service.Stop()

			logger.Info.Println("关闭控制流程结束")
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
