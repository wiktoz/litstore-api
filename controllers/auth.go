package controllers

import (
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/utils"
	"log"
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
	csrfToken, errCsrfToken := utils.GenerateCSRFToken()

	if errAccessToken != nil || errRefreshToken != nil || errCsrfToken != nil {
		log.Fatalf(errAccessToken.Error())
		log.Fatalf(errRefreshToken.Error())
		log.Fatalf(errCsrfToken.Error())

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	// Set Tokens to Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_refresh_token",
		Value:    jwtRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(utils.JwtRefreshExp),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt_access_token",
		Value:    jwtAccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(utils.JwtAccessExp),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(utils.CsrfExp),
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
