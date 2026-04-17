package env

import (
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadPins_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")

	pins := []PinEntry{
		{Key: "DB_PASS", Value: "secret", PinnedAt: time.Now().UTC()},
	}
	if err := SavePins(path, pins); err != nil {
		t.Fatalf("SavePins: %v", err)
	}
	loaded, err := LoadPins(path)
	if err != nil {
		t.Fatalf("LoadPins: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Key != "DB_PASS" || loaded[0].Value != "secret" {
		t.Fatalf("unexpected loaded pins: %+v", loaded)
	}
}

func TestLoadPins_MissingFileReturnsNil(t *testing.T) {
	pins, err := LoadPins("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if pins != nil {
		t.Fatalf("expected nil pins")
	}
}

func TestApplyPins_OverridesSecrets(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "original", "API_KEY": "key"}
	pins := []PinEntry{{Key: "DB_PASS", Value: "pinned", PinnedAt: time.Now()}}
	out := ApplyPins(secrets, pins)
	if out["DB_PASS"] != "pinned" {
		t.Fatalf("expected pinned value, got %s", out["DB_PASS"])
	}
	if out["API_KEY"] != "key" {
		t.Fatalf("non-pinned key should be unchanged")
	}
}

func TestApplyPins_ExpiredPinIgnored(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "original"}
	pins := []PinEntry{
		{Key: "DB_PASS", Value: "pinned", PinnedAt: time.Now().Add(-2 * time.Hour), ExpiresAt: time.Now().Add(-1 * time.Hour)},
	}
	out := ApplyPins(secrets, pins)
	if out["DB_PASS"] != "original" {
		t.Fatalf("expired pin should not override, got %s", out["DB_PASS"])
	}
}

func TestApplyPins_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"X": "1"}
	pins := []PinEntry{{Key: "X", Value: "2", PinnedAt: time.Now()}}
	ApplyPins(secrets, pins)
	if secrets["X"] != "1" {
		t.Fatal("input map was mutated")
	}
}
