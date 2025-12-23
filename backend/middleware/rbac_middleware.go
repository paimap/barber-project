package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "role not found in token",
			})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient role",
		})
		c.Abort()
	}
}
