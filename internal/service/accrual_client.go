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

var (
	// ErrAccrualNoContent is returned when the accrual system returns 204 No Content.
	// This indicates the order is not found in the accrual system yet.
	ErrAccrualNoContent = errors.New("accrual: no content")

	// ErrAccrualTooManyRequests is returned when the accrual system returns 429 Too Many Requests.
	// This indicates rate limiting is in effect. See AccrualTooManyRequestsError for retry details.
	ErrAccrualTooManyRequests = errors.New("accrual: too many requests")
)

// AccrualTooManyRequestsError represents a rate limiting error from the accrual system.
//
// It includes the Retry-After duration indicating when the next request can be made.
type AccrualTooManyRequestsError struct {
	// RetryAfter is the duration to wait before making another request.
	RetryAfter time.Duration
}

// Error implements the error interface for AccrualTooManyRequestsError.
func (e *AccrualTooManyRequestsError) Error() string {
	return ErrAccrualTooManyRequests.Error()
}

// AccrualClient is an HTTP client for communicating with the external accrual system.
//
// It handles order status queries and parses responses according to the accrual API
// specification. The client includes a 3-second timeout for all requests.
type AccrualClient struct {
	// baseURL is the base URL of the accrual system API.
	baseURL string

	// client is the HTTP client with configured timeout.
	client *http.Client
}

// AccrualOrderResponse represents the response from the accrual system API.
//
// It contains the order number, processing status, and accrued points amount
// (if the order has been processed).
type AccrualOrderResponse struct {
	// Order is the order number.
	Order string `json:"order"`

	// Status is the processing status of the order.
	Status model.OrderStatus `json:"status"`

	// Accrual is the accrued points amount (nil if not yet determined).
	Accrual *float64 `json:"accrual,omitempty"`
}

// NewAccrualClient creates a new AccrualClient instance.
//
// Parameters:
//   - baseURL: Base URL of the accrual system API (e.g., "https://accrual.example.com")
//
// Returns:
//   - *AccrualClient: New accrual client with 3-second request timeout
//
// Example usage:
//
//	client := service.NewAccrualClient("http://localhost:8081")
//	resp, err := client.GetOrder(ctx, "12345678901")
func NewAccrualClient(baseURL string) *AccrualClient {
	return &AccrualClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// GetOrder queries the accrual system for the status of a specific order.
//
// The function makes an HTTP GET request to the accrual system's /api/orders/{number}
// endpoint and handles the following response codes:
//   - 200 OK: Returns the order status and accrual amount
//   - 204 No Content: Order not found in accrual system (returns ErrAccrualNoContent)
//   - 429 Too Many Requests: Rate limited (returns *AccrualTooManyRequestsError)
//   - Other: Returns an error with the status code
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - number: Order number to query
//
// Returns:
//   - *AccrualOrderResponse: Order status and accrual information
//   - error: ErrAccrualNoContent, *AccrualTooManyRequestsError, or HTTP error
//
// Example usage:
//
//	resp, err := client.GetOrder(ctx, "12345678901")
//	if err != nil {
//	    if errors.Is(err, service.ErrAccrualNoContent) {
//	        // Order not yet in accrual system, retry later
//	    }
//	    var rateLimitErr *service.AccrualTooManyRequestsError
//	    if errors.As(err, &rateLimitErr) {
//	        // Wait for rateLimitErr.RetryAfter before retrying
//	    }
//	}
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
