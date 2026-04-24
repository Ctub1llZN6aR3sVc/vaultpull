package sync

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_MaskWriteRedactsSensitiveKeys(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "supersecret", "HOST": "localhost"},
	})

	tmp := tempEnvFile(t)

	s := New(client, Options{
		Paths: []string{"secret/app"},
		OutputFile: tmp,
		MaskWrite: &env.MaskWriteOptions{
			AutoDetect: true,
		},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	content := string(data)

	if contains(content, "supersecret") {
		t.Errorf("expected API_KEY to be masked, but found plaintext value")
	}
	if !contains(content, "***") {
		t.Errorf("expected masked placeholder in output")
	}
	if !contains(content, "localhost") {
		t.Errorf("expected HOST=localhost to be unmasked")
	}
}

func TestRun_MaskWriteNilSkipsMasking(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "supersecret"},
	})

	tmp := tempEnvFile(t)

	s := New(client, Options{
		Paths: []string{"secret/app"},
		OutputFile: tmp,
		MaskWrite: nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !contains(string(data), "supersecret") {
		t.Errorf("expected plaintext value when MaskWrite is nil")
	}
}
