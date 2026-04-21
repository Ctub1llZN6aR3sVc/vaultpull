package env

import (
	"strings"
	"testing"
)

func TestSquash_NoDuplicates(t *testing.T) {
	a := map[string]string{"FOO": "foo", "BAR": "bar"}
	b := map[string]string{"BAZ": "baz"}
	out, res := Squash([]map[string]string{a, b}, nil)
	if len(res.Squashed) != 0 {
		t.Fatalf("expected no squashed keys, got %v", res.Squashed)
	}
	if out["FOO"] != "foo" || out["BAR"] != "bar" || out["BAZ"] != "baz" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestSquash_JoinsDuplicates(t *testing.T) {
	a := map[string]string{"KEY": "first"}
	b := map[string]string{"KEY": "second"}
	out, res := Squash([]map[string]string{a, b}, nil)
	if len(res.Squashed) != 1 || res.Squashed[0] != "KEY" {
		t.Fatalf("expected KEY in squashed, got %v", res.Squashed)
	}
	if out["KEY"] != "first,second" {
		t.Fatalf("expected 'first,second', got %q", out["KEY"])
	}
}

func TestSquash_CustomSeparator(t *testing.T) {
	a := map[string]string{"K": "a"}
	b := map[string]string{"K": "b"}
	out, _ := Squash([]map[string]string{a, b}, &SquashOptions{Separator: "|"})
	if out["K"] != "a|b" {
		t.Fatalf("expected 'a|b', got %q", out["K"])
	}
}

func TestSquash_KeepFirst(t *testing.T) {
	a := map[string]string{"K": "first"}
	b := map[string]string{"K": "second"}
	out, _ := Squash([]map[string]string{a, b}, &SquashOptions{KeepFirst: true})
	if out["K"] != "first" {
		t.Fatalf("expected 'first', got %q", out["K"])
	}
}

func TestSquash_KeepLast(t *testing.T) {
	a := map[string]string{"K": "first"}
	b := map[string]string{"K": "second"}
	out, _ := Squash([]map[string]string{a, b}, &SquashOptions{KeepLast: true})
	if out["K"] != "second" {
		t.Fatalf("expected 'second', got %q", out["K"])
	}
}

func TestSquash_DoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"K": "a"}
	b := map[string]string{"K": "b"}
	Squash([]map[string]string{a, b}, nil)
	if a["K"] != "a" || b["K"] != "b" {
		t.Fatal("Squash mutated input maps")
	}
}

func TestSquash_Summary(t *testing.T) {
	a := map[string]string{"X": "1"}
	b := map[string]string{"X": "2"}
	_, res := Squash([]map[string]string{a, b}, nil)
	if !strings.Contains(res.Summary(), "X") {
		t.Fatalf("expected summary to mention X, got %q", res.Summary())
	}
}

func TestSquash_SummaryNoSquashed(t *testing.T) {
	a := map[string]string{"A": "1"}
	_, res := Squash([]map[string]string{a}, nil)
	if !strings.Contains(res.Summary(), "no duplicate") {
		t.Fatalf("expected 'no duplicate' in summary, got %q", res.Summary())
	}
}
