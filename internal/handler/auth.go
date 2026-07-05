package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-musthave-diploma/internal/service"
)

// AuthService defines the interface for user authentication operations.
//
// Implementations handle user registration and login, returning
// user data on success or appropriate errors on failure.
type AuthService interface {
	// Register creates a new user account with the given login and password.
	// Returns the created user or an error if registration fails.
	Register(ctx gin.Context, login string, password string)
}

// authRequest represents the JSON payload for authentication endpoints.
//
// Used for both registration and login requests. Both fields are required
// and must be non-empty for the request to be valid.
type authRequest struct {
	// Login is the user's chosen username (must be unique).
	Login string `json:"login"`

	// Password is the user's password (stored as a hash).
	Password string `json:"password"`
}

// authResponse represents the JSON response for successful authentication.
//
// Contains the JWT token to be used in subsequent authenticated requests.
type authResponse struct {
	// Token is the JWT authentication token.
	// Include this in the Authorization header as "Bearer <token>".
	Token string `json:"token"`
}

// Register handles user registration requests.
//
// Expected request body:
//
//	{"login": "username", "password": "secret"}
//
// Response codes:
//   - 200 OK: Registration successful, returns JWT token
//   - 400 Bad Request: Invalid or missing request body
//   - 409 Conflict: Login already taken
//   - 500 Internal Server Error: Server error (token generation failed)
//
// On success, the JWT token is returned in the response body and also
// set in the Authorization header as "Bearer <token>".
func (h *Handler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Login == "" || req.Password == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrLoginAlreadyTaken) {
			c.Status(http.StatusConflict)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.Generate(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, authResponse{Token: token})
}

// Login handles user authentication requests.
//
// Expected request body:
//
//	{"login": "username", "password": "secret"}
//
// Response codes:
//   - 200 OK: Login successful, returns JWT token
//   - 400 Bad Request: Invalid or missing request body
//   - 401 Unauthorized: Invalid credentials (wrong login or password)
//   - 500 Internal Server Error: Server error (token generation failed)
//
// On success, the JWT token is returned in the response body and also
// set in the Authorization header as "Bearer <token>".
func (h *Handler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Login == "" || req.Password == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.Generate(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, authResponse{Token: token})
}
