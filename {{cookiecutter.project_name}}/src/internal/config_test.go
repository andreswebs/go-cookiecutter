package config_test

import (
	"testing"

	"github.com/andreswebs/lexbrasilis-api/internal/config"
)

func TestLoad_DefaultPort(t *testing.T) {
	t.Setenv("PORT", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Load().Port = %q, want %q", cfg.Port, "8080")
	}
}

func TestLoad_CustomPort(t *testing.T) {
	t.Setenv("PORT", "9090")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("Load().Port = %q, want %q", cfg.Port, "9090")
	}
}
