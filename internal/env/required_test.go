package env

import (
	"testing"
)

func TestRequired_AllPresent(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
	res, err := Required(secrets, RequiredOptions{Keys: []string{"DB_HOST", "DB_PORT", "API_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK() {
		t.Fatalf("expected OK, got %s", res.Summary())
	}
}

func TestRequired_MissingKey(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
	}
	res, err := Required(secrets, RequiredOptions{Keys: []string{"DB_HOST", "DB_PASS"}})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_PASS" {
		t.Fatalf("expected DB_PASS in Missing, got %v", res.Missing)
	}
}

func TestRequired_EmptyValue(t *testing.T) {
	secrets := map[string]string{
		"API_KEY": "   ",
	}
	res, err := Required(secrets, RequiredOptions{Keys: []string{"API_KEY"}})
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	if len(res.Empty) != 1 || res.Empty[0] != "API_KEY" {
		t.Fatalf("expected API_KEY in Empty, got %v", res.Empty)
	}
}

func TestRequired_MixedViolations(t *testing.T) {
	secrets := map[string]string{
		"A": "",
	}
	res, err := Required(secrets, RequiredOptions{Keys: []string{"A", "B"}})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(res.Empty) != 1 {
		t.Fatalf("expected 1 empty, got %d", len(res.Empty))
	}
	if len(res.Missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(res.Missing))
	}
}

func TestRequired_EmptyKeyList(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	res, err := Required(secrets, RequiredOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK() {
		t.Fatal("expected OK for empty key list")
	}
}

func TestRequired_SummaryOK(t *testing.T) {
	res := RequiredResult{}
	if res.Summary() != "all required keys present" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestRequired_SummaryWithViolations(t *testing.T) {
	res := RequiredResult{Missing: []string{"X"}, Empty: []string{"Y"}}
	s := res.Summary()
	if s == "" || s == "all required keys present" {
		t.Fatalf("unexpected summary: %s", s)
	}
}

func TestRequired_NilSecrets(t *testing.T) {
	res, err := Required(nil, RequiredOptions{Keys: []string{"DB_HOST"}})
	if err == nil {
		t.Fatal("expected error for nil secrets map")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_HOST" {
		t.Fatalf("expected DB_HOST in Missing, got %v", res.Missing)
	}
}
