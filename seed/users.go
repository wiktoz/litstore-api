package seed

import (
	"litstore/api/config"
	"litstore/api/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedDefaultRoles(db *gorm.DB) {
	// Create Admin Role
	var adminRole models.Role = models.Role{
		Name:        "admin",
		Permissions: getAdminPermissions(db),
	}

	// Insert the role into the database
	result := db.Create(&adminRole)
	if result.Error != nil {
		log.Println("Error creating admin role:", result.Error)
	} else {
		log.Println("Admin role created successfully")
	}
}

func SeedDefaultUsers(db *gorm.DB) {
	// Get Admin Role
	var adminRole models.Role
	result := db.Where("name = ?", "admin").First(&adminRole)

	if result.Error != nil {
		log.Println("Error retrieving Admin role:", result.Error)
		return
	}

	// Hash password
	var password string = "root2137"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		log.Printf("Cannot hash password\n")
	}

	// Create User
	var user models.User = models.User{Email: "wiktoz05@icloud.com", Password: string(hash), Confirmed: true, Roles: []models.Role{adminRole}}
	result = db.Create(&user)

	if result.Error != nil {
		log.Println("Error creating user:", result.Error)
	} else {
		log.Println("User created successfully:", user.Email)
	}
}

func getAdminPermissions(db *gorm.DB) []models.Permission {
	var permissions []models.Permission

	result := db.Where("name IN ?", config.AllPermissions).Find(&permissions)

	if result.Error != nil {
		log.Println("Error retrieving permissions:", result.Error)
	}

	return permissions
}
