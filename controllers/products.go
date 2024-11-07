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
	name := c.Param("name")

	var product = models.Product{Name: name}

	initializers.DB.Find(&product)
	c.IndentedJSON(http.StatusOK, product)
}

func EditProductById(c *gin.Context) {

}

func DeleteProductById(c *gin.Context) {

}

func GetProductsBySearch(c *gin.Context) {

}

func InsertProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "You have access!",
	})

	//var body models.Product

	//c.Bind(&body)

	//initializers.DB.Create(&body)
}
