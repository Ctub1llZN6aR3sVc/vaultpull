package sync

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/env"
	"github.com/yourusername/vaultpull/internal/vault"
)

func TestRun_ValidationFailsOnEmptySensitiveValue(t *testing.T) {
	mockSecrets := map[string]map[string]string{
		"secret/app": {
			"API_TOKEN": "",
		},
	}
	client := vault.NewMockClient(mockSecrets)

	tmpFile := tempEnvFile(t)

	s := New(client, []string{"secret/app"}, tmpFile)
	s.Validate = true

	err := s.Run()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestRun_ValidationPassesWithValidSecrets(t *testing.T) {
	mockSecrets := map[string]map[string]string{
		"secret/app": {
			"APP_NAME": "vaultpull",
			"API_TOKEN": "abc123xyz",
		},
	}
	client := vault.NewMockClient(mockSecrets)

	tmpFile := tempEnvFile(t)

	s := New(client, []string{"secret/app"}, tmpFile)
	s.Validate = true

	err := s.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ValidationSkippedWhenDisabled(t *testing.T) {
	mockSecrets := map[string]map[string]string{
		"secret/app": {
			"API_TOKEN": "",
		},
	}
	client := vault.NewMockClient(mockSecrets)

	tmpFile := tempEnvFile(t)

	s := New(client, []string{"secret/app"}, tmpFile)
	s.Validate = false

	err := s.Run()
	if err != nil {
		t.Fatalf("expected no error when validation disabled, got: %v", err)
	}
}

func TestValidate_DirectCall(t *testing.T) {
	secrets := map[string]string{
		"APP_ENV": "production",
	}
	result := env.Validate(secrets)
	if !result.IsValid() {
		t.Fatalf("expected valid secrets, got: %s", result.Summary())
	}
}
