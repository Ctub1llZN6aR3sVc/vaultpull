package env

import (
	"testing"
)

func TestCollect_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := Collect(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Out["FOO"] != "bar" || res.Out["BAZ"] != "qux" {
		t.Errorf("expected passthrough, got %v", res.Out)
	}
}

func TestCollect_GroupByKey(t *testing.T) {
	secrets := map[string]string{
		"ENV":    "prod",
		"DB_URL": "postgres://localhost",
		"PORT":   "5432",
	}
	res, err := Collect(secrets, &CollectOptions{GroupBy: "ENV"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Out["PROD_DB_URL"] != "postgres://localhost" {
		t.Errorf("expected PROD_DB_URL, got keys: %v", res.Out)
	}
	if res.Out["PROD_PORT"] != "5432" {
		t.Errorf("expected PROD_PORT, got keys: %v", res.Out)
	}
	if res.Out["ENV"] != "prod" {
		t.Errorf("group-by key should be preserved, got %v", res.Out)
	}
}

func TestCollect_CustomSeparator(t *testing.T) {
	secrets := map[string]string{
		"NS":  "app",
		"KEY": "value",
	}
	res, err := Collect(secrets, &CollectOptions{GroupBy: "NS", Separator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Out["APP.KEY"] != "value" {
		t.Errorf("expected APP.KEY, got %v", res.Out)
	}
}

func TestCollect_MissingGroupByKey(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	_, err := Collect(secrets, &CollectOptions{GroupBy: "MISSING"})
	if err == nil {
		t.Fatal("expected error for missing group-by key")
	}
}

func TestCollect_GroupsPopulated(t *testing.T) {
	secrets := map[string]string{
		"STAGE": "staging",
		"HOST":  "localhost",
	}
	res, err := Collect(secrets, &CollectOptions{GroupBy: "STAGE"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Groups) != 1 || res.Groups[0] != "STAGING" {
		t.Errorf("expected groups [STAGING], got %v", res.Groups)
	}
}

func TestCollect_SummaryString(t *testing.T) {
	secrets := map[string]string{
		"ENV": "dev",
		"A":   "1",
		"B":   "2",
	}
	res, err := Collect(secrets, &CollectOptions{GroupBy: "ENV"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCollect_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"ENV": "test", "X": "1"}
	orig := map[string]string{"ENV": "test", "X": "1"}
	_, err := Collect(secrets, &CollectOptions{GroupBy: "ENV"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range orig {
		if secrets[k] != v {
			t.Errorf("input mutated at key %q", k)
		}
	}
}
