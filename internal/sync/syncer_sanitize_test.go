package sync

import (
	"os"
	"testing"

	"github.com/fmartingr/vaultpull/internal/env"
	"github.com/fmartingr/vaultpull/internal/vault"
)

func TestRun_SanitizeStripsControlChars(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "val\x01ue"},
	})

	f, _ := os.CreateTemp(t.TempDir(), "*.env")
	f.Close()

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: f.Name(),
		Sanitize: &env.SanitizeOptions{
			StripControlChars: true,
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f.Name())
	if result["API_KEY"] != "value" {
		t.Errorf("expected control chars stripped, got %q", result["API_KEY"])
	}
}

func TestRun_SanitizeNilSkipsSanitization(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "val\x01ue"},
	})

	f, _ := os.CreateTemp(t.TempDir(), "*.env")
	f.Close()

	s := New(client, Options{
		Paths:    []string{"secret/app"},
		EnvFile:  f.Name(),
		Sanitize: nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f.Name())
	if result["API_KEY"] != "val\x01ue" {
		t.Errorf("expected value unchanged, got %q", result["API_KEY"])
	}
}
