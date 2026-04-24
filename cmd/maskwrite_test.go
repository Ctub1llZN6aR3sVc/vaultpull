package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeMaskEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestMaskWriteCmd_IsRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "maskwrite <file>" {
			return
		}
	}
	t.Error("maskwrite command not registered")
}

func TestMaskWriteCmd_MasksAutoDetected(t *testing.T) {
	p := writeMaskEnv(t, "API_KEY=supersecret\nHOST=localhost\n")

	rootCmd.SetArgs([]string{"maskwrite", "--auto", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	content := string(data)
	if strings.Contains(content, "supersecret") {
		t.Errorf("expected API_KEY to be masked")
	}
	if !strings.Contains(content, "localhost") {
		t.Errorf("expected HOST to be preserved")
	}
}

func TestMaskWriteCmd_DryRunDoesNotWrite(t *testing.T) {
	p := writeMaskEnv(t, "TOKEN=abc123\n")

	rootCmd.SetArgs([]string{"maskwrite", "--keys", "TOKEN", "--dry-run", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "abc123") {
		t.Errorf("dry-run should not have modified the file")
	}
}

func TestMaskWriteCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"maskwrite", "/nonexistent/.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}
