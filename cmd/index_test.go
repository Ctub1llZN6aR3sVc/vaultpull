package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeIndexEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestIndexCmd_IsRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "index [file]" {
			return
		}
	}
	t.Fatal("index command not registered")
}

func TestIndexCmd_ShowsKeys(t *testing.T) {
	p := writeIndexEnv(t, "API_KEY=abc\nDB_PASS=secret\n")
	rootCmd.SetArgs([]string{"index", p})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIndexCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"index", "/nonexistent/.env"})
	// should not hard-fail; missing file is warned and treated as empty
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
