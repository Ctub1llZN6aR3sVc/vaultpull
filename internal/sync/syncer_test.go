package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

func tempEnvFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, ".env")
}

func TestRun_WritesSecretsToEnvFile(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	})

	envPath := tempEnvFile(t)
	writer, err := env.NewWriter(envPath)
	if err != nil {
		t.Fatalf("failed to create writer: %v", err)
	}

	s := sync.New(client, writer)
	if err := s.Run([]string{"secret/app"}); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	got, err := env.Read(envPath)
	if err != nil {
		t.Fatalf("failed to read env file: %v", err)
	}

	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestRun_MergesMultiplePaths(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/db": {"DB_NAME": "mydb"},
		"secret/api": {"API_KEY": "abc123"},
	})

	envPath := tempEnvFile(t)
	writer, err := env.NewWriter(envPath)
	if err != nil {
		t.Fatalf("failed to create writer: %v", err)
	}

	s := sync.New(client, writer)
	if err := s.Run([]string{"secret/db", "secret/api"}); err != nil {
		t.Fatalf("Run() returned error: %v", err)
	}

	got, err := env.Read(envPath)
	if err != nil {
		t.Fatalf("failed to read env file: %v", err)
	}

	if got["DB_NAME"] != "mydb" {
		t.Errorf("expected DB_NAME=mydb, got %q", got["DB_NAME"])
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", got["API_KEY"])
	}
}

func TestRun_ReturnsErrorOnMissingPath(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{})

	envPath := filepath.Join(t.TempDir(), ".env")
	writer, err := env.NewWriter(envPath)
	if err != nil {
		t.Fatalf("failed to create writer: %v", err)
	}

	s := sync.New(client, writer)
	if err := s.Run([]string{"secret/nonexistent"}); err == nil {
		t.Error("expected error for missing vault path, got nil")
	}

	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		t.Error("expected env file to not be created on error")
	}
}
