package middleware

import (
	"litstore/api/config"
	"litstore/api/initializers"
	"litstore/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
			cookieCsrfToken, err := c.Cookie("csrf-token")

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
		tokenString, err := c.Cookie("jwt_access_token")

		if err != nil {
			c.JSON(401, gin.H{"error": "No token passed", "success": false})
			c.Abort()
			return
		}

		token, err := utils.ParseJWT(tokenString)

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token", "success": false})
			c.Abort()
			return
		}

		userID, err := token.Claims.GetSubject()

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token", "success": false})
			c.Abort()
			return
		}

		// check for permissions by userID
		var count int

		// Perform a join query to count how many times the permission exists for the user
		result := initializers.DB.Table("users").
			Select("count(*)").
			Joins("JOIN user_permissions ON user_permissions.user_id = users.id").
			Joins("JOIN permissions ON permissions.id = user_permissions.permission_id").
			Where("users.id = ? AND permissions.name = ?", userID, requiredPermission).
			Group("users.id").
			Scan(&count)

		if result.Error != nil || count <= 0 {
			c.JSON(401, gin.H{"error": "Insufficient permissions", "success": false})
			c.Abort()
			return
		}

		c.Next()
	}
}
