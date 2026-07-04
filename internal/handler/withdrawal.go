package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-musthave-diploma/internal/middleware"
)

type withdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.balanceService.GetWithdrawals(c.Request.Context(), userID)
	if err != nil {
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
