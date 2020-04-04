package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"shop/validates"
)

type Order struct {
	gorm.Model
	//Model
	Title    string `gorm:"type:varchar(255)" json:"title"`
	Username string `gorm:"type:varchar(255)" json:"username"`
	Description string `gorm:"type:varchar(128);not null" json:"description"`
	Price       float32  `gorm:"type:int(255)" json:"price"`
	ImagePath 	string `gorm:"type:varchar(255)" json:"imagepath"`
}

func (o Order) TableName() string {
	return "gorm_order"
}

//func CreateUser(aul *validates.CreateUpdateUserRequest) (user *User) {
func CreateOrder(aul *validates.CreateOrderRequest) {

	order := &Order{
		Username: aul.Username,
		Title: aul.Title,
	}

	if err := DB.Create(order).Error; err != nil {
		color.Red(fmt.Sprintf("CreateOrderErr:%s \n ", err))
	}

	//addRoles(aul, user)

	return
}

func OrderFindByUser(username string) []*Order {
	var orders []*Order
	//DB.Model(&Order{}).Find(&orders)
	tx := DB.Model(&Order{}).Where("username =?", username).Find(&orders)
	if tx.Error == nil {
		return orders
	} else {
		return nil
	}
}
func OrderFindById(id int) *Order {
	//var orders *Order
	order := new(Order)
	//DB.Model(&Order{}).Find(&orders)
	tx := DB.Model(&Order{}).Where("id =?", id).Find(order)
	if tx.Error == nil {
		return order
	} else {
		return nil
	}
}
