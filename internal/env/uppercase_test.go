package env

import (
	"testing"
)

func TestUppercase_NoOptions(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "db_port": "5432"}
	out, res := Uppercase(secrets, UppercaseOptions{})
	if len(res.Changed) != 0 {
		t.Fatalf("expected no changes, got %v", res.Changed)
	}
	if out["db_host"] != "localhost" {
		t.Fatalf("unexpected value: %s", out["db_host"])
	}
}

func TestUppercase_UppercasesKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	out, res := Uppercase(secrets, UppercaseOptions{Keys: true})
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatal("expected DB_HOST key")
	}
	if len(res.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(res.Changed))
	}
}

func TestUppercase_UppercasesValues(t *testing.T) {
	secrets := map[string]string{"mode": "debug"}
	out, _ := Uppercase(secrets, UppercaseOptions{Values: true})
	if out["mode"] != "DEBUG" {
		t.Fatalf("expected DEBUG, got %s", out["mode"])
	}
}

func TestUppercase_OnlyKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost", "app_env": "prod"}
	out, res := Uppercase(secrets, UppercaseOptions{Keys: true, OnlyKeys: []string{"db_host"}})
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatal("expected DB_HOST")
	}
	if _, ok := out["app_env"]; !ok {
		t.Fatal("expected app_env unchanged")
	}
	if len(res.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(res.Changed))
	}
}

func TestUppercase_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"key": "value"}
	Uppercase(secrets, UppercaseOptions{Keys: true, Values: true})
	if _, ok := secrets["key"]; !ok {
		t.Fatal("input map was mutated")
	}
}

func TestUppercase_SummaryNoChanges(t *testing.T) {
	_, res := Uppercase(map[string]string{"A": "B"}, UppercaseOptions{})
	if res.Summary() != "uppercase: no changes" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestUppercase_SummaryWithChanges(t *testing.T) {
	_, res := Uppercase(map[string]string{"a": "b"}, UppercaseOptions{Keys: true})
	s := res.Summary()
	if s == "uppercase: no changes" {
		t.Fatal("expected non-empty summary")
	}
}
