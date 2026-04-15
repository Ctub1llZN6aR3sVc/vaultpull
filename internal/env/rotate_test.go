package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRotate_CreatesFileWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{"KEY": "value"}
	res, err := Rotate(path, secrets, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BackupPath != "" {
		t.Errorf("expected no backup for new file, got %s", res.BackupPath)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected env file to exist: %v", err)
	}
}

func TestRotate_BacksUpExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Create initial file.
	if err := os.WriteFile(path, []byte("OLD=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	secrets := map[string]string{"NEW": "2"}
	res, err := Rotate(path, secrets, RotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.BackupPath == "" {
		t.Fatal("expected a backup path")
	}
	if _, err := os.Stat(res.BackupPath); err != nil {
		t.Errorf("backup file should exist: %v", err)
	}

	// New file should contain NEW, not OLD.
	got, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if got["NEW"] != "2" {
		t.Errorf("expected NEW=2, got %q", got["NEW"])
	}
	if _, ok := got["OLD"]; ok {
		t.Error("OLD key should not be present after rotation")
	}
}

func TestRotate_PrunesOldBackups(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	// Perform 3 rotations with MaxBackups=2.
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(path, []byte("K=v\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := Rotate(path, map[string]string{"K": "v"}, RotateOptions{MaxBackups: 2})
		if err != nil {
			t.Fatalf("rotation %d failed: %v", i, err)
		}
	}

	matches, err := filepath.Glob(filepath.Join(dir, ".env.*.bak"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) > 2 {
		t.Errorf("expected at most 2 backups, got %d", len(matches))
	}
}

func TestRotate_CustomBackupDir(t *testing.T) {
	dir := t.TempDir()
	backupDir := t.TempDir()
	path := filepath.Join(dir, ".env")

	if err := os.WriteFile(path, []byte("A=1\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	res, err := Rotate(path, map[string]string{"A": "2"}, RotateOptions{BackupDir: backupDir})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Dir(res.BackupPath) != backupDir {
		t.Errorf("expected backup in %s, got %s", backupDir, res.BackupPath)
	}
}
