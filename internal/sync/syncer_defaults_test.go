package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_DefaultsAppliedForMissingKeys(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "vault-value"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: f,
		Defaults: map[string]string{
			"API_KEY":  "ignored-default",
			"LOG_LEVEL": "info",
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, _ := env.Read(f)
	if secrets["API_KEY"] != "vault-value" {
		t.Errorf("vault value should win, got %s", secrets["API_KEY"])
	}
	if secrets["LOG_LEVEL"] != "info" {
		t.Errorf("default should be applied, got %s", secrets["LOG_LEVEL"])
	}
}

func TestRun_NilDefaultsSkipsApplication(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "vault-value"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: f,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, _ := env.Read(f)
	if len(secrets) != 1 {
		t.Errorf("expected 1 key, got %d", len(secrets))
	}
}
