package sync

import (
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_ResolveAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_URL": "https://example.com"},
	})

	s := New(client, []string{"secret/app"}, output)
	s.Resolve = &env.ResolveOptions{
		Defaults: map[string]string{"LOG_LEVEL": "info"},
		Required: []string{"API_URL", "LOG_LEVEL"},
	}

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := env.Read(output)
	if got["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info from defaults, got %q", got["LOG_LEVEL"])
	}
	if got["API_URL"] != "https://example.com" {
		t.Errorf("expected API_URL from vault, got %q", got["API_URL"])
	}
}

func TestRun_ResolveFailsOnMissingRequired(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"FOO": "bar"},
	})

	s := New(client, []string{"secret/app"}, output)
	s.Resolve = &env.ResolveOptions{
		Required: []string{"MISSING_KEY"},
	}

	if err := s.Run(); err == nil {
		t.Fatal("expected error for missing required key")
	}
}

func TestRun_ResolveNilSkipsResolution(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "value"},
	})

	s := New(client, []string{"secret/app"}, output)
	// s.Resolve is nil by default

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := env.Read(output)
	if got["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", got["KEY"])
	}
}
