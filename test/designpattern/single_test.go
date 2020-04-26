package designpattern

import (
	"sync"
	"testing"
)

var m *Single
var once sync.Once

type Single struct{
	Name string
}
func GetInstance()*Single{
	once.Do(func() {
		m = &Single{}
	})
	return m
}

func TestSingle(t *testing.T){
	single := GetInstance()
	println(single)
}
