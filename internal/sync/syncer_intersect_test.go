package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_IntersectKeepsOnlyCommonKeys(t *testing.T) {
	secrets := map[string]string{
		"KEEP":   "vault_val",
		"REMOVE": "vault_val2",
	}

	f := tempEnvFile(t, "KEEP=existing\nLOCAL_ONLY=local\n")

	s := New(vault.NewMockClient(map[string]map[string]string{
		"secret/app": secrets,
	}), Options{
		Paths:   []string{"secret/app"},
		OutFile: f,
		Intersect: &env.IntersectOptions{KeepValue: "b"},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f)
	if _, ok := result["KEEP"]; !ok {
		t.Error("expected KEEP in output")
	}
	if _, ok := result["REMOVE"]; ok {
		t.Error("REMOVE should have been excluded by intersect")
	}
}

func TestRun_IntersectNilSkipsIntersection(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}

	f := tempEnvFile(t, "")

	s := New(vault.NewMockClient(map[string]map[string]string{
		"secret/app": secrets,
	}), Options{
		Paths:   []string{"secret/app"},
		OutFile: f,
		Intersect: nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, _ := env.Read(f)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	_ = os.Remove(f)
}
