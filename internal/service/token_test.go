package service

import "testing"

func TestTokenService_GenerateAndParse(t *testing.T) {
	tokenService := NewTokenService("secret")

	token, err := tokenService.Generate(42)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	userID, err := tokenService.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if userID != 42 {
		t.Fatalf("Parse() userID = %d, want %d", userID, 42)
	}
}

func TestTokenService_ParseInvalidToken(t *testing.T) {
	tokenService := NewTokenService("secret")

	_, err := tokenService.Parse("invalid-token")
	if err == nil {
		t.Fatal("Parse() expected error")
	}
}
