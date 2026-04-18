package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_CastAppliesTypes(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"PORT": "8080.0", "ENABLED": "1"},
	})

	tmp, _ := os.CreateTemp(t.TempDir(), "*.env")
	tmp.Close()

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp.Name(),
		Cast: &env.CastOptions{
			Types: map[string]string{
				"PORT":    "int",
				"ENABLED": "bool",
			},
		},
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Secrets["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", result.Secrets["PORT"])
	}
	if result.Secrets["ENABLED"] != "true" {
		t.Errorf("expected ENABLED=true, got %s", result.Secrets["ENABLED"])
	}
}

func TestRun_NoCastByDefault(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"PORT": "8080.0"},
	})

	tmp, _ := os.CreateTemp(t.TempDir(), "*.env")
	tmp.Close()

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp.Name(),
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Secrets["PORT"] != "8080.0" {
		t.Errorf("expected original value 8080.0, got %s", result.Secrets["PORT"])
	}
}
