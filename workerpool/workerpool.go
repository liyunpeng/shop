package workerpool

import "shop/logger"
type WorkerPool struct {
	workerlen int
	//WorkerPool的Job channel，用于调用者把具体的数据写入到这里，
	JobQueue    chan Job
	WorkerQueue chan Worker
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan Worker, workerlen),
	}
}

/*
当数据无限多的时候func (wp *WorkerPool) Run() 会无限创建协程，这里需要做一些处理，
这里是为了让所有的请求不等待，并且体现一下最大峰值时的协程数。具体因项目而异。
*/
func (wp *WorkerPool) Run() {
	logger.Info.Println("------------开始创建routine池，就是把所有的routine启动好， 不是请求数据来了，才启动routine----------")
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		logger.Info.Println("创建的worker=", worker)
		/*
			在数据没有到来前， 就启动了所有的routine, 构成一个协成池
		*/
		worker.Run(wp.WorkerQueue)
	}
	logger.Info.Println("------------------------- routine池的创建完成 -----------------\n\n ")

	go func() {
		for {
			select {
			/*
				用routine池的job队列来接收请求的数据
			*/
			case job := <-wp.JobQueue:
				/*
					在工作池中去一个空闲的Worker去执行该Job
					读到一个数据时, 就获取一个可用的Worker，并将Job对象传递到该Worker的chan通道
				*/
				logger.Info.Println("读到一个数据后, 就从worker池中获取一个可用的Worker， 这是启动好的routine， 请求没有等待")
				// 从worker通道里面取出一个worker， 这时worker通道里就少了一个worker
				worker := <-wp.WorkerQueue
				logger.Info.Println("将Job对象传递到该Worker的chan通道")
				worker.JobQueue <- job
			}
		}
	}()
}
