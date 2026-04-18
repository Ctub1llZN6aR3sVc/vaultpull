package env

import (
	"testing"
)

func TestPromote_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}
	out, res := Promote(dst, src, PromoteOptions{})
	if out["B"] != "2" {
		t.Fatalf("expected B=2, got %q", out["B"])
	}
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Fatalf("unexpected Added: %v", res.Added)
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new"}
	out, res := Promote(dst, src, PromoteOptions{Overwrite: false})
	if out["A"] != "original" {
		t.Fatalf("expected original value, got %q", out["A"])
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %v", res.Skipped)
	}
}

func TestPromote_OverwritesWhenEnabled(t *testing.T) {
	dst := map[string]string{"A": "old"}
	src := map[string]string{"A": "new"}
	out, res := Promote(dst, src, PromoteOptions{Overwrite: true})
	if out["A"] != "new" {
		t.Fatalf("expected new value, got %q", out["A"])
	}
	if len(res.Overwrote) != 1 {
		t.Fatalf("expected 1 overwrote, got %v", res.Overwrote)
	}
}

func TestPromote_DryRunDoesNotMutate(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}
	out, res := Promote(dst, src, PromoteOptions{DryRun: true})
	if _, ok := dst["B"]; ok {
		t.Fatal("dst should not be mutated in dry-run")
	}
	if _, ok := out["B"]; ok {
		t.Fatal("out should not contain B in dry-run")
	}
	if len(res.Added) != 1 {
		t.Fatalf("expected 1 added in result, got %v", res.Added)
	}
}

func TestPromote_DoesNotMutateInputDst(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"B": "2"}
	Promote(dst, src, PromoteOptions{})
	if _, ok := dst["B"]; ok {
		t.Fatal("original dst map should not be mutated")
	}
}

func TestPromote_EmptySrcReturnsUnchanged(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{}
	out, res := Promote(dst, src, PromoteOptions{})
	if len(out) != 1 || out["A"] != "1" {
		t.Fatalf("expected unchanged output, got %v", out)
	}
	if len(res.Added) != 0 || len(res.Skipped) != 0 || len(res.Overwrote) != 0 {
		t.Fatalf("expected empty result, got %v", res)
	}
}

func TestPromoteResult_String(t *testing.T) {
	r := PromoteResult{Added: []string{"X"}, Skipped: []string{"Y", "Z"}, Overwrote: nil}
	s := r.String()
	if s != "added=1 skipped=2 overwrote=0" {
		t.Fatalf("unexpected string: %q", s)
	}
}
