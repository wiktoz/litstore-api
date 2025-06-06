package models

import (
	"litstore/api/models/enums"
)

type Item struct {
	Base
	Price      float64    `gorm:"type:numeric(6,2);not null" json:"price"`
	PromoPrice float64    `gorm:"type:numeric(6,2);default:null" json:"promo_price"`
	Stock      uint       `gorm:"not null" json:"stock"`
	Unit       enums.Unit `gorm:"type:unit_type;default:'pc.'" json:"unit"`
	SKU        string     `gorm:"size:30;not null;unique" json:"sku"`

	ProductID       uint          `gorm:"not null" json:"product_id"`
	VariantOptionID uint          `gorm:"not null" json:"variant_option_id"`
	Product         Product       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VariantOption   VariantOption `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Deliveries []Delivery `gorm:"many2many:items_deliveries" json:"deliveries"`
}
