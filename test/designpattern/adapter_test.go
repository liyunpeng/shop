package designpattern

import (
	"shop/logger"
	"testing"
)

type Target interface {
	request()
}
type Adaptee struct {

}
func(it *Adaptee)specificrequeset(){
	logger.Info.Println("asdf")
}

type Adapter struct {
	adaptee *Adaptee
}
func(it *Adapter)setAdaptee(adaptee *Adaptee){
	it.adaptee = adaptee
}
func(it *Adapter)request(){
	it.adaptee.specificrequeset()
}

func TestAdapter(t *testing.T){
	logger.InitCustLogger()
	target := new(Adapter)
	adaptee := new(Adaptee)
	target.setAdaptee(adaptee)
	target.request()
}
