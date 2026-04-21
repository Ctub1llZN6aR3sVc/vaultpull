package sync

import (
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_TypeCheckPassesWithValidTypes(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"PORT": "8080", "DEBUG": "true"},
	})

	tmp := tempEnvFile(t)

	rules := []env.TypeRule{
		{Key: "PORT", Expected: "int"},
		{Key: "DEBUG", Expected: "bool"},
	}

	s := New(client, []string{"secret/app"}, tmp)
	s.TypeRules = rules

	if err := s.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_TypeCheckFailsOnBadValue(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"PORT": "not-a-port"},
	})

	tmp := tempEnvFile(t)

	rules := []env.TypeRule{
		{Key: "PORT", Expected: "int"},
	}

	s := New(client, []string{"secret/app"}, tmp)
	s.TypeRules = rules
	s.StrictTypeCheck = true

	if err := s.Run(); err == nil {
		t.Fatal("expected error for type violation, got nil")
	}
}

func TestRun_TypeCheckNilSkipsCheck(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"PORT": "not-a-port"},
	})

	tmp := tempEnvFile(t)

	s := New(client, []string{"secret/app"}, tmp)
	s.TypeRules = nil
	s.StrictTypeCheck = true

	if err := s.Run(); err != nil {
		t.Fatalf("expected no error when TypeRules is nil, got: %v", err)
	}
}
