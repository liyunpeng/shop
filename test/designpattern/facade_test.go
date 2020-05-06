package designpattern

import (
	"fmt"
	"testing"
)

// 虽然把蔬菜定义在厨房的采购部并不合理
// 但是我还是这么做了
const (
	STATUS_BUYED  = iota
	STATUS_STORED
	STATUS_CUTED
	STATUS_COOKED
	STATUS_EATED
)

type VegStatus int

type Vegetable struct {
	name   string
	status VegStatus
}

func BuyVegetable(name string) *Vegetable {
	return &Vegetable{name, STATUS_BUYED}
}

func Eat(veg *Vegetable) {
	veg.status = STATUS_EATED
}

var storage []*Vegetable = make([]*Vegetable, 0)

func SaveVegetables(veg ...*Vegetable) {
	storage = append(storage,veg...)
	for _,v := range veg{
		v.status=STATUS_STORED
	}
}

func GetVegetables() *Vegetable {
	l := len(storage)
	res := storage[l-1]
	storage = storage[:l-1]
	return res
}


//  切菜
//  具体拿去做什么不知道
//  只管切
func CutVegtable(veg ...*Vegetable)[]*Vegetable{
	for _,v := range veg{
		v.status=STATUS_CUTED
	}
	return veg
}

//  炒菜
//  给谁吃不知道
//  炒就是了
func CookVegtable(vec ...*Vegetable)[]*Vegetable{
	for _,v := range vec{
		v.status=STATUS_COOKED
	}
	return vec
}

//  菜单
//  客户随意选择
//  制作方式可以由主厨决定
//  只要客户喜欢
func SauteVegtable()[]*Vegetable{
	qc := BuyVegetable("青菜")
	suan := BuyVegetable("蒜")
	jiang := BuyVegetable("姜")
	SaveVegetables(qc,suan,jiang)

	vegs := CookVegtable(CutVegtable(storage...)...)
	return vegs
}

func EatVegtables(veg ...*Vegetable){
	for _,v :=range veg{
		Eat(v)
	}
}

func TestFacade(t *testing.T){

	// No Facade
	bc := BuyVegetable("白菜")
	SaveVegetables(bc)
	vecs := CookVegtable(GetVegetables())
	EatVegtables(vecs...)
	for _,v := range vecs{
		fmt.Println(*v)
	}

	// Favade
	sauteVegtable := SauteVegtable()
	EatVegtables(sauteVegtable...)
	for _,v := range sauteVegtable{
		fmt.Println(*v)
	}
}