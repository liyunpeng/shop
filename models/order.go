package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	Title  string `gorm:"type:varchar(255)" json:"title"`
	Username  string `gorm:"type:varchar(255)" json:"username"`
}
func (o Order) TableName() string {
	return "gorm_order"
}

func OrderFindByUser( username string) ( []*Order){
	var orders []*Order
	//DB.Model(&Order{}).Find(&orders)
	tx := DB.Model(&Order{}).Where("username =?", username).Find(&orders)
	if ( tx.Error == nil){
		return orders
	}else{
		return nil
	}
}
func OrderFindById( id int) ( *Order){
	//var orders *Order
	order := new(Order)
	//DB.Model(&Order{}).Find(&orders)
	tx := DB.Model(&Order{}).Where("id =?", id).Find(order)
	if ( tx.Error == nil){
		return order
	}else{
		return nil
	}
}

