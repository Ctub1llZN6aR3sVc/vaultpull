package env

import (
	"testing"
)

func TestUnion_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{}
	r := Union(a, b, UnionOptions{})
	if r.Secrets["FOO"] != "1" {
		t.Fatalf("expected FOO=1, got %q", r.Secrets["FOO"])
	}
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "FOO" {
		t.Fatalf("unexpected OnlyInA: %v", r.OnlyInA)
	}
}

func TestUnion_OnlyInB(t *testing.T) {
	a := map[string]string{}
	b := map[string]string{"BAR": "2"}
	r := Union(a, b, UnionOptions{})
	if r.Secrets["BAR"] != "2" {
		t.Fatalf("expected BAR=2, got %q", r.Secrets["BAR"])
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BAR" {
		t.Fatalf("unexpected OnlyInB: %v", r.OnlyInB)
	}
}

func TestUnion_PreferA(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	r := Union(a, b, UnionOptions{PreferA: true})
	if r.Secrets["KEY"] != "from-a" {
		t.Fatalf("expected from-a, got %q", r.Secrets["KEY"])
	}
	if len(r.InBoth) != 1 {
		t.Fatalf("expected 1 shared key")
	}
}

func TestUnion_PreferB(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	r := Union(a, b, UnionOptions{PreferA: false})
	if r.Secrets["KEY"] != "from-b" {
		t.Fatalf("expected from-b, got %q", r.Secrets["KEY"])
	}
}

func TestUnion_DoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	Union(a, b, UnionOptions{})
	if len(a) != 1 || len(b) != 1 {
		t.Fatal("inputs were mutated")
	}
}

func TestUnion_Summary(t *testing.T) {
	a := map[string]string{"A": "1", "C": "3"}
	b := map[string]string{"B": "2", "C": "x"}
	r := Union(a, b, UnionOptions{PreferA: true})
	s := r.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestUnion_BothEmpty(t *testing.T) {
	r := Union(map[string]string{}, map[string]string{}, UnionOptions{})
	if len(r.Secrets) != 0 {
		t.Fatal("expected empty secrets")
	}
}
