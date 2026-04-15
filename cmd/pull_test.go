package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, dir string) string {
	t.Helper()
	content := `vault:
  profiles:
    default:
      address: "http://127.0.0.1:8200"
      token: "test-token"
      env_file: ".env"
      paths:
        - "secret/data/app"
`
	path := filepath.Join(dir, "vaultpull.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	return path
}

func TestPullCmd_IsRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "pull" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'pull' command to be registered on root")
	}
}

func TestPullCmd_MissingConfig(t *testing.T) {
	dir := t.TempDir()
	cfgFile = filepath.Join(dir, "nonexistent.yaml")
	t.Cleanup(func() { cfgFile = "vaultpull.yaml" })

	err := runPull(pullCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing config file, got nil")
	}
}

func TestPullCmd_MissingToken(t *testing.T) {
	dir := t.TempDir()
	cfgFile = writeTestConfig(t, dir)
	t.Cleanup(func() { cfgFile = "vaultpull.yaml" })

	os.Unsetenv("VAULT_TOKEN")

	// Override token in config to empty to trigger missing token error
	content := `vault:
  profiles:
    default:
      address: "http://127.0.0.1:8200"
      token: ""
      env_file: ".env"
      paths:
        - "secret/data/app"
`
	path := filepath.Join(dir, "vaultpull.yaml")
	os.WriteFile(path, []byte(content), 0644)
	cfgFile = path

	err := runPull(pullCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing vault token, got nil")
	}
}
