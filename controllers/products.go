package controllers

import (
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	initializers.DB.Find(&products)

	c.IndentedJSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")

	// Validate the ID format (e.g., UUID with 36 characters)
	if len(id) != 36 || !utils.ValidateUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID format!",
		})
		return
	}

	var product models.Product

	result := initializers.DB.Where("id = ?", id).Find(&product)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot find product",
		})

		return
	}

	c.JSON(http.StatusOK, product)
}

func GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var product models.Product

	err := initializers.DB.
		Preload("ProductImages", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("ProductImages.Image").
		Preload("Descriptions").
		Preload("Variants").
		Preload("Items").
		First(&product, "slug = ?", slug).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Product not found",
			})
			return
		}

		// Some other DB error
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func EditProductById(c *gin.Context) {
	id := c.Param("id")

	// Validate the ID format (e.g., UUID with 36 characters)
	if len(id) != 36 || !utils.ValidateUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID format!",
		})
		return
	}
}

func DeleteProductById(c *gin.Context) {
	id := c.Param("id")

	// Validate the ID format (e.g., UUID with 36 characters)
	if len(id) != 36 || !utils.ValidateUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID format!",
		})
		return
	}

	// Perform the delete operation
	result := initializers.DB.Where("ID = ?", id).Delete(&models.Product{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Product not found!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product deleted successfully!",
	})
}

func GetProductsBySearch(c *gin.Context) {
	phrase := c.Param("phrase")

	var product models.Product

	result := initializers.DB.Where("Email LIKE ?", "%"+strings.ToLower(phrase)+"%").Find(&product)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Product not found",
		})

		return
	}

	c.JSON(http.StatusOK, product)
}

func InsertProduct(c *gin.Context) {
	type CreateProductInput struct {
		Name          string     `json:"name" binding:"required,min=3"`
		Manufacturer  string     `json:"manufacturer" binding:"required,min=3"`
		New           bool       `json:"new"`
		Active        bool       `json:"active"`
		CategoryID    *uuid.UUID `json:"category_id"`
		SubcategoryID *uuid.UUID `json:"subcategory_id"`

		ImageIDs []uuid.UUID `json:"image_ids" binding:"required"`
	}

	var input CreateProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Create product instance
	product := models.Product{
		Name:          input.Name,
		Manufacturer:  input.Manufacturer,
		New:           input.New,
		Active:        input.Active,
		CategoryID:    input.CategoryID,
		SubcategoryID: input.SubcategoryID,
	}

	tx := initializers.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Prepare ProductImages join records
	var productImages []models.ProductImage
	for i := range input.ImageIDs {
		productImages = append(productImages, models.ProductImage{
			ProductID:  product.ID,
			ImageID:    input.ImageIDs[i],
			OrderIndex: i,
		})
	}

	if len(productImages) > 0 {
		if err := tx.Create(&productImages).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": product})
}
