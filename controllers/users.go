package controllers

import (
	"litstore/api/config"
	"litstore/api/dto/requests"
	"litstore/api/dto/responses"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GetSelfUser godoc
// @Summary      Get Self User
// @Description  Get Currently Logged User by JWT from Cookies
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.GetUserResponse
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/me [get]
func GetUserSelf(c *gin.Context) {
	// Get Access Token from Cookies
	tokenString, err := c.Cookie(config.JwtAccessName)

	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "No token provided"})
		c.Abort()
		return
	}

	// Try to Parse JWT from Cookie
	token, err := utils.ParseJWT(tokenString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "Invalid token"})
		c.Abort()
		return
	}

	// Try to Get UserID from JWT
	userID, err := token.Claims.GetSubject()

	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "Invalid token"})
		c.Abort()
		return
	}

	// Fetch User from DB
	var user models.User

	result := initializers.DB.Preload("Roles.Permissions").Preload("Permissions").Preload("Addresses").Where("ID = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "User not found"})
		c.Abort()
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
// @Success      200  {array}   responses.GetUserResponse
// @Failure      401  {object}  responses.Error
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User

	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary      Get User by ID
// @Description  Get User by their ID
// @Tags         user
// @Param id path string true "User ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}   models.User
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/id/{id} [get]
func GetUserById(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	result := initializers.DB.Preload("Roles.Permissions").Preload("Permissions").Preload("Addresses").Where("ID = ?", id).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "User not found"})

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
// @Success      200  {array}   responses.GetUserResponse
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/search/{phrase} [get]
func GetUsersBySearch(c *gin.Context) {
	phrase := c.Param("phrase")

	var user models.User

	result := initializers.DB.Where("Email LIKE ?", "%"+strings.ToLower(phrase)+"%").First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "User not found"})

		return
	}

	c.JSON(http.StatusOK, user)
}

// EditUserById godoc
// @Summary      Updates user data
// @Description  Finds user by ID and updates with values provided in body
// @Tags         user
// @Param id path string true "User ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}   responses.GetUserResponse
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /users/id/{id} [put]
func EditUserById(c *gin.Context) {
	id := c.Param("id")

	var body models.User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error{Message: "Invalid body data"})
		return
	}

	var user models.User

	findUser := initializers.DB.Where("ID = ?", id).First(&user)

	if findUser.Error != nil {
		c.JSON(http.StatusNotFound, responses.Error{Message: "User not found"})
		return
	}

	result := initializers.DB.Model(&user).Updates(&body)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Message: "Cannot update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUserById godoc
// @Summary      Deletes user
// @Description  Finds user by ID and deletes
// @Tags         user
// @Param id path string true "User ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.Success
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/id/{id} [delete]
func DeleteUserById(c *gin.Context) {
	id := c.Param("id")

	// Perform the delete operation
	result := initializers.DB.Where("ID = ?", id).Delete(&models.User{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, responses.Error{Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, responses.Error{Message: "User deleted"})
}

// InsertUserAddress godoc
// @Summary      Inserts user's address
// @Description  Inserts address for currently logged user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.GetUserAddressResponse
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/address/new [post]
func InsertUserAddress(c *gin.Context) {
	var body requests.InsertUserAddressRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error{Message: "Invalid body data"})
		return
	}

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "UserID not provided"})
		return
	}

	var count int64
	err := initializers.DB.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Message: "Failed to insert address"})
		return
	}

	address := models.Address{
		UserID:     userID.(uuid.UUID),
		Name:       body.Name,
		Surname:    body.Surname,
		Street:     body.Street,
		House:      body.House,
		Flat:       body.Flat,
		PostCode:   body.PostCode,
		City:       body.City,
		Phone:      body.Phone,
		Country:    body.Country,
		OrderIndex: uint(count),
	}

	result := initializers.DB.Create(&address)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Message: "Failed to insert address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Address inserted successfully",
		"address": address,
	})
}

// GetUserAddresses godoc
// @Summary      Fetches all user's addresses
// @Description  Fetches all addresses for currently logged user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {array}  responses.GetUserAddressResponse
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/address/all [get]
func GetUserAddresses(c *gin.Context) {
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "UserID not provided"})
		return
	}

	var addresses []models.Address

	result := initializers.DB.Where("user_id = ?", userID).Order("order_index").Find(&addresses)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Message: "Failed to fetch addresses"})
		return
	}

	// map to DTO
	response := make([]responses.GetUserAddressResponse, 0, len(addresses))

	for _, addr := range addresses {
		response = append(response, responses.GetUserAddressResponse{
			ID:       addr.ID.String(),
			Name:     addr.Name,
			Surname:  addr.Surname,
			Street:   addr.Street,
			House:    addr.House,
			Flat:     addr.Flat,
			PostCode: addr.PostCode,
			City:     addr.City,
			Phone:    addr.Phone,
			Country:  addr.Country,
		})
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUserAddress godoc
// @Summary      Delete user address by ID
// @Description  Deletes user's address by specified ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {array}  responses.Success
// @Failure      401  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Router       /users/address/{id} [delete]
func DeleteUserAddress(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "UserID not provided"})
		return
	}

	result := initializers.DB.Where("user_id = ? AND id = ?", userID, id).Delete(&models.Address{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, responses.Error{Message: "Address not found"})
		return
	}

	c.JSON(http.StatusOK, responses.Error{Message: "Address deleted successfully"})
}

// ChangePassword godoc
// @Summary      Change Password
// @Description  Change Password for currently logged user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /user/password/change [post]
func ChangePassword(c *gin.Context) {
	var body requests.ChangePasswordRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	// Get user from context
	userID, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, responses.Error{Message: "UserID not provided"})
		return
	}

	// Check old password
	var user models.User

	result := initializers.DB.Where("ID = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Cannot fetch user",
		})

		c.Abort()
		return
	}

	// Compare hashes
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid old password",
		})

		return
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	// Update password in DB
	user.Password = string(hash)
	result = initializers.DB.Save(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}
