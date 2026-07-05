package config

import (
	"flag"
	"testing"
	"time"
)

func TestLoadFromDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	cfg := LoadFrom(fs, func(string) string {
		return ""
	}, nil)

	if cfg.RunAddress != ":8080" {
		t.Fatalf("RunAddress = %q", cfg.RunAddress)
	}

	if cfg.JWTSecret != "super-secret-key" {
		t.Fatalf("JWTSecret = %q", cfg.JWTSecret)
	}

	if cfg.AccrualPollInterval != time.Second {
		t.Fatalf("AccrualPollInterval = %v", cfg.AccrualPollInterval)
	}

	if cfg.AccrualBatchSize != 10 {
		t.Fatalf("AccrualBatchSize = %d", cfg.AccrualBatchSize)
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

	cfg := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)

	if cfg.RunAddress != ":9090" {
		t.Fatal("RunAddress")
	}

	if cfg.DatabaseURI != "postgres://db" {
		t.Fatal("DatabaseURI")
	}

	if cfg.AccrualSystemAddress != "http://accrual" {
		t.Fatal("AccrualSystemAddress")
	}

	if cfg.JWTSecret != "secret" {
		t.Fatal("JWTSecret")
	}

	if cfg.AccrualPollInterval != 5*time.Second {
		t.Fatal("AccrualPollInterval")
	}

	if cfg.AccrualBatchSize != 100 {
		t.Fatal("AccrualBatchSize")
	}
}

func TestLoadFromInvalidDuration(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"ACCRUAL_POLL_INTERVAL": "invalid",
	}

	cfg := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)

	if cfg.AccrualPollInterval != time.Second {
		t.Fatalf("AccrualPollInterval = %v, want 1s", cfg.AccrualPollInterval)
	}
}

func TestLoadFromInvalidBatchSize(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"ACCRUAL_BATCH_SIZE": "invalid",
	}

	cfg := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)

	if cfg.AccrualBatchSize != 10 {
		t.Fatalf("AccrualBatchSize = %d, want 10", cfg.AccrualBatchSize)
	}
}

func TestLoadFromNegativeBatchSize(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"ACCRUAL_BATCH_SIZE": "-5",
	}

	cfg := LoadFrom(fs, func(key string) string {
		return env[key]
	}, nil)

	if cfg.AccrualBatchSize != 10 {
		t.Fatalf("AccrualBatchSize = %d, want 10", cfg.AccrualBatchSize)
	}
}

func TestLoadFromFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	args := []string{"-a", ":7070", "-d", "postgres://flag", "-r", "http://flag", "-s", "flag-secret"}

	cfg := LoadFrom(fs, func(string) string {
		return ""
	}, args)

	if cfg.RunAddress != ":7070" {
		t.Fatal("RunAddress")
	}

	if cfg.DatabaseURI != "postgres://flag" {
		t.Fatal("DatabaseURI")
	}

	if cfg.AccrualSystemAddress != "http://flag" {
		t.Fatal("AccrualSystemAddress")
	}

	if cfg.JWTSecret != "flag-secret" {
		t.Fatal("JWTSecret")
	}
}

func TestLoadFromFlagsOverride(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	env := map[string]string{
		"RUN_ADDRESS":  ":9090",
		"JWT_SECRET":   "env-secret",
	}

	args := []string{"-a", ":7070", "-s", "flag-secret"}

	cfg := LoadFrom(fs, func(key string) string {
		return env[key]
	}, args)

	if cfg.RunAddress != ":7070" {
		t.Fatal("RunAddress")
	}

	if cfg.JWTSecret != "flag-secret" {
		t.Fatal("JWTSecret")
	}
}

func TestLoad(t *testing.T) {
	cfg := Load()
	if cfg == nil {
		t.Fatal("Load() returned nil")
	}
}
