package workerpool

import (
	"shop/logger"
)
type Worker struct {
	//Worker的Job channel，当WorkerPool读取到Job，
	// 并拿到可用的Worker的时候，
	// 会将Job实例写入该Worker的Job channel，用来直接执行Do()方法。
	JobQueue chan Job
}

func NewWorker() Worker {
	return Worker{JobQueue: make(chan Job)}
}

//每一个被初始化的worker都会在后期单独占用一个协程
//初始化的时候会先把自己的JobQueue传递到Worker通道中，
func (w Worker) Run(wq chan Worker) {
	go func() {
		for {
			logger.Info.Println("worker阻塞读取自己的JobQueue")
			// 向worker通道里增加一个worker, 因为这个时间点，本worker已经干完活了， 手上没有活了，
			// 所以把这个worker放到worker池里面
			wq <- w
			select {
			// 阻塞在取job, 即是worker待命的时候， 只要来job， 立即干活
			case job := <-w.JobQueue:
				logger.Info.Println("读到一个Job就执行Job对象的Do()方法")
				job.Do()
			}
		}
	}()
}
