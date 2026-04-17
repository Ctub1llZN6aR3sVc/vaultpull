package sync

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/densestvoid/vaultpull/internal/env"
	"github.com/densestvoid/vaultpull/internal/vault"
)

func TestRun_PinOverridesVaultValue(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	pinFile := filepath.Join(dir, "pins.json")

	pins := []env.PinEntry{
		{Key: "DB_PASS", Value: "pinned-value", PinnedAt: time.Now()},
	}
	if err := env.SavePins(pinFile, pins); err != nil {
		t.Fatalf("SavePins: %v", err)
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "vault-value", "OTHER": "other"},
	})
	s := New(client, []PathMapping{{VaultPath: "secret/app", EnvFile: envFile}},
		WithPinFile(pinFile),
	)
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	result, _ := env.Read(envFile)
	if result["DB_PASS"] != "pinned-value" {
		t.Fatalf("expected pinned-value, got %s", result["DB_PASS"])
	}
	if result["OTHER"] != "other" {
		t.Fatalf("expected other, got %s", result["OTHER"])
	}
}

func TestRun_NoPinFileWritesVaultValue(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "vault-value"},
	})
	s := New(client, []PathMapping{{VaultPath: "secret/app", EnvFile: envFile}})
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	result, _ := env.Read(envFile)
	if result["DB_PASS"] != "vault-value" {
		t.Fatalf("expected vault-value, got %s", result["DB_PASS"])
	}
}
