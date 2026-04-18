package env

import (
	"testing"
)

func TestClone_CopiesAllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	out, res := Clone(src, dst, CloneOptions{})
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("expected all keys copied, got %v", out)
	}
	if len(res.Cloned) != 2 {
		t.Fatalf("expected 2 cloned, got %d", len(res.Cloned))
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res := Clone(src, dst, CloneOptions{Overwrite: false})
	if out["A"] != "old" {
		t.Fatalf("expected old value preserved, got %s", out["A"])
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped")
	}
}

func TestClone_OverwriteReplacesExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res := Clone(src, dst, CloneOptions{Overwrite: true})
	if out["A"] != "new" {
		t.Fatalf("expected new value, got %s", out["A"])
	}
	if len(res.Cloned) != 1 {
		t.Fatalf("expected 1 cloned")
	}
}

func TestClone_SelectiveKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, res := Clone(src, dst, CloneOptions{Keys: []string{"A", "C"}})
	if _, ok := out["B"]; ok {
		t.Fatal("B should not be cloned")
	}
	if len(res.Cloned) != 2 {
		t.Fatalf("expected 2 cloned, got %d", len(res.Cloned))
	}
}

func TestClone_DryRunDoesNotMutate(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	out, res := Clone(src, dst, CloneOptions{DryRun: true})
	if _, ok := out["A"]; ok {
		t.Fatal("dry run should not write key")
	}
	if len(res.Cloned) != 1 {
		t.Fatalf("expected 1 reported as cloned even in dry run")
	}
}

func TestClone_MissingSourceKeySkipped(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	_, res := Clone(src, dst, CloneOptions{Keys: []string{"A", "MISSING"}})
	if len(res.Skipped) != 1 || res.Skipped[0] != "MISSING" {
		t.Fatalf("expected MISSING in skipped, got %v", res.Skipped)
	}
}

func TestClone_SummaryMessage(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	_, res := Clone(src, dst, CloneOptions{})
	s := res.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
