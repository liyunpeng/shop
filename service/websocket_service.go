package service

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func StartWebSocketService() {
	fmt.Println("启动 websocket 服务")
	http.Handle("/ws", websocket.Handler(WebSocketHandle))
	err := http.ListenAndServe(":88", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("websocket 启动异常")
	}else{
		fmt.Println("websocket 监听服务")
	}
}
