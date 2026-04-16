package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	secrets := map[string]string{"KEY": "value"}

	if err := SaveSnapshot(path, "dev", secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestSaveAndLoadSnapshot_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc"}

	if err := SaveSnapshot(path, "staging", secrets); err != nil {
		t.Fatalf("save: %v", err)
	}
	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if snap.Profile != "staging" {
		t.Errorf("expected profile staging, got %s", snap.Profile)
	}
	if snap.Secrets["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret")
	}
}

func TestLoadSnapshot_MissingFileReturnsNil(t *testing.T) {
	snap, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Errorf("expected nil snapshot")
	}
}

func TestDiffSnapshot_NilSnapshot(t *testing.T) {
	current := map[string]string{"NEW_KEY": "val"}
	result := DiffSnapshot(nil, current)
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added key, got %d", len(result.Added))
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	old := map[string]string{"A": "1", "B": "2"}
	_ = SaveSnapshot(path, "dev", old)
	snap, _ := LoadSnapshot(path)

	current := map[string]string{"A": "changed", "C": "new"}
	result := DiffSnapshot(snap, current)

	if _, ok := result.Changed["A"]; !ok {
		t.Errorf("expected A to be changed")
	}
	if _, ok := result.Added["C"]; !ok {
		t.Errorf("expected C to be added")
	}
	if _, ok := result.Removed["B"]; !ok {
		t.Errorf("expected B to be removed")
	}
}
