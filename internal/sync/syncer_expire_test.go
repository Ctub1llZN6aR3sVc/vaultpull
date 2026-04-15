package sync

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_ExpiryWarningOnExpiredKey(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASSWORD": "hunter2"},
	})

	tmpFile := tempEnvFile(t, "")

	s := New(mock, []string{"secret/app"}, tmpFile)

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets := map[string]string{"DB_PASSWORD": "hunter2"}
	entries := env.BuildExpiryEntries(secrets, time.Now().Add(-48*time.Hour), 24*time.Hour)

	expired := env.CheckExpiry(secrets, entries)
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired key, got %d", len(expired))
	}
	if expired[0] != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD to be expired, got %s", expired[0])
	}

	warnings := env.ExpiryWarnings(expired)
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(warnings))
	}
}

func TestRun_NoExpiryWarningWhenFresh(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "freshkey"},
	})

	tmpFile := tempEnvFile(t, "")
	s := New(mock, []string{"secret/app"}, tmpFile)

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets := map[string]string{"API_KEY": "freshkey"}
	entries := env.BuildExpiryEntries(secrets, time.Now(), 24*time.Hour)

	expired := env.CheckExpiry(secrets, entries)
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
}
