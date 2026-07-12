// Package repository provides data access layer implementations for the GopherMart application.
//
// It includes:
//   - MemoryRepository: In-memory implementation for testing
//   - PostgresRepository: PostgreSQL-backed implementation for production
//   - Storage interface: Composite interface combining all repository interfaces
//
// The package defines common errors that may be returned by repository operations
// and provides thread-safe implementations for concurrent access.
package repository

import "errors"

var (
	// ErrNotFound is returned when a requested entity does not exist.
	//
	// This error is returned by Get* methods when the specified entity
	// (user, order, etc.) cannot be found in the repository.
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is returned when attempting to create a duplicate entity.
	//
	// This error is returned by Create* methods when trying to create
	// an entity with a unique identifier (login, order number) that
	// already exists in the repository.
	ErrAlreadyExists = errors.New("already exists")

	ErrNotEnoughBalance = errors.New("not enough balance")
)
