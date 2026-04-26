package sync

import (
	"os"
	"testing"

	"github.com/eliziario/vaultpull/internal/vault"
)

func TestRun_AliasCreatesNewKey(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "hunter2"},
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Alias: map[string]string{
			"DATABASE_PASSWORD": "DB_PASS",
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := readEnvFile(t, tmp)
	if env["DATABASE_PASSWORD"] != "hunter2" {
		t.Fatalf("expected alias key, got %q", env["DATABASE_PASSWORD"])
	}
	if env["DB_PASS"] != "hunter2" {
		t.Fatal("expected source key to remain")
	}
}

func TestRun_AliasNilSkipsAliasing(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"FOO": "bar"},
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Alias:   nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := readEnvFile(t, tmp)
	if env["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %q", env["FOO"])
	}
}

// TestRun_AliasUnknownSourceKeyIgnored verifies that an alias referencing a
// source key that does not exist in the fetched secrets is silently skipped,
// leaving the output file without the alias key.
func TestRun_AliasUnknownSourceKeyIgnored(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"FOO": "bar"},
	})

	tmp := tempEnvFile(t)
	defer os.Remove(tmp)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Alias: map[string]string{
			"ALIAS_KEY": "NONEXISTENT",
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := readEnvFile(t, tmp)
	if _, ok := env["ALIAS_KEY"]; ok {
		t.Fatal("expected alias key to be absent when source key does not exist")
	}
	if env["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %q", env["FOO"])
	}
}
