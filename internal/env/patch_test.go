package env

import (
	"testing"
)

func TestPatch_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	patch := map[string]string{"B": "2"}
	out, res := Patch(dst, patch, PatchOptions{})
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %q", out["B"])
	}
	if len(res.Patched) != 1 || res.Patched[0] != "B" {
		t.Errorf("unexpected patched: %v", res.Patched)
	}
}

func TestPatch_UpdatesExistingKeys(t *testing.T) {
	dst := map[string]string{"A": "old"}
	out, res := Patch(dst, map[string]string{"A": "new"}, PatchOptions{})
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %q", out["A"])
	}
	if len(res.Patched) != 1 {
		t.Errorf("expected 1 patched key")
	}
}

func TestPatch_ExistingOnlySkipsNewKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	out, res := Patch(dst, map[string]string{"B": "2"}, PatchOptions{ExistingOnly: true})
	if _, ok := out["B"]; ok {
		t.Error("B should not be added with ExistingOnly")
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Errorf("expected B in skipped, got %v", res.Skipped)
	}
}

func TestPatch_DryRunDoesNotMutate(t *testing.T) {
	dst := map[string]string{"A": "1"}
	out, res := Patch(dst, map[string]string{"A": "99", "B": "2"}, PatchOptions{DryRun: true})
	if out["A"] != "1" {
		t.Errorf("dry run should not change A")
	}
	if _, ok := out["B"]; ok {
		t.Error("dry run should not add B")
	}
	if len(res.Patched) != 2 {
		t.Errorf("dry run should still report patched keys, got %v", res.Patched)
	}
}

func TestPatch_DoesNotMutateInput(t *testing.T) {
	dst := map[string]string{"A": "1"}
	Patch(dst, map[string]string{"A": "2"}, PatchOptions{})
	if dst["A"] != "1" {
		t.Error("original dst was mutated")
	}
}
