package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_FilterIncludePrefix(t *testing.T) {
	mockSecrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_PORT":    "8080",
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": mockSecrets,
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, []string{"secret/app"}, tmp)
	s.FilterOpts = env.FilterOptions{
		IncludePrefix: []string{"DB_"},
	}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(tmp)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
	if _, ok := result["APP_PORT"]; ok {
		t.Error("APP_PORT should be filtered out")
	}
}

func TestRun_FilterExcludeKeys(t *testing.T) {
	mockSecrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "topsecret",
		"APP_PORT":    "8080",
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": mockSecrets,
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, []string{"secret/app"}, tmp)
	s.FilterOpts = env.FilterOptions{
		ExcludeKeys: []string{"DB_PASSWORD"},
	}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(tmp)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if _, ok := result["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should be excluded")
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in output")
	}
}

func TestRun_NoFilterWritesAll(t *testing.T) {
	mockSecrets := map[string]string{
		"KEY_A": "val_a",
		"KEY_B": "val_b",
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": mockSecrets,
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, []string{"secret/app"}, tmp)

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(tmp)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}
