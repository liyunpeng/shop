package models

import "shop/util"

type Model struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt util.JSONTime `json:"createdAt"`
	UpdatedAt util.JSONTime `json:"updatedAt"`
	DeletedAt util.JSONTime `json:"deletedAt"`
}

