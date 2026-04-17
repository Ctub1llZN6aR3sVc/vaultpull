package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLintCmd_IsRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "lint" {
			found = true
		}
	}
	if !found {
		t.Error("lint command not registered")
	}
}

func TestLintCmd_MissingConfig(t *testing.T) {
	rootCmd.SetArgs([]string{"lint", "--config", "/nonexistent/vaultpull.yaml", "--token", "tok"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func TestLintCmd_MissingToken(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "vaultpull.yaml")
	content := `default_profile: default
profiles:
  - name: default
    address: http://127.0.0.1:8200
    paths: ["secret/app"]
    output: .env
`
	os.WriteFile(cfgPath, []byte(content), 0644)
	os.Unsetenv("VAULT_TOKEN")
	rootCmd.SetArgs([]string{"lint", "--config", cfgPath})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when token is missing")
	}
}
