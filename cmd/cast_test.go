package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCastEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestCastCmd_IsRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "cast [file]" {
			return
		}
	}
	t.Fatal("cast command not registered")
}

func TestCastCmd_CastsIntValue(t *testing.T) {
	file := writeCastEnv(t, "PORT=9000.0\nNAME=app\n")

	buf := new(strings.Builder)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"cast", file, "--type", "PORT=int"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "PORT=9000") {
		t.Errorf("expected PORT=9000 in output, got: %s", buf.String())
	}
}

func TestCastCmd_StrictFailsOnBadValue(t *testing.T) {
	file := writeCastEnv(t, "PORT=notanumber\n")

	rootCmd.SetArgs([]string{"cast", file, "--type", "PORT=int", "--strict"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
}

func TestCastCmd_MissingFile(t *testing.T) {
	rootCmd.SetArgs([]string{"cast", "/nonexistent/.env"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}
