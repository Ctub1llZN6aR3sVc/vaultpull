package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_InvertSwapsKeysAndValues(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"localhost": "HOST",
			"5432":      "PORT",
		},
	})

	tmp := tempEnvFile(t)

	s := New(client, &Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Invert:  &env.InvertOptions{},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	got := readEnvFile(t, tmp)
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost after invert, got %q", got["HOST"])
	}
	if got["PORT"] != "5432" {
		t.Errorf("expected PORT=5432 after invert, got %q", got["PORT"])
	}
}

func TestRun_InvertNilSkipsInversion(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"HOST": "localhost",
		},
	})

	tmp := tempEnvFile(t)

	s := New(client, &Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Invert:  nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	got := readEnvFile(t, tmp)
	if got["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost without inversion, got %q", got["HOST"])
	}
}
