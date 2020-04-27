package designpattern

import (
	"fmt"
	"testing"
)

type IHello interface {
	Hello(name string)
}

type A struct {
}
func (*A) Hello(name string) {
	logger.Info.Println("hello " + name + ", i am a")
}

type B struct {
	IHello
}
func (*B) Hello(name string) {
	logger.Info.Println("hello " + name + ", i am b")
}

type C struct {
	IHello
}


type D struct {
}
func (*D) Hello(name string) {
	logger.Info.Println("hello " + name + ", i am d")
}

func TestZuheInteface(t *testing.T) {
	name := "Lee"
	a := A{}
	a.Hello(name) //hello Lee, i am a

	b := B{&A{}}
	b.Hello(name) //hello Lee, i am b

	/*
	B中写了一个与IHello同名的方法Hello，
	此时直接访问b.Hello是访问的b的方法，
	想访问A的方法需要b.IHello.Hello(name)。
	我们可以把组合方式直接访问被嵌入类型方法(这里为A的方法)看做一个语法糖
	 */
	b.IHello.Hello(name) //hello Lee, i am a

	c := C{&A{}}
	c.Hello(name) //hello Lee, i am a

	/*
	根据运行时上下文指定其他具体实现，比如D，更加灵活
	 */
	c.IHello = &D{}
	c.Hello(name) //hello Lee, i am d
}

