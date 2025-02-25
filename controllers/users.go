package controllers

import (
	"litstore/api/config"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserSelf(c *gin.Context) {
	tokenString, err := c.Cookie(config.JwtAccessName)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token passed", "success": false})
		c.Abort()
		return
	}

	token, err := utils.ParseJWT(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "success": false})
		c.Abort()
		return
	}

	userID, err := token.Claims.GetSubject()

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "success": false})
		c.Abort()
		return
	}

	var user models.User

	result := initializers.DB.Preload("Roles.Permissions").Preload("Permissions").Preload("Addresses").Where("ID = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User not found",
		})

		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsers(c *gin.Context) {

}

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	result := initializers.DB.Preload("Roles.Permissions").Preload("Permissions").Preload("Addresses").Where("ID = ?", id).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User not found",
		})

		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsersBySearch(c *gin.Context) {
	phrase := c.Param("phrase")

	var user models.User

	result := initializers.DB.Where("Email LIKE ?", "%"+strings.ToLower(phrase)+"%").First(&user)

	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "User not found",
		})

		return
	}

	c.JSON(http.StatusOK, user)
}

func EditUserById(c *gin.Context) {

}

func DeleteUserById(c *gin.Context) {

}
