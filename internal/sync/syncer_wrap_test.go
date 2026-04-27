package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_WrapAddsAffixes(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "abc123"},
	})

	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	s := New(client, []string{"secret/app"}, tmp.Name())
	s.Wrap = &env.WrapOptions{Prefix: "[", Suffix: "]"}

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := env.Read(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	if result["API_KEY"] != "[abc123]" {
		t.Errorf("expected [abc123], got %q", result["API_KEY"])
	}
}

func TestRun_WrapNilSkipsWrapping(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"TOKEN": "raw"},
	})

	tmp, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	s := New(client, []string{"secret/app"}, tmp.Name())
	s.Wrap = nil

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := env.Read(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	if result["TOKEN"] != "raw" {
		t.Errorf("expected raw, got %q", result["TOKEN"])
	}
}
