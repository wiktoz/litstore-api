package models

import "gorm.io/gorm"

type Subcategory struct {
	gorm.Model
	CategoryID     uint   `gorm:"not null" json:"category_id"`
	Name           string `gorm:"size:50;not null" json:"name"`
	Description    string `gorm:"type:text" json:"description"`
	SeoDescription string `gorm:"type:text" json:"seo_description"`
	Img            string `gorm:"size:100" json:"img"`
	BgImg          string `gorm:"size:100" json:"bg_img"`
	DisplayNavbar  bool   `gorm:"default:true" json:"display_navbar"`
	DisplayFooter  bool   `gorm:"default:true" json:"display_footer"`
	Active         bool   `gorm:"default:true" json:"active"`
	Slug           string `gorm:"size:60;not null" json:"slug"`
}
