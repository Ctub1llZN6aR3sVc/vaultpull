package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_SwapRenamesKey(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"OLD_TOKEN": "abc123"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: f,
		Swap: &env.SwapOptions{
			Pairs: map[string]string{"OLD_TOKEN": "APP_TOKEN"},
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(f)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if _, ok := result["OLD_TOKEN"]; ok {
		t.Error("expected OLD_TOKEN to be removed after swap")
	}
	if result["APP_TOKEN"] != "abc123" {
		t.Errorf("expected APP_TOKEN=abc123, got %q", result["APP_TOKEN"])
	}
}

func TestRun_SwapNilSkipsSwap(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"MY_KEY": "value"},
	})

	f := tempEnvFile(t)
	defer os.Remove(f)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		EnvFile: f,
		Swap:    nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(f)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if result["MY_KEY"] != "value" {
		t.Errorf("expected MY_KEY=value, got %q", result["MY_KEY"])
	}
}
