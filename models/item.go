package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ProductID       int     `gorm:"not null" json:"product_id"`
	VariantOptionID int     `gorm:"not null" json:"variant_option_id"`
	Price           float64 `gorm:"type:numeric(6,2);not null" json:"price"`
	Stock           int     `gorm:"not null" json:"stock"`
	Unit            string  `gorm:"type:item_unit;default:'pcs.';not null" json:"unit"`
	SKU             string  `gorm:"size:50;not null;unique" json:"sku"`
}
