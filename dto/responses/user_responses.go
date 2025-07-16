package responses

import "litstore/api/models"

type GetUserAddressResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Street   string `json:"street"`
	House    string `json:"house"`
	Flat     string `json:"flat"`
	PostCode string `json:"post_code"`
	City     string `json:"city"`
	Phone    string `json:"phone"`
	Country  string `json:"country"`
}

type GetUserResponse struct {
	models.Base
	Email       string              `json:"email"`
	Confirmed   bool                `json:"confirmed"`
	Blocked     bool                `json:"blocked"`
	Addresses   []models.Address    `gorm:"foreignKey:UserID;references:ID" json:"addresses"`
	Roles       []models.Role       `gorm:"many2many:users_roles" json:"roles"`
	Permissions []models.Permission `gorm:"many2many:users_permissions" json:"permissions"`
}
