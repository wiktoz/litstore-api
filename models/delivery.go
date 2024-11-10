package models

import (
	"litstore/api/utils"

	"gorm.io/gorm"
)

type Delivery struct {
	Base
	Name            string  `gorm:"size:60;not null" json:"name"`
	Description     string  `gorm:"size:150" json:"description"`
	ImgURL          string  `gorm:"size:150;not null" json:"img"`
	Price           float64 `gorm:"type:numeric(6,2);not null" json:"price"`
	FreeFrom        float64 `gorm:"type:numeric(6,2);not null" json:"free_from"`
	PersonalCollect bool    `gorm:"default:false;not null" json:"personal_collect"`
	CashOnDelivery  bool    `gorm:"default:false;not null" json:"cash_on_delivery"`
	Active          bool    `gorm:"default:true;not null" json:"active"`
	Slug            string  `gorm:"size:70;not null" json:"slug"`
}

func (p *Delivery) BeforeCreate(tx *gorm.DB) (err error) {
	p.Slug = utils.GenerateUniqueSlug(tx, p, "Name")
	return nil
}
