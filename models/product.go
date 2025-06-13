package models

import (
	"litstore/api/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name         string `json:"name" binding:"required,min=3"`
	Manufacturer string `json:"manufacturer" binding:"required,min=3"`
	New          bool   `json:"new"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`

	Descriptions  []ProductDescription `gorm:"foreignKey:ProductID" json:"descriptions"`
	ProductImages []ProductImage       `gorm:"foreignKey:ProductID" json:"images"`
	Variants      []Variant            `gorm:"many2many:products_variants" json:"variants"`
	Items         []Item               `gorm:"foreignKey:ProductID" json:"items"`

	CategoryID    *uuid.UUID `json:"category_id" binding:"omitempty"`
	SubcategoryID *uuid.UUID `json:"subcategory_id" binding:"omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.Slug = utils.GenerateUniqueSlug(tx, p, "Name")
	return nil
}
