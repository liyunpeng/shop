package stub

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/prashantv/gostub"

	. "github.com/smartystreets/goconvey/convey"
)

var counter = 100

func stubGlobalVariable() {
	stubs := gostub.Stub(&counter, 200)
	defer stubs.Reset()
	fmt.Println("Counter:", counter)
}


// Reset方法将全局变量的值恢复为原值。
var configFile = "config.json"
func GetConfig() ([]byte, error) {
	return ioutil.ReadFile(configFile)
}
func stubGlobalVariableA() {
	stubs := gostub.Stub(&configFile, "/tmp/test.config")
	defer stubs.Reset()
	/// 返回tmp/test.config文件的内容
	data, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}


// 为函数打桩
var timeNow = time.Now
var osHostname = os.Hostname

func getDate() int {
	return timeNow().Day()
}
func getHostName() (string, error) {
	return osHostname()
}

func StubTimeNowFunction() {
	stubs := gostub.Stub(&timeNow, func() time.Time {
		return time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC)
	})
	fmt.Println(getDate())
	defer stubs.Reset()
}

func StubHostNameFunction() {
	stubs := gostub.StubFunc(&osHostname, "LocalHost", nil)
	defer stubs.Reset()
	fmt.Println(getHostName())
}


//为过程打桩
//没有返回值的函数称为过程。通常将资源清理类函数定义为过程。
var CleanUp = cleanUp

func cleanUp(val string) {
	fmt.Println(val)
}

func StubProcess(){
	stubs := gostub.StubFunc(&CleanUp)
	CleanUp("Hello go")
	defer stubs.Reset()
}

func TestStub(t *testing.T) {
	stubGlobalVariable()
	stubGlobalVariableA()



	StubTimeNowFunction()
	StubHostNameFunction()


	StubProcess()
}

//var counter = 100
//var CleanUp = cleanUp

//func cleanUp(val string) {
//	fmt.Println(val)
//}
var TimeNow = time.Now
func TestFuncDemo(t *testing.T) {
	Convey("TestFuncDemo", t, func() {
		Convey("for succ", func() {
			stubs := gostub.Stub(&counter, 200)
			defer stubs.Reset()
			stubs.Stub(&TimeNow, func() time.Time {
				return time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC)
			})
			stubs.StubFunc(&CleanUp)
			fmt.Println(counter)
			fmt.Println(TimeNow().Day())
			CleanUp("Hello go")
		})
	})
}