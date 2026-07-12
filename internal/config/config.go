// Package config provides application configuration loading and management.
//
// It supports loading configuration from command-line flags and environment variables,
// with sensible defaults for all settings. The package is designed to be testable
// through the LoadFrom function which accepts injectable dependencies.
//
// Configuration precedence: explicitly provided command-line flags take priority
// over environment variables, which in turn override hardcoded defaults.
package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration parameters.
//
// Each field corresponds to either a command-line flag or an environment variable.
// Default values are set in Load/LoadFrom, and can be overridden via
// environment variables (RUN_ADDRESS, DATABASE_URI, ACCRUAL_SYSTEM_ADDRESS,
// JWT_SECRET, ACCRUAL_POLL_INTERVAL, ACCRUAL_BATCH_SIZE) or command-line flags.
type Config struct {
	// RunAddress is the HTTP server listen address (e.g., ":8080").
	RunAddress string

	// DatabaseURI is the PostgreSQL connection string (e.g., "postgres://user:pass@host/db").
	DatabaseURI string

	// AccrualSystemAddress is the base URL of the accrual system API (e.g., "https://accrual.example.com").
	AccrualSystemAddress string

	// JWTSecret is the secret key used for signing and verifying JWT tokens.
	JWTSecret string

	// AccrualPollInterval is the interval between withdrawal polling cycles.
	AccrualPollInterval time.Duration

	// AccrualBatchSize is the number of orders to process in a single accrual batch.
	AccrualBatchSize int
}

// Load reads configuration from the default flag set and environment variables.
//
// It uses os.Args[1:] for command-line arguments and the default flag.FlagSet.
// This is the primary entry point for production code. For testing, use LoadFrom
// to inject custom flag sets and environment variable mocks.
//
// Defaults:
//   - RunAddress:          ":8080"
//   - JWTSecret:           "super-secret-key"
//   - AccrualPollInterval: 1s
//   - AccrualBatchSize:    10
func Load() (*Config, error) {
	return LoadFrom(flag.CommandLine, os.Getenv, os.Args[1:])
}

var ErrJWTSecretNotSet = errors.New("JWT secret is not set")

// LoadFrom reads configuration from the provided flag set, environment variable
// getter function, and command-line arguments.
//
// This function enables dependency injection for testing. The getenv parameter
// should be a function that returns environment variable values (typically
// os.Getenv in production, a mock map lookup in tests). The fs parameter
// defines which flag.FlagSet to parse, and args provides the command-line
// arguments to process.
//
// The function follows this precedence order:
//  1. Hardcoded defaults
//  2. Environment variables (if set and non-empty)
//  3. Command-line flags (parsed last, highest precedence)
//
// Environment variables recognized:
//   - RUN_ADDRESS:           HTTP listen address
//   - DATABASE_URI:          PostgreSQL connection string
//   - ACCRUAL_SYSTEM_ADDRESS: Accrual system API base URL
//   - JWT_SECRET:            JWT signing secret
//   - ACCRUAL_POLL_INTERVAL: Polling interval (Go duration string, e.g., "5s")
//   - ACCRUAL_BATCH_SIZE:    Number of orders per batch (must be positive)
//
// Command-line flags recognized:
//   - -a: Run address
//   - -d: Database URI
//   - -r: Accrual system address
//   - -s: JWT secret
func LoadFrom(
	fs *flag.FlagSet,
	getenv func(string) string,
	args []string,
) (*Config, error) {
	cfg := &Config{
		RunAddress:          ":8080",
		AccrualPollInterval: time.Second,
		AccrualBatchSize:    10,
	}

	if v := getenv("RUN_ADDRESS"); v != "" {
		cfg.RunAddress = v
	}

	if v := getenv("DATABASE_URI"); v != "" {
		cfg.DatabaseURI = v
	}

	if v := getenv("ACCRUAL_SYSTEM_ADDRESS"); v != "" {
		cfg.AccrualSystemAddress = v
	}

	if v := getenv("JWT_SECRET"); v != "" {
		cfg.JWTSecret = v
	}

	if v := getenv("ACCRUAL_POLL_INTERVAL"); v != "" {
		duration, err := time.ParseDuration(v)
		if err == nil {
			cfg.AccrualPollInterval = duration
		}
	}

	if v := getenv("ACCRUAL_BATCH_SIZE"); v != "" {
		size, err := strconv.Atoi(v)
		if err == nil && size > 0 {
			cfg.AccrualBatchSize = size
		}
	}

	fs.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "address and port")
	fs.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database URI")
	fs.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "accrual address")
	fs.StringVar(&cfg.JWTSecret, "s", cfg.JWTSecret, "JWT secret")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.JWTSecret == "" {
		return nil, ErrJWTSecretNotSet
	}

	return cfg, nil
}
