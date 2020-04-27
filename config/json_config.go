package config

import (
	"encoding/json"
	"io/ioutil"
	"shop/logger"
)

//type Config struct {
//	Server ServerConf `yaml:"server"`
//	Mysql  MysqlConf  `yaml:"mysql"`
//	Redis  RedisConf  `yaml:"redis"`
//}
//
//type ServerConf struct {
//	Port int      `yaml:"port"`
//	List []string `yaml:"list,flow"`
//}
//
//type RedisConf struct {
//	Enable bool `yaml:"enable"`
//}
//
//type MysqlConf struct {
//	MaxIdle  int    `yaml:"maxidle"`
//	MaxOpen  int    `yaml:"maxopen"`
//	User     string `yaml:"user"`
//	Host     string `yaml:"host"`
//	Password string `yaml:"password"`
//	Port     string `yaml:"port"`
//	Name     string `yaml:"name"`
//}

type LogConfig struct {
	Topic    string `json:"topic"`
	LogPath  string `json:"log_path"`
	Service  string `json:"service"`
	SendRate int    `json:"send_rate"`
}

type http struct {
	Host string `json:"host1"`
	Port int    `json:"port1"`
}

type db struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Db     string `json:"db"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type amqp struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Vhost  string `json:"vhost"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type baseConfig struct {
	http
	db   `json:"dbConfig"`
	amqp `json:"rabbitmqConfig"`
}

var (
	HttpConfig *http
	DBConfig   *db
	AmqpConfig *amqp
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    baseConfig
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		logger.Info.Println(err)
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		logger.Info.Println(err)
		return
	}

	HttpConfig = &conf.http
	DBConfig = &conf.db
	AmqpConfig = &conf.amqp
	return
}
