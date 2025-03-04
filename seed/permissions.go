package seed

import (
	"litstore/api/config"
	"litstore/api/models"
	"log"

	"gorm.io/gorm"
)

/*
Fill default permissions on startup
*/
func SeedDefaultPermissions(db *gorm.DB) {
	defaultPermissions := config.AllPermissions

	for _, permission := range defaultPermissions {
		var existingPermission models.Permission
		if err := db.Where("name = ?", string(permission)).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				db.Create(&models.Permission{Name: string(permission)})
				log.Printf("Permission '%s' added to the database.\n", permission)
			} else {
				log.Printf("Error checking permission '%s': %v\n", permission, err)
			}
		}
	}
}
