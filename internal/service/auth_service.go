package service

import (
	"context"
	"errors"

	"go-musthave-diploma/internal/model"
	"go-musthave-diploma/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginAlreadyTaken  = errors.New("login already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) *AuthService {
	return &AuthService{
		users: users,
	}
}

func (s *AuthService) Register(ctx context.Context, login string, password string) (*model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.users.CreateUser(ctx, login, string(passwordHash))
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return nil, ErrLoginAlreadyTaken
		}

		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, login string, password string) (*model.User, error) {
	user, err := s.users.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
