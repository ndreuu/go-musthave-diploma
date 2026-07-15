package repository

import (
	"context"
	"go-musthave-diploma/internal/model"
)

// Storage is the composite interface that combines all repository interfaces.
//
// It represents the complete data access layer for the application, providing
// access to users, orders, withdrawals, accrual operations, and balance queries.
// Implementations include MemoryRepository for testing and PostgresRepository
// for production use.
type Storage interface {
	UserRepository
	OrderRepository
	WithdrawalRepository
	AccrualRepository
	BalanceRepository
}

var _ Storage = (*PostgresRepository)(nil)
var _ Storage = (*MemoryRepository)(nil)

// AccrualRepository defines the interface for accrual system operations.
//
// Implementations handle fetching orders that need accrual processing
// and updating order status and accrual amounts after processing.
type AccrualRepository interface {
	// GetOrdersForAccrual retrieves orders that need accrual processing.
	// Returns orders with status NEW or PROCESSING, up to the specified limit.
	GetOrdersForAccrual(ctx context.Context, limit int) ([]model.Order, error)

	// UpdateOrderAccrual updates the status and accrual amount for an order.
	// Used after the accrual system has processed an order.
	UpdateOrderAccrual(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
}

// BalanceRepository defines the interface for balance-related queries.
//
// Implementations provide methods to calculate total accrued points
// and total withdrawn amounts for a user.
type BalanceRepository interface {
	GetBalance(ctx context.Context, userID int64) (*model.Balance, error)
}
