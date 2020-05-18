package main

import (
	"fmt"
)

//   chan<- //只写
func producer(out chan<- int) {
	defer close(out)  // 在最后一个写通道动作后，close通道
	for i := 0; i < 5; i++ {
		fmt.Println("produce: ", i)
		out <- i //如果对方不读 会阻塞
	}
}

//   <-chan //只读
func consumer(in <-chan int) {

	for num := range in {   // range无缓存通道, 要求必须在最后一个写通道的后面，close通道
		fmt.Println("consume: ", num)
	}
}

func main() {

	c := make(chan int) //   chan   //读写

	go producer(c) //生产者

	consumer(c) //消费者

	fmt.Println("done")
}
