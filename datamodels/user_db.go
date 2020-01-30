package datamodels

import (
	//"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"time"
)

func getDb()(db *gorm.DB){
	defer func(){
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()

	//TODO？ 全局db没用
	/*
		链接localhost数据库， 用户名root, 密码root
	*/
	db, err := gorm.Open(
		"mysql", "root:root@/gotest?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		fmt.Println("open db sucess", db)

	} else {
		fmt.Println("open db error ", err)
		panic("数据库错误")
	}

	return db
}

type UserG struct {
	gorm.Model
	Salt      string `gorm:"type:varchar(255)" json:"salt"`
	Username  string `gorm:"type:varchar(32)" json:"username"`
	Password  string `gorm:"type:varchar(200);column:password" json:"-"`
	Languages string `gorm:"type:varchar(200);column:languages" json:"languages"`
}

func (u UserG) TableName() string {
	return "gorm_user"
}

func (u *UserG) Insert() (err error){
	db := getDb()
	defer func() { db.Close() }()

	/*
		这里插入的post记录， 没有append comment记录， 所有post的related方法不会得到comment记录
	*/
	//db.Create(post)

	/*
			等同于，
		    INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;
	*/
	db.Set("gorm:insert_option", "ON CONFLICT").Create(u)

	return nil
}