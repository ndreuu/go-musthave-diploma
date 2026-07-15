// Package logger provides structured logging utilities for the GopherMart application.
//
// It wraps the zap logging library to provide a simple interface for creating
// production-ready loggers with configurable log levels.
//
// Example usage:
//
//	log, err := logger.NewLogger("info")
//	if err != nil {
//	    panic(err)
//	}
//	defer log.Sync()
//
//	log.Info("server started", zap.String("addr", ":8080"))
//	log.Error("database error", zap.Error(err))
package logger

import (
	"go.uber.org/zap"
)

// NewLogger creates a new zap logger with the specified log level.
//
// The logger is configured with production settings (JSON output, structured fields)
// and uses the provided level to filter log messages. Valid level strings are:
// "debug", "info", "warn", "error", "dpanic", "panic", "fatal".
//
// Parameters:
//   - level: Minimum log level to output (e.g., "info" logs info and above)
//
// Returns:
//   - *zap.Logger: Configured logger instance
//   - error: Error if level string is invalid or logger creation fails
//
// Example usage:
//
//	log, err := logger.NewLogger("debug")
//	if err != nil {
//	    return err
//	}
//	defer log.Sync()
//
//	log.Debug("debug message")  // Output: enabled
//	log.Info("info message")    // Output: enabled
//	log.Warn("warn message")    // Output: enabled
//	log.Error("error message")  // Output: enabled
func NewLogger(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl, nil
}
