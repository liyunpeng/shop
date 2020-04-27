package designpattern

import (
	"fmt"
	"testing"
	"time"
)

type Option func(options *Config)

type Config struct {
	Name    string
	Timeout time.Duration
}

type Services struct {
	conf Config
}

func SetTimeout(t time.Duration) Option {
	return func(options *Config) {
		options.Timeout = t
	}
}

func SetName(name string) Option {
	return func(options *Config) {
		options.Name = name
	}
}

/*
	类似go-micro中的opts， 实现默认的多参数
 */
func NewServices(opts ...Option) Services {
	c := Config{
		Name: "default name",
		Timeout: time.Second,
	}
	// 每个函数设置一个参数
	for _, op := range opts {
		op(&c)
	}
	s := Services{}
	s.conf = c
	return s
}


func TestDefaultParams(t *testing.T) {
	s := NewServices(
		//SetName("peter"),   // 没有设置name参数， name参数去默认值为default name
		SetTimeout(time.Second*5),
	)

	fmt.Println("name:", s.conf.Name)
	fmt.Println("time", s.conf.Timeout)
}
