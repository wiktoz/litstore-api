package controllers

import (
	"fmt"
	"litstore/api/initializers"
	"litstore/api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	// Generate JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Print(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})

		return
	}

	// Set token to cookie

	c.SetCookie("token", tokenString, 86400, "", "localhost", true, true)

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
