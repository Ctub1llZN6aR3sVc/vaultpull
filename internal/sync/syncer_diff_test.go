package sync

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_DiffReportsNewKeys(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Pre-populate with an existing key.
	if err := os.WriteFile(envFile, []byte("EXISTING=old\n"), 0600); err != nil {
		t.Fatal(err)
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"EXISTING": "new",
			"BRAND_NEW": "value",
		},
	})

	profile := config.Profile{
		EnvFile: envFile,
		Paths:   []string{"secret/app"},
	}

	var buf bytes.Buffer
	s := New(client, profile, &buf)
	diff, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff.IsEmpty() {
		t.Fatal("expected non-empty diff")
	}

	foundAdded, foundChanged := false, false
	for _, c := range diff.Changes {
		if c.Key == "BRAND_NEW" && c.Type == "added" {
			foundAdded = true
		}
		if c.Key == "EXISTING" && c.Type == "changed" {
			foundChanged = true
		}
	}
	if !foundAdded {
		t.Error("expected BRAND_NEW to be reported as added")
	}
	if !foundChanged {
		t.Error("expected EXISTING to be reported as changed")
	}
}

func TestRun_DiffEmptyWhenNoChanges(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Pre-populate with the same value that Vault will return.
	if err := os.WriteFile(envFile, []byte("EXISTING=same\n"), 0600); err != nil {
		t.Fatal(err)
	}

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"EXISTING": "same",
		},
	})

	profile := config.Profile{
		EnvFile: envFile,
		Paths:   []string{"secret/app"},
	}

	var buf bytes.Buffer
	s := New(client, profile, &buf)
	diff, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !diff.IsEmpty() {
		t.Errorf("expected empty diff, got %d change(s)", len(diff.Changes))
	}
}
