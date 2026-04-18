package sync

import (
	"testing"

	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

func TestRun_GroupPartitionsSecrets(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"DB_HOST":  "localhost",
			"DB_PORT":  "5432",
			"APP_NAME": "vaultpull",
		},
	})

	secrets, err := client.GetSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := env.Group(secrets, []string{"DB", "APP"}, "_")

	if res.Groups["DB"]["HOST"] != "localhost" {
		t.Errorf("expected DB_HOST in DB group")
	}
	if res.Groups["APP"]["NAME"] != "vaultpull" {
		t.Errorf("expected APP_NAME in APP group")
	}
	if len(res.Ungrouped) != 0 {
		t.Errorf("expected no ungrouped keys, got %d", len(res.Ungrouped))
	}
}

func TestRun_GroupUngroupedFallback(t *testing.T) {
	client := vault.NewMockClient(map[string]map[string]string{
		"secret/app": {
			"UNRELATED": "val",
		},
	})

	secrets, err := client.GetSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := env.Group(secrets, []string{"DB"}, "_")
	if res.Ungrouped["UNRELATED"] != "val" {
		t.Errorf("expected UNRELATED in ungrouped")
	}
	if len(res.Groups["DB"]) != 0 {
		t.Errorf("expected DB group to be empty")
	}
}
