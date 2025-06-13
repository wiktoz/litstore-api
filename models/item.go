package models

import (
	"litstore/api/models/enums"

	"github.com/google/uuid"
)

type Item struct {
	Base
	Price      float64    `gorm:"type:numeric(6,2);not null" json:"price"`
	PromoPrice float64    `gorm:"type:numeric(6,2);default:null" json:"promo_price"`
	Stock      uint       `gorm:"not null" json:"stock"`
	Unit       enums.Unit `gorm:"type:unit_type;default:'pc.'" json:"unit"`
	SKU        string     `gorm:"size:30;not null;unique" json:"sku"`
	Active     bool       `gorm:"default:false" json:"active"`

	ProductID uuid.UUID `gorm:"not null" json:"product_id"`
	Product   Product   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	VariantOptions []VariantOption `gorm:"many2many:item_variant_options" json:"variant_options"`

	Deliveries []Delivery `gorm:"many2many:items_deliveries" json:"deliveries"`
}
