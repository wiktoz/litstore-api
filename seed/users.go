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

func SeedDefaultCategories(db *gorm.DB) {
	// Create Categories
	categories := []models.Category{
		{Name: "Shoes", Description: "All types of shoes", ImgURL: "https://example.com/shoes.jpg", BgImgURL: "https://example.com/bg_shoes.jpg", DisplayNavbar: true, DisplayFooter: true, Active: true},
		{Name: "Clothing", Description: "All types of clothing", ImgURL: "https://example.com/clothing.jpg", BgImgURL: "https://example.com/bg_clothing.jpg", DisplayNavbar: true, DisplayFooter: true, Active: true},
		{Name: "Accessories", Description: "All types of accessories", ImgURL: "https://example.com/accessories.jpg", BgImgURL: "https://example.com/bg_accessories.jpg", DisplayNavbar: true, DisplayFooter: true, Active: true},
	}

	for _, category := range categories {
		result := db.Create(&category)
		if result.Error != nil {
			log.Println("Error creating category:", result.Error)
		} else {
			log.Println("Category created successfully:", category.Name)
		}
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
