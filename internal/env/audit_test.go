package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWriteAuditLog_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "logs", "audit.log")

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Profile:   "staging",
		EnvFile:   ".env",
		Added:     3,
		Removed:   1,
		Changed:   2,
	}

	if err := WriteAuditLog(logPath, entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("expected audit log file to exist")
	}
}

func TestWriteAuditLog_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	entry := AuditEntry{Timestamp: time.Now().UTC(), Profile: "prod"}
	if err := WriteAuditLog(logPath, entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}

func TestReadAuditLog_ReturnsEntries(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	for i, profile := range []string{"dev", "staging", "prod"} {
		e := AuditEntry{
			Timestamp: time.Now().UTC(),
			Profile:   profile,
			Added:     i + 1,
		}
		if err := WriteAuditLog(logPath, e); err != nil {
			t.Fatalf("write error: %v", err)
		}
	}

	entries, err := ReadAuditLog(logPath)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Profile != "staging" {
		t.Errorf("expected staging, got %s", entries[1].Profile)
	}
}

func TestReadAuditLog_MissingFileReturnsEmpty(t *testing.T) {
	entries, err := ReadAuditLog("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestWriteAuditLog_AppendPreservesExistingEntries(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	first := AuditEntry{Timestamp: time.Now().UTC(), Profile: "dev", Added: 1}
	second := AuditEntry{Timestamp: time.Now().UTC(), Profile: "prod", Added: 2}

	if err := WriteAuditLog(logPath, first); err != nil {
		t.Fatalf("first write error: %v", err)
	}
	if err := WriteAuditLog(logPath, second); err != nil {
		t.Fatalf("second write error: %v", err)
	}

	entries, err := ReadAuditLog(logPath)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after append, got %d", len(entries))
	}
	if entries[0].Profile != "dev" {
		t.Errorf("expected first entry profile 'dev', got %s", entries[0].Profile)
	}
	if entries[1].Profile != "prod" {
		t.Errorf("expected second entry profile 'prod', got %s", entries[1].Profile)
	}
}
