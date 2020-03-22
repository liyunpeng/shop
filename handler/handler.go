package handler

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"shop/util"

	//"shop/services"
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
var WebsocketChan chan  string
func Handle(conn *websocket.Conn) {
	defer conn.Close()

	go Push(conn)

	for {
		jsonHandler := websocket.JSON
		userInfo := &UserInfo{}
		err := jsonHandler.Receive(conn, userInfo)
		if err != nil {
			fmt.Println(err)
			break
		}
		jsonData, _ := json.Marshal(userInfo)

		fmt.Println("receive frontend data:", string(jsonData[:]))
		res := &Response{
			Code: 1,
			Msg:  "success",
			Data:  "from web socket ",
		}
		err = jsonHandler.Send(conn, res)
		if err != nil {
			fmt.Println(err)
			break
		}


	}
}

func Push(conn *websocket.Conn) {
	jsonHandler := websocket.JSON

	for {
		//err := jsonHandler.Send(conn, res)
		//if err != nil {
		//	fmt.Println(err)
		//	break
		//}
		//time.Sleep(time.Millisecond * 500)
		select {
		case msg := <- WebsocketChan:
			fmt.Println("向前端发送数据=",msg)
			res := &Response{
				Code: 1,
				Msg:  "success",
				Data: msg,
			}
			jsonHandler.Send(conn, res)
		case <- util.ChanStop:
			fmt.Println("websocket执行发送的routine结束")
			return
		}
	}
}
