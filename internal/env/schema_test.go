package env

import (
	"strings"
	"testing"
)

func TestValidateSchema_AllPresent(t *testing.T) {
	fields := []SchemaField{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PASS", Required: true},
	}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PASS": "secret"}
	r := ValidateSchema(secrets, fields, false)
	if r.HasError {
		t.Errorf("expected no error, got missing: %v", r.Missing)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	fields := []SchemaField{
		{Key: "DB_HOST", Required: true},
		{Key: "DB_PASS", Required: true},
	}
	secrets := map[string]string{"DB_HOST": "localhost"}
	r := ValidateSchema(secrets, fields, false)
	if !r.HasError {
		t.Fatal("expected error for missing key")
	}
	if len(r.Missing) != 1 || r.Missing[0] != "DB_PASS" {
		t.Errorf("unexpected missing: %v", r.Missing)
	}
}

func TestValidateSchema_StrictDetectsExtra(t *testing.T) {
	fields := []SchemaField{
		{Key: "DB_HOST", Required: true},
	}
	secrets := map[string]string{"DB_HOST": "localhost", "UNKNOWN_KEY": "val"}
	r := ValidateSchema(secrets, fields, true)
	if len(r.Extra) != 1 || r.Extra[0] != "UNKNOWN_KEY" {
		t.Errorf("expected UNKNOWN_KEY in extra, got %v", r.Extra)
	}
}

func TestValidateSchema_NonStrictIgnoresExtra(t *testing.T) {
	fields := []SchemaField{{Key: "DB_HOST", Required: true}}
	secrets := map[string]string{"DB_HOST": "localhost", "EXTRA": "x"}
	r := ValidateSchema(secrets, fields, false)
	if len(r.Extra) != 0 {
		t.Errorf("expected no extras in non-strict mode")
	}
}

func TestSchemaResult_Summary_OK(t *testing.T) {
	r := SchemaResult{}
	if !strings.Contains(r.Summary(), "OK") {
		t.Errorf("expected OK in summary")
	}
}

func TestSchemaResult_Summary_WithErrors(t *testing.T) {
	r := SchemaResult{Missing: []string{"API_KEY"}, HasError: true}
	s := r.Summary()
	if !strings.Contains(s, "API_KEY") {
		t.Errorf("expected API_KEY in summary, got: %s", s)
	}
}
