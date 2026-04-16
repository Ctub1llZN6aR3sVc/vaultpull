package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/vaultpull/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "vaultpull.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
default_profile: dev
profiles:
  dev:
    name: dev
    vault_addr: http://127.0.0.1:8200
    vault_path: secret/data/myapp/dev
    output_file: .env
    auth_method: token
`
	path := writeTempConfig(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultProfile != "dev" {
		t.Errorf("expected default_profile=dev, got %q", cfg.DefaultProfile)
	}
	p, err := cfg.GetProfile("")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if p.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected vault_addr: %q", p.VaultAddr)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/vaultpull.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	content := `
default_profile: dev
profiles:
  dev:
    name: dev
    vault_addr: http://127.0.0.1:8200
    vault_path: secret/data/myapp/dev
    output_file: .env
    auth_method: token
`
	path := writeTempConfig(t, content)
	cfg, _ := config.Load(path)
	_, err := cfg.GetProfile("staging")
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestGetProfile_ExplicitName(t *testing.T) {
	content := `
default_profile: dev
profiles:
  dev:
    name: dev
    vault_addr: http://127.0.0.1:8200
    vault_path: secret/data/myapp/dev
    output_file: .env
    auth_method: token
  staging:
    name: staging
    vault_addr: http://staging.example.com:8200
    vault_path: secret/data/myapp/staging
    output_file: .env.staging
    auth_method: token
`
	path := writeTempConfig(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := cfg.GetProfile("staging")
	if err != nil {
		t.Fatalf("GetProfile(staging): %v", err)
	}
	if p.VaultAddr != "http://staging.example.com:8200" {
		t.Errorf("unexpected vault_addr: %q", p.VaultAddr)
	}
}
