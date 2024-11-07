package middleware

import (
	"fmt"
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

func Authorization(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_access_token")

		if err != nil {
			c.JSON(401, gin.H{"error": "No token", "success": false})
			c.Abort()
			return
		}

		token, err := utils.ParseJWT(tokenString)

		if err != nil {
			c.JSON(401, gin.H{"error": err, "success": false})
			c.Abort()
			return
		}

		userID, err := token.Claims.GetSubject()

		if err != nil {
			c.JSON(401, gin.H{"error": "No subject", "success": false})
			c.Abort()
			return
		}

		// check for permissions by userID

		fmt.Println(userID)

		c.Next()
	}
}
