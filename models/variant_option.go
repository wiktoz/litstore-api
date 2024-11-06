package models

import "gorm.io/gorm"

type VariantOption struct {
	gorm.Model
	VariantID  uint   `gorm:"not null" json:"variant_id"`
	Name       string `gorm:"size:60;not null" json:"name"`
	OrderIndex uint   `gorm:"not null" json:"order_index"`

	Items []Item `gorm:"foreignKey:VariantOptionID"`
}
