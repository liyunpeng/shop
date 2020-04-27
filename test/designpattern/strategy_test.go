package designpattern

import (
	"fmt"
	"testing"
)

// 策略模式里的策略就是指算法

/**
 * 策略接口
 */
type Strategier interface {
	Compute(num1, num2 int) int
}

/**
* 一系列的策略，这些策略不要依赖高层模块的实现, 即import不能用业务的包，也避免了包的循环依赖
 */
type Division struct {}
func (p Division) Compute(num1, num2 int) int {
	defer func() {
		if f := recover(); f != nil {
			logger.Info.Println(f)
			return
		}
	}()

	if num2 == 0 {
		panic("num2 must not be 0!")
	}

	return num1 / num2
}

type Subtraction struct {}
func (p Subtraction) Compute(num1, num2 int) int {
	return num1-num2
}

type Multiplication struct {}
func (p Multiplication) Compute(num1, num2 int) int {
	return num1*num2
}

type Addition struct {}
func (p Addition) Compute(num1, num2 int) int {
	return num1 + num2
}

/**
* 工厂方法，根据不用的type来返回不同的策略
* 这里体现了工厂方法， 产品类不会出现在业务程序里， 增加新产品类，业务程序不用增加代码
 */
func NewStrategy(t string) (res Strategier) {
	switch t {
	case "s": // 减法
		res = Subtraction{}
	case "m": // 乘法
		res = Multiplication{}
	case "d": // 除法
		res = Division{}
	case "a": // 加法
		fallthrough
	default:
		res = Addition{}
	}

	return
}


type Computer struct {
	Num1, Num2 int
	strate Strategier
}

func (p *Computer) SetStrategy(strate Strategier) {
	p.strate = strate
}

func (p Computer) Do() int {
	defer func() {
		if f := recover(); f != nil {
			logger.Info.Println(f)
		}
	}()

	if p.strate == nil {
		panic("Strategier is null")
	}

	return p.strate.Compute(p.Num1, p.Num2)
}

//var stra *string = flag.String("type", "a", "input the strategy")
//var num1 *int = flag.Int("num1", 1, "input num1")
//var num2 *int = flag.Int("num2", 1, "input num2")
//
//func init() {
//	flag.Parse()
//}

func Strategy() {
	// 用户提交的参数
	num1 := 100
	num2 := 200
	stratey := "a"

	// 服务期按照指定策略算法运算
	com := Computer{Num1: num1, Num2: num2}
	strate := NewStrategy(stratey)
	com.SetStrategy(strate)
	logger.Info.Println(com.Do())
}

func TestStrategy(t *testing.T){
	Strategy()
}
