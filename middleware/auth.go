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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cannot connect to redis", "success": false})
			c.Abort()
			return
		}

		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "success": false})
			c.Abort()
			return
		}

		// check for permissions by userID
		var count int

		query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_permissions
			JOIN permissions ON user_permissions.permission_id = permissions.id
			WHERE user_permissions.user_id = (SELECT id FROM users WHERE username = 'john_doe') 
			AND permissions.name = 'product_read'
		)
		OR EXISTS (
			SELECT 1
			FROM user_roles
			JOIN roles ON user_roles.role_id = roles.id
			JOIN role_permissions ON role_permissions.role_id = roles.id
			JOIN permissions ON permissions.id = role_permissions.permission_id
			WHERE user_roles.user_id = (SELECT id FROM users WHERE username = 'john_doe') 
			AND permissions.name = 'product_read'
		);`

		result := initializers.DB.Raw(query, userID, requiredPermission, userID, requiredPermission).Scan(&count)

		if result.Error != nil || count == 0 {
			c.JSON(401, gin.H{"error": "Insufficient permissions", "success": false})
			c.Abort()
			return
		}

		c.Next()
	}
}
