package handler

import (
	"errors"
	"net/http"

	"go-musthave-diploma/internal/middleware"
	"go-musthave-diploma/internal/service"

	"github.com/gin-gonic/gin"
)

type balanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *Handler) GetBalance(c *gin.Context) {
	userID, ok := middleware.UserIDFromGin(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	balance, err := h.balanceService.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balanceResponse{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	})
}

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
		case errors.Is(err, service.ErrNotEnoughBalance):
			c.Status(http.StatusPaymentRequired)
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.Status(http.StatusOK)
}
