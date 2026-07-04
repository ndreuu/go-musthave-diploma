package repository

import (
	"context"

	"go-musthave-diploma/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, userID int64, number string) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error)
}
