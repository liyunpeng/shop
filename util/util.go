package util

import (
	"fmt"
	"github.com/kataras/iris/v12/sessions"
	"runtime"
	"sync"
)

var ChanStop = make( chan bool)
var WaitGroup sync.WaitGroup

var  COOKEI_NAME = "sessioncookiename"
const SessionUserName = "serssion_user_name"
func PrintFuncName() {
	funcName,file,line,ok := runtime.Caller(1)
	if(ok){
		fmt.Println("func name: " + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s, line: %d\n",file,line)
	}
}



func GetCurrentUserID( session *sessions.Session) int64 {
	userID := session.GetInt64Default(SessionUserName, 0)
	return userID
}

func  IsLoggedIn( session *sessions.Session) bool {
	return GetCurrentUserID(session) > 0
}

func Logout(session *sessions.Session) {
	session.Destroy()
}

