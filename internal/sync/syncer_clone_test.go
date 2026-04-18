package sync

import (
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestClone_CopiesVaultSecretsToSecondFile(t *testing.T) {
	src := tempEnvFile(t, "")
	dst := tempEnvFile(t, "EXISTING=keep\n")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc"},
	})

	s := New(client, []string{"secret/app"}, src)
	if err := s.Run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	srcSecrets, _ := env.Read(src)
	dstSecrets, _ := env.Read(dst)

	out, res := env.Clone(srcSecrets, dstSecrets, env.CloneOptions{Overwrite: false})
	if out["DB_PASS"] != "s3cr3t" {
		t.Fatalf("expected DB_PASS cloned, got %v", out)
	}
	if out["EXISTING"] != "keep" {
		t.Fatal("expected EXISTING preserved")
	}
	if len(res.Cloned) == 0 {
		t.Fatal("expected at least one cloned key")
	}
}

func TestClone_OverwriteReplacesDestination(t *testing.T) {
	src := tempEnvFile(t, "")
	dst := tempEnvFile(t, "DB_PASS=old\n")

	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "new"},
	})

	s := New(client, []string{"secret/app"}, src)
	if err := s.Run(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	srcSecrets, _ := env.Read(src)
	dstSecrets, _ := env.Read(dst)

	out, _ := env.Clone(srcSecrets, dstSecrets, env.CloneOptions{Overwrite: true})
	if out["DB_PASS"] != "new" {
		t.Fatalf("expected DB_PASS=new, got %s", out["DB_PASS"])
	}
}
