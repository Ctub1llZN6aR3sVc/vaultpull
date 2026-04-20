package sync

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_TruncateShortensLongValues(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "supersecretlongtoken12345"},
	})

	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	s := New(client, []string{"secret/app"}, f.Name())
	s.Truncate = &env.TruncateOptions{MaxLength: 10}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	val := result["API_KEY"]
	if len(val) > 10 {
		t.Fatalf("expected value truncated to <=10 chars, got %q (len %d)", val, len(val))
	}
	if val != "supersec..." {
		t.Fatalf("expected 'supersec...', got %q", val)
	}
}

func TestRun_TruncateNilSkipsTruncation(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"LONG_KEY": "verylongvaluethatexceedslimit"},
	})

	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	s := New(client, []string{"secret/app"}, f.Name())
	// s.Truncate is nil by default

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	result, err := env.Read(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if result["LONG_KEY"] != "verylongvaluethatexceedslimit" {
		t.Fatalf("expected full value, got %q", result["LONG_KEY"])
	}
}
