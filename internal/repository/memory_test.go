package repository

import (
	"context"
	"errors"
	"testing"

	"go-musthave-diploma/internal/model"
)

func TestMemoryRepository_CreateAndGetUser(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	user, err := repo.CreateUser(ctx, "andrew", "hash")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("ID = %d, want 1", user.ID)
	}

	got, err := repo.GetUserByLogin(ctx, "andrew")
	if err != nil {
		t.Fatalf("GetUserByLogin() error = %v", err)
	}

	if got.Login != "andrew" {
		t.Fatalf("Login = %q", got.Login)
	}
}

func TestMemoryRepository_DuplicateUser(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_, _ = repo.CreateUser(ctx, "andrew", "hash")

	_, err := repo.CreateUser(ctx, "andrew", "hash2")
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("error = %v, want %v", err, ErrAlreadyExists)
	}
}

func TestMemoryRepository_UserNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_, err := repo.GetUserByLogin(ctx, "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrNotFound)
	}
}

func TestMemoryRepository_OrderAccrual(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	if err := repo.CreateOrder(ctx, 1, "9278923470"); err != nil {
		t.Fatalf("CreateOrder() error = %v", err)
	}

	accrual := 500.0
	if err := repo.UpdateOrderAccrual(ctx, "9278923470", model.OrderStatusProcessed, &accrual); err != nil {
		t.Fatalf("UpdateOrderAccrual() error = %v", err)
	}

	order, err := repo.GetOrderByNumber(ctx, "9278923470")
	if err != nil {
		t.Fatalf("GetOrderByNumber() error = %v", err)
	}

	if order.Status != model.OrderStatusProcessed {
		t.Fatalf("status = %q", order.Status)
	}

	sum, err := repo.GetAccrualSumByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetAccrualSumByUserID() error = %v", err)
	}

	if sum != 500 {
		t.Fatalf("sum = %v, want 500", sum)
	}
}

func TestMemoryRepository_OrderNotFound(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_, err := repo.GetOrderByNumber(ctx, "nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrNotFound)
	}

	err = repo.UpdateOrderAccrual(ctx, "nonexistent", model.OrderStatusProcessed, nil)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want %v", err, ErrNotFound)
	}
}

func TestMemoryRepository_GetOrdersByUserID(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_ = repo.CreateOrder(ctx, 1, "9278923470")
	_ = repo.CreateOrder(ctx, 1, "2377225624")
	_ = repo.CreateOrder(ctx, 2, "12345678901")

	orders, err := repo.GetOrdersByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetOrdersByUserID() error = %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("len = %d, want 2", len(orders))
	}
}

func TestMemoryRepository_GetOrdersForAccrual(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_ = repo.CreateOrder(ctx, 1, "9278923470")
	_ = repo.CreateOrder(ctx, 1, "2377225624")

	orders, err := repo.GetOrdersForAccrual(ctx, 10)
	if err != nil {
		t.Fatalf("GetOrdersForAccrual() error = %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("len = %d, want 2", len(orders))
	}
}

func TestMemoryRepository_GetOrdersForAccrual_Limit(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_ = repo.CreateOrder(ctx, 1, "9278923470")
	_ = repo.CreateOrder(ctx, 1, "2377225624")
	_ = repo.CreateOrder(ctx, 1, "12345678901")

	orders, err := repo.GetOrdersForAccrual(ctx, 2)
	if err != nil {
		t.Fatalf("GetOrdersForAccrual() error = %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("len = %d, want 2", len(orders))
	}
}

func TestMemoryRepository_Withdrawals(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	if err := repo.CreateWithdrawal(ctx, 1, "2377225624", 200); err != nil {
		t.Fatalf("CreateWithdrawal() error = %v", err)
	}

	withdrawals, err := repo.GetWithdrawalsByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetWithdrawalsByUserID() error = %v", err)
	}

	if len(withdrawals) != 1 {
		t.Fatalf("len = %d, want 1", len(withdrawals))
	}

	sum, err := repo.GetWithdrawalSumByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetWithdrawalSumByUserID() error = %v", err)
	}

	if sum != 200 {
		t.Fatalf("sum = %v, want 200", sum)
	}
}

func TestMemoryRepository_EmptyWithdrawals(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	withdrawals, err := repo.GetWithdrawalsByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetWithdrawalsByUserID() error = %v", err)
	}

	if len(withdrawals) != 0 {
		t.Fatalf("len = %d, want 0", len(withdrawals))
	}

	sum, err := repo.GetWithdrawalSumByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetWithdrawalSumByUserID() error = %v", err)
	}

	if sum != 0 {
		t.Fatalf("sum = %v, want 0", sum)
	}
}

func TestMemoryRepository_EmptyAccrual(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	sum, err := repo.GetAccrualSumByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("GetAccrualSumByUserID() error = %v", err)
	}

	if sum != 0 {
		t.Fatalf("sum = %v, want 0", sum)
	}
}

func TestMemoryRepository_OrderOverwrite(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryRepository()

	_ = repo.CreateOrder(ctx, 1, "9278923470")
	_ = repo.CreateOrder(ctx, 2, "9278923470")

	order, err := repo.GetOrderByNumber(ctx, "9278923470")
	if err != nil {
		t.Fatalf("GetOrderByNumber() error = %v", err)
	}

	if order.UserID != 2 {
		t.Fatalf("UserID = %d, want 2", order.UserID)
	}
}

func TestNewMemoryRepository(t *testing.T) {
	repo := NewMemoryRepository()
	if repo == nil {
		t.Fatal("NewMemoryRepository() returned nil")
	}
}
