package workerpool

import (
	"runtime"
	"shop/logger"
	"time"
)

type Score struct {
	Num int
}

func (s *Score) Do() {
	logger.Info.Println("num:", s.Num)
	//这里延迟是模拟处理数据的耗时
	time.Sleep(1 * 1 * time.Second)
}

func GopollMain() {

	/*
		num := 100 * 100 * 20
		这句表示worker协程池提前准备 20万个worker 生产环境需要这么多的workder
		开发环境为了调试方便， 用准备2个worker
	*/
	num := 2
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	p := NewWorkerPool(num)
	p.Run()

	/*
		datanum := 100 * 100 * 100 * 100
		模拟生产环境短时间发一亿个请求
	*/
	datanum := 3
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
			p.JobQueue <- sc
			//p.JobQueue <- sc
			//p.JobQueue <- sc
			//p.JobQueue <- sc
		}
	}()
	timer := time.NewTimer(10 * time.Second)

loopa:

	for {

		logger.Info.Println("启动的routine个数统计： runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
		select {
		case <-timer.C:
			break loopa

		}
	}

	logger.Info.Println("启动的routine个数统计： runtime.NumGoroutine() :", runtime.NumGoroutine())
	logger.Info.Println("GopollMain 结束")
}
