package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeChainEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestChainCmd_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "chain" {
			found = true
			break
		}
	}
	if !found {
		t.Error("chain command not registered")
	}
}

func TestChainCmd_NoFilesReturnsError(t *testing.T) {
	rootCmd.SetArgs([]string{"chain"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no files provided")
	}
}

func TestChainCmd_MergesFiles(t *testing.T) {
	dir := t.TempDir()
	f1 := writeChainEnv(t, dir, "a.env", "FOO=from_a\nSHARED=a_wins\n")
	f2 := writeChainEnv(t, dir, "b.env", "BAR=from_b\nSHARED=b_loses\n")

	rootCmd.SetArgs([]string{"chain", "--files", f1 + "," + f2})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
