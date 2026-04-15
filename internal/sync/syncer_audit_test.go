package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_WritesAuditLogWhenEnabled(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	logPath := filepath.Join(dir, "audit.log")

	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "value"},
	})

	s := New(mock, []string{"secret/app"}, envFile)
	s.AuditLogPath = logPath
	s.Profile = "test"

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := env.ReadAuditLog(logPath)
	if err != nil {
		t.Fatalf("read audit log error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(entries))
	}
	if entries[0].Profile != "test" {
		t.Errorf("expected profile 'test', got '%s'", entries[0].Profile)
	}
	if entries[0].Added != 1 {
		t.Errorf("expected 1 added key, got %d", entries[0].Added)
	}
}

func TestRun_NoAuditLogWhenPathEmpty(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"KEY": "value"},
	})

	s := New(mock, []string{"secret/app"}, envFile)
	// AuditLogPath intentionally left empty

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No audit file should be created anywhere in dir
	matches, _ := filepath.Glob(filepath.Join(dir, "*.log"))
	if len(matches) != 0 {
		t.Errorf("expected no log files, found: %v", matches)
	}

	_ = os.Remove(envFile)
}
