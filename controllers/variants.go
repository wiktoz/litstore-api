package controllers

import (
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/models/enums"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVariants(c *gin.Context) {
	var variants []models.Variant

	result := initializers.DB.Preload("Options").Find(&variants)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot find variants",
		})

		return
	}

	c.JSON(http.StatusOK, variants)
}

func GetVariantById(c *gin.Context) {
	id := c.Param("id")

	var variant models.Variant

	result := initializers.DB.Where("id = ?", id).Find(&variant)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot find product",
		})

		return
	}

	c.JSON(http.StatusOK, variant)
}

func EditVariantById(c *gin.Context) {

}

func DeleteVariantById(c *gin.Context) {

}

func InsertVariant(c *gin.Context) {
	type VariantInput struct {
		Name        string           `json:"name" binding:"required"`
		DisplayName string           `json:"display_name" binding:"required"`
		SelectType  enums.SelectType `json:"select_type" binding:"required"`
		Options     []string         `json:"options" binding:"required,dive,required"` // must be non-empty
	}

	var input VariantInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Options) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one variant option is required"})
		return
	}

	if !enums.IsValidSelectType(string(input.SelectType)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid select_type"})
		return
	}

	variant := models.Variant{
		Name:        input.Name,
		DisplayName: input.DisplayName,
		SelectType:  input.SelectType,
	}

	// Build and assign options
	var options []models.VariantOption
	for i, opt := range input.Options {
		options = append(options, models.VariantOption{
			VariantID:  &variant.ID,
			Name:       opt,
			OrderIndex: uint(i),
		})
	}
	variant.Options = options

	if err := initializers.DB.Create(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create variant"})
		return
	}

	c.JSON(http.StatusCreated, variant)
}
