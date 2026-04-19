package env

import (
	"testing"
)

func TestPick_ReturnsOnlyRequestedKeys(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	res, err := Pick(secrets, PickOptions{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Picked) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(res.Picked))
	}
	if res.Picked["A"] != "1" || res.Picked["C"] != "3" {
		t.Errorf("unexpected values: %v", res.Picked)
	}
}

func TestPick_MissingKeyNotInResult(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	res, err := Pick(secrets, PickOptions{Keys: []string{"A", "Z"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Picked["Z"]; ok {
		t.Error("Z should not be in picked result")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "Z" {
		t.Errorf("expected missing=[Z], got %v", res.Missing)
	}
}

func TestPick_FailOnMissingReturnsError(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	_, err := Pick(secrets, PickOptions{Keys: []string{"A", "MISSING"}, FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPick_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	_, _ = Pick(secrets, PickOptions{Keys: []string{"A"}})
	if len(secrets) != 2 {
		t.Error("input map was mutated")
	}
}

func TestPick_EmptyKeysReturnsEmptyMap(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	res, err := Pick(secrets, PickOptions{Keys: []string{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Picked) != 0 {
		t.Errorf("expected empty map, got %v", res.Picked)
	}
}

func TestPick_SummaryNoMissing(t *testing.T) {
	res := PickResult{Picked: map[string]string{"A": "1", "B": "2"}}
	s := res.Summary()
	if s != "picked 2 key(s)" {
		t.Errorf("unexpected summary: %s", s)
	}
}
