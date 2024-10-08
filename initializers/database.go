package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func UNUSED(x ...interface{}) {}

var DB *gorm.DB

func ConnectDB() {
	DB, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})

	UNUSED(DB)

	if err != nil {
		log.Fatal("Failed to connect to the database")
	}
}
