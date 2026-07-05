package repository

import (
	"context"

	"go-musthave-diploma/internal/model"
)

// OrderRepository defines the interface for order data access operations.
//
// Implementations handle order management including creation, retrieval,
// and lookup operations for the GopherMart loyalty system.
type OrderRepository interface {
	// CreateOrder creates a new order for the specified user.
	// The order is created with status NEW. Returns ErrAlreadyExists
	// if an order with the same number already exists.
	CreateOrder(ctx context.Context, userID int64, number string) error

	// GetOrderByNumber retrieves an order by its number.
	// Returns the order if found, or ErrNotFound if no order with that number exists.
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)

	// GetOrdersByUserID retrieves all orders for a specific user.
	// Returns an empty slice if the user has no orders.
	GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error)
}
