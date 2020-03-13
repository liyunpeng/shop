package services

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"strings"
	"sync/atomic"
	"time"
)

type SecondLimit struct {
	unixSecond int64
	curCount   int32
	Limit      int32
}

// NewSecondLimit to init a SecondLimit obj
func NewSecondLimit(limit int32) *SecondLimit {
	secLimit := &SecondLimit{
		unixSecond: time.Now().Unix(),
		curCount:   0,
		Limit:      limit,
	}

	return secLimit
}

// Add is func to
func (s *SecondLimit) Add(count int) {
	sec := time.Now().Unix()
	if sec == s.unixSecond {
		atomic.AddInt32(&s.curCount, int32(count))
		return
	}

	atomic.StoreInt64(&s.unixSecond, sec)
	atomic.StoreInt32(&s.curCount, int32(count))
}

// Wait to Limit num
func (s *SecondLimit) Wait() bool {
	for {
		sec := time.Now().Unix()
		if (sec == atomic.LoadInt64(&s.unixSecond)) && s.curCount >= s.Limit {
			time.Sleep(time.Millisecond)
			logs.Debug("Limit is runing, Limit: %d s.curCount:%d", s.Limit, s.curCount)
			continue
		}

		if sec != atomic.LoadInt64(&s.unixSecond) {
			atomic.StoreInt64(&s.unixSecond, sec)
			atomic.StoreInt32(&s.curCount, 0)
		}
		logs.Debug("Limit is exited")
		return false
	}
}

// LogConfig is struct from etcd log conf
type LogConfig struct {
	Topic    string `json:"topic"`
	LogPath  string `json:"log_path"`
	Service  string `json:"service"`
	SendRate int    `json:"send_rate"`
}

type TailWithConf struct {
	tail     *tail.Tail
	offset   int64
	logConf  LogConfig
	secLimit *SecondLimit
	exitChan chan bool
}

func (t *TailWithConf) readLog(fileName string) {

	for line := range t.tail.Lines {
		if line.Err != nil {
			logs.Error("read line error:%v ", line.Err)
			continue
		}

		lineStr := strings.TrimSpace(line.Text)
		fmt.Println("从被监控的文件", fileName, "中读到的字符串=", lineStr)

		if len(lineStr) == 0 || lineStr[0] == '\n' {
			continue
		}
		fmt.Println("向kafka生产者数据通道发送消息 消息字符串=",
			line.Text, "消息的topic=", t.logConf.Topic)
		kafkaSender.addMessage(line.Text, t.logConf.Topic)
		t.secLimit.Add(1)
		t.secLimit.Wait()

		select {
		case <-t.exitChan:
			logs.Warn("tail obj is exited: config:", t.logConf)
			return
		default:
		}
	}
	waitGroup.Done()
}
