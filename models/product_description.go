package models

import (
	"litstore/api/models/enums"
)

type ProductDescription struct {
	Base
	ProductID uint       `gorm:"not null" json:"product_id"`
	Lang      enums.Lang `gorm:"type:lang_type;default:'pl'" json:"lang"`
	Content   string     `gorm:"size:350;not null" json:"content"`
}
