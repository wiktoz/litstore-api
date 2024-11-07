package initializers

import (
	"litstore/api/models"
	"log"
)

func SyncDatabase() {
	DB.Exec("CREATE TYPE lang_type AS ENUM ('pl', 'en', 'fr', 'de');")
	DB.Exec("CREATE TYPE select_type AS ENUM ('button', 'select');")
	DB.Exec("CREATE TYPE unit_type AS ENUM ('pc.', 'l', 'kg', 'set');")

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
}
