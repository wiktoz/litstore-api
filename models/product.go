package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name         string               `json:"name"`
	Manufacturer string               `json:"manufacturer"`
	New          bool                 `json:"new"`
	Active       bool                 `json:"active"`
	Slug         string               `json:"slug"`
	Descriptions []ProductDescription `gorm:"foreignKey:ProductID" json:"descriptions"`
	Photos       []ProductPhoto       `gorm:"foreignKey:ProductID" json:"photos"`
	Variants     []Variant            `gorm:"many2many:products_variants" json:"variants"`
	Items        []Item               `gorm:"foreignKey:ProductID"`

	CategoryID    uint        `json:"category_id"`
	Category      Category    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SubcategoryID uint        `json:"subcategory_id"`
	Subcategory   Subcategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
