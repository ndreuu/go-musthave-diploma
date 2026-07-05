package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-musthave-diploma/internal/model"
)

func TestAccrualClient_GetOrderOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/orders/9278923470" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write([]byte(`{
			"order":"9278923470",
			"status":"PROCESSED",
			"accrual":500
		}`))
	}))
	defer server.Close()

	client := NewAccrualClient(server.URL)

	resp, err := client.GetOrder(context.Background(), "9278923470")
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}

	if resp.Order != "9278923470" {
		t.Fatalf("Order = %q", resp.Order)
	}

	if resp.Status != model.OrderStatusProcessed {
		t.Fatalf("Status = %q", resp.Status)
	}

	if resp.Accrual == nil {
		t.Fatal("Accrual is nil")
	}

	if *resp.Accrual != 500 {
		t.Fatalf("Accrual = %v", *resp.Accrual)
	}
}

func TestAccrualClient_NoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewAccrualClient(server.URL)

	_, err := client.GetOrder(context.Background(), "1")

	if err != ErrAccrualNoContent {
		t.Fatalf("got %v", err)
	}
}

func TestAccrualClient_TooManyRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewAccrualClient(server.URL)

	_, err := client.GetOrder(context.Background(), "1")

	var rateLimitErr *AccrualTooManyRequestsError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("expected AccrualTooManyRequestsError, got %v", err)
	}

	if rateLimitErr.RetryAfter != 60*time.Second {
		t.Fatalf("RetryAfter = %v", rateLimitErr.RetryAfter)
	}
}

