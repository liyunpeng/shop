// file: main.go

package main

import (
	stdContext "context"
	"fmt"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"golang.org/x/net/context"
	"io"
	_ "net/http/pprof"
	"shop/custchan"
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

	"shop/client"
	"shop/config"
	"shop/models"
	"shop/transformer"
	"shop/web/routes"
	"time"

	_ "shop/validates"
)

//var Conf *config.Config


func init() {

	util.Ctx, util.Cancel = context.WithCancel(context.Background())
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
func startService(transformConfiguration *transformer.Conf) {
	go service.StartTailService()

	go service.StartWebSocketService()

	util.WaitGroup.Add(1)
	go service.StartGrpcService(transformConfiguration.Grpc)

	go service.StartMicroService()

	go service.StartOauth2Service()
}

func startClient(transformConfiguration *transformer.Conf) {
	etcdKeys := GetEtcdKeys()
	client.EtcdClientInsance = client.NewEtcdClientWrap(
		[]string{transformConfiguration.Etcd.Addr}, 5 * time.Second)
	go func() {
		fmt.Println("到etcd服务器，按指定的键遍历键值对")
		for _, key := range etcdKeys {
			resp := client.EtcdClientInsance.Get(key)
			if resp != nil || resp.Count < 1 {
				continue
			}
			for _, ev := range resp.Kvs {
				custchan.ConfChan <- string(ev.Value)
				fmt.Printf("etcdkey = %s \n etcdvalue = %s \n", ev.Key, ev.Value)
			}
		}
	}()

	// 启动对etcd的监听服务，有新的键值对会被监听到
	util.WaitGroup.Add(1)

	go client.EtcdClientInsance.EtcdWatch( util.Ctx , etcdKeys)


	util.WaitGroup.Add(1)
	go client.StartKafkaProducer(
		transformConfiguration.Kafka.Addr, 1, true)

	util.WaitGroup.Add(1)
	go client.StartKafkaConsumer(transformConfiguration.Kafka.Addr)

	go client.StartOauth2Client()

	go client.StartGrpcClient()
}


// 按天生成日志文件
func todayFilename() string {
	today := time.Now().Format("20060102")
	return today + ".log"
}

// 创建打开文件
func newLogFile() *os.File {
	filename := todayFilename()
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}


func main() {
	//defer fmt.Println("主routine完全退出")
	//defer fmt.Println("主routine内存分析完毕")
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

	err := config.InitConfig("./config/conf.json")
	if err != nil{
		return
	}
	// 分别打印http db rabbitmq配置
	fmt.Println("host=", config.HttpConfig.Host)
	fmt.Println("port=", config.DBConfig.Port)
	fmt.Println("vhost=", config.AmqpConfig.Vhost)

	f := newLogFile()
	defer f.Close()

	irisConfiguration := iris.TOML("./config/conf.tml")
	config.TransformConfiguration = config.GetTransformConfiguration(irisConfiguration)

	client.InitRedis()

	go func () {
		err := http.ListenAndServe(":9909", nil )
		if err != nil {
			panic(err)
		}
	}()

	app := iris.New()

	util.Logger = app.Logger()
	//outputlog := &client.LoggerOutput{}
	//util.Logger.SetOutput(outputlog)


	//f.Write()

	/*
	输出远程请求的log：
	[INFO] 2020/04/23 22:37 302 18.4357ms ::1 POST /user/login Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.122 Safari/537.36
	*/
	requestLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,
		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},
		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(requestLogger)

	util.Logger.SetOutput(io.MultiWriter(f, os.Stdout ))  //, outputlog))
	util.Logger.SetLevel("debug")


	//util.Logger.SetOutput()
	service.WebsocketChan = make( chan string, 10)

	startService(config.TransformConfiguration)
	models.Register(config.TransformConfiguration)
	models.DB.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
		&models.Order{},
		&models.Goods{},
	)



	startClient(config.TransformConfiguration)
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
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})

	controllers.RegisterControllers(app)
	routes.RegisterApi(app)
	sys.CreateSystemData(app)  //apiRoutes)

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
	util.Logger.Debug("启动iris服务完毕")


	cookieGet(app)

	control(app)


	util.Logger.Info("等待所有routine关闭动作完成")
	util.WaitGroup.Wait()
	util.Logger.Info("所有routine的关闭动作已全部完成")
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

			util.Cancel()
			close(util.ChanStop)
			close(custchan.KafkaProducerMsgChan)

			timeout := 5 * time.Second
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
			cancel()
			util.Logger.Println("关闭iris 服务")
			app.Shutdown(ctx)

			util.Logger.Println("关闭grpc 服务")
			service.GrpcSever.Stop()

			util.Logger.Println("关闭go-micro 微服务")
			service.Stop()

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
