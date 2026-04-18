package env

import (
	"testing"
)

func TestApplyDefaults_AddsNewKeys(t *testing.T) {
	dst := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2", "C": "3"}
	out, res := ApplyDefaults(dst, defs, DefaultsOptions{})
	if out["B"] != "2" || out["C"] != "3" {
		t.Fatalf("expected B=2, C=3, got %v", out)
	}
	if len(res.Applied) != 2 {
		t.Fatalf("expected 2 applied, got %d", len(res.Applied))
	}
}

func TestApplyDefaults_SkipsExistingWithoutOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	defs := map[string]string{"A": "default"}
	out, res := ApplyDefaults(dst, defs, DefaultsOptions{})
	if out["A"] != "original" {
		t.Fatalf("expected original, got %s", out["A"])
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Fatalf("expected A skipped, got %v", res.Skipped)
	}
}

func TestApplyDefaults_OverwriteReplacesExisting(t *testing.T) {
	dst := map[string]string{"A": "original"}
	defs := map[string]string{"A": "default"}
	out, res := ApplyDefaults(dst, defs, DefaultsOptions{Overwrite: true})
	if out["A"] != "default" {
		t.Fatalf("expected default, got %s", out["A"])
	}
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied")
	}
}

func TestApplyDefaults_DryRunDoesNotMutate(t *testing.T) {
	dst := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2"}
	out, res := ApplyDefaults(dst, defs, DefaultsOptions{DryRun: true})
	if _, ok := out["B"]; ok {
		t.Fatal("dry run should not write B")
	}
	if len(res.Applied) != 1 || res.Applied[0] != "B" {
		t.Fatalf("expected B in applied, got %v", res.Applied)
	}
}

func TestApplyDefaults_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2"}
	ApplyDefaults(dst, defs, DefaultsOptions{})
	if _, ok := dst["B"]; ok {
		t.Fatal("original dst should not be mutated")
	}
}

func TestApplyDefaults_SummaryNoApplied(t *testing.T) {
	res := DefaultsResult{}
	if res.Summary() != "defaults: no keys applied" {
		t.Fatalf("unexpected summary: %s", res.Summary())
	}
}
