package repository

import (
	"context"

	"go-musthave-diploma/internal/model"
)

// UserRepository defines the interface for user data access operations.
//
// Implementations handle user registration and authentication by providing
// methods to create new users and retrieve existing users by their login.
type UserRepository interface {
	// CreateUser creates a new user with the given login and password hash.
	// Returns the created user with assigned ID, or ErrAlreadyExists if
	// a user with the same login already exists.
	CreateUser(ctx context.Context, login string, passwordHash string) (*model.User, error)

	// GetUserByLogin retrieves a user by their login username.
	// Returns the user if found, or ErrNotFound if no user with that login exists.
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}
