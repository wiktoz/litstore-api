package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	DB, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	return DB
}

func GetDB() *gorm.DB {
	return DB
}
