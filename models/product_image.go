package models

import "github.com/google/uuid"

type ProductImage struct {
	ProductID  uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	ImageID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	OrderIndex int       `json:"order_index"`

	Image Image `gorm:"foreignKey:ImageID" json:"image"`
}
