package util

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"io"
	"os"
	"time"
)

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

func SetLogger(app *iris.Application) {
	Logger = app.Logger()

	f := newLogFile()
	//defer f.Close()
	Logger.SetOutput(io.MultiWriter(f, os.Stdout ))  //, outputlog))
	Logger.SetLevel("debug")

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
}
