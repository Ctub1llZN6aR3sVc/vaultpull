package sync

import (
	"os"
	"testing"

	"github.com/elizaos/vaultpull/internal/env"
	"github.com/elizaos/vaultpull/internal/vault"
)

var encKey = []byte("12345678901234567890123456789012")

func TestRun_EncryptsSecretsWhenEnabled(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASSWORD": "s3cr3t"},
	})
	tmp := tempEnvFile(t)
	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
		EncryptKey: encKey,
	})
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	existing, _ := env.Read(tmp)
	if existing["DB_PASSWORD"] == "s3cr3t" {
		t.Fatal("expected encrypted value, got plaintext")
	}
	dec, err := env.Decrypt(existing["DB_PASSWORD"], encKey)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(dec) != "s3cr3t" {
		t.Fatalf("expected s3cr3t, got %s", dec)
	}
}

func TestRun_SkipsEncryptionWhenKeyNil(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "plainvalue"},
	})
	tmp := tempEnvFile(t)
	s := New(client, Options{
		Paths:      []string{"secret/app"},
		OutputFile: tmp,
	})
	if err := s.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}
	existing, _ := env.Read(tmp)
	if existing["API_KEY"] != "plainvalue" {
		t.Fatalf("expected plainvalue, got %s", existing["API_KEY"])
	}
	os.Remove(tmp)
}
