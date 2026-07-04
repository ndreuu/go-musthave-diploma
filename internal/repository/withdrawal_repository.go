package repository

import (
	"context"

	"go-musthave-diploma/internal/model"
)

type WithdrawalRepository interface {
	CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error
	GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]model.Withdrawal, error)
}
