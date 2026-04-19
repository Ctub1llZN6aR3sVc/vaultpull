package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_ExtractByPrefix(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"APP_HOST": "localhost", "APP_PORT": "9000", "DB_URL": "postgres"},
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: tmp,
		Extract: &env.ExtractOptions{Prefixes: []string{"APP_"}},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	got, _ := env.Read(tmp)
	if got["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", got["APP_HOST"])
	}
	if _, ok := got["DB_URL"]; ok {
		t.Error("DB_URL should have been excluded by extract")
	}
}

func TestRun_ExtractNilSkipsExtraction(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"APP_HOST": "localhost", "DB_URL": "postgres"},
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: tmp,
		Extract: nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	got, _ := env.Read(tmp)
	if got["APP_HOST"] != "localhost" || got["DB_URL"] != "postgres" {
		t.Errorf("expected all keys present, got %v", got)
	}
}
