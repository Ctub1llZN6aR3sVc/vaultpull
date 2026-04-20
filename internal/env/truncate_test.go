package env

import (
	"strings"
	"testing"
)

func TestTruncate_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "hello", "B": "world"}
	out, res := Truncate(secrets, TruncateOptions{})
	if out["A"] != "hello" || out["B"] != "world" {
		t.Fatalf("expected unchanged map, got %v", out)
	}
	if len(res.Truncated) != 0 {
		t.Fatalf("expected no truncations, got %v", res.Truncated)
	}
}

func TestTruncate_ShortensLongValue(t *testing.T) {
	secrets := map[string]string{"KEY": "abcdefghij"}
	out, res := Truncate(secrets, TruncateOptions{MaxLength: 5})
	if out["KEY"] != "ab..." {
		t.Fatalf("expected 'ab...', got %q", out["KEY"])
	}
	if len(res.Truncated) != 1 || res.Truncated[0] != "KEY" {
		t.Fatalf("expected KEY in truncated, got %v", res.Truncated)
	}
}

func TestTruncate_CustomSuffix(t *testing.T) {
	secrets := map[string]string{"X": "hello world"}
	out, _ := Truncate(secrets, TruncateOptions{MaxLength: 7, Suffix: "~~"})
	if out["X"] != "hello~~" {
		t.Fatalf("expected 'hello~~', got %q", out["X"])
	}
}

func TestTruncate_RestrictedToKeys(t *testing.T) {
	secrets := map[string]string{"A": "longvalue", "B": "longvalue"}
	out, res := Truncate(secrets, TruncateOptions{MaxLength: 4, Keys: []string{"A"}})
	if out["B"] != "longvalue" {
		t.Fatalf("B should be untouched, got %q", out["B"])
	}
	if len(res.Truncated) != 1 || res.Truncated[0] != "A" {
		t.Fatalf("only A should be truncated, got %v", res.Truncated)
	}
}

func TestTruncate_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"K": "verylongvalue"}
	out, res := Truncate(secrets, TruncateOptions{MaxLength: 5, DryRun: true})
	if out["K"] != "verylongvalue" {
		t.Fatalf("dry run should not change value, got %q", out["K"])
	}
	if len(res.Truncated) != 1 {
		t.Fatalf("dry run should still report truncation, got %v", res.Truncated)
	}
}

func TestTruncate_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"K": "longvalue"}
	Truncate(secrets, TruncateOptions{MaxLength: 3})
	if secrets["K"] != "longvalue" {
		t.Fatal("input map was mutated")
	}
}

func TestTruncate_SummaryWithTruncations(t *testing.T) {
	secrets := map[string]string{"A": "toolongvalue"}
	_, res := Truncate(secrets, TruncateOptions{MaxLength: 5})
	if !strings.Contains(res.Summary(), "1 value(s) truncated") {
		t.Fatalf("unexpected summary: %q", res.Summary())
	}
}

func TestTruncate_SummaryNoTruncations(t *testing.T) {
	secrets := map[string]string{"A": "hi"}
	_, res := Truncate(secrets, TruncateOptions{MaxLength: 10})
	if !strings.Contains(res.Summary(), "no values exceeded") {
		t.Fatalf("unexpected summary: %q", res.Summary())
	}
}
