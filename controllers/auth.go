package controllers

import (
	"litstore/api/config"
	"litstore/api/dto/requests"
	"litstore/api/dto/responses"
	"litstore/api/emails"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/models/enums"
	"litstore/api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary      Login user
// @Description  Login user by email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param credentials body requests.LoginRequest true "Login credentials"
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var body requests.LoginRequest

	// No data sent
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, responses.Error{Message: "Invalid email or password"})

		return
	}

	// Check if user exists
	var userToFind = models.User{Email: body.Email}
	var user models.User

	result := initializers.DB.Where(&userToFind).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, responses.Error{Message: "Invalid email or password"})

		return
	}

	// Compare hashes
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.Error{Message: "Invalid email or password"})

		return
	}

	// Generate JWT tokens
	jwtAccessToken, errAccessToken := utils.GenerateJWT(user.ID.String(), "access")
	jwtRefreshToken, errRefreshToken := utils.GenerateJWT(user.ID.String(), "refresh")
	csrfToken, errCsrfToken := utils.GenerateToken()

	if errAccessToken != nil || errRefreshToken != nil || errCsrfToken != nil {
		c.JSON(http.StatusBadRequest, responses.Error{Message: "Error. Try again later"})

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

	c.JSON(http.StatusOK, responses.Success{Message: "Successfully logged in"})
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param credentials body requests.RegisterRequest true "Register credentials"
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var body requests.RegisterRequest

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

	sendVerificationEmail(user.Email)

	// Response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully registered, verification email sent",
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout user by revoking JWT Token and destroying Cookies
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.Success
// @Failure      500  {object}  responses.Error
// @Router       /auth/logout [post]
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

// VerifyEmail godoc
// @Summary      Verify email
// @Description  Verify account by clicking link from email message
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param token body requests.VerifyEmail true "Token from email message for email verification"
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /auth/email/verify [post]
func VerifyEmail(c *gin.Context) {
	var body requests.VerifyEmail

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	// Get token from DB
	secret, err := utils.ReadHMACSecret(config.HMACSecretPath)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})
		return
	}

	hashedToken := utils.ComputeHMACToken(secret, body.Token)

	var actionToken models.ActionToken
	result := initializers.DB.Where("token_hash = ?", hashedToken).First(&actionToken)

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

	// Update user confirmation status
	var user models.User
	result = initializers.DB.Where("ID = ?", actionToken.UserID).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})
		return
	}
	user.Confirmed = true
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
		"message": "Email verified successfully",
	})
}

func sendVerificationEmail(email string) {
	// Attempt to find user silently
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)

	// If user exists and is not confirmed, proceed
	if result.Error == nil && !user.Confirmed {
		// Revoke old, unused, unexpired tokens
		initializers.DB.Model(&models.ActionToken{}).
			Where("user_id = ? AND action = ? AND used_at IS NULL AND expires_at > ?", user.ID, enums.EmailVerification, time.Now()).
			Update("used_at", time.Now())

		// Generate new token
		actionToken, token, err := models.GenerateActionToken(user.ID, enums.EmailVerification)
		if err != nil {
			return // silent failure
		}

		// Store new token
		if createErr := initializers.DB.Create(actionToken).Error; createErr != nil {
			return // silent failure
		}

		// Send email
		_ = emails.SendVerificationEmail(user.Email, token) // silently ignore errors
	}
}

// ResendVerificationEmail godoc
// @Summary      Resend Verification Email
// @Description  Resend Email with Verification Token to activate user's account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param email body requests.ResendVerificationEmailRequest true "Email address to resend verification token"
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /auth/email/resend [post]
func ResendVerificationEmail(c *gin.Context) {
	var body requests.ResendVerificationEmailRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	sendVerificationEmail(body.Email)

	// Always return success regardless of internal logic
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "If the email is registered and not verified, a verification email will be sent.",
	})
}

// DemandResetPassword godoc
// @Summary      Demand Password Reset
// @Description  Send Token on email when user forgot their password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param email body requests.DemandResetPasswordRequest true "Email address to send password reset token"
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /auth/password/forgot [post]
func DemandResetPassword(c *gin.Context) {
	var body requests.DemandResetPasswordRequest

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

// ResetPassword godoc
// @Summary      Password Reset
// @Description  Reset password with Token obtained from Email and set a new password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param tokenPassword body requests.ResetPasswordRequest true "Token from email message and a new password"
// @Success      200  {object}  responses.Success
// @Failure      400  {object}  responses.Error
// @Router       /auth/password/forgot/reset [post]
func ResetPassword(c *gin.Context) {
	var body requests.ResetPasswordRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect data provided",
		})
		return
	}

	// Get token from DB
	secret, err := utils.ReadHMACSecret(config.HMACSecretPath)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Internal error",
		})
		return
	}

	hashedToken := utils.ComputeHMACToken(secret, body.Token)

	var actionToken models.ActionToken
	result := initializers.DB.Where("token_hash = ?", hashedToken).First(&actionToken)

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
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)

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

	// Mark token as used
	now := time.Now()
	actionToken.UsedAt = &now
	initializers.DB.Save(&actionToken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset successfully",
	})
}
