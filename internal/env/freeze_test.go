package env

import (
	"os"
	"path/filepath"
	"testing"
)

func freezePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "freeze.json")
}

func TestFreeze_WritesEntries(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc123"}
	path := freezePath(t)
	res, err := Freeze(secrets, []string{"DB_PASS", "API_KEY"}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Frozen) != 2 {
		t.Errorf("expected 2 frozen, got %d", len(res.Frozen))
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("freeze file not created: %v", err)
	}
}

func TestFreeze_SkipsMissingKeys(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret"}
	path := freezePath(t)
	res, err := Freeze(secrets, []string{"DB_PASS", "MISSING"}, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestCheckFrozen_DetectsChanged(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "original"}
	path := freezePath(t)
	if _, err := Freeze(secrets, []string{"DB_PASS"}, path); err != nil {
		t.Fatal(err)
	}
	modified := map[string]string{"DB_PASS": "changed"}
	chang(modified, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changed) != 1 || changed[0] != "DB_PASS" {
		t.Errorf("expected DB_PASS in changed, got %v", changed)
	}
}

func TestCheckFrozen_NoChanges(t *testing.T) {
	secrets := map[string]string{"API_KEY": "stable"}
	path := freezePath(t)
	if _, err := Freeze(secrets, []string{"API_KEY"}, path); err != nil {
		t.Fatal(err)
	}
	changed, err := CheckFrozen(secrets, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestCheckFrozen_MissingFileReturnsNil(t *testing.T) {
	changed, err := CheckFrozen(map[string]string{}, "/nonexistent/freeze.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed != nil {
		t.Errorf("expected nil, got %v", changed)
	}
}
