package middleware

import (
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
		// verify token, check in db permissions

		c.Next()
	}
}
