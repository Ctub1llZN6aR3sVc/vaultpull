package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_SavesSnapshotWhenEnabled(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	snapFile := filepath.Join(dir, "snap.json")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"TOKEN": "abc123"},
	})
	s := New(client, Options{
		Paths:        []string{"secret/app"},
		OutputFile:   envFile,
		SnapshotFile: snapFile,
	})
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, err := os.Stat(snapFile); err != nil {
		t.Errorf("snapshot file not created: %v", err)
	}
	snap, err := env.LoadSnapshot(snapFile)
	if err != nil || snap == nil {
		t.Fatalf("could not load snapshot: %v", err)
	}
	if snap.Secrets["TOKEN"] != "abc123" {
		t.Errorf("expected TOKEN=abc123 in snapshot")
	}
}

func TestRun_NoSnapshotWhenPathEmpty(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "val"},
	})
	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: envFile,
	})
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// No snapshot path — just ensure no panic and no stray file
	matches, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	if len(matches) != 0 {
		t.Errorf("expected no JSON files, found %v", matches)
	}
}
