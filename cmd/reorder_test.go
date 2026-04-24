package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeReorderEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestReorderCmd_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "reorder" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("reorder command not registered")
	}
}

func TestReorderCmd_AlphabeticalDryRun(t *testing.T) {
	p := writeReorderEnv(t, "C=3\nA=1\nB=2\n")

	out := &strings.Builder{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"reorder", "--file", p, "--alpha", "--dry-run"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	result := out.String()
	idxA := strings.Index(result, "A=")
	idxB := strings.Index(result, "B=")
	idxC := strings.Index(result, "C=")
	if idxA > idxB || idxB > idxC {
		t.Fatalf("expected A < B < C order, got:\n%s", result)
	}
}

func TestReorderCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"reorder", "--file", "/nonexistent/.env", "--alpha"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}
