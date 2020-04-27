package designpattern

import (
	"fmt"
	"testing"
)

type Base struct {
	FirstName, LastName string
	Age float32
}

func (base *Base) HasFeet() {
	fmt.Println(base.FirstName + base.LastName + "has feet! Base")
}

func (base *Base) Flying() {
	fmt.Println("Base Can flying!")
}

type Sub struct {
	Base
	Area string
}

func (sub *Sub) Flying() {
	/*
	如果“子类”重写了“基类”的成员方法，
	需要在子类的成员方法中调用基类的同名成员方法，
	一定要以sub.Base.Flying()这样显式的方法调用，
	而不是使用sub.Flying()这种调用继承方法的方式调用，
	这样会出现无限循环，即一直在调用子类的方法。
	 */
	sub.Base.Flying()
	fmt.Println("Sub flying")
}

func TestZuhe( t *testing.T) {
	chk := new(Sub)
	chk.Flying()
	chk2 := &Sub{Base{"Bob", "Steven", 2.0}, "China"}
	fmt.Println(chk2.Area)
}

