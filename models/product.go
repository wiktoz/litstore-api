package models

import (
	"litstore/api/utils"

	"gorm.io/gorm"
)

type Product struct {
	Base
	Name         string               `json:"name" binding:"required,min=3"`
	Manufacturer string               `json:"manufacturer" binding:"required,min=3"`
	New          bool                 `json:"new"`
	Active       bool                 `json:"active"`
	Slug         string               `json:"slug"`
	Descriptions []ProductDescription `gorm:"foreignKey:ProductID" json:"descriptions"`
	Photos       []ProductPhoto       `gorm:"foreignKey:ProductID" json:"photos"`
	Variants     []Variant            `gorm:"many2many:products_variants" json:"variants"`
	Items        []Item               `gorm:"foreignKey:ProductID" json:"items"`

	CategoryID    *uint `json:"category_id" binding:"omitempty"`
	SubcategoryID *uint `json:"subcategory_id" binding:"omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.Slug = utils.GenerateUniqueSlug(tx, p, "Name")
	return nil
}
