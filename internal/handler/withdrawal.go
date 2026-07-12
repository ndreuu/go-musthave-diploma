package handler

import (
	"net/http"
	"time"

	"go-musthave-diploma/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// withdrawalResponse represents a single withdrawal record in the JSON response.
//
// Contains the order number, withdrawal amount, and processing timestamp.
type withdrawalResponse struct {
	// Order is the order number associated with the withdrawal.
	Order string `json:"order"`

	// Sum is the amount of points withdrawn.
	Sum float64 `json:"sum"`

	// ProcessedAt is the timestamp when the withdrawal was processed.
	ProcessedAt time.Time `json:"processed_at"`
}

// GetWithdrawals handles requests to retrieve the user's withdrawal history.
//
// Requires authentication via JWT token in the Authorization header.
// Returns all withdrawals made by the authenticated user with their
// associated order numbers, amounts, and processing timestamps.
//
// Response codes:
//   - 200 OK: Returns list of withdrawals
//   - 204 No Content: User has no withdrawal history
//   - 401 Unauthorized: Missing or invalid authentication token
//   - 500 Internal Server Error: Database or service error
//
// Response body (on success): Array of withdrawal objects
//
//	[
//	  {"order": "12345678901", "sum": 100.5, "processed_at": "2024-01-01T00:00:00Z"},
//	  {"order": "98765432109", "sum": 200.0, "processed_at": "2024-01-02T00:00:00Z"}
//	]
func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
		h.log.Error(
			"failed to get user withdrawals",
			zap.Error(err),
			zap.Int64("user_id", userID),
		)

		c.Status(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	resp := make([]withdrawalResponse, 0, len(withdrawals))
	for _, w := range withdrawals {
		resp = append(resp, withdrawalResponse{
			Order:       w.Order,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt,
		})
	}

	c.JSON(http.StatusOK, resp)
}
