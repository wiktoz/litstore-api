package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name string `gorm:"size:60;not null" json:"name"`
}
