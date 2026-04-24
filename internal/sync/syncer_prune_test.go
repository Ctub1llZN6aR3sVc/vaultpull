package sync_test

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_PruneRemovesEmptyValues(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_HOST": "localhost", "DB_PASS": ""},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := sync.New(client, sync.Options{
		Paths:  []string{"secret/app"},
		Output: f,
		Prune:  &env.PruneOptions{RemoveEmpty: true},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := readEnvMap(t, f)
	if _, ok := got["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if _, ok := got["DB_PASS"]; ok {
		t.Error("expected DB_PASS to be pruned (empty value)")
	}
}

func TestRun_PruneNilSkipsPruning(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_HOST": "localhost", "DB_PASS": ""},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := sync.New(client, sync.Options{
		Paths:  []string{"secret/app"},
		Output: f,
		Prune:  nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := readEnvMap(t, f)
	if _, ok := got["DB_PASS"]; !ok {
		t.Error("expected DB_PASS to be written when prune is nil")
	}
}
