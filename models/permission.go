package models

type Permission struct {
	Base
	Name string `gorm:"size:60;not null" json:"name"`
}
