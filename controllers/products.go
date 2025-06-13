package controllers

import (
	"fmt"
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
		Preload("Variants.Options").
		Preload("Items.VariantOptions").
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
		Name          string      `json:"name" binding:"required,min=3"`
		Manufacturer  string      `json:"manufacturer" binding:"required,min=3"`
		New           bool        `json:"new"`
		Active        bool        `json:"active"`
		CategoryID    *uuid.UUID  `json:"category_id"`
		SubcategoryID *uuid.UUID  `json:"subcategory_id"`
		ImageIDs      []uuid.UUID `json:"image_ids" binding:"required"`
		VariantIDs    []uuid.UUID `json:"variant_ids" binding:"required"`
	}

	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

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

	// Associate variants with product
	var variants []models.Variant
	if len(input.VariantIDs) > 0 {
		if err := tx.Where("id IN ?", input.VariantIDs).Find(&variants).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
		if err := tx.Model(&product).Association("Variants").Append(variants); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}
	}

	// Associate images with product
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

	// Generate cartesian product items
	if err := generateItemsForProduct(tx, &product); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": product})
}

func generateItemsForProduct(tx *gorm.DB, product *models.Product) error {
	var variants []models.Variant

	// Load variants with options
	if err := tx.Preload("Options").Model(product).Association("Variants").Find(&variants); err != nil {
		return err
	}

	// If no variants, create one item without variant options
	if len(variants) == 0 {
		item := models.Item{
			ProductID: product.ID,
			Stock:     0,
			Price:     0,
			SKU:       fmt.Sprintf("SKU-%s", uuid.NewString()[:8]),
			Active:    false,
		}
		return tx.Create(&item).Error
	}

	// Otherwise: compute cartesian product of variant options
	optionGroups := [][]models.VariantOption{}
	for _, variant := range variants {
		if len(variant.Options) == 0 {
			continue
		}
		optionGroups = append(optionGroups, variant.Options)
	}

	if len(optionGroups) == 0 {
		// No options in selected variants
		item := models.Item{
			ProductID: product.ID,
			Stock:     0,
			Price:     0,
			SKU:       fmt.Sprintf("SKU-%s", uuid.NewString()[:8]),
			Active:    false,
		}
		return tx.Create(&item).Error
	}

	// Recursive function to build combinations
	var combinations [][]models.VariantOption
	var build func(idx int, current []models.VariantOption)
	build = func(idx int, current []models.VariantOption) {
		if idx == len(optionGroups) {
			// Make a copy
			combination := make([]models.VariantOption, len(current))
			copy(combination, current)
			combinations = append(combinations, combination)
			return
		}
		for _, option := range optionGroups[idx] {
			build(idx+1, append(current, option))
		}
	}
	build(0, []models.VariantOption{})

	// Insert all combinations as items
	for _, opts := range combinations {
		item := models.Item{
			ProductID:      product.ID,
			Stock:          0,
			Price:          0,
			SKU:            fmt.Sprintf("SKU-%s", uuid.NewString()[:8]),
			VariantOptions: opts,
			Active:         false,
		}
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}
