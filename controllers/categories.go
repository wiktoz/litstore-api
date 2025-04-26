package controllers

import (
	"litstore/api/initializers"
	"litstore/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertCategory(c *gin.Context) {
	var body models.Category

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"err":     err.Error(),
		})

		return
	}

	result := initializers.DB.Create(&body)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot create category",
			"err":     result.Error,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully created category",
		"product": result,
	})
}

func EditCategoryById(c *gin.Context) {

}

func GetCategories(c *gin.Context) {
	var categories []models.Category

	result := initializers.DB.Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "Categories not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategoryById(c *gin.Context) {
	var category models.Category

	result := initializers.DB.Where("ID = ?", c.Param("id")).First(&category)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func GetCategoryBySlug(c *gin.Context) {
	var category models.Category

	result := initializers.DB.Where("slug = ?", c.Param("slug")).First(&category)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)

}

func DeleteCategoryById(c *gin.Context) {

}
