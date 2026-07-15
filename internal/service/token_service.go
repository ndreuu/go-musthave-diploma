package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// tokenTTL is the time-to-live for generated JWT tokens.
// Tokens expire after this duration from the time of issuance.
const tokenTTL = 24 * time.Hour

// ErrInvalidToken is returned when a JWT token cannot be parsed or is invalid.
//
// This error may occur when:
//   - The token signature doesn't match
//   - The token has expired
//   - The token format is malformed
//   - The subject claim cannot be parsed as a user ID
var ErrInvalidToken = errors.New("invalid token")

// TokenService handles JWT token generation and validation.
//
// It uses HS256 (HMAC with SHA-256) for signing tokens. The service
// generates tokens with a 24-hour expiration time and includes standard
// JWT claims (subject, issued at, expires at).
type TokenService struct {
	// secret is the HMAC signing key for JWT operations.
	secret []byte
}

// NewTokenService creates a new TokenService instance.
//
// Parameters:
//   - secret: Secret key for JWT signing and verification (should be kept secure)
//
// Returns:
//   - *TokenService: New token service instance
//
// Example usage:
//
//	tokenService := service.NewTokenService("my-secret-key")
//	token, err := tokenService.Generate(userID)
func NewTokenService(secret string) *TokenService {
	return &TokenService{
		secret: []byte(secret),
	}
}

// Generate creates a new JWT token for the specified user ID.
//
// The token includes the following claims:
//   - sub (subject): The user ID as a string
//   - iat (issued at): Current timestamp
//   - exp (expires at): Current timestamp + tokenTTL (24 hours)
//
// The token is signed using HS256 with the configured secret key.
//
// Parameters:
//   - userID: Unique identifier of the user to generate token for
//
// Returns:
//   - string: The encoded JWT token string
//   - error: Error if token generation fails
//
// Example usage:
//
//	token, err := tokenService.Generate(123)
//	if err != nil {
//	    return err
//	}
//	// Use token in Authorization header: "Bearer " + token
func (s *TokenService) Generate(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

// Parse validates and extracts the user ID from a JWT token.
//
// The function performs the following validation:
//   - Verifies the token signature using the secret key
//   - Checks token expiration
//   - Parses the subject claim as an int64 user ID
//
// Parameters:
//   - tokenString: The JWT token string to parse (from Authorization header)
//
// Returns:
//   - int64: The user ID extracted from the token's subject claim
//   - error: ErrInvalidToken if token is invalid or cannot be parsed
//
// Example usage:
//
//	userID, err := tokenService.Parse(tokenString)
//	if err != nil {
//	    if errors.Is(err, service.ErrInvalidToken) {
//	        // Handle invalid token (return 401 Unauthorized)
//	    }
//	}
func (s *TokenService) Parse(tokenString string) (int64, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, ErrInvalidToken
	}

	return userID, nil
}
