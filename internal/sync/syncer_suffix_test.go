package sync

import (
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_SuffixAddsToKeys(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_HOST": "localhost", "DB_PORT": "5432"},
	})

	tmp := tempEnvFile(t)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Suffix:  &env.SuffixOptions{Add: "_V2"},
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	got, err := env.Read(tmp)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if _, ok := got["DB_HOST_V2"]; !ok {
		t.Error("expected DB_HOST_V2 in output")
	}
	if _, ok := got["DB_PORT_V2"]; !ok {
		t.Error("expected DB_PORT_V2 in output")
	}
	if _, ok := got["DB_HOST"]; ok {
		t.Error("original key DB_HOST should not be present")
	}
}

func TestRun_SuffixNilSkipsSuffixing(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "secret123"},
	})

	tmp := tempEnvFile(t)

	s := New(client, Options{
		Paths:   []string{"secret/app"},
		OutFile: tmp,
		Suffix:  nil,
	})

	if err := s.Run(); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	got, err := env.Read(tmp)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if _, ok := got["API_KEY"]; !ok {
		t.Error("expected API_KEY unchanged when suffix is nil")
	}
}
