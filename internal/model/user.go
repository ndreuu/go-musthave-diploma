// Package model defines the core data structures used throughout the application.
//
// It includes domain models for users, orders, and withdrawals that represent
// the business entities of the GopherMart loyalty system.
package model

// User represents a registered user in the GopherMart system.
//
// Users can upload order numbers to earn loyalty points, check their balance,
// and withdraw points for various rewards.
type User struct {
	// ID is the unique identifier for the user (auto-generated).
	ID int64

	// Login is the user's unique username used for authentication.
	Login string

	// PasswordHash is the bcrypt hash of the user's password.
	// The original password is never stored.
	PasswordHash string
}
