package service

import (
	"context"
	"errors"
	"sort"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
)

var (
	ErrInvalidOrderNumber     = errors.New("invalid order number")
	ErrOrderUploadedByAnother = errors.New("order uploaded by another user")
	ErrOrderAlreadyUploaded   = errors.New("order already uploaded")
)

type OrderService struct {
	orders repository.OrderRepository
}

func NewOrderService(orders repository.OrderRepository) *OrderService {
	return &OrderService{orders: orders}
}

func (s *OrderService) UploadOrder(ctx context.Context, userID int64, number string) error {
	if !IsValidLuhn(number) {
		return ErrInvalidOrderNumber
	}

	order, err := s.orders.GetOrderByNumber(ctx, number)
	if err == nil {
		if order.UserID == userID {
			return ErrOrderAlreadyUploaded
		}
		return ErrOrderUploadedByAnother
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	return s.orders.CreateOrder(ctx, userID, number)
}

func (s *OrderService) GetOrders(ctx context.Context, userID int64) ([]model.Order, error) {
	orders, err := s.orders.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].UploadedAt.After(orders[j].UploadedAt)
	})

	return orders, nil
}
