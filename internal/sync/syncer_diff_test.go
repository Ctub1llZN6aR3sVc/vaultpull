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

	// Pre-populate env file with an existing key
	err := os.WriteFile(envFile, []byte("EXISTING=value\n"), 0600)
	if err != nil {
		t.Fatalf("failed to write initial env file: %v", err)
	}

	mockClient := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"EXISTING": "value",
			"NEW_KEY":  "new_value",
		},
	})

	syncer := New(mockClient, envFile)
	diff, err := syncer.Run([]string{"secret/app"})
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if diff == nil {
		t.Fatal("expected a diff result, got nil")
	}

	if diff.IsEmpty() {
		t.Fatal("expected non-empty diff")
	}

	found := false
	for _, c := range diff.Changes {
		if c.Key == "NEW_KEY" && c.Type == env.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW_KEY to appear as added in diff")
	}
}
