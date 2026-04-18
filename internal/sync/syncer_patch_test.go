package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_PatchExistingOnlySkipsNewVaultKeys(t *testing.T) {
	path := tempEnvFile(t, "EXISTING=old\n")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"EXISTING": "updated", "NEW_KEY": "value"},
	})

	s := New(client, []PathMapping{{VaultPath: "secret/app", EnvFile: path}},
		Options{
			Patch: &env.PatchOptions{ExistingOnly: true},
		},
	)
	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !containsStr(content, "EXISTING=updated") {
		t.Errorf("expected EXISTING to be updated, got:\n%s", content)
	}
	if containsStr(content, "NEW_KEY") {
		t.Errorf("NEW_KEY should be skipped with ExistingOnly, got:\n%s", content)
	}
}

func TestRun_PatchNilSkipsPatching(t *testing.T) {
	path := tempEnvFile(t, "")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"FOO": "bar"},
	})

	s := New(client, []PathMapping{{VaultPath: "secret/app", EnvFile: path}},
		Options{Patch: nil},
	)
	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !containsStr(string(data), "FOO=bar") {
		t.Errorf("expected FOO=bar in output")
	}
}
