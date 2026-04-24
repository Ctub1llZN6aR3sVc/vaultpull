package env

import (
	"testing"
)

func TestReorder_NoOptions(t *testing.T) {
	secrets := map[string]string{"B": "2", "A": "1", "C": "3"}
	out, res := Reorder(secrets, ReorderOptions{})
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
	// With no Keys and no Alphabetical, all keys are unlisted (appended sorted)
	if len(res.Ordered) != 3 {
		t.Fatalf("expected ordered length 3, got %d", len(res.Ordered))
	}
}

func TestReorder_Alphabetical(t *testing.T) {
	secrets := map[string]string{"C": "3", "A": "1", "B": "2"}
	_, res := Reorder(secrets, ReorderOptions{Alphabetical: true})
	if res.Ordered[0] != "A" || res.Ordered[1] != "B" || res.Ordered[2] != "C" {
		t.Fatalf("expected A,B,C got %v", res.Ordered)
	}
	if len(res.Unlisted) != 0 {
		t.Fatalf("expected no unlisted keys, got %v", res.Unlisted)
	}
}

func TestReorder_ExplicitKeys(t *testing.T) {
	secrets := map[string]string{"C": "3", "A": "1", "B": "2"}
	_, res := Reorder(secrets, ReorderOptions{Keys: []string{"C", "A"}})
	if res.Ordered[0] != "C" || res.Ordered[1] != "A" {
		t.Fatalf("expected C,A first, got %v", res.Ordered)
	}
	if res.Ordered[2] != "B" {
		t.Fatalf("expected B last, got %v", res.Ordered)
	}
	if len(res.Unlisted) != 1 || res.Unlisted[0] != "B" {
		t.Fatalf("expected B in unlisted, got %v", res.Unlisted)
	}
}

func TestReorder_MissingExplicitKeyIgnored(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	_, res := Reorder(secrets, ReorderOptions{Keys: []string{"A", "MISSING", "B"}})
	if len(res.Ordered) != 2 {
		t.Fatalf("expected 2 ordered keys, got %d", len(res.Ordered))
	}
}

func TestReorder_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	orig := map[string]string{"A": "1", "B": "2"}
	Reorder(secrets, ReorderOptions{Alphabetical: true})
	for k, v := range orig {
		if secrets[k] != v {
			t.Fatalf("input mutated at key %s", k)
		}
	}
}

func TestReorder_SummaryNoUnlisted(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	_, res := Reorder(secrets, ReorderOptions{Keys: []string{"A", "B"}})
	if res.Summary() != "all keys reordered as specified" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}

func TestReorder_SummaryWithUnlisted(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	_, res := Reorder(secrets, ReorderOptions{Keys: []string{"A"}})
	s := res.Summary()
	if s == "all keys reordered as specified" {
		t.Fatalf("expected unlisted summary, got: %s", s)
	}
}
