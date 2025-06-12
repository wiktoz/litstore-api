package controllers

import (
	"litstore/api/config"
	"litstore/api/emails"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/models/enums"
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
	jwtAccessToken, errAccessToken := utils.GenerateJWT(user.ID.String(), "access")
	jwtRefreshToken, errRefreshToken := utils.GenerateJWT(user.ID.String(), "refresh")
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
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
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
			"err":     err,
		})
		return
	}

	// Revoking access token
	accessToken := utils.Token{Name: config.JwtAccessName, ExpTime: config.JwtAccessExpTime}

	err = utils.RevokeToken(c, initializers.RDB, accessToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Cannot logout at this moment",
			"err":     err,
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out",
	})
}

func ChangePassword(c *gin.Context) {
	var body struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	// Get user from context
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword))

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

func DemandResetPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	// Attempt to find user
	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)

	// If user found, proceed to generate token and send email
	if result.Error == nil {
		// Generate password reset token
		actionToken, token, err := models.GenerateActionToken(user.ID, enums.PasswordReset)
		if err == nil {
			// Save token to DB
			if createErr := initializers.DB.Create(actionToken).Error; createErr == nil {
				// Send reset email (ignore errors to avoid revealing info)
				_ = emails.SendPasswordResetEmail(user.Email, token)
			}
		}
	}

	// Always respond with the same message regardless of whether user exists
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "If a user with that email exists, they will receive a password reset email.",
	})
}

func ResetPassword(c *gin.Context) {
	var body struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	// Get token from DB
	var actionToken models.ActionToken
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Token), 12)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})
		return
	}

	result := initializers.DB.Where("token_hash = ?", string(hash)).First(&actionToken)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token not valid",
		})

		return
	}

	// Verify token
	err = models.VerifyActionToken(&actionToken)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token not valid",
		})
		return
	}

	// Hash new password
	hash, err = bcrypt.GenerateFromPassword([]byte(body.Password), 12)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})
		return
	}

	// Update user password
	var user models.User
	result = initializers.DB.Where("ID = ?", actionToken.UserID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Internal error",
		})
		return
	}

	user.Password = string(hash)
	initializers.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successfully",
	})
}
