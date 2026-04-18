package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_DedupeKeepsFirstWhenEnabled(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_URL": "vault-value"},
	})

	tmp := tempEnvFile(t)
	// pre-populate the env file with an existing value for the same key
	if err := env.NewWriter(tmp).Write(map[string]string{"DB_URL": "local-value"}); err != nil {
		t.Fatal(err)
	}

	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
		Dedupe:     true,
		KeepFirst:  true,
	})
	if err := s.Run(); err != nil {
		t.Fatal(err)
	}

	secrets, _ := env.Read(tmp)
	if secrets["DB_URL"] != "local-value" {
		t.Fatalf("expected local-value (keepFirst=true), got %s", secrets["DB_URL"])
	}
}

func TestRun_DedupeKeepsLastWhenDisabled(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_URL": "vault-value"},
	})

	tmp := tempEnvFile(t)
	if err := env.NewWriter(tmp).Write(map[string]string{"DB_URL": "local-value"}); err != nil {
		t.Fatal(err)
	}

	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
		Dedupe:     true,
		KeepFirst:  false,
	})
	if err := s.Run(); err != nil {
		t.Fatal(err)
	}

	secrets, _ := env.Read(tmp)
	if secrets["DB_URL"] != "vault-value" {
		t.Fatalf("expected vault-value (keepFirst=false), got %s", secrets["DB_URL"])
	}
}
