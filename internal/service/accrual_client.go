package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-musthave-diploma/internal/model"
)

var ErrAccrualNoContent = errors.New("accrual: no content")
var ErrAccrualTooManyRequests = errors.New("accrual: too many requests")

type AccrualTooManyRequestsError struct {
	RetryAfter time.Duration
}

func (e *AccrualTooManyRequestsError) Error() string {
	return ErrAccrualTooManyRequests.Error()
}

type AccrualClient struct {
	baseURL string
	client  *http.Client
}

type AccrualOrderResponse struct {
	Order   string            `json:"order"`
	Status  model.OrderStatus `json:"status"`
	Accrual *float64          `json:"accrual,omitempty"`
}

func NewAccrualClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *AccrualClient) GetOrder(ctx context.Context, number string) (*AccrualOrderResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var result AccrualOrderResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}
		return &result, nil

	case http.StatusNoContent:
		return nil, ErrAccrualNoContent

	case http.StatusTooManyRequests:
		retryAfter := time.Second

		if value := resp.Header.Get("Retry-After"); value != "" {
			seconds, err := strconv.Atoi(value)
			if err == nil && seconds > 0 {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}

		return nil, &AccrualTooManyRequestsError{
			RetryAfter: retryAfter,
		}

	default:
		return nil, fmt.Errorf("accrual unexpected status: %d", resp.StatusCode)
	}
}
