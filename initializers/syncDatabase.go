package initializers

import (
	"litstore/api/models"
	"litstore/api/seed"
	"log"
)

func SyncDatabase() {
	err := DB.AutoMigrate(
		&models.User{},
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
		&models.VariantOption{},
		&models.Variant{},
	)

	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	log.Println("Database migration complete!")

	seed.SeedDefaultPermissions(DB)

	log.Println("Permission seed complete!")
}
