package models

type ProductPhoto struct {
	Base
	ProductID  uint   `gorm:"not null" json:"product_id"`
	ImgURL     string `gorm:"size:100" json:"img_url"`
	OrderIndex uint   `gorm:"not null" json:"order_index"`
}
