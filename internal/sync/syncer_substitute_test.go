package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_SubstituteResolvesReferences(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"BASE_URL": "https://api.example.com",
			"FULL_URL": "${BASE_URL}/v1",
		},
	})

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	s := New(client, []string{"secret/app"}, envPath)
	s.Substitute = &env.SubstituteOptions{}

	if err := s.Run(); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	got, err := env.Read(envPath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if got["FULL_URL"] != "https://api.example.com/v1" {
		t.Errorf("expected substituted URL, got %q", got["FULL_URL"])
	}
}

func TestRun_SubstituteNilSkipsSubstitution(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"BASE": "hello",
			"VAL":  "${BASE}_world",
		},
	})

	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	s := New(client, []string{"secret/app"}, envPath)
	s.Substitute = nil

	if err := s.Run(); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	got, err := env.Read(envPath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	// without substitution, raw template string is preserved
	if got["VAL"] != "${BASE}_world" {
		t.Errorf("expected raw value, got %q", got["VAL"])
	}

	_ = os.Remove(envPath)
}
