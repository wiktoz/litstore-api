package models

type Role struct {
	Base
	Name        string       `gorm:"size:60;not null" json:"name"`
	Permissions []Permission `gorm:"many2many:roles_permissions" json:"permissions"`
}
