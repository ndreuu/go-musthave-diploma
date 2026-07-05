package repository

import (
	"context"

	"go-musthave-diploma/internal/model"
)

// WithdrawalRepository defines the interface for withdrawal data access operations.
//
// Implementations handle withdrawal record creation and retrieval for tracking
// points withdrawn from user balances in the GopherMart loyalty system.
type WithdrawalRepository interface {
	// CreateWithdrawal creates a new withdrawal record for the specified user.
	// The withdrawal is recorded with the current timestamp as ProcessedAt.
	CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error

	// GetWithdrawalsByUserID retrieves all withdrawals for a specific user.
	// Returns an empty slice if the user has no withdrawal history.
	GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]model.Withdrawal, error)
}
