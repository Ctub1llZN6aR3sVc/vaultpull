package env

import (
	"testing"
)

func TestExtract_ByKeys(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	res, err := Extract(secrets, ExtractOptions{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["A"] != "1" || res.Secrets["C"] != "3" {
		t.Errorf("unexpected secrets: %v", res.Secrets)
	}
	if _, ok := res.Secrets["B"]; ok {
		t.Error("B should not be present")
	}
}

func TestExtract_ByPrefix(t *testing.T) {
	secrets := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres"}
	res, err := Extract(secrets, ExtractOptions{Prefixes: []string{"APP_"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["APP_HOST"] != "localhost" || res.Secrets["APP_PORT"] != "8080" {
		t.Errorf("unexpected secrets: %v", res.Secrets)
	}
	if _, ok := res.Secrets["DB_URL"]; ok {
		t.Error("DB_URL should not be present")
	}
}

func TestExtract_StripPrefix(t *testing.T) {
	secrets := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	res, err := Extract(secrets, ExtractOptions{Prefixes: []string{"APP_"}, StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["HOST"] != "localhost" || res.Secrets["PORT"] != "8080" {
		t.Errorf("unexpected secrets: %v", res.Secrets)
	}
}

func TestExtract_MissingKeyRecorded(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	res, err := Extract(secrets, ExtractOptions{Keys: []string{"A", "MISSING"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING" {
		t.Errorf("expected missing=[MISSING], got %v", res.Missing)
	}
}

func TestExtract_FailOnMissing(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	_, err := Extract(secrets, ExtractOptions{Keys: []string{"A", "GONE"}, FailOnMissing: true})
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestExtract_SummaryOK(t *testing.T) {
	res := ExtractResult{Secrets: map[string]string{"A": "1"}, Missing: nil}
	if res.Summary() != "extract: ok" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestExtract_SummaryWithMissing(t *testing.T) {
	res := ExtractResult{Secrets: map[string]string{}, Missing: []string{"X", "Y"}}
	s := res.Summary()
	if s != "extract: missing keys: X, Y" {
		t.Errorf("unexpected summary: %s", s)
	}
}
