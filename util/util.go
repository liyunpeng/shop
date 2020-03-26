package util

import (
	"fmt"
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
