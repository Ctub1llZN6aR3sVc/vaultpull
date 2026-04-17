package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_LintPassesWithCleanSecrets(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "abc123", "DB_HOST": "localhost"},
	})
	tmp := tempEnvFile(t)
	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
		Lint:       true,
	})
	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.LintResult.IsClean() {
		t.Errorf("expected clean lint, got: %v", result.LintResult.Issues)
	}
}

func TestRun_LintDetectsIssue(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"BAD KEY": "value"},
	})
	tmp := tempEnvFile(t)
	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
		Lint:       true,
	})
	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LintResult.IsClean() {
		t.Error("expected lint issue for key with spaces")
	}
}

func TestRun_LintSkippedWhenDisabled(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"BAD KEY": "value"},
	})
	tmp := tempEnvFile(t)
	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
		Lint:       false,
	})
	result, err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.LintResult.Issues) != 0 {
		t.Error("expected no lint issues when lint is disabled")
	}
}
