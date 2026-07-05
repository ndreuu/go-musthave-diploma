package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-musthave-diploma/internal/repository"
	"go-musthave-diploma/internal/service"
	"go.uber.org/zap"
)

func TestAccrualPoller_Start(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order":"9278923470","status":"PROCESSED","accrual":500}`))
	}))
	defer server.Close()

	repo := repository.NewMemoryRepository()
	client := service.NewAccrualClient(server.URL)
	logger, _ := zap.NewDevelopment()

	poller := service.NewAccrualPoller(repo, client, logger, 100*time.Millisecond, 10)

	_ = repo.CreateOrder(context.Background(), 1, "9278923470")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	poller.Start(ctx)
}

func TestAccrualPoller_NoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	repo := repository.NewMemoryRepository()
	client := service.NewAccrualClient(server.URL)
	logger, _ := zap.NewDevelopment()

	poller := service.NewAccrualPoller(repo, client, logger, 100*time.Millisecond, 10)

	_ = repo.CreateOrder(context.Background(), 1, "9278923470")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	poller.Start(ctx)
}

func TestAccrualPoller_RateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	repo := repository.NewMemoryRepository()
	client := service.NewAccrualClient(server.URL)
	logger, _ := zap.NewDevelopment()

	poller := service.NewAccrualPoller(repo, client, logger, 50*time.Millisecond, 10)

	_ = repo.CreateOrder(context.Background(), 1, "9278923470")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	poller.Start(ctx)
}

func TestNewAccrualPoller(t *testing.T) {
	repo := repository.NewMemoryRepository()
	client := service.NewAccrualClient("http://test")
	logger, _ := zap.NewDevelopment()

	poller := service.NewAccrualPoller(repo, client, logger, time.Second, 10)
	if poller == nil {
		t.Fatal("NewAccrualPoller() returned nil")
	}
}
