package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-musthave-diploma/internal/service"
)

type AuthService interface {
	Register(ctx gin.Context, login string, password string)
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Login == "" || req.Password == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrLoginAlreadyTaken) {
			c.Status(http.StatusConflict)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.Generate(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, authResponse{Token: token})
}

func (h *Handler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Login == "" || req.Password == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := h.tokenService.Generate(user.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, authResponse{Token: token})
}
