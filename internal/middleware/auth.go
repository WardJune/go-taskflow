package middleware

import (
	"errors"
	"strings"

	"github.com/WardJune/taskflow/pkg/response"
	"github.com/WardJune/taskflow/pkg/token"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, errors.New("authorization header required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, errors.New("invalid authorization format"))
			c.Abort()
			return
		}

		claims, err := token.Verify(parts[1], jwtSecret)

		if err != nil {
			response.Unauthorized(c, errors.New("invalid or expired token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
