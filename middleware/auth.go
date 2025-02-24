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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Cannot connect to redis", "success": false})
			c.Abort()
			return
		}

		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Blacklisted token", "success": false})
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
