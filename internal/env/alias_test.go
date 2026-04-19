package env

import (
	"testing"
)

func TestAlias_CopiesValue(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret"}
	out, res := Alias(secrets, AliasOptions{Aliases: map[string]string{"DATABASE_PASSWORD": "DB_PASS"}})
	if out["DATABASE_PASSWORD"] != "secret" {
		t.Fatalf("expected alias to be set, got %q", out["DATABASE_PASSWORD"])
	}
	if out["DB_PASS"] != "secret" {
		t.Fatal("expected source key to remain")
	}
	if len(res.Aliased) != 1 {
		t.Fatalf("expected 1 aliased, got %d", len(res.Aliased))
	}
}

func TestAlias_RemoveSource(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "val"}
	out, _ := Alias(secrets, AliasOptions{
		Aliases:      map[string]string{"NEW_KEY": "OLD_KEY"},
		RemoveSource: true,
	})
	if _, ok := out["OLD_KEY"]; ok {
		t.Fatal("expected source key to be removed")
	}
	if out["NEW_KEY"] != "val" {
		t.Fatal("expected new key to have value")
	}
}

func TestAlias_MissingSourceKey(t *testing.T) {
	secrets := map[string]string{"OTHER": "x"}
	_, res := Alias(secrets, AliasOptions{Aliases: map[string]string{"NEW": "MISSING"}})
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING" {
		t.Fatalf("expected missing key reported, got %v", res.Missing)
	}
	if len(res.Aliased) != 0 {
		t.Fatal("expected no aliased keys")
	}
}

func TestAlias_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	out, res := Alias(secrets, AliasOptions{
		Aliases: map[string]string{"B": "A"},
		DryRun:  true,
	})
	if _, ok := out["B"]; ok {
		t.Fatal("dry run should not add alias key")
	}
	if len(res.Aliased) != 1 {
		t.Fatal("dry run should still report what would be aliased")
	}
}

func TestAlias_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	Alias(secrets, AliasOptions{Aliases: map[string]string{"Z": "X"}, RemoveSource: true})
	if _, ok := secrets["X"]; !ok {
		t.Fatal("input map should not be mutated")
	}
}

func TestAlias_SummaryNoAliases(t *testing.T) {
	res := AliasResult{}
	if res.Summary() != "alias: no keys aliased" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}
