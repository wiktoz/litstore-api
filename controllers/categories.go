package controllers

import (
	"litstore/api/dto/responses"
	"litstore/api/initializers"
	"litstore/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InsertCategory godoc
// @Summary      Insert Category
// @Description  Insert a new category into the database
// @Param        category  body  models.Category  true  "Category object"
// @Tags         category
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.Category
// @Failure      401  {object}  responses.Error
// @Router       /categories/new [post]
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

// EditCategoryById godoc
// @Summary      Edit Category by ID
// @Description  Finds category by ID and updates with values provided in body
// @Tags         category
// @Param id path string true "Category ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.Category
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /categories/id/{id} [put]
func EditCategoryById(c *gin.Context) {

}

// GetCategories godoc
// @Summary      Get Categories
// @Description  Fetches all categories from DB
// @Tags         category
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Category
// @Failure      401  {object}  responses.Error
// @Router       /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category

	result := initializers.DB.Find(&categories)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "Categories not found"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryById godoc
// @Summary      Get Category by ID
// @Description  Get Category by their ID
// @Tags         category
// @Param id path string true "Category ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.Category
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /categories/id/{id} [get]
func GetCategoryById(c *gin.Context) {
	var category models.Category

	result := initializers.DB.Where("ID = ?", c.Param("id")).First(&category)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategoryBySlug godoc
// @Summary      Get Category by Slug
// @Description  Get Category by their Slug
// @Tags         category
// @Param slug path string true "Category Slug"
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.Category
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /categories/slug/{slug} [get]
func GetCategoryBySlug(c *gin.Context) {
	var category models.Category

	result := initializers.DB.Where("slug = ?", c.Param("slug")).First(&category)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)

}

// DeleteCategoryById godoc
// @Summary      Delete Category by ID
// @Description  Delete Category by their ID
// @Tags         category
// @Param id path string true "Category ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.Error
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /categories/id/{id} [delete]
func DeleteCategoryById(c *gin.Context) {

}
