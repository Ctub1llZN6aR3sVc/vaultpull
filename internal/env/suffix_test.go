package env

import (
	"testing"
)

func TestSuffix_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Suffix(secrets, SuffixOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Affected) != 0 {
		t.Errorf("expected no affected keys, got %v", res.Affected)
	}
	if res.Secrets["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", res.Secrets["FOO"])
	}
}

func TestSuffix_AddSuffix(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	res, err := Suffix(secrets, SuffixOptions{Add: "_V2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Secrets["DB_HOST_V2"]; !ok {
		t.Error("expected key DB_HOST_V2")
	}
	if _, ok := res.Secrets["DB_PORT_V2"]; !ok {
		t.Error("expected key DB_PORT_V2")
	}
	if len(res.Affected) != 2 {
		t.Errorf("expected 2 affected, got %d", len(res.Affected))
	}
}

func TestSuffix_StripSuffix(t *testing.T) {
	secrets := map[string]string{"HOST_OLD": "localhost", "PORT_OLD": "5432", "NAME": "db"}
	res, err := Suffix(secrets, SuffixOptions{Strip: "_OLD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Secrets["HOST"]; !ok {
		t.Error("expected key HOST after strip")
	}
	if _, ok := res.Secrets["PORT"]; !ok {
		t.Error("expected key PORT after strip")
	}
	if res.Secrets["NAME"] != "db" {
		t.Error("expected NAME to be unchanged")
	}
}

func TestSuffix_FailOnConflict(t *testing.T) {
	secrets := map[string]string{"FOO": "a", "FOO_V2": "b"}
	_, err := Suffix(secrets, SuffixOptions{Add: "_V2", FailOnConflict: true})
	if err == nil {
		t.Error("expected conflict error, got nil")
	}
}

func TestSuffix_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	_, err := Suffix(secrets, SuffixOptions{Add: "_NEW", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := secrets["KEY"]; !ok {
		t.Error("dry run should not mutate original map")
	}
	if _, ok := secrets["KEY_NEW"]; ok {
		t.Error("dry run should not add new key to original map")
	}
}

func TestSuffix_RestrictedToKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "a", "BAR": "b", "BAZ": "c"}
	res, err := Suffix(secrets, SuffixOptions{Add: "_X", Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Secrets["FOO_X"]; !ok {
		t.Error("expected FOO_X")
	}
	if _, ok := res.Secrets["BAR"]; !ok {
		t.Error("expected BAR to remain unchanged")
	}
	if _, ok := res.Secrets["BAR_X"]; ok {
		t.Error("BAR should not be suffixed")
	}
}

func TestSuffix_Summary(t *testing.T) {
	res := SuffixResult{Affected: []string{"A", "B"}}
	if res.Summary() == "suffix: no keys affected" {
		t.Error("expected non-empty summary")
	}
	empty := SuffixResult{}
	if empty.Summary() != "suffix: no keys affected" {
		t.Errorf("unexpected summary: %s", empty.Summary())
	}
}
