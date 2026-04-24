package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_ReorderAlphabetical(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"ZEBRA": "z", "APPLE": "a", "MANGO": "m"},
	})

	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: f.Name(),
		Reorder: &env.ReorderOptions{Alphabetical: true},
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if result.ReorderResult == nil {
		t.Fatal("expected ReorderResult to be set")
	}
	if len(result.ReorderResult.Ordered) != 3 {
		t.Fatalf("expected 3 ordered keys, got %d", len(result.ReorderResult.Ordered))
	}
	if result.ReorderResult.Ordered[0] != "APPLE" {
		t.Fatalf("expected APPLE first, got %s", result.ReorderResult.Ordered[0])
	}
}

func TestRun_ReorderNilSkipsReorder(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"B": "2", "A": "1"},
	})

	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: f.Name(),
		Reorder: nil,
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if result.ReorderResult != nil {
		t.Fatal("expected ReorderResult to be nil when Reorder is nil")
	}
}
