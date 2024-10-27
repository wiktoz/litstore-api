package initializers

import "litstore/api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
