package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-musthave-diploma/internal/service"
)

const UserIDKey = "userID"

func Auth(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, prefix)

		userID, err := tokenService.Parse(token)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func UserIDFromGin(c *gin.Context) (int64, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}
