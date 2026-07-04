package repository

import (
	"context"
	"sync"
	"time"

	"go-musthave-diploma/internal/model"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	nextID int64

	users       map[string]*model.User
	orders      map[string]*model.Order
	withdrawals []model.Withdrawal
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		users:  make(map[string]*model.User),
		orders: make(map[string]*model.Order),
	}
}

func (r *MemoryRepository) CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.withdrawals = append(r.withdrawals, model.Withdrawal{
		Order:       order,
		UserID:      userID,
		Sum:         sum,
		ProcessedAt: time.Now(),
	})

	return nil
}

func (r *MemoryRepository) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Withdrawal

	for _, withdrawal := range r.withdrawals {
		if withdrawal.UserID == userID {
			result = append(result, withdrawal)
		}
	}

	return result, nil
}

func (r *MemoryRepository) CreateUser(ctx context.Context, login string, passwordHash string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[login]; exists {
		return nil, ErrAlreadyExists
	}

	user := &model.User{
		ID:           r.nextID,
		Login:        login,
		PasswordHash: passwordHash,
	}

	r.nextID++
	r.users[login] = user

	return user, nil
}

func (r *MemoryRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[login]
	if !exists {
		return nil, ErrNotFound
	}

	return user, nil
}

func (r *MemoryRepository) CreateOrder(ctx context.Context, userID int64, number string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[number] = &model.Order{
		Number:     number,
		UserID:     userID,
		Status:     model.OrderStatusNew,
		UploadedAt: time.Now(),
	}

	return nil
}

func (r *MemoryRepository) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[number]
	if !ok {
		return nil, ErrNotFound
	}

	return order, nil
}

func (r *MemoryRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Order

	for _, order := range r.orders {
		if order.UserID == userID {
			result = append(result, *order)
		}
	}

	return result, nil
}

func (r *MemoryRepository) GetOrdersForAccrual(ctx context.Context, limit int) ([]model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]model.Order, 0, limit)

	for _, order := range r.orders {
		if order.Status == model.OrderStatusNew ||
			order.Status == model.OrderStatusProcessing {
			orders = append(orders, *order)

			if len(orders) >= limit {
				break
			}
		}
	}

	return orders, nil
}

func (r *MemoryRepository) UpdateOrderAccrual(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[number]
	if !ok {
		return ErrNotFound
	}

	order.Status = status
	order.Accrual = accrual

	r.orders[number] = order

	return nil
}
