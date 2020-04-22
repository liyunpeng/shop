package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type OrderDetail struct {
	gorm.Model
	//Model
	Name       string  `gorm:"type:varchar(255)" json:"name"`
	Description string  `gorm:"type:varchar(128);not null" json:"description"`
	Type string  `gorm:"type:varchar(64);not null" json:"type"`
	Price       float32 `gorm:"type:decimal(7,2)" json:"price"`
	ImagePath   string  `gorm:"type:varchar(255)" json:"imagepath"`
	Stock  int   `gorm:"type:int(10)" json:"stock"`

}

func (o OrderDetail) TableName() string {
	return "gorm_OrderDetail"
}

func CreaterrderDetail(orderDetail *OrderDetail) {

	if err := DB.Create(orderDetail).Error; err != nil {
		color.Red(fmt.
			Sprintf("CreateOrderDetailErr:%s \n ", err))
	}
}

func OrderDetailFindByName(name string) []*OrderDetail {
	var orderDetail []*OrderDetail
	tx := DB.Model(&OrderDetail{}).Where("name =?", name).Find(&orderDetail)
	if tx.Error == nil {
		return orderDetail
	} else {
		return nil
	}
}
func OrderDetailFindAll() []*OrderDetail {
	var orderDetail []*OrderDetail
	tx := DB.Model(&OrderDetail{}).Find(&orderDetail)
	if tx.Error == nil {
		return orderDetail
	} else {
		return nil
	}
}
func OrderDetailFindById(id int) *OrderDetail {
	orderDetail := new(OrderDetail)
	//DB.Model(&OrderDetail{}).Find(&OrderDetails)
	tx := DB.Model(&OrderDetail{}).Where("id =?", id).Find(orderDetail)
	if tx.Error == nil {
		return orderDetail
	} else {
		return nil
	}
}
