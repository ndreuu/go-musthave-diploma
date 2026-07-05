package logger

import "testing"

func TestNewLogger(t *testing.T) {
	log, err := NewLogger("info")
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	if log == nil {
		t.Fatal("NewLogger() returned nil logger")
	}

	_ = log.Sync()
}

func TestNewLogger_InvalidLevel(t *testing.T) {
	log, err := NewLogger("bad-level")
	if err == nil {
		t.Fatal("expected error")
	}

	if log != nil {
		t.Fatal("expected nil logger")
	}
}
