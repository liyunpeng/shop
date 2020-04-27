package service

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"shop/logger"
	"shop/util"
	//"shop/client"
)

type UserInfo struct {
	UserId int32
	Name   string
	Age    string
}

type Response struct {
	Code int32
	Msg  string
	Data interface{}
}

var WebsocketChan chan string

func WebSocketHandle(conn *websocket.Conn) {
	defer conn.Close()
	util.WaitGroup.Add(1)

	logger.Info.Println("websocket与客户端建立连接，启动接收和发送数据的服务")
	go sendToClient(conn)
	for {
		jsonHandler := websocket.JSON
		userInfo := &UserInfo{}
		err := jsonHandler.Receive(conn, userInfo)
		if err != nil {
			logger.Info.Println(err)
			break
		}
		jsonData, _ := json.Marshal(userInfo)

		logger.Info.Println("receive frontend data:", string(jsonData[:]))
		res := &Response{
			Code: 1,
			Msg:  "success",
			Data: "from web socket ",
		}
		err = jsonHandler.Send(conn, res)
		if err != nil {
			logger.Info.Println(err)
			break
		}

	}
}

func sendToClient(conn *websocket.Conn) {
	defer util.WaitGroup.Done()
	defer util.PrintFuncName()
	jsonHandler := websocket.JSON

	for {
		select {
		case msg := <-WebsocketChan:
			logger.Info.Println("向前端发送数据=", msg)
			res := &Response{
				Code: 1,
				Msg:  "success",
				Data: msg,
			}
			err := jsonHandler.Send(conn, res)
			if err != nil {
				logger.Info.Println(err)

				//go StartWebSocketService()
				//panic("websockt 向客户端发送数据错误")
				break
			}
			//time.Sleep(time.Millisecond * 500)
		case <-util.ChanStop:
			logger.Info.Println("websocket执行发送的routine结束")
			return
		}
	}
}
