package env

import (
	"strings"
	"testing"
)

func TestTypeCheck_AllValid(t *testing.T) {
	secrets := map[string]string{
		"PORT":    "8080",
		"RATIO":   "0.75",
		"ENABLED": "true",
		"NAME":    "vaultpull",
	}
	rules := []TypeRule{
		{Key: "PORT", Expected: "int"},
		{Key: "RATIO", Expected: "float"},
		{Key: "ENABLED", Expected: "bool"},
		{Key: "NAME", Expected: "string"},
	}
	res := TypeCheck(secrets, rules)
	if !res.OK() {
		t.Fatalf("expected no violations, got: %s", res.Summary())
	}
}

func TestTypeCheck_IntViolation(t *testing.T) {
	secrets := map[string]string{"PORT": "not-a-number"}
	rules := []TypeRule{{Key: "PORT", Expected: "int"}}
	res := TypeCheck(secrets, rules)
	if res.OK() {
		t.Fatal("expected violation for non-int PORT")
	}
	if res.Violations[0].Key != "PORT" {
		t.Errorf("expected violation key PORT, got %q", res.Violations[0].Key)
	}
}

func TestTypeCheck_BoolViolation(t *testing.T) {
	secrets := map[string]string{"ENABLED": "yes-please"}
	rules := []TypeRule{{Key: "ENABLED", Expected: "bool"}}
	res := TypeCheck(secrets, rules)
	if res.OK() {
		t.Fatal("expected violation for invalid bool")
	}
}

func TestTypeCheck_MissingKeySkipped(t *testing.T) {
	secrets := map[string]string{"OTHER": "value"}
	rules := []TypeRule{{Key: "PORT", Expected: "int"}}
	res := TypeCheck(secrets, rules)
	if !res.OK() {
		t.Fatalf("expected no violations for missing key, got: %s", res.Summary())
	}
}

func TestTypeCheck_MultipleViolations(t *testing.T) {
	secrets := map[string]string{
		"PORT":    "abc",
		"TIMEOUT": "xyz",
	}
	rules := []TypeRule{
		{Key: "PORT", Expected: "int"},
		{Key: "TIMEOUT", Expected: "float"},
	}
	res := TypeCheck(secrets, rules)
	if len(res.Violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(res.Violations))
	}
}

func TestTypeCheck_SummaryContainsKey(t *testing.T) {
	secrets := map[string]string{"PORT": "bad"}
	rules := []TypeRule{{Key: "PORT", Expected: "int"}}
	res := TypeCheck(secrets, rules)
	if !strings.Contains(res.Summary(), "PORT") {
		t.Errorf("expected summary to mention PORT, got: %s", res.Summary())
	}
}

func TestTypeCheck_SummaryOK(t *testing.T) {
	res := TypeCheckResult{}
	if !strings.Contains(res.Summary(), "all values match") {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
