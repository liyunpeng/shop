package services

import (
	//"errors"
	"fmt"
	"shop/datamodels"
	"github.com/jinzhu/gorm"
	//"shop/repositories"
)

// UserGService handles CRUID operations of a user datamodel,
// it depends on a user repository for its actions.
// It's here to decouple the data source from the higher level compoments.
// As a result a different repository type can be used with the same logic without any aditional changes.
// It's an interface and it's used as interface everywhere
// because we may need to change or try an experimental different domain logic at the future.
type UserGService interface {
	//GetAll() []datamodels.UserG
	//GetByID(id int64) (datamodels.UserG, bool)
	//GetByUsernameAndPassword(username, userPassword string) (datamodels.UserG, bool)
	//DeleteByID(id int64) bool
	//
	//Update(id int64, user datamodels.UserG) (datamodels.UserG, error)
	//UpdatePassword(id int64, newPassword string) (datamodels.UserG, error)
	//UpdateUsername(id int64, newUsername string) (datamodels.UserG, error)
	//
	//Create(userPassword string, user datamodels.UserG) (datamodels.UserG, error)

	CreateUsergTable()

	InsertUserg(userg datamodels.UserG)
}

type userGService struct {
	//repo repositories.UserRepository,
	//userg datamodels.UserG
	db *gorm.DB
}

// NewUserGService returns the default user service.
func NewUserGService() UserGService {

	return &userGService{
		db: getDb(),
	}
}

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

func (u *userGService) CreateUsergTable(){
	userg := datamodels.UserG{}

	userg.CreateTable(u.db)

}

func (u *userGService) InsertUserg(userg datamodels.UserG){
	//userg := datamodels.UserG{}

	userg.Insert(u.db)

}
