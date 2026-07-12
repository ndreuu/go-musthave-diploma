package config

import (
	"errors"
	"flag"
	"os"
	"testing"
	"time"
)

func TestLoadFromWithoutJWTSecret(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	cfg, err := LoadFrom(fs, func(string) string {
		return ""
	}, nil)

	if !errors.Is(err, ErrJWTSecretNotSet) {
		t.Fatalf("LoadFrom() error = %v, want %v", err, ErrJWTSecretNotSet)
	}

	if cfg != nil {
		t.Fatalf("LoadFrom() config = %#v, want nil", cfg)
	}
}

func TestLoadFromDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"JWT_SECRET": "test-secret",
	}

	cfg, err := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.RunAddress != ":8080" {
		t.Fatalf("RunAddress = %q, want %q", cfg.RunAddress, ":8080")
	}

	if cfg.JWTSecret != "test-secret" {
		t.Fatalf("JWTSecret = %q, want %q", cfg.JWTSecret, "test-secret")
	}

	if cfg.AccrualPollInterval != time.Second {
		t.Fatalf(
			"AccrualPollInterval = %v, want %v",
			cfg.AccrualPollInterval,
			time.Second,
		)
	}

	if cfg.AccrualBatchSize != 10 {
		t.Fatalf("AccrualBatchSize = %d, want %d", cfg.AccrualBatchSize, 10)
	}
}

func TestLoadFromEnv(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"RUN_ADDRESS":            ":9090",
		"DATABASE_URI":           "postgres://db",
		"ACCRUAL_SYSTEM_ADDRESS": "http://accrual",
		"JWT_SECRET":             "secret",
		"ACCRUAL_POLL_INTERVAL":  "5s",
		"ACCRUAL_BATCH_SIZE":     "100",
	}

	cfg, err := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.RunAddress != ":9090" {
		t.Fatalf("RunAddress = %q, want %q", cfg.RunAddress, ":9090")
	}

	if cfg.DatabaseURI != "postgres://db" {
		t.Fatalf(
			"DatabaseURI = %q, want %q",
			cfg.DatabaseURI,
			"postgres://db",
		)
	}

	if cfg.AccrualSystemAddress != "http://accrual" {
		t.Fatalf(
			"AccrualSystemAddress = %q, want %q",
			cfg.AccrualSystemAddress,
			"http://accrual",
		)
	}

	if cfg.JWTSecret != "secret" {
		t.Fatalf("JWTSecret = %q, want %q", cfg.JWTSecret, "secret")
	}

	if cfg.AccrualPollInterval != 5*time.Second {
		t.Fatalf(
			"AccrualPollInterval = %v, want %v",
			cfg.AccrualPollInterval,
			5*time.Second,
		)
	}

	if cfg.AccrualBatchSize != 100 {
		t.Fatalf("AccrualBatchSize = %d, want %d", cfg.AccrualBatchSize, 100)
	}
}

func TestLoadFromInvalidDuration(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"JWT_SECRET":             "test-secret",
		"ACCRUAL_POLL_INTERVAL": "invalid",
	}

	cfg, err := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.AccrualPollInterval != time.Second {
		t.Fatalf(
			"AccrualPollInterval = %v, want %v",
			cfg.AccrualPollInterval,
			time.Second,
		)
	}
}

func TestLoadFromInvalidBatchSize(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"JWT_SECRET":         "test-secret",
		"ACCRUAL_BATCH_SIZE": "invalid",
	}

	cfg, err := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.AccrualBatchSize != 10 {
		t.Fatalf("AccrualBatchSize = %d, want %d", cfg.AccrualBatchSize, 10)
	}
}

func TestLoadFromNegativeBatchSize(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"JWT_SECRET":         "test-secret",
		"ACCRUAL_BATCH_SIZE": "-5",
	}

	cfg, err := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.AccrualBatchSize != 10 {
		t.Fatalf("AccrualBatchSize = %d, want %d", cfg.AccrualBatchSize, 10)
	}
}

func TestLoadFromFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	args := []string{
		"-a", ":7070",
		"-d", "postgres://flag",
		"-r", "http://flag",
		"-s", "flag-secret",
	}

	cfg, err := LoadFrom(fs, func(string) string {
		return ""
	}, args)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.RunAddress != ":7070" {
		t.Fatalf("RunAddress = %q, want %q", cfg.RunAddress, ":7070")
	}

	if cfg.DatabaseURI != "postgres://flag" {
		t.Fatalf(
			"DatabaseURI = %q, want %q",
			cfg.DatabaseURI,
			"postgres://flag",
		)
	}

	if cfg.AccrualSystemAddress != "http://flag" {
		t.Fatalf(
			"AccrualSystemAddress = %q, want %q",
			cfg.AccrualSystemAddress,
			"http://flag",
		)
	}

	if cfg.JWTSecret != "flag-secret" {
		t.Fatalf("JWTSecret = %q, want %q", cfg.JWTSecret, "flag-secret")
	}
}

func TestLoadFromFlagsOverrideEnv(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"RUN_ADDRESS": ":9090",
		"JWT_SECRET":  "env-secret",
	}

	args := []string{
		"-a", ":7070",
		"-s", "flag-secret",
	}

	cfg, err := LoadFrom(fs, func(key string) string {
		return env[key]
	}, args)
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}

	if cfg.RunAddress != ":7070" {
		t.Fatalf("RunAddress = %q, want %q", cfg.RunAddress, ":7070")
	}

	if cfg.JWTSecret != "flag-secret" {
		t.Fatalf("JWTSecret = %q, want %q", cfg.JWTSecret, "flag-secret")
	}
}

func TestLoadFromFlagParseError(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	cfg, err := LoadFrom(
		fs,
		func(string) string {
			return ""
		},
		[]string{"-unknown"},
	)

	if err == nil {
		t.Fatal("LoadFrom() error = nil, want flag parsing error")
	}

	if cfg != nil {
		t.Fatalf("LoadFrom() config = %#v, want nil", cfg)
	}
}

func TestLoad(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	oldArgs := os.Args
	os.Args = []string{"gophermart"}
	t.Cleanup(func() {
		os.Args = oldArgs
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.JWTSecret != "test-secret" {
		t.Fatalf("JWTSecret = %q, want %q", cfg.JWTSecret, "test-secret")
	}
}
