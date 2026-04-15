package env

import (
	"os"
	"path/filepath"
	"testing"
)

func tempEnvFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".env")
}

func TestWrite_CreatesFile(t *testing.T) {
	path := tempEnvFile(t)
	w := NewWriter(path)

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	if err := w.Write(secrets); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	for k, v := range secrets {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestWrite_FilePermissions(t *testing.T) {
	path := tempEnvFile(t)
	w := NewWriter(path)

	if err := w.Write(map[string]string{"KEY": "value"}); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions: got %o, want 0600", perm)
	}
}

func TestMerge_PreservesExistingKeys(t *testing.T) {
	path := tempEnvFile(t)
	w := NewWriter(path)

	if err := w.Write(map[string]string{"EXISTING": "old", "KEEP": "me"}); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	if err := w.Merge(map[string]string{"EXISTING": "new", "ADDED": "yes"}); err != nil {
		t.Fatalf("Merge() error: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	cases := map[string]string{
		"EXISTING": "new",
		"KEEP":     "me",
		"ADDED":    "yes",
	}
	for k, v := range cases {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
}

func TestEscapeValue_QuotesSpecialChars(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"with#hash", `"with#hash"`},
	}
	for _, tc := range cases {
		got := escapeValue(tc.input)
		if got != tc.want {
			t.Errorf("escapeValue(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
