package repository

import (
	"context"

	"go-musthave-diploma/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login string, passwordHash string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
}
