package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func Load() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database connection URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "accrual system address")

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

	return cfg
}
