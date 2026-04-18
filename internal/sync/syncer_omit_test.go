package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_OmitRemovesKeys(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "DEBUG": "true", "API_KEY": "abc"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: f,
		Omit: &env.OmitOptions{
			Keys: []string{"DEBUG"},
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f)
	if _, ok := result["DEBUG"]; ok {
		t.Fatal("DEBUG should have been omitted")
	}
	if result["DB_PASS"] != "s3cr3t" {
		t.Fatal("DB_PASS should be present")
	}
}

func TestRun_OmitNilSkipsOmission(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"A": "1", "B": "2"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: f,
		Omit:    nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f)
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}
