package models

import (
	"litstore/api/models/enums"
)

type Variant struct {
	Base
	Name        string           `gorm:"size:60;not null" json:"name"`
	DisplayName string           `gorm:"size:60;not null" json:"display_name"`
	SelectType  enums.SelectType `gorm:"type:select_type;default:'select'" json:"select_type"`
	Options     []VariantOption  `gorm:"foreignKey:VariantID;references:ID" json:"options"`
}
