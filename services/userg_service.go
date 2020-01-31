package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	//"errors"
	//"fmt"
	"github.com/liyunpeng/shop/datamodels"
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
	GetByUsernameAndPassword(username, userPassword string) (datamodels.UserG, bool)
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
func NewUserGService(db1 *gorm.DB) UserGService {

	return &userGService{
		db: db1,
	}
}



func (u *userGService) CreateUsergTable(){
	userg := datamodels.UserG{}

	userg.CreateTable(u.db)

}

func (u *userGService) InsertUserg(userg datamodels.UserG){
	//userg := datamodels.UserG{}

	userg.Insert(u.db)
}

func (u *userGService) GetByUsernameAndPassword(username, userPassword string) (datamodels.UserG, bool){
	userg := &datamodels.UserG{}



	userg.FindByName(u.db, username)
	if len(userg.Username) > 0 {
		fmt.Println("找到用户名=", userg.Username)
		return  *userg, true
	}else{
		fmt.Println("没找到用户名=", username)
		return  *userg, false
	}
}