package models

import "github.com/google/uuid"

type ProductPhoto struct {
	Base
	ProductID  *uuid.UUID `gorm:"not null" json:"product_id"`
	ImgURL     string     `gorm:"size:100" json:"img_url"`
	OrderIndex uint       `gorm:"not null" json:"order_index"`
}
