package handler

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"go-musthave-diploma/internal/middleware"
	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/service"

	"github.com/gin-gonic/gin"
)

// UploadOrder handles requests to upload a new order number for processing.
//
// Requires authentication via JWT token in the Authorization header.
// The request body should contain only the order number as plain text.
//
// Request body: Plain text order number (e.g., "12345678901")
//
// Response codes:
//   - 202 Accepted: Order uploaded successfully, awaiting processing
//   - 200 OK: Order was already uploaded by this user (no action taken)
//   - 400 Bad Request: Empty or malformed request body
//   - 401 Unauthorized: Missing or invalid authentication token
//   - 409 Conflict: Order was uploaded by another user
//   - 422 Unprocessable Entity: Invalid order number format (failed Luhn check)
//   - 500 Internal Server Error: Database or service error
func (h *Handler) UploadOrder(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	number := strings.TrimSpace(string(body))
	if number == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.orderService.UploadOrder(c.Request.Context(), userID, number)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrderNumber):
			c.Status(http.StatusUnprocessableEntity)
		case errors.Is(err, service.ErrOrderAlreadyUploaded):
			c.Status(http.StatusOK)
		case errors.Is(err, service.ErrOrderUploadedByAnother):
			c.Status(http.StatusConflict)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusAccepted)
}

// orderResponse represents a single order in the JSON response.
//
// Contains the order number, current status, accrued points (if any),
// and the timestamp when the order was uploaded.
type orderResponse struct {
	// Number is the unique order identifier.
	Number string `json:"number"`

	// Status is the current processing status of the order.
	Status model.OrderStatus `json:"status"`

	// Accrual is the amount of points accrued for this order (nil if not yet processed).
	Accrual *float64 `json:"accrual,omitempty"`

	// UploadedAt is the timestamp when the order was uploaded by the user.
	UploadedAt time.Time `json:"uploaded_at"`
}

// GetOrders handles requests to retrieve the user's uploaded orders.
//
// Requires authentication via JWT token in the Authorization header.
// Returns all orders uploaded by the authenticated user with their
// current processing status and accrued points.
//
// Response codes:
//   - 200 OK: Returns list of orders
//   - 204 No Content: User has no uploaded orders
//   - 401 Unauthorized: Missing or invalid authentication token
//   - 500 Internal Server Error: Database or service error
//
// Response body (on success): Array of order objects
//
//	[
//	  {"number": "12345678901", "status": "PROCESSING", "uploaded_at": "2024-01-01T00:00:00Z"},
//	  {"number": "98765432109", "status": "PROCESSED", "accrual": 500.5, "uploaded_at": "2024-01-02T00:00:00Z"}
//	]
func (h *Handler) GetOrders(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetOrders(c.Request.Context(), userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resp := make([]orderResponse, 0, len(orders))
	for _, order := range orders {
		resp = append(resp, orderResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}
