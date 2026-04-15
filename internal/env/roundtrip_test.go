package env

import (
	"path/filepath"
	"testing"
)

// TestRoundTrip verifies that values written by Writer can be read back
// correctly by, preserving special characters and quoted values.
func TestRoundTrip_WriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path)

	secrets := map[string]string{
		"SIMPLE":    "value",
		"WITH_SPACE": "hello world",
		"WITH_HASH": "abc#def",
		"EMPTY":     "",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	for key, want := range secrets {
		val, ok := got[key]
		if !ok {
			t.Errorf("key %q not found after round-trip", key)
			continue
		}
		if val != want {
			t.Errorf("key %q: got %q, want %q", key, val, want)
		}
	}
}

// TestRoundTrip_MergePreservesExisting verifies that merging new secrets
// does not overwrite keys already present in the file.
func TestRoundTrip_MergePreservesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := NewWriter(path)

	initial := map[string]string{"EXISTING": "keep_me", "OLD": "value"}
	if err := w.Write(initial); err != nil {
		t.Fatalf("initial Write failed: %v", err)
	}

	newSecrets := map[string]string{"NEW_KEY": "new_value"}
	if err := w.Merge(newSecrets); err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read after merge failed: %v", err)
	}

	if got["EXISTING"] != "keep_me" {
		t.Errorf("EXISTING was overwritten: got %q", got["EXISTING"])
	}
	if got["NEW_KEY"] != "new_value" {
		t.Errorf("NEW_KEY not written: got %q", got["NEW_KEY"])
	}
}
