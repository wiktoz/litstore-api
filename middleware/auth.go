package middleware

import (
	"litstore/api/config"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		writeMethods := map[string]bool{
			http.MethodPost:   true,
			http.MethodPut:    true,
			http.MethodPatch:  true,
			http.MethodDelete: true,
		}

		if _, ok := writeMethods[c.Request.Method]; ok {
			csrfToken := c.GetHeader("X-CSRF-Token")
			cookieCsrfToken, err := c.Cookie(config.CsrfName)

			if err != nil || csrfToken != cookieCsrfToken {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch or missing"})
				return
			}
		}

		c.Next()
	}
}

func Authorization(requiredPermission config.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// check if token is not revoked
		blacklisted, err := utils.IsBlacklisted(c, initializers.RDB, tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Cannot connect to redis",
			})

			c.Abort()
			return
		}

		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Blacklisted token",
			})

			c.Abort()
			return
		}

		// Save userID in context for further use
		userIDObj, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Failed to parse token",
			})

			c.Abort()
			return
		}
		c.Set("userID", userIDObj)

		// Get user details
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

		// check if user is blocked
		if user.Blocked {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "User is blocked by administrator",
			})

			c.Abort()
			return
		}

		// check if user has confirmed email
		if !user.Confirmed {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "User has not confirmed an email",
			})

			c.Abort()
			return
		}

		// check for permissions by userID

		if requiredPermission != "" {
			var exists bool

			query := `
			SELECT EXISTS (
				SELECT 1
				FROM users_permissions
				JOIN permissions ON users_permissions.permission_id = permissions.id
				WHERE users_permissions.user_id = ?
				AND permissions.name = ?
			)
			OR EXISTS (
				SELECT 1
				FROM users_roles
				JOIN roles ON users_roles.role_id = roles.id
				JOIN roles_permissions ON roles_permissions.role_id = roles.id
				JOIN permissions ON permissions.id = roles_permissions.permission_id
				WHERE users_roles.user_id = ?
				AND permissions.name = ?
			);`

			result := initializers.DB.Raw(query, userID, requiredPermission, userID, requiredPermission).Scan(&exists)

			if result.Error != nil || !exists {
				c.JSON(401, gin.H{"error": "Insufficient permissions", "success": false})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
