package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_DiffReportsNewKeys(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Write an existing .env file with one key
	existing := []byte("EXISTING_KEY=old_value\n")
	if err := os.WriteFile(envFile, existing, 0600); err != nil {
		t.Fatalf("failed to write existing env file: %v", err)
	}

	mockSecrets := map[string]string{
		"EXISTING_KEY": "new_value",
		"BRAND_NEW_KEY": "hello",
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": mockSecrets,
	})

	s := New(client, []string{"secret/app"}, envFile)
	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	// Read the resulting file and verify diff expectations
	result, err := env.Read(envFile)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if result["BRAND_NEW_KEY"] != "hello" {
		t.Errorf("expected BRAND_NEW_KEY=hello, got %q", result["BRAND_NEW_KEY"])
	}
	if result["EXISTING_KEY"] != "new_value" {
		t.Errorf("expected EXISTING_KEY=new_value, got %q", result["EXISTING_KEY"])
	}

	// Compute diff between original and final to confirm change detection
	original := map[string]string{"EXISTING_KEY": "old_value"}
	diff := env.Diff(original, result)

	if _, ok := diff.Added["BRAND_NEW_KEY"]; !ok {
		t.Error("expected BRAND_NEW_KEY to appear in diff as added")
	}
	if _, ok := diff.Changed["EXISTING_KEY"]; !ok {
		t.Error("expected EXISTING_KEY to appear in diff as changed")
	}
}
