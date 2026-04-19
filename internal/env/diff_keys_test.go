package env

import (
	"testing"
)

func TestDiffKeys_IdenticalMaps(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "x", "BAR": "y"}
	d := DiffKeys(a, b)
	if len(d.OnlyInA) != 0 || len(d.OnlyInB) != 0 {
		t.Errorf("expected no exclusive keys, got onlyInA=%v onlyInB=%v", d.OnlyInA, d.OnlyInB)
	}
	if len(d.InBoth) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(d.InBoth))
	}
}

func TestDiffKeys_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1", "EXTRA": "x"}
	b := map[string]string{"FOO": "1"}
	d := DiffKeys(a, b)
	if len(d.OnlyInA) != 1 || d.OnlyInA[0] != "EXTRA" {
		t.Errorf("expected EXTRA only in A, got %v", d.OnlyInA)
	}
	if len(d.OnlyInB) != 0 {
		t.Errorf("expected nothing only in B, got %v", d.OnlyInB)
	}
}

func TestDiffKeys_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"FOO": "1", "NEW": "2"}
	d := DiffKeys(a, b)
	if len(d.OnlyInB) != 1 || d.OnlyInB[0] != "NEW" {
		t.Errorf("expected NEW only in B, got %v", d.OnlyInB)
	}
}

func TestDiffKeys_EmptyMaps(t *testing.T) {
	d := DiffKeys(map[string]string{}, map[string]string{})
	if len(d.OnlyInA) != 0 || len(d.OnlyInB) != 0 || len(d.InBoth) != 0 {
		t.Error("expected all empty slices for empty maps")
	}
}

func TestDiffKeys_SummaryIdentical(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"A": "2"}
	d := DiffKeys(a, b)
	got := KeyDiffSummary(d)
	if got != "key sets are identical" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestDiffKeys_SummaryWithDifferences(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "2", "C": "3"}
	d := DiffKeys(a, b)
	got := KeyDiffSummary(d)
	if got == "" || got == "key sets are identical" {
		t.Errorf("expected non-identical summary, got: %s", got)
	}
}
