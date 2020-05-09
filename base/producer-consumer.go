package main

import (
	"fmt"
	"sync"
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

func producerConsumer() {

	c := make(chan int) //   chan   //读写

	go producer(c) //生产者

	consumer(c) //消费者

	fmt.Println("done")
}

func writeTwoChan() {
	ch := make(chan string)
	go func() {
		for m := range ch {
			fmt.Println("processed:", m)

		}

	}()

	/*
		这是bug2
		一般在主程序里都是读通道的动作， 现在把两个写通道放在主程序里，是为了验证个结论，
		我们故意把这个函数放在main函数最后执行， 最后程序只输出了
		processed: cmd.1
		没有输出 processed: cmd.2
		原因是在主程序在写了管道之后， routine的读管道的阻塞的解除是需要一点时间的，
		而主程序在写完管道后，输出一条语句后，就直接退除了，也就失去了对控制台的输出权利
		这时routine再向控制台输出时，控制台是接收不到的。
		想看到routine的完整输出打印，有两个办法
		一是在最后一个写管道后等待两秒，即添加
		time.Sleep(2*time.Second)
		另一种办法是把输出的log打印到文件里，即在routine里用下面语句打印：
		global.fmt.Println("processed:", m)
		bug1和bug2的原因相同，都是主程序在routine还没结束就退出了
	*/
	ch <- "cmd.1"
	ch <- "cmd.2" //won't be processed
	//time.Sleep(2*time.Second)
	close(ch)
	fmt.Println("writeTwoChan end")
}



func closeChan() {
	done := make(chan struct{})
	go func() {
		/*
			这是bug1,
			原因：从<- done routine exit没有打印出来，
			可以断定此routine在主程序退出前没有结束，
			通过实验得知， close(done)能够解除所有读done通道的阻塞操作，
			但是稍微会晚一点时间，1秒左右的时间， 但是主程序退出了
			导致1秒后routine的读done通道阻塞解除时， 主程序已经结束了，
			对控制台的输出就结束了，导致routine不能向控制台输出。
			所以close通道， 能解除routine里读通道的阻塞，只是时间稍微延后
			为了能从log判断出routine是否正常退出，
			每个routine的结尾都要打印一个结束log info，不是只有在debug时才打开的log.
		*/
		<-done
		fmt.Println("<- done closeChan routine exit")
	}()
	close(done)
	fmt.Println("closeChan  programe  exit")
}

func closeChan1() {
	done := make(chan struct{})
	go func() {
		for {
			select {

			case <-done:
				fmt.Println("<- done closeChan1 routine exit")
				return
			}
		}

	}()
	close(done)
	fmt.Println("closeChan1  programe  exit")
}
/*
	管道操作的样本程序
	用于生产者消费者的wq管道
	用于让每个routine开始执行退出动作的done管道
	用于等待每个routine都执行完的wg信号
*/
func chanStddPrograme() {

	/*
		命名习惯:
		用于信号通知的通道命名为done,  因为不存放具体数据， 所以用空结构体
		用于生产者消费者模式的通道命名为wq， 因为要存放数据，而数据类型不确定， 所以用interface类型，即可以是任意类型数据
		等待每个routine结束的wg信号
	*/
	done := make(chan struct{})
	wq := make(chan interface{})
	var wg sync.WaitGroup
	workerCount := 2
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go consumer2(i, wq, done, &wg)
	}

	/*
			通道被用在两种情况：
		   	一，用于生产者消费者模式，
		    判断标准： 有连续写的动作就是被用于生产者消费者模式， 这里连续写两次，这里的wq即是， wq是work queue缩写
			二，用于信号通知，
		 	判断标准： 只有一次写，或没有写，只有一个close动作， done通道即是
	*/
	for i := 0; i < workerCount; i++ {
		wq <- i
	}

	/*
	 解除所有的读通道， 不管是几次读，
	 这里读通道阻塞的routine被启动了两次， 就有两个读通道阻塞，close将这两个阻塞全部解除
	*/
	close(done)
	/*
		用信号等待所有routine结束的原因
		本来上面的close动作能够解锁所有routine里的done通道的阻塞，然后routine就结束了，
		但是routine做这些事的时候是磨磨蹭蹭的， 主程序是很麻利的往下执行了，
		主程序下面事也少， 主程序就麻利的退出了
		而这时候routine还没磨蹭完，routine需要用的log文件以及控制台，都被主程序给关闭了，
		routine就执行异常了
		所以主程序必须用个信号等待所有的routine的磨蹭的把事情做完，主程序才退出
		所以，在主程序的末尾都会有这种标准写法：
		close(done)
		wg.wait()
		两个动作挨在一起
		close让每个routine开始执行退出动作
		wg.wait()等待每个routine执行完退出动作
	*/
	wg.Wait()
	fmt.Println("all done!")
}
func consumer2(routineId int, wq <-chan interface{}, done <-chan struct{}, wg *sync.WaitGroup) {
	fmt.Printf("routine[%v]  is running\n", routineId)
	fmt.Printf("[%v] is running\n", routineId)
	/*
	 defer 在本函数的结束前的最后一步调用， 但很多时候， 不能确定本函数在哪里结束
	*/
	defer wg.Done()
	/*
			接收管道的标准写法：
			go 起个routine
			在rouinte里 用
		    for {
				select {
			建立等待消息的循环， 一般都不是只接收一次消息，所以用for无限循环
			有两个case
			一个case, 做消费者， 即读通道wq
			一个case, 做接收结束信号用， 即读通道done,  整个routine就在这结束，即return
	*/
	for {
		select {
		/*
			对wq通道读操作，可以判断本routine充当消费者角色
		*/
		case product := <-wq:
			fmt.Printf("routine[%v] product[%v] is consumed \n", routineId, product)
			fmt.Printf("routine[%v] product[%v] is consumed \n", routineId, product)
		case <-done:
			fmt.Printf("[%v] is done\n", routineId)
			fmt.Printf("[%v] is done\n", routineId)
			return
		}
	}
}

func main() {
	fmt.Println("<------------------------- Chan begin -------------------->")
	producerConsumer()

	closeChan1()

	chanStddPrograme()

	closeChan()

	writeTwoChan()
	fmt.Println("<------------------------- Chan end -------------------->")
}
