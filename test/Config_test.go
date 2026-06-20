package test

import (
	"testing"

	"user/pkg/config"
)

func TestConfig_Get_ExistingKey(t *testing.T) {
	v := config.Get("config.mysql.driveName")
	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s != "mysql" {
		t.Errorf("expected 'mysql', got '%s'", s)
	}
}

func TestConfig_Get_DataSourceName(t *testing.T) {
	v := config.Get("config.mysql.dataSourceName")
	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s == "" {
		t.Error("expected non-empty dataSourceName")
	}
}

func TestConfig_Get_NonExistent(t *testing.T) {
	v := config.Get("config.nonexistent.key")
	if v != nil {
		t.Errorf("expected nil, got %v", v)
	}
}

func TestConfig_Get_JWTSecret(t *testing.T) {
	v := config.Get("config.jwt.secret")
	s, ok := v.(string)
	if !ok {
		t.Fatalf("expected string, got %T", v)
	}
	if s == "" {
		t.Error("expected non-empty jwt secret")
	}
}
