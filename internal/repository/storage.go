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
}

type AccrualRepository interface {
	GetOrdersForAccrual(ctx context.Context, limit int) ([]model.Order, error)
	UpdateOrderAccrual(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error
}
