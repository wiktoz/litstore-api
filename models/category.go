package models

import (
	"litstore/api/utils"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name           string `gorm:"size:50;not null" json:"name"`
	Description    string `gorm:"type:text" json:"description"`
	SeoDescription string `gorm:"type:text" json:"seo_description"`
	ImgURL         string `gorm:"size:100" json:"img_url"`
	BgImgURL       string `gorm:"size:100" json:"bg_img_url"`
	DisplayNavbar  bool   `gorm:"default:true" json:"display_navbar"`
	DisplayFooter  bool   `gorm:"default:true" json:"display_footer"`
	Active         bool   `gorm:"default:true" json:"active"`
	Slug           string `gorm:"size:60;not null; unique" json:"slug"`

	Products      []Product     `gorm:"foreignKey:CategoryID"`
	Subcategories []Subcategory `gorm:"foreignKey:CategoryID"`
}

func (p *Category) BeforeCreate(tx *gorm.DB) (err error) {
	p.Slug = utils.GenerateUniqueSlug(tx, p, "Name")
	return nil
}
