package env

import (
	"testing"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"C": "3"}
	out, res := Dedupe(base, incoming, true)
	if len(res.Removed) != 0 {
		t.Fatalf("expected no dupes, got %v", res.Removed)
	}
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	incoming := map[string]string{"KEY": "override"}
	out, res := Dedupe(base, incoming, true)
	if len(res.Removed) != 1 || res.Removed[0] != "KEY" {
		t.Fatalf("expected KEY in dupes, got %v", res.Removed)
	}
	if out["KEY"] != "original" {
		t.Fatalf("expected original value, got %s", out["KEY"])
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	base := map[string]string{"KEY": "original"}
	incoming := map[string]string{"KEY": "override"}
	out, res := Dedupe(base, incoming, false)
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 dupe, got %v", res.Removed)
	}
	if out["KEY"] != "override" {
		t.Fatalf("expected override value, got %s", out["KEY"])
	}
}

func TestDedupe_MultipleDuplicates(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	incoming := map[string]string{"A": "x", "C": "y", "D": "4"}
	_, res := Dedupe(base, incoming, true)
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 dupes, got %v", res.Removed)
	}
	if res.Removed[0] != "A" || res.Removed[1] != "C" {
		t.Fatalf("unexpected dupe keys: %v", res.Removed)
	}
}

func TestDedupe_DoesNotMutateInput(t *testing.T) {
	base := map[string]string{"K": "v1"}
	incoming := map[string]string{"K": "v2"}
	Dedupe(base, incoming, false)
	if base["K"] != "v1" {
		t.Fatal("base was mutated")
	}
}

func TestDedupeResult_Summary(t *testing.T) {
	res := DedupeResult{}
	if res.Summary() != "no duplicate keys found" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
	res.Removed = []string{"A", "B"}
	s := res.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
