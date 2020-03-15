package handler

import (
	//"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
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
	jsonHandler := websocket.JSON
	userInfo := &UserInfo{}
	//res := &Response{
	//	Code: 1,
	//	Msg:  "success",
	//}

	go Push(conn)

	for {
		err := jsonHandler.Receive(conn, userInfo)
		if err != nil {
			fmt.Println(err)
			break
		}
		//jsonData, _ := json.Marshal(userInfo)
		//fmt.Println("receive data:", string(jsonData[:]))
		//err = jsonHandler.Send(conn, res)
		//if err != nil {
		//	fmt.Println(err)
		//	break
		//}


	}
}

func Push(conn *websocket.Conn) {
	jsonHandler := websocket.JSON
	//res := &Response{
	//	Code: 1,
	//	Msg:  "success",
	//	Data: "hello client",
	//}
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
			jsonHandler.Send(conn, msg)
		}
	}
}
