package designpattern

import (
	"fmt"
	"testing"
)

// Bridge模式 属于单一职责，主要解决在继承的情况下，子类急剧膨胀，难以维护。
// 而Bridge模式，改继承为组合
// 在现实生活中，某些类具有两个或多个维度的变化，如图形既可按形状分，又可按颜色分。
// 如何设计类似于 Photoshop 这样的软件，能画不同形状和不同颜色的图形呢？
// 如果用继承方式，m 种形状和 n 种颜色的图形就有 m×n 种，不但对应的子类很多，
// 而且扩展困难。

// 包是是一个抽象化的角色
type Bag interface {
	GetName()
}

type BridgeBag struct {
	color Color
}

// 扩展抽象化角色 手包
type HandBag struct {
	//BridgeBag
	color Color
	Name string
}

func (h *HandBag) GetName(){
	logger.Info.Printf("包的名字是%s,颜色是%s \n", h.Name,h.color)
}

// 扩展抽象化角色 钱包
type Wallet struct {
	//BridgeBag
	color Color
	Name string
}

func (w Wallet) GetName(){
	logger.Info.Printf("包的名字是%s,颜色是%s \n", w.Name,w.color)
}

type Color interface {
	GetColor()
}

// 具体实现的角色
type Blue struct {
	Name string
}

func (b *Blue) GetColor(){
	logger.Info.Printf("颜色是%s \n", b.Name)
}

type Red struct {
	Name string
}

func (r *Red) GetColor(){
	logger.Info.Printf("颜色是%s \n", r.Name)
}


func BridgePattern(){
	// 红色
	r := Red{Name:"红色"}
	// 蓝色
	b := Blue{Name:"蓝色"}
	// 红色的手包
	h := HandBag{Name:"手包"}
	h.color = &r
	h.GetName()

	// 蓝色手包
	h.color = &b
	h.GetName()

	w := Wallet{Name:"钱包"}
	// 红色钱包
	w.color = &r
	w.GetName()

	// 蓝色钱包
	w.color = &b
	w.GetName()
}

func TestBridgePattern( t *testing.T)  {
	BridgePattern()
}
