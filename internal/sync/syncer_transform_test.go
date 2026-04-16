package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/toughtackle/vaultpull/internal/env"
	"github.com/toughtackle/vaultpull/internal/vault"
)

func TestRun_TransformUppercase(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"db_host": "localhost"},
	})

	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	s := New(client, []string{"secret/app"}, outFile)
	s.Transform = env.TransformOptions{Uppercase: true}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "DB_HOST") {
		t.Fatalf("expected DB_HOST in output, got: %s", data)
	}
}

func TestRun_TransformPrefix(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"TOKEN": "abc123"},
	})

	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	s := New(client, []string{"secret/app"}, outFile)
	s.Transform = env.TransformOptions{Prefix: "MYAPP_"}

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "MYAPP_TOKEN") {
		t.Fatalf("expected MYAPP_TOKEN in output, got: %s", data)
	}
}

func TestRun_NoTransformByDefault(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"my_key": "val"},
	})

	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	s := New(client, []string{"secret/app"}, outFile)

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "my_key") {
		t.Fatalf("expected my_key unchanged, got: %s", data)
	}
}
