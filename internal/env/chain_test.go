package env

import (
	"testing"
)

func TestChain_FirstSourceWins(t *testing.T) {
	a := map[string]string{"FOO": "from_a", "BAR": "bar_a"}
	b := map[string]string{"FOO": "from_b", "BAZ": "baz_b"}
	out := Chain([]string{"FOO", "BAR", "BAZ"}, a, b)
	if out["FOO"] != "from_a" {
		t.Errorf("expected from_a, got %s", out["FOO"])
	}
	if out["BAR"] != "bar_a" {
		t.Errorf("expected bar_a, got %s", out["BAR"])
	}
	if out["BAZ"] != "baz_b" {
		t.Errorf("expected baz_b, got %s", out["BAZ"])
	}
}

func TestChain_FallsBackWhenEmpty(t *testing.T) {
	a := map[string]string{"FOO": ""}
	b := map[string]string{"FOO": "fallback"}
	out := Chain([]string{"FOO"}, a, b)
	if out["FOO"] != "fallback" {
		t.Errorf("expected fallback, got %s", out["FOO"])
	}
}

func TestChain_MissingKeyNotInResult(t *testing.T) {
	a := map[string]string{"FOO": "val"}
	out := Chain([]string{"FOO", "MISSING"}, a)
	if _, ok := out["MISSING"]; ok {
		t.Error("expected MISSING to be absent from result")
	}
}

func TestChainAll_MergesAllSources(t *testing.T) {
	a := map[string]string{"A": "1", "SHARED": "from_a"}
	b := map[string]string{"B": "2", "SHARED": "from_b"}
	out := ChainAll(a, b)
	if out["A"] != "1" {
		t.Errorf("expected 1, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected 2, got %s", out["B"])
	}
	if out["SHARED"] != "from_a" {
		t.Errorf("expected from_a, got %s", out["SHARED"])
	}
}

func TestChainAll_EmptySources(t *testing.T) {
	out := ChainAll()
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
