package env

import (
	"testing"
)

func TestValidate_ValidSecrets(t *testing.T) {
	secrets := map[string]string{
		"APP_NAME":     "vaultpull",
		"DATABASE_URL": "postgres://localhost/db",
	}
	result := Validate(secrets)
	if !result.IsValid() {
		t.Fatalf("expected valid, got errors: %s", result.Summary())
	}
}

func TestValidate_EmptyValueForSensitiveKey(t *testing.T) {
	secrets := map[string]string{
		"API_SECRET": "",
	}
	result := Validate(secrets)
	if result.IsValid() {
		t.Fatal("expected validation error for empty sensitive value")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "API_SECRET" {
		t.Errorf("unexpected error key: %s", result.Errors[0].Key)
	}
}

func TestValidate_KeyWithWhitespace(t *testing.T) {
	secrets := map[string]string{
		"BAD KEY": "value",
	}
	result := Validate(secrets)
	if result.IsValid() {
		t.Fatal("expected validation error for key with whitespace")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	secrets := map[string]string{
		"DB PASSWORD": "",
		"API_TOKEN":   "",
	}
	result := Validate(secrets)
	if result.IsValid() {
		t.Fatal("expected validation errors")
	}
	if len(result.Errors) < 2 {
		t.Fatalf("expected at least 2 errors, got %d", len(result.Errors))
	}
}

func TestValidate_SummaryNoErrors(t *testing.T) {
	result := ValidationResult{}
	if result.Summary() != "all secrets valid" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestValidate_SummaryWithErrors(t *testing.T) {
	secrets := map[string]string{
		"API_SECRET": "",
	}
	result := Validate(secrets)
	summary := result.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	if result.IsValid() {
		t.Error("expected invalid result")
	}
}
