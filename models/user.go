package models

type User struct {
	Base
	Email       string       `gorm:"size:60;not null" json:"email"`
	Password    string       `gorm:"size:64;not null" json:"password"`
	Confirmed   bool         `gorm:"default:false" json:"confirmed"`
	Blocked     bool         `gorm:"default:false" json:"blocked"`
	Addresses   []Address    `gorm:"foreignKey:UserID;references:ID" json:"addresses"`
	Roles       []Role       `gorm:"many2many:users_roles" json:"roles"`
	Permissions []Permission `gorm:"many2many:users_permissions" json:"permissions"`
}
