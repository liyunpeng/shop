package test

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

const (
	_CHAN_SIZE  = 10
	_GUARD_SIZE = 10

	_TEST_CNT = 32
)

type Obj struct {
	flag int64
	c    chan interface{}
}

func (obj *Obj) readLoop() error {
	counter := _TEST_CNT
	for {
		time.Sleep(5 * time.Millisecond)
		if len(obj.c) > _CHAN_SIZE {
			return errors.New(fmt.Sprintf("Chan overflow, len: %v.", len(obj.c)))
		} else if len(obj.c) > 0 {
			<-obj.c
			counter--
		}
		if counter <= 0 {
			return nil
		}
	}
}

func (obj *Obj) writeMsg(idx int, v interface{}) (err error) {
	for {
		if len(obj.c) < _CHAN_SIZE {
			obj.c <- v
			logger.Info.Printf("R(%v)+1 ", idx)
			return nil
		}
	}
}

func (obj *Obj) writeMsgWithCASCheck(idx int, v interface{}) (err error) {
	for {
		if atomic.CompareAndSwapInt64(&obj.flag, 0, 1) {
			if len(obj.c) < _CHAN_SIZE {
				obj.c <- v
				atomic.StoreInt64(&obj.flag, 0)
				logger.Info.Printf("R(%v)+1 ", idx)
				return nil
			} else {
				atomic.StoreInt64(&obj.flag, 0)
			}
		}
	}

	return nil
}

func StartRoutineWriteChan(useCAS bool) {
	//useCAS := false
	//if len(os.Args) > 1 && os.Args[1] == "cas" {
	//	useCAS = true
	//}
	routineCnt := 4
	tryCnt := _TEST_CNT / routineCnt
	var obj = &Obj{c: make(chan interface{}, _CHAN_SIZE+_GUARD_SIZE)}

	for idx := 0; idx < routineCnt; idx++ {
		go func(nameIdx int) {
			for tryIdx := 0; tryIdx < tryCnt; tryIdx++ {
				if useCAS {
					obj.writeMsgWithCASCheck(nameIdx, nil)
				} else {
					obj.writeMsg(nameIdx, nil)
				}
			}
		}(idx)
	}

	// logger.Info.Println(casObj.readLoop())
	logger.Info.Println(obj.readLoop())
	logger.Info.Println("quit.")
}

func TestStartRoutineWriteChan( t *testing.T){
	StartRoutineWriteChan(false)

	StartRoutineWriteChan(true)
}
