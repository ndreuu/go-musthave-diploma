// Package middleware provides HTTP middleware functions for the Gin web framework.
//
// It includes authentication middleware that validates JWT tokens and extracts
// user identity from requests, making it available to downstream handlers.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go-musthave-diploma/internal/service"
)

// UserIDKey is the context key used to store the authenticated user's ID.
//
// Handlers can retrieve the user ID using middleware.UserIDFromGin(c) or
// directly via c.Get(middleware.UserIDKey). The value is stored as int64.
const UserIDKey = "userID"

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
// On successful validation, the user ID is stored in the Gin context
// under the key UserIDKey and can be retrieved by downstream handlers.
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

// UserIDFromGin extracts the authenticated user's ID from the Gin context.
//
// This helper function safely retrieves the user ID that was set by the
// Auth middleware. It returns the user ID and a boolean indicating whether
// the value was present and of the correct type.
//
// Parameters:
//   - c: Gin context containing the request state
//
// Returns:
//   - userID: The authenticated user's ID (0 if not found or invalid type)
//   - exists: True if the user ID was found and is valid, false otherwise
//
// Example usage:
//
//	func handler(c *gin.Context) {
//	    userID, ok := middleware.UserIDFromGin(c)
//	    if !ok {
//	        c.Status(http.StatusUnauthorized)
//	        return
//	    }
//	    // userID is now available for use
//	}
func UserIDFromGin(c *gin.Context) (int64, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int64)
	return userID, ok
}
