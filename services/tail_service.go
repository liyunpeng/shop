package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
)

type TailService interface {
	RunServer()
}

type tailService struct {

}

func NewTailService() *tailService{
	return &tailService{}
}

func (t *tailService ) RunServer() {
	tailManager = NewTailManager()
	tailManager.Process()
	waitGroup.Wait()
}

var tailManager *TailManager

type TailManager struct {
	tailWithConfMap map[string]*TailWithConf
	lock            sync.Mutex
}

// NewTailManager init TailManager obj
func NewTailManager() *TailManager {
	return &TailManager{
		tailWithConfMap: make(map[string]*TailWithConf, 16),
	}
}

func (t *TailManager) NewTailWithConf(logConfig LogConfig) (*TailWithConf,  error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	fmt.Println("在Tail服务里增加一个被监控的文件：", logConfig)
	_, ok := t.tailWithConfMap[logConfig.LogPath]
	if ok {
		return nil, errors.New("map中已存在该键值")
	}

	tail, err := tail.TailFile(logConfig.LogPath, tail.Config{
		ReOpen:   true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // read to tail
		MustExist: false,                                //file does not exist, it does not return an error
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return nil, err
	}

	tailWithConf := &TailWithConf{
		tail:     tail,
		offset:   0,
		logConf:  logConfig,
		secLimit: NewSecondLimit(int32(logConfig.SendRate)),
		exitChan: make(chan bool, 1),
	}

	return tailWithConf, err
}

func (t *TailManager) reloadConfig(logConfArr []LogConfig) (err error) {
	fmt.Println("hdcloud tail管理者重新加载配置")
	for _, logConfArrValue := range logConfArr {
		tailWithConf, ok := t.tailWithConfMap[logConfArrValue.LogPath]

		if !ok {
			fmt.Println("hdcloud tail监控的文件不存在，则为此文件新建一个tail对象")
			tailWithConf, err = t.NewTailWithConf(logConfArrValue)
			if err != nil {
				logs.Error("add log file failed:%v", err)
				continue
			}

			t.tailWithConfMap[logConfArrValue.LogPath] = tailWithConf

			waitGroup.Add(1)

			fmt.Println("新的监控文件对应一个新的tail对象， ")
			go tailWithConf.readLog(logConfArrValue.LogPath)

			continue
		}
		tailWithConf.logConf = logConfArrValue
		tailWithConf.secLimit.Limit = int32(logConfArrValue.SendRate)
		t.tailWithConfMap[logConfArrValue.LogPath] = tailWithConf
		fmt.Println("tailWithConf:", tailWithConf)
	}

	for key, tailWithConf := range t.tailWithConfMap {
		var found = false
		for _, newValue := range logConfArr {
			if key == newValue.LogPath {
				found = true
				break
			}
		}
		if found == false {
			logs.Warn("log path :%s is remove", key)
			tailWithConf.exitChan <- true
			delete(t.tailWithConfMap, key)
		}
	}
	return
}

func (t *TailManager) Process() {
	for etcdConfValue := range ConfChan {
		logs.Debug("log etcdConfValue: %v", etcdConfValue)

		var logConfArr []LogConfig

		err := json.Unmarshal([]byte(etcdConfValue), &logConfArr)
		fmt.Println("从etcd得到的配置字符串解析出的配置对象: ", logConfArr)

		if err != nil {
			logs.Error("unmarshal failed, err: %v etcdConfValue :%s", err, etcdConfValue)
			fmt.Println("unmarshal failed, err: %v etcdConfValue :%s", err, etcdConfValue)
			continue
		}

		err = t.reloadConfig(logConfArr)
		if err != nil {
			logs.Error("reload config from etcd failed: %v", err)
			continue
		}
	}
}
