package models

import (
	"bytes"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Salt      string `gorm:"type:varchar(255)" json:"salt"`
	Username  string `gorm:"unique_index" json:"username"`
	Password  string `gorm:"type:varchar(200);column:password" json:"-"`
	Phonenumber  string `gorm:"type:varchar(200);column:phonenumber" json:"phonenumber"`
	Level string `gorm:"type:varchar(200);column:level" json:"level"`
}

func (u User) TableName() string {
	return "gorm_user"
}

func UserCreateTable() (s string) {
	var buffer bytes.Buffer
	/*
		gorm创建的表名默认为小写开头, 出现大写字符， 则会_分割， 以复数结尾， 可能加s,也可能加es
	*/
	if !DB.HasTable("gorm_user") {
		DB.CreateTable(&User{})
		buffer.WriteString("gorm_user表创建成功\n")
	} else {
		buffer.WriteString("gorm_user表已存在，不再次创建\n")
	}

	return buffer.String()
}

func UserInsert(user *User){
	DB.Create(user)
}

func UserFindByName(name string) *User{
	user := new(User)
	DB.Where("username =?", name).First(user)
	return  user
}

func UserUpdate(user *User) (err error){
	DB.Model(&User{}).Update(user)
	return nil
}

func UserDeleteByName(username string) {
	user := &User{
		Username: username,
	}
	DB.Delete(user)
}

func UserFindById(id uint) *User{
	user := new(User)
	DB.Where("id =?", id).First(user)
	return  user
}

func UserFindAll() ( []*User){
	//var users []*User
	DB.DB().Ping()
	usersa := make([]*User, 100)

	for i := 0; i < 100; i++ {
		usersa[i] = new(User)
	}
	//TODO: find 查找全部咋用， DB.Find(usersa)
	DB.Model(&User{}).First(usersa[0])
	return usersa
}

