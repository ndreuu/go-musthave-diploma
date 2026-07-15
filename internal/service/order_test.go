package service

import (
	"context"
	"errors"
	"testing"

	"go-musthave-diploma/internal/repository"
)

func TestOrderService_UploadOrderNew(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	orderService := NewOrderService(repo)

	err := orderService.UploadOrder(ctx, 1, "9278923470")
	if err != nil {
		t.Fatalf("UploadOrder() error = %v", err)
	}
}

func TestOrderService_UploadOrderInvalidNumber(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	orderService := NewOrderService(repo)

	err := orderService.UploadOrder(ctx, 1, "9278923471")
	if !errors.Is(err, ErrInvalidOrderNumber) {
		t.Fatalf("UploadOrder() error = %v, want %v", err, ErrInvalidOrderNumber)
	}
}

func TestOrderService_UploadOrderSameUser(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	orderService := NewOrderService(repo)

	err := orderService.UploadOrder(ctx, 1, "9278923470")
	if err != nil {
		t.Fatalf("UploadOrder() first call error = %v", err)
	}

	err = orderService.UploadOrder(ctx, 1, "9278923470")
	if !errors.Is(err, ErrOrderAlreadyUploaded) {
		t.Fatalf("UploadOrder() error = %v, want %v", err, ErrOrderAlreadyUploaded)
	}
}

func TestOrderService_UploadOrderOtherUser(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	orderService := NewOrderService(repo)

	err := orderService.UploadOrder(ctx, 1, "9278923470")
	if err != nil {
		t.Fatalf("UploadOrder() first call error = %v", err)
	}

	err = orderService.UploadOrder(ctx, 2, "9278923470")
	if !errors.Is(err, ErrOrderUploadedByAnother) {
		t.Fatalf("UploadOrder() error = %v, want %v", err, ErrOrderUploadedByAnother)
	}
}

func TestOrderService_GetOrders(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	orderService := NewOrderService(repo)

	if err := orderService.UploadOrder(ctx, 1, "9278923470"); err != nil {
		t.Fatalf("UploadOrder() error = %v", err)
	}

	orders, err := orderService.GetOrders(ctx, 1)
	if err != nil {
		t.Fatalf("GetOrders() error = %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("GetOrders() len = %d, want %d", len(orders), 1)
	}

	if orders[0].Number != "9278923470" {
		t.Fatalf("GetOrders()[0].Number = %q, want %q", orders[0].Number, "9278923470")
	}
}
