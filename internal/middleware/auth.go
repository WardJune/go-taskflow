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
		tokenString := ""

		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 || parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			response.Unauthorized(c, errors.New("authorization required"))
			return
		}

		claims, err := token.Verify(tokenString, jwtSecret)

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
