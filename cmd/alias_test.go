package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeAliasEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestAliasCmd_IsRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "alias" {
			return
		}
	}
	t.Fatal("alias command not registered")
}

func TestAliasCmd_CreatesAlias(t *testing.T) {
	p := writeAliasEnv(t, "DB_PASS=secret\n")

	rootCmd.SetArgs([]string{"alias", "--file", p, "--map", "DATABASE_PASSWORD=DB_PASS"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if !strings.Contains(string(data), "DATABASE_PASSWORD=secret") {
		t.Fatalf("expected alias in file, got:\n%s", data)
	}
}

func TestAliasCmd_MissingMapReturnsError(t *testing.T) {
	p := writeAliasEnv(t, "FOO=bar\n")
	rootCmd.SetArgs([]string{"alias", "--file", p})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no --map provided")
	}
}

func TestAliasCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"alias", "--file", "/nonexistent/.env", "--map", "B=A"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
