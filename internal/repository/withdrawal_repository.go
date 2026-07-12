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
	Withdraw(ctx context.Context, userID int64, order string, sum float64) error

	GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]model.Withdrawal, error)
}
