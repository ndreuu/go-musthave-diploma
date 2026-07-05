package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string

	JWTSecret string

	AccrualPollInterval time.Duration
	AccrualBatchSize    int
}

func Load() *Config {
	cfg := &Config{
		RunAddress:           ":8080",
		JWTSecret:            "super-secret-key",
		AccrualPollInterval:  time.Second,
		AccrualBatchSize:     10,
	}

	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database connection URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "accrual system address")
	flag.StringVar(&cfg.JWTSecret, "s", "super-secret-key", "JWT secret")

	flag.Parse()

	if v := os.Getenv("RUN_ADDRESS"); v != "" {
		cfg.RunAddress = v
	}

	if v := os.Getenv("DATABASE_URI"); v != "" {
		cfg.DatabaseURI = v
	}

	if v := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); v != "" {
		cfg.AccrualSystemAddress = v
	}

	if value := os.Getenv("JWT_SECRET"); value != "" {
		cfg.JWTSecret = value
	}

	if value := os.Getenv("ACCRUAL_POLL_INTERVAL"); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			cfg.AccrualPollInterval = duration
		}
	}

	if value := os.Getenv("ACCRUAL_BATCH_SIZE"); value != "" {
		size, err := strconv.Atoi(value)
		if err == nil && size > 0 {
			cfg.AccrualBatchSize = size
		}
	}

	return cfg
}
