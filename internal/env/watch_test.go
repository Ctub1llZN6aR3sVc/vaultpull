package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeWatchFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeWatchFile: %v", err)
	}
	return p
}

func TestHashFile_ReturnsEmptyForMissing(t *testing.T) {
	h, err := HashFile("/nonexistent/path/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h != "" {
		t.Errorf("expected empty hash, got %q", h)
	}
}

func TestHashFile_DeterministicHash(t *testing.T) {
	dir := t.TempDir()
	p := writeWatchFile(t, dir, ".env", "KEY=value\n")

	h1, err := HashFile(p)
	if err != nil {
		t.Fatalf("hash1: %v", err)
	}
	h2, err := HashFile(p)
	if err != nil {
		t.Fatalf("hash2: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected same hash, got %q vs %q", h1, h2)
	}
}

func TestWatchState_ChangedAfterWrite(t *testing.T) {
	dir := t.TempDir()
	p := writeWatchFile(t, dir, ".env", "KEY=old\n")

	w, err := NewWatchState(p)
	if err != nil {
		t.Fatalf("NewWatchState: %v", err)
	}

	changed, err := w.Changed()
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if changed {
		t.Error("expected no change immediately after init")
	}

	if err := os.WriteFile(p, []byte("KEY=new\n"), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	changed, err = w.Changed()
	if err != nil {
		t.Fatalf("Changed after write: %v", err)
	}
	if !changed {
		t.Error("expected change after file update")
	}
}

func TestWatchState_RefreshClearsChanged(t *testing.T) {
	dir := t.TempDir()
	p := writeWatchFile(t, dir, ".env", "KEY=v1\n")

	w, err := NewWatchState(p)
	if err != nil {
		t.Fatalf("NewWatchState: %v", err)
	}

	_ = os.WriteFile(p, []byte("KEY=v2\n"), 0600)

	if err := w.Refresh(); err != nil {
		t.Fatalf("Refresh: %v", err)
	}

	changed, err := w.Changed()
	if err != nil {
		t.Fatalf("Changed: %v", err)
	}
	if changed {
		t.Error("expected no change after refresh")
	}
}
