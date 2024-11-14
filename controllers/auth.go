package controllers

import (
	"litstore/api/config"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// No data sent
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})

		return
	}

	// Check if user exists
	var userToFind = models.User{Email: body.Email}
	var user models.User

	result := initializers.DB.Where(&userToFind).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})

		return
	}

	// Compare hashes
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})

		return
	}

	// Generate JWT tokens
	jwtAccessToken, errAccessToken := utils.GenerateJWT(user.ID, "access")
	jwtRefreshToken, errRefreshToken := utils.GenerateJWT(user.ID, "refresh")
	csrfToken, errCsrfToken := utils.GenerateToken()

	if errAccessToken != nil || errRefreshToken != nil || errCsrfToken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	// Set Tokens to Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.JwtRefreshName,
		Value:    jwtRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(config.JwtRefreshExpTime),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.JwtAccessName,
		Value:    jwtAccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(config.JwtAccessExpTime),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.CsrfName,
		Value:    csrfToken,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(config.CsrfExpTime),
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully logged in",
	})
}

func Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// No data sent
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})

		return
	}

	// Check if user exists
	var userToFind = models.User{Email: body.Email}
	var user models.User

	result := initializers.DB.Where(&userToFind).First(&user)

	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User already exists",
		})

		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	// Insert user into DB
	user = models.User{Email: body.Email, Password: string(hash)}
	result = initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully registered",
	})
}

func Logout(c *gin.Context) {
	// Revoking refresh token
	refreshToken := utils.Token{Name: config.JwtRefreshName, ExpTime: config.JwtRefreshExpTime}

	err := utils.RevokeToken(c, initializers.RDB, refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Cannot logout at this moment",
		})
	}

	// Revoking access token
	accessToken := utils.Token{Name: config.JwtAccessName, ExpTime: config.JwtAccessExpTime}

	err = utils.RevokeToken(c, initializers.RDB, accessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Cannot logout at this moment",
		})
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out",
	})
}
