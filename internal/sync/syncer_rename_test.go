package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_RenameAppliesMapping(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"OLD_DB_PASS": "secret123"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: f,
		Rename:  map[string]string{"OLD_DB_PASS": "DB_PASSWORD"},
	})

	_, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f)
	if result["DB_PASSWORD"] != "secret123" {
		t.Errorf("expected DB_PASSWORD=secret123, got %q", result["DB_PASSWORD"])
	}
	if _, ok := result["OLD_DB_PASS"]; ok {
		t.Error("OLD_DB_PASS should have been renamed away")
	}
}

func TestRun_NoRenameByDefault(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "abc"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: f,
	})

	_, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f)
	if result["API_KEY"] != "abc" {
		t.Errorf("expected API_KEY=abc, got %q", result["API_KEY"])
	}
}
