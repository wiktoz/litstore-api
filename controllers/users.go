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

// GetSelfUser godoc
// @Summary      Get Self User
// @Description  Get Currently Logged User by JWT from Cookies
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.User
// @Failure      401  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Router       /users/me [get]
func GetUserSelf(c *gin.Context) {

	// Get Access Token from Cookies
	tokenString, err := c.Cookie(config.JwtAccessName)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{Message: "No token provided"})
		c.Abort()
		return
	}

	// Try to Parse JWT from Cookie
	token, err := utils.ParseJWT(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{Message: "Invalid token"})
		c.Abort()
		return
	}

	// Try to Get UserID from JWT
	userID, err := token.Claims.GetSubject()

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{Message: "Invalid token"})
		c.Abort()
		return
	}

	// Fetch User from DB
	var user models.User

	result := initializers.DB.Preload("Roles.Permissions").Preload("Permissions").Preload("Addresses").Where("ID = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "User not found"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsers godoc
// @Summary      Get Users
// @Description  Fetches all users from DB
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      401  {object}  models.Error
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User

	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary      Get User by ID
// @Description  Get User by their ID
// @Tags         user
// @Param id path int true "User ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.User
// @Failure      401  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Router       /users/id/{id} [get]
func GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	result := initializers.DB.Preload("Roles.Permissions").Preload("Permissions").Preload("Addresses").Where("ID = ?", id).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "User not found"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUsersBySearch godoc
// @Summary      Get Users by Search
// @Description  Finds users by a search phrase
// @Tags         user
// @Param phrase path string true "Search Phrase"
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      401  {object}  models.Error
// @Failure      404  {object}  models.Error
// @Router       /users/search/{phrase} [get]
func GetUsersBySearch(c *gin.Context) {
	phrase := c.Param("phrase")

	var user models.User

	result := initializers.DB.Where("Email LIKE ?", "%"+strings.ToLower(phrase)+"%").First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.Error{Message: "User not found"})

		return
	}

	c.JSON(http.StatusOK, user)
}

func EditUserById(c *gin.Context) {

}

func DeleteUserById(c *gin.Context) {

}
