package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name           string `json:"name"`
	Category_ID    int    `json:"category_id"`
	Subcategory_ID int    `json:"subcategory_id"`
	Manufacturer   string `json:"manufacturer"`
	New            bool   `json:"new"`
	Active         bool   `json:"active"`
	Slug           string `json:"slug"`
}
