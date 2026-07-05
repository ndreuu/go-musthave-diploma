package repository

import (
	"context"
	"go-musthave-diploma/internal/model"
)

type Storage interface {
	UserRepository
	OrderRepository
	WithdrawalRepository
	AccrualRepository
	BalanceRepository
}

type AccrualRepository interface {
	GetOrdersForAccrual(ctx context.Context, limit int) ([]model.Order, error)
	UpdateOrderAccrual(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
}

type BalanceRepository interface {
	GetAccrualSumByUserID(ctx context.Context, userID int64) (float64, error)
	GetWithdrawalSumByUserID(ctx context.Context, userID int64) (float64, error)
}
