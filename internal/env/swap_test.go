package env

import (
	"testing"
)

func TestSwap_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, res, err := Swap(secrets, SwapOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if len(res.Swapped) != 0 {
		t.Fatalf("expected no swaps, got %v", res.Swapped)
	}
}

func TestSwap_RenamesKey(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "secret_value", "KEEP": "yes"}
	out, res, err := Swap(secrets, SwapOptions{Pairs: map[string]string{"OLD_KEY": "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if out["NEW_KEY"] != "secret_value" {
		t.Errorf("expected NEW_KEY=secret_value, got %q", out["NEW_KEY"])
	}
	if out["KEEP"] != "yes" {
		t.Error("expected KEEP to be preserved")
	}
	if len(res.Swapped) != 1 || res.Swapped[0] != "OLD_KEY" {
		t.Errorf("unexpected swapped list: %v", res.Swapped)
	}
}

func TestSwap_MissingKeyNoFail(t *testing.T) {
	secrets := map[string]string{"KEEP": "yes"}
	out, res, err := Swap(secrets, SwapOptions{Pairs: map[string]string{"MISSING": "NEW"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["NEW"]; ok {
		t.Error("NEW should not be present when source is missing")
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 missing, got %v", res.Missing)
	}
}

func TestSwap_MissingKeyFailOnMissing(t *testing.T) {
	secrets := map[string]string{"KEEP": "yes"}
	_, _, err := Swap(secrets, SwapOptions{
		Pairs:         map[string]string{"GHOST": "SPIRIT"},
		FailOnMissing: true,
	})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSwap_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	out, res, err := Swap(secrets, SwapOptions{
		Pairs:  map[string]string{"A": "Z"},
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["A"]; !ok {
		t.Error("dry run should not remove original key")
	}
	if _, ok := out["Z"]; ok {
		t.Error("dry run should not create new key")
	}
	if len(res.Swapped) != 1 {
		t.Errorf("expected 1 in swapped list even on dry run, got %v", res.Swapped)
	}
}

func TestSwap_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"X": "val"}
	Swap(secrets, SwapOptions{Pairs: map[string]string{"X": "Y"}})
	if _, ok := secrets["X"]; !ok {
		t.Error("Swap must not mutate the input map")
	}
}

func TestSwapResult_Summary(t *testing.T) {
	r := SwapResult{Swapped: []string{"A", "B"}, Missing: []string{"C"}}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
