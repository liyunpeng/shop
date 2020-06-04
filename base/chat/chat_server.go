package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)


type ClientInfo struct {
	message chan string
	Name    string
}

var (
	messages    = make(chan string)
	clientChans = make(map[ClientInfo]bool)
)

func broadcast() {
	for { //  所有不能一次select就退出的都用for
		select {
		case msg := <-messages:
			for cli := range clientChans {
				debugLog(msg)
				cli.message <- msg
			}
		}
	}

	fmt.Println("broadcast routine exit ")

}

func InputTimeout(c net.Conn, timeout time.Duration, input func(chan struct{}) ()) {
	// 一个回车后写done通道， 一个超时后写done通道
	done := make(chan struct{})
	// 客户端每次输入前，写设个sig通道，
	// 这样读sig通道的地方会重新计时，
	// 达到每次输入等待时间限定指定时间内
	sig := make(chan struct{})
	go func() {
		timer := time.NewTimer(timeout)
		for {
			// 不想select一次退出，必须在前有for循环
			select { // select 下所有的case都是在等着读的通道

			// 问：是否有可能主程序已经退出了， routine还在这等待，
			// 答：不会，因为本routine有自超时的处理， 这是避免routine在出程序退出还存在等待的好办法
			case <-sig:
				timer.Reset(timeout)

			case <-timer.C:
				done <- struct{}{} // 通道作为信号使用的写法， 空结构体
				return             // 不是退出for, 而是退出了整个go routine函数
			}
		}
	}() // 启动routine，都是在调用函数，必须有这个圆括号

	go func() {
		input(sig)
		done <- struct{}{}
	}()

	<-done // 写操作的routine和超时通道的routine 都会写这个done通道， 一个写了，这里就解除阻塞
}

func debugLog(s string) {
	//  fmt.Println("[Debug] : ", s)  统一关闭调试打印
}

func handleConn1(c net.Conn) {

	clientinfo := ClientInfo{make(chan string), ""} //不能用new的范式， 不然clientChans[clientinfo]  找不到

	debugLog(c.RemoteAddr().String())

	clientinfo.message = make(chan string) //  只有make 才创建管道， 用var xxx chan stirng只是声明

	go ouputToConnection(c, clientinfo.message)

	clientinfo.message <- "input your name:"
	/*
	    问： NewScanner和newReader都是bufio构造了一个输入流， 两者有何不同
		答： newscanner的源是随时输入的， 如标准输入，和客户端的输入， 客户端的输入就是连接客户端端的网络连接
		这里就是net.conn.  这些都是输入源， 而且阻塞的，input.scan就会阻塞在这里， 有输入并且输入阻塞的，
		用户按下回车键， input.text就返回
		newreader的源就是固定， 一般是一个文件， 和一个字符串。因为没有客户输入， 所以没有scan阻塞函数， 而有直接读取数据的函数ureadline
	*/
	input := bufio.NewScanner(c)
	inputC := func(sig chan struct{}) {
		for i := 0; input.Scan(); i++ {
			sig <- struct{}{}
			if i == 0 {
				clientinfo.Name = input.Text()

				messages <- c.RemoteAddr().String() + " " + clientinfo.Name + " 11111111111 enter chat room"

				clientChans[clientinfo] = true // 作为索引的结构体不能用new的方式，因为new放回的是指针地址
			} else {
				s := input.Text()
				debugLog(s)
				messages <- clientinfo.Name + " : " + s
			}
		}
	}

	InputTimeout(c, 30*time.Second, inputC)

	close(clientinfo.message)

	delete(clientChans, clientinfo)

	messages <- c.RemoteAddr().String() + " leave chat room"

	c.Close()

}
/*
package main
import (
  "fmt"
)

type State uint8
type Event uint8
type Trans struct {
  sourceState State
  event Event
  targetState State
}

const (
  opened State = iota
  closed
  locked
  unlocked
)

const (
  openDoor Event = iota
  closeDoor
  lockDoor
  unlockDoor
)

type Door struct {
  state State
}

func (d *Door) ChangeState(transArray []Trans, sourceState State, event Event){
  for _, v := range transArray {
    if v.sourceState == sourceState && v.event == event {
      d.state = v.targetState
      break
    }
  }
}

func main() {
  transArray := []Trans{
    {opened,  closeDoor, closed},
    {closed,  lockDoor, locked},
    {locked,  unlockDoor, unlocked},
    {unlocked, openDoor, opened},
  }

  d :=  &Door{
    state: opened,
  }

  d.ChangeState(transArray, opened, closeDoor)
  fmt.Println("door state is changed to ", d.state)

  d.ChangeState(transArray, closed, lockDoor)
  fmt.Println("door state is changed to ", d.state)


}
 */
func ouputToConnection(c net.Conn, data <-chan string) {
	for v := range data {
		fmt.Fprintf(c, " %s \r\n", v) //  需要加深理解
	}
}

func main() {
	listener, _ := net.Listen("tcp", ":7890")

	go broadcast() //  需要加深理解

	for { //  需要加深理解
		conn, _ := listener.Accept()
		go handleConn1(conn)
	}
}