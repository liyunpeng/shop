package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"strconv"
	"sync"
)

type Goods struct {
	gorm.Model
	//Model
	Name       string  `gorm:"type:varchar(255)" json:"name"`
	Description string  `gorm:"type:varchar(128);not null" json:"description"`
	Type string  `gorm:"type:varchar(64);not null" json:"type"`
	Price       float32 `gorm:"type:decimal(7,2)" json:"price"`
	ImagePath   string  `gorm:"type:varchar(255)" json:"imagepath"`
	Stock  int   `gorm:"type:int(10)" json:"stock"`

}

func (o Goods) TableName() string {
	return "gorm_goods"
}

func GoodsDelete() {
	DB.Delete(&Goods{})
}

func CreateGoods(goods *Goods) {
	if err := DB.Create(goods).Error; err != nil {
		color.Red(fmt.Sprintf("CreateGoodsErr:%s \n ", err))
	}
}

func BuyGood(id int ) {
	//DB.Model(&Goods{}).Where("id =?", id).Update(“”)
	//DB.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
	var mutex sync.Mutex
	mutex.Lock()
	DB.Raw("UPDATE gorm_goods SET stock = stock - 1  WHERE id = (" + strconv.Itoa(id) + ")")
	mutex.Unlock()
}

func GoodsFindByName(name string) []*Goods {
	var goods []*Goods
	tx := DB.Model(&Goods{}).Where("name =?", name).Find(&goods)
	if tx.Error == nil {
		return goods
	} else {
		return nil
	}
}
func GoodsFindAll() []*Goods {
	var goods []*Goods
	tx := DB.Model(&Goods{}).Find(&goods)
	if tx.Error == nil {
		return goods
	} else {
		return nil
	}
}
func GoodsFindById(id int64) *Goods {
	goods := new(Goods)
	//DB.Model(&Goods{}).Find(&Goodss)
	tx := DB.Model(&Goods{}).Where("id =?", id).Find(goods)
	if tx.Error == nil {
		return goods
	} else {
		return nil
	}
}
