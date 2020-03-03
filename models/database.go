package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"shop/config"
	_ "shop/config"
	"time"
)

var DB *gorm.DB
func Register(conf *config.Config) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()
	mysqlConf := conf.Mysql
	var err error

	DB, err = gorm.Open(
		"mysql", "root:123456@/gotest?charset=utf8&parseTime=True&loc=Local")
	//"mysql", "root:password@tcp(192.168.0.220:31111)/gotest?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		DB.DB().SetMaxIdleConns(mysqlConf.MaxIdle)
		DB.DB().SetMaxOpenConns(mysqlConf.MaxOpen)
		DB.DB().SetConnMaxLifetime(time.Duration(300) * time.Minute)
		err = DB.DB().Ping()
		fmt.Println("成功连接数据库")
	} else {
		fmt.Println("没有连接到数据库 err= ", err)
		panic("数据库错误")
	}
}

