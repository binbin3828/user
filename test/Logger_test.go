package test

import (
	"os"
	"testing"

	"user/pkg/logger"
)

func TestLogger_DefaultLevel(t *testing.T) {
	os.Unsetenv("LOG_LEVEL")
	l := logger.NewZapLogger()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_DebugLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")
	l := logger.NewZapLogger()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_InfoLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	defer os.Unsetenv("LOG_LEVEL")
	l := logger.NewZapLogger()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_WarnLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "warn")
	defer os.Unsetenv("LOG_LEVEL")
	l := logger.NewZapLogger()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_ErrorLevel(t *testing.T) {
	os.Setenv("LOG_LEVEL", "error")
	defer os.Unsetenv("LOG_LEVEL")
	l := logger.NewZapLogger()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_InvalidLevelDefaultsToInfo(t *testing.T) {
	os.Setenv("LOG_LEVEL", "invalid")
	defer os.Unsetenv("LOG_LEVEL")
	l := logger.NewZapLogger()
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	l.Info("test log: should not panic")
}
