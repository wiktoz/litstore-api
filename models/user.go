package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `gorm:"size:60;not null" json:"email"`
	Password      string `gorm:"size:64;not null" json:"password"`
	MainAddressID int    `json:"main_address_id"`
	RoleID        int    `json:"role_id"`
	Active        bool   `gorm:"default:true" json:"active"`
}
