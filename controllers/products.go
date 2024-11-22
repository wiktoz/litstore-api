package controllers

import (
	"litstore/api/initializers"
	"litstore/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	initializers.DB.Find(&products)

	c.IndentedJSON(http.StatusOK, products)
}

func GetProductById(c *gin.Context) {
	id := c.Param("id")

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

	result := initializers.DB.Where("slug = ?", slug).Find(&product)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot find product",
		})

		return
	}

	c.JSON(http.StatusOK, product)
}

func EditProductById(c *gin.Context) {

}

func DeleteProductById(c *gin.Context) {

}

func GetProductsBySearch(c *gin.Context) {

}

func InsertProduct(c *gin.Context) {
	var body models.Product

	c.Bind(&body)

	result := initializers.DB.Create(&body)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Cannot create product",
			"err":     result.Error,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully created product",
		"product": body.ID,
	})
}
