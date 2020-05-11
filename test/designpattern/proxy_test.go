package designpattern

import (
	"fmt"
	"testing"
)

type NetWork interface {
	Conn()
	Close()
}

// 声明服务器对象
type Server struct {
	localIP string
}

// 真实的网络连接器，连接服务器
type IpNetWork struct {
	server *Server
}

func (c *IpNetWork) Conn(serverIp string) {
	c.server = &Server{localIP: serverIp}
	fmt.Println(c.server.localIP + "已连接")
}

func (c *IpNetWork) Close() {
	fmt.Println(c.server.localIP + "已关闭连接")
}

// 代理的网络连接器，代理连接器
type ProxyNetWork struct {
	Ip *IpNetWork
}

func (c *ProxyNetWork) Conn(serverIp string) {
	// 代理中实际使用的是真实的网络连接器
	c.Ip = &IpNetWork{}
	c.Ip.Conn(serverIp)
}

func (c *ProxyNetWork) Close() {
	c.Ip.Close()
}

func TestProxy(t *testing.T) {
	proxy := &ProxyNetWork{}
	proxy.Conn("192.168.51.471")
	proxy.Close()
}