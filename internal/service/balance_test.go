package service

import (
	"context"
	"errors"
	"testing"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
)

func TestBalanceService_GetBalanceEmpty(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	balanceService := NewBalanceService(repo, repo)

	balance, err := balanceService.GetBalance(ctx, 1)
	if err != nil {
		t.Fatalf("GetBalance() error = %v", err)
	}

	if balance.Current != 0 {
		t.Fatalf("Current = %v, want 0", balance.Current)
	}

	if balance.Withdrawn != 0 {
		t.Fatalf("Withdrawn = %v, want 0", balance.Withdrawn)
	}
}

func TestBalanceService_WithdrawNotEnoughBalance(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	balanceService := NewBalanceService(repo, repo)

	err := balanceService.Withdraw(ctx, 1, "9278923470", 100)
	if !errors.Is(err, ErrNotEnoughBalance) {
		t.Fatalf("Withdraw() error = %v, want %v", err, ErrNotEnoughBalance)
	}
}

func TestBalanceService_WithdrawInvalidOrder(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	balanceService := NewBalanceService(repo, repo)

	err := balanceService.Withdraw(ctx, 1, "9278923471", 100)
	if !errors.Is(err, ErrInvalidOrderNumber) {
		t.Fatalf("Withdraw() error = %v, want %v", err, ErrInvalidOrderNumber)
	}
}

func TestBalanceService_WithdrawSuccess(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	balanceService := NewBalanceService(repo, repo)
	orderService := NewOrderService(repo)

	err := orderService.UploadOrder(ctx, 1, "9278923470")
	if err != nil {
		t.Fatalf("UploadOrder() error = %v", err)
	}

	accrual := 500.0
	err = repo.UpdateOrderAccrual(ctx, "9278923470", model.OrderStatusProcessed, &accrual)
	if err != nil {
		t.Fatalf("UpdateOrderAccrual() error = %v", err)
	}

	err = balanceService.Withdraw(ctx, 1, "2377225624", 200)
	if err != nil {
		t.Fatalf("Withdraw() error = %v", err)
	}

	balance, err := balanceService.GetBalance(ctx, 1)
	if err != nil {
		t.Fatalf("GetBalance() error = %v", err)
	}

	if balance.Current != 300 {
		t.Fatalf("Current = %v, want 300", balance.Current)
	}

	if balance.Withdrawn != 200 {
		t.Fatalf("Withdrawn = %v, want 200", balance.Withdrawn)
	}

	withdrawals, err := balanceService.GetWithdrawals(ctx, 1)
	if err != nil {
		t.Fatalf("GetWithdrawals() error = %v", err)
	}

	if len(withdrawals) != 1 {
		t.Fatalf("withdrawals len = %d, want 1", len(withdrawals))
	}

	if withdrawals[0].Order != "2377225624" {
		t.Fatalf("withdrawal order = %q, want %q", withdrawals[0].Order, "2377225624")
	}
}
