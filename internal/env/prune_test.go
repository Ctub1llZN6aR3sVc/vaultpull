package env_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
)

func TestPrune_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, res := env.Prune(secrets, env.PruneOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if len(res.Pruned) != 0 {
		t.Fatalf("expected no pruned keys, got %v", res.Pruned)
	}
}

func TestPrune_RemoveEmpty(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "EMPTY": "", "ALSO_EMPTY": ""}
	out, res := env.Prune(secrets, env.PruneOptions{RemoveEmpty: true})
	if _, ok := out["FOO"]; !ok {
		t.Error("expected FOO to be kept")
	}
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if len(res.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(res.Pruned))
	}
}

func TestPrune_RemoveKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	out, res := env.Prune(secrets, env.PruneOptions{RemoveKeys: []string{"FOO", "BAZ"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key remaining, got %d", len(out))
	}
	if out["BAR"] != "2" {
		t.Error("expected BAR to remain")
	}
	if len(res.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(res.Pruned))
	}
}

func TestPrune_RemovePrefix(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_NAME": "myapp"}
	out, res := env.Prune(secrets, env.PruneOptions{RemovePrefix: []string{"DB_"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key remaining, got %d", len(out))
	}
	if out["APP_NAME"] != "myapp" {
		t.Error("expected APP_NAME to remain")
	}
	if len(res.Pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(res.Pruned))
	}
}

func TestPrune_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"FOO": "", "BAR": "val"}
	orig := map[string]string{"FOO": "", "BAR": "val"}
	env.Prune(secrets, env.PruneOptions{RemoveEmpty: true})
	for k, v := range orig {
		if secrets[k] != v {
			t.Errorf("input mutated at key %s", k)
		}
	}
}

func TestPrune_SummaryWithPruned(t *testing.T) {
	secrets := map[string]string{"FOO": "", "BAR": "val"}
	_, res := env.Prune(secrets, env.PruneOptions{RemoveEmpty: true})
	s := res.Summary()
	if s == "prune: nothing removed" {
		t.Error("expected a non-empty summary")
	}
}

func TestPrune_SummaryNoPruned(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	_, res := env.Prune(secrets, env.PruneOptions{})
	if res.Summary() != "prune: nothing removed" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
