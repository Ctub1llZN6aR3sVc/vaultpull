package env

import (
	"testing"
)

func TestOmit_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	out, res := Omit(secrets, OmitOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected no removals")
	}
}

func TestOmit_ByKey(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	out, res := Omit(secrets, OmitOptions{Keys: []string{"B"}})
	if _, ok := out["B"]; ok {
		t.Fatal("B should have been removed")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if len(res.Removed) != 1 || res.Removed[0] != "B" {
		t.Fatalf("unexpected removed list: %v", res.Removed)
	}
}

func TestOmit_ByPrefix(t *testing.T) {
	secrets := map[string]string{"DEV_HOST": "h", "DEV_PORT": "p", "PROD_HOST": "ph"}
	out, res := Omit(secrets, OmitOptions{Prefixes: []string{"DEV_"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["PROD_HOST"]; !ok {
		t.Fatal("PROD_HOST should remain")
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestOmit_EmptyValues(t *testing.T) {
	secrets := map[string]string{"A": "val", "B": "", "C": ""}
	out, res := Omit(secrets, OmitOptions{Empty: true})
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestOmit_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	Omit(secrets, OmitOptions{Keys: []string{"A"}})
	if len(secrets) != 2 {
		t.Fatal("input map was mutated")
	}
}

func TestOmit_Summary(t *testing.T) {
	res := OmitResult{Removed: []string{"X", "Y"}}
	s := res.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
