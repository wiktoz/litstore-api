package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:60;not null" json:"name"`
	Permissions []Permission `gorm:"many2many:roles_permissions" json:"permissions"`
}
