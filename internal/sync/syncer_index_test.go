package sync

import (
	"testing"

	"github.com/your/vaultpull/internal/env"
	"github.com/your/vaultpull/internal/vault"
)

func TestRun_IndexBuiltFromSecrets(t *testing.T) {
	tmp := tempEnvFile(t)
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "secret", "API_KEY": "key"},
	})
	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
	})
	res, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idx := env.Index(res.Secrets, &env.IndexOptions{Source: "vault"})
	if idx.Total != 2 {
		t.Fatalf("expected 2 indexed entries, got %d", idx.Total)
	}
	for _, e := range idx.Entries {
		if e.Source != "vault" {
			t.Fatalf("expected source vault, got %s", e.Source)
		}
	}
}

func TestRun_IndexEmptyWhenNoSecrets(t *testing.T) {
	tmp := tempEnvFile(t)
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/empty": {},
	})
	s := New(client, Options{
		Paths:   []string{"secret/empty"},
		OutFile: tmp,
	})
	res, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idx := env.Index(res.Secrets, nil)
	if idx.Total != 0 {
		t.Fatal("expected empty index")
	}
}
