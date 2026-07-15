package service

import (
	"context"
	"errors"
	"testing"

	"go-musthave-diploma/internal/repository"
)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	authService := NewAuthService(repo)

	user, err := authService.Register(ctx, "andrew", "password")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.ID == 0 {
		t.Fatal("Register() user ID should not be zero")
	}

	if user.Login != "andrew" {
		t.Fatalf("Register() login = %q, want %q", user.Login, "andrew")
	}

	if user.PasswordHash == "password" {
		t.Fatal("Register() should store password hash, not raw password")
	}

	loggedInUser, err := authService.Login(ctx, "andrew", "password")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if loggedInUser.ID != user.ID {
		t.Fatalf("Login() user ID = %d, want %d", loggedInUser.ID, user.ID)
	}
}

func TestAuthService_RegisterDuplicateLogin(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	authService := NewAuthService(repo)

	_, err := authService.Register(ctx, "andrew", "password")
	if err != nil {
		t.Fatalf("Register() first call error = %v", err)
	}

	_, err = authService.Register(ctx, "andrew", "another-password")
	if !errors.Is(err, ErrLoginAlreadyTaken) {
		t.Fatalf("Register() error = %v, want %v", err, ErrLoginAlreadyTaken)
	}
}

func TestAuthService_LoginInvalidCredentials(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewMemoryRepository()
	authService := NewAuthService(repo)

	_, err := authService.Login(ctx, "missing", "password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}

	_, err = authService.Register(ctx, "andrew", "password")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, err = authService.Login(ctx, "andrew", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}
