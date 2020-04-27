package service

import (
	"golang.org/x/net/websocket"
	"shop/logger"
	"net/http"
)

func StartWebSocketService() {
	logger.Info.Println("启动 websocket 服务")
	http.Handle("/ws", websocket.Handler(WebSocketHandle))
	err := http.ListenAndServe(":88", nil)
	if err != nil {
		logger.Info.Println(err)
		logger.Info.Println("websocket 启动异常")
	} else {
		logger.Info.Println("websocket 监听服务")
	}
}
