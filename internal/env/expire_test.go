package env

import (
	"testing"
	"time"
)

func TestIsExpired_ReturnsTrueWhenPast(t *testing.T) {
	e := ExpiryEntry{
		Key:       "DB_PASSWORD",
		SyncedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if !e.IsExpired() {
		t.Error("expected entry to be expired")
	}
}

func TestIsExpired_ReturnsFalseWhenFuture(t *testing.T) {
	e := ExpiryEntry{
		Key:       "API_KEY",
		SyncedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if e.IsExpired() {
		t.Error("expected entry to not be expired")
	}
}

func TestCheckExpiry_DetectsExpiredKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "secret",
		"API_KEY":     "key123",
	}
	entries := []ExpiryEntry{
		{Key: "DB_PASSWORD", ExpiresAt: time.Now().Add(-1 * time.Hour)},
		{Key: "API_KEY", ExpiresAt: time.Now().Add(1 * time.Hour)},
	}
	expired := CheckExpiry(secrets, entries)
	if len(expired) != 1 || expired[0] != "DB_PASSWORD" {
		t.Errorf("expected [DB_PASSWORD], got %v", expired)
	}
}

func TestCheckExpiry_EmptyWhenNoneExpired(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc"}
	entries := []ExpiryEntry{
		{Key: "TOKEN", ExpiresAt: time.Now().Add(24 * time.Hour)},
	}
	expired := CheckExpiry(secrets, entries)
	if len(expired) != 0 {
		t.Errorf("expected no expired keys, got %v", expired)
	}
}

func TestBuildExpiryEntries_SetsCorrectExpiry(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "val"}
	now := time.Now()
	ttl := 24 * time.Hour
	entries := BuildExpiryEntries(secrets, now, ttl)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Key != "SECRET_KEY" {
		t.Errorf("expected key SECRET_KEY, got %s", e.Key)
	}
	expected := now.Add(ttl)
	if e.ExpiresAt.Unix() != expected.Unix() {
		t.Errorf("unexpected ExpiresAt: %v", e.ExpiresAt)
	}
}

func TestBuildExpiryEntries_ZeroTTLReturnsNil(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	entries := BuildExpiryEntries(secrets, time.Now(), 0)
	if entries != nil {
		t.Errorf("expected nil for zero TTL, got %v", entries)
	}
}

func TestExpiryWarnings_FormatsMessages(t *testing.T) {
	warnings := ExpiryWarnings([]string{"DB_PASS", "API_TOKEN"})
	if len(warnings) != 2 {
		t.Fatalf("expected 2 warnings, got %d", len(warnings))
	}
	for _, w := range warnings {
		if len(w) == 0 {
			t.Error("expected non-empty warning message")
		}
	}
}
