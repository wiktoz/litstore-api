package initializers

import "litstore/api/models"

func SyncDatabase() {
	DB.AutoMigrate(
		&models.Address{},
		&models.Category{},
		&models.Delivery{},
		&models.ProductDescription{},
		&models.Item{},
		&models.Permission{},
		&models.ProductPhoto{},
		&models.Product{},
		&models.Role{},
		&models.Subcategory{},
		&models.User{},
		&models.VariantOption{},
		&models.Variant{},
	)
}
