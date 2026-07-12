// Package middleware provides HTTP middleware functions for the Gin web framework.
//
// It includes authentication middleware that validates JWT tokens and extracts
// user identity from requests, making it available to downstream handlers.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey struct{}

// userIDKey is a private typed key used to store the authenticated user ID
// in the request context.
var userIDKey contextKey

type TokenParser interface {
	Parse(token string) (int64, error)
}

// Auth creates a Gin middleware function that validates JWT tokens from the
// Authorization header and sets the user ID in the request context.
//
// The middleware expects the Authorization header in the format:
//
//	Authorization: Bearer <token>
//
// If the header is missing, malformed, or contains an invalid token,
// the middleware responds with 401 Unauthorized and aborts the request.
//
// On successful validation, the user ID is stored in the request context
// and can be retrieved using UserIDFromGin.
//
// Parameters:
//   - tokenService: Service instance for JWT token parsing and validation
//
// Returns:
//   - gin.HandlerFunc that performs authentication on each request
//
// Example usage:
//
//	r := gin.New()
//	auth := r.Group("/api/user")
//	auth.Use(middleware.Auth(tokenService))
//	auth.GET("/profile", getProfileHandler)
func Auth(tokenService TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, prefix)
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := tokenService.Parse(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			c.Request.Context(),
			userIDKey,
			userID,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// UserIDFromContext extracts the authenticated user's ID from context.
func userIDFromContext(ctx context.Context) (int64, bool) {
	if ctx == nil {
		return 0, false
	}

	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

// UserIDFromGin extracts the authenticated user's ID from the Gin request context.
func UserIDFromGin(c *gin.Context) (int64, bool) {
	if c == nil || c.Request == nil {
		return 0, false
	}

	return userIDFromContext(c.Request.Context())
}
