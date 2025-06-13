package models

import "github.com/google/uuid"

type Address struct {
	Base
	UserID     uuid.UUID `gorm:"not null" json:"user_id"`
	Name       string    `gorm:"size:60;not null" json:"name"`
	Surname    string    `gorm:"size:60;not null" json:"surname"`
	Street     string    `gorm:"size:60;not null" json:"street"`
	House      string    `gorm:"size:20;not null" json:"house"`
	Flat       string    `gorm:"size:20" json:"flat"`
	PostCode   string    `gorm:"size:10;not null" json:"post_code"`
	City       string    `gorm:"size:60;not null" json:"city"`
	Phone      string    `gorm:"size:20;not null" json:"phone"`
	Country    string    `gorm:"size:60;not null" json:"country"`
	OrderIndex uint      `gorm:"not null" json:"order_index"`
}
