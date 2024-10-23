package models

import "gorm.io/gorm"

type Delivery struct {
	gorm.Model
	Name            string  `gorm:"size:60;not null" json:"name"`
	Description     string  `gorm:"type:text" json:"description"`
	Img             string  `gorm:"size:150;not null" json:"img"`
	Price           float64 `gorm:"type:numeric(6,2);not null" json:"price"`
	FreeFrom        float64 `gorm:"type:numeric(6,2);not null" json:"free_from"`
	PersonalCollect bool    `gorm:"default:false;not null" json:"personal_collect"`
	CashOnDelivery  bool    `gorm:"default:false;not null" json:"cash_on_delivery"`
	Active          bool    `gorm:"default:true;not null" json:"active"`
	Slug            string  `gorm:"size:70;not null" json:"slug"`
}
