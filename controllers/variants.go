package controllers

import (
	"litstore/api/initializers"
	"litstore/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVariants(c *gin.Context) {
	var variants []models.Variant

	result := initializers.DB.Find(&variants)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot find product",
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

}
