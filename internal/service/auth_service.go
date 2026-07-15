package service

import (
	"context"
	"errors"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrLoginAlreadyTaken is returned when attempting to register with a login
	// that is already in use by another user.
	ErrLoginAlreadyTaken = errors.New("login already taken")

	// ErrInvalidCredentials is returned when login or password is incorrect
	// during authentication attempts.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService handles user registration and authentication operations.
//
// It uses bcrypt for password hashing and verification, ensuring secure
// password storage. The service interacts with the UserRepository for
// data persistence.
type AuthService struct {
	// users is the repository for user data access.
	users repository.UserRepository
}

// NewAuthService creates a new AuthService instance.
//
// Parameters:
//   - users: UserRepository implementation for data access
//
// Returns:
//   - *AuthService: New authentication service instance
//
// Example usage:
//
//	authService := service.NewAuthService(userRepo)
//	user, err := authService.Register(ctx, "username", "password")
func NewAuthService(users repository.UserRepository) *AuthService {
	return &AuthService{
		users: users,
	}
}

// Register creates a new user account with the given login and password.
//
// The password is hashed using bcrypt with the default cost factor before
// being stored. If a user with the same login already exists, returns
// ErrLoginAlreadyTaken.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - login: Unique username for the new account
//   - password: Plain text password (will be hashed)
//
// Returns:
//   - *model.User: Pointer to the created user with assigned ID
//   - error: ErrLoginAlreadyTaken if login exists, or other error
//
// Example usage:
//
//	user, err := authService.Register(ctx, "newuser", "secretpassword")
//	if err != nil {
//	    if errors.Is(err, service.ErrLoginAlreadyTaken) {
//	        // Handle duplicate login
//	    }
//	}
func (s *AuthService) Register(ctx context.Context, login string, password string) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.users.CreateUser(ctx, login, string(passwordHash))
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, ErrLoginAlreadyTaken
		}

		return nil, err
	}

	return user, nil
}

// Login authenticates a user with the provided login and password.
//
// The function retrieves the user by login and verifies the password
// against the stored bcrypt hash. Returns ErrInvalidCredentials if
// the user doesn't exist or the password doesn't match.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - login: Username to authenticate
//   - password: Plain text password to verify
//
// Returns:
//   - *model.User: Pointer to the authenticated user
//   - error: ErrInvalidCredentials if authentication fails, or other error
//
// Example usage:
//
//	user, err := authService.Login(ctx, "username", "password")
//	if err != nil {
//	    if errors.Is(err, service.ErrInvalidCredentials) {
//	        // Handle invalid credentials
//	    }
//	}
func (s *AuthService) Login(ctx context.Context, login string, password string) (*model.User, error) {
	user, err := s.users.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
