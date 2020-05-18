package main

import (
	"io"
	"log"
	"net"
	"os"
	"fmt"
	"sync"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":7890")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	/*
		每起一个routine, 就往里传送一个wg信号，
		routine里把这个wg信号add 1

	 */
	go recvNetwork(conn , &wg)

	sendtoNework(conn)

	conn.CloseWrite()
	wg.Wait()

	fmt.Println("client exit")
}

func recvNetwork(c *net.TCPConn, wg *sync.WaitGroup){
	wg.Add(1)

	/*
		将tcp连接（简称连接）中的数据发给标准输出
		连接里没有数据, io.Cory就阻塞
	 */
	if _, err := io.Copy(os.Stdout, c); err != nil {
		log.Fatal(err)
	}

	fmt.Println("receive data from network connection finished")

	c.CloseRead()

	wg.Done()
}

func sendtoNework(c net.Conn){
	/*
		将标准输入的数据拷贝给tcp连接
		tcp连接是个writer，也是reader, 同时实现了这两个接口
		标准输入没有数据就阻塞，除非ctl+d退出输入
	 */
	if _, err := io.Copy(c, os.Stdin); err != nil {
		log.Fatal(err)
	}
	fmt.Println("send data to network connection finished \n")
}