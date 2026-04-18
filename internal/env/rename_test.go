package env

import (
	"testing"
)

func TestRename_BasicRename(t *testing.T) {
	secrets := map[string]string{"OLD_KEY": "value1", "KEEP": "value2"}
	out, res, err := Rename(secrets, RenameOptions{Map: map[string]string{"OLD_KEY": "NEW_KEY"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", out["NEW_KEY"])
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if out["KEEP"] != "value2" {
		t.Error("KEEP should be preserved")
	}
	if len(res.Renamed) != 1 || res.Renamed[0] != "OLD_KEY" {
		t.Errorf("unexpected renamed list: %v", res.Renamed)
	}
}

func TestRename_MissingKeyNoFail(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	_, res, err := Rename(secrets, RenameOptions{Map: map[string]string{"MISSING": "NEW"}, FailOnMissing: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missed) != 1 || res.Missed[0] != "MISSING" {
		t.Errorf("expected MISSING in missed list, got %v", res.Missed)
	}
}

func TestRename_MissingKeyFailOnMissing(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	_, _, err := Rename(secrets, RenameOptions{Map: map[string]string{"MISSING": "NEW"}, FailOnMissing: true})
	if err == nil {
		t.Error("expected error for missing key with FailOnMissing")
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"OLD": "val"}
	Rename(secrets, RenameOptions{Map: map[string]string{"OLD": "NEW"}})
	if _, ok := secrets["OLD"]; !ok {
		t.Error("input map should not be mutated")
	}
}

func TestRename_SummaryNoRenames(t *testing.T) {
	res := RenameResult{}
	if res.Summary() != "no renames applied" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestRename_SummaryWithCounts(t *testing.T) {
	res := RenameResult{Renamed: []string{"A", "B"}, Missed: []string{"C"}}
	if res.Summary() != "2 renamed, 1 not found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
