package handler

import (
	"errors"
	"net/http"

	"go-musthave-diploma/internal/middleware"
	"go-musthave-diploma/internal/repository"
	"go-musthave-diploma/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// balanceResponse represents the JSON response for balance queries.
//
// Contains the current balance and total withdrawn amount for the user.
type balanceResponse struct {
	// Current is the current available balance in points.
	Current float64 `json:"current"`

	// Withdrawn is the total amount of points withdrawn by the user.
	Withdrawn float64 `json:"withdrawn"`
}

// withdrawRequest represents the JSON payload for withdrawal requests.
//
// Specifies the order number and amount to withdraw from the user's balance.
type withdrawRequest struct {
	// Order is the order number associated with the withdrawal.
	Order string `json:"order"`

	// Sum is the amount of points to withdraw (must be positive).
	Sum float64 `json:"sum"`
}

// GetBalance handles requests to retrieve the user's current balance.
//
// Requires authentication via JWT token in the Authorization header.
//
// Response codes:
//   - 200 OK: Returns current balance and withdrawn amount
//   - 401 Unauthorized: Missing or invalid authentication token
//   - 500 Internal Server Error: Database or service error
//
// Response body (on success):
//
//	{"current": 100.5, "withdrawn": 50.0}
func (h *Handler) GetBalance(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	balance, err := h.balanceService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		h.log.Error(
			"failed to get user balance",
			zap.Error(err),
			zap.Int64("user_id", userID),
		)

		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	})
}

// Withdraw handles requests to withdraw points from the user's balance.
//
// Requires authentication via JWT token in the Authorization header.
//
// Expected request body:
//
//	{"order": "12345678901", "sum": 100.5}
//
// Response codes:
//   - 200 OK: Withdrawal successful
//   - 400 Bad Request: Invalid request body (missing order or non-positive sum)
//   - 401 Unauthorized: Missing or invalid authentication token
//   - 422 Unprocessable Entity: Invalid order number format (failed Luhn check)
//   - 402 Payment Required: Insufficient balance for withdrawal
//   - 500 Internal Server Error: Database or service error
func (h *Handler) Withdraw(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	var req withdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Order == "" || req.Sum <= 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	err := h.balanceService.Withdraw(c.Request.Context(), userID, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOrderNumber):
			c.Status(http.StatusUnprocessableEntity)
		case errors.Is(err, repository.ErrNotEnoughBalance):
			c.Status(http.StatusPaymentRequired)
		default:
			h.log.Error(
				"failed to withdraw points",
				zap.Error(err),
				zap.Int64("user_id", userID),
				zap.String("order", req.Order),
				zap.Float64("sum", req.Sum),
			)

			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
