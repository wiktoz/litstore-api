package models

import (
	"litstore/api/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subcategory struct {
	Base
	Name           string `gorm:"size:50;not null" json:"name"`
	Description    string `gorm:"type:text" json:"description"`
	SeoDescription string `gorm:"type:text" json:"seo_description"`
	Img            string `gorm:"size:100" json:"img"`
	BgImg          string `gorm:"size:100" json:"bg_img"`
	DisplayNavbar  bool   `gorm:"default:true" json:"display_navbar"`
	DisplayFooter  bool   `gorm:"default:true" json:"display_footer"`
	Active         bool   `gorm:"default:true" json:"active"`
	Slug           string `gorm:"size:60;not null" json:"slug"`

	CategoryID *uuid.UUID `gorm:"not null" json:"category_id"`
	Category   Category   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Products []Product `gorm:"foreignKey:SubcategoryID"`
}

func (p *Subcategory) BeforeCreate(tx *gorm.DB) (err error) {
	p.Slug = utils.GenerateUniqueSlug(tx, p, "Name")
	return nil
}
