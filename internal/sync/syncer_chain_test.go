package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestChain_VaultOverridesLocal(t *testing.T) {
	vaultSecrets := map[string]string{"API_KEY": "vault_val", "DB_PASS": "vault_db"}
	local := map[string]string{"API_KEY": "local_val", "LOCAL_ONLY": "local"}

	out := env.ChainAll(vaultSecrets, local)

	if out["API_KEY"] != "vault_val" {
		t.Errorf("vault should win, got %s", out["API_KEY"])
	}
	if out["LOCAL_ONLY"] != "local" {
		t.Errorf("local-only key should be present, got %s", out["LOCAL_ONLY"])
	}
	if out["DB_PASS"] != "vault_db" {
		t.Errorf("expected vault_db, got %s", out["DB_PASS"])
	}
}

func TestRun_ChainFallbackFromLocal(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")

	// Write a local .env with a fallback value
	_ = os.WriteFile(envFile, []byte("FALLBACK_KEY=local_fallback\n"), 0600)

	mock := vault.NewMockClient(map[string]string{
		"NEW_KEY": "from_vault",
	})

	s := New(mock, envFile)
	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := env.Read(envFile)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if result["NEW_KEY"] != "from_vault" {
		t.Errorf("expected from_vault, got %s", result["NEW_KEY"])
	}
}
