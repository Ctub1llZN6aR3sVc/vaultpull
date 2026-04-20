package env

import (
	"testing"
)

func TestPrefix_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "1", "BAR": "2"}
	res, err := Prefix(secrets, PrefixOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Out))
	}
	if len(res.Renamed) != 0 {
		t.Errorf("expected no renames, got %v", res.Renamed)
	}
}

func TestPrefix_AddPrefix(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	res, err := Prefix(secrets, PrefixOptions{Add: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Out["APP_FOO"]; !ok {
		t.Errorf("expected key APP_FOO in output, got %v", res.Out)
	}
	if len(res.Renamed) != 1 {
		t.Errorf("expected 1 rename, got %d", len(res.Renamed))
	}
}

func TestPrefix_StripPrefix(t *testing.T) {
	secrets := map[string]string{"APP_FOO": "bar", "OTHER": "val"}
	res, err := Prefix(secrets, PrefixOptions{Strip: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Out["FOO"]; !ok {
		t.Errorf("expected key FOO in output, got %v", res.Out)
	}
	if _, ok := res.Out["OTHER"]; !ok {
		t.Errorf("expected key OTHER preserved in output")
	}
}

func TestPrefix_StripThenAdd(t *testing.T) {
	secrets := map[string]string{"OLD_FOO": "1"}
	res, err := Prefix(secrets, PrefixOptions{Strip: "OLD_", Add: "NEW_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Out["NEW_FOO"]; !ok {
		t.Errorf("expected key NEW_FOO, got %v", res.Out)
	}
}

func TestPrefix_FailOnConflict(t *testing.T) {
	secrets := map[string]string{"FOO": "1", "APP_FOO": "2"}
	_, err := Prefix(secrets, PrefixOptions{Add: "APP_", FailOnConflict: true})
	if err == nil {
		t.Error("expected conflict error, got nil")
	}
}

func TestPrefix_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"FOO": "1"}
	res, err := Prefix(secrets, PrefixOptions{Add: "DRY_", DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Out["DRY_FOO"]; !ok {
		t.Errorf("expected DRY_FOO in dry-run result")
	}
	if _, ok := secrets["FOO"]; !ok {
		t.Errorf("original secrets map should be unchanged")
	}
}

func TestPrefix_SummaryNoChanges(t *testing.T) {
	res := PrefixResult{}
	if res.Summary() != "no changes" {
		t.Errorf("expected 'no changes', got %q", res.Summary())
	}
}

func TestPrefix_SummaryWithRenames(t *testing.T) {
	res := PrefixResult{
		Renamed: []string{"A -> B"},
		Skipped: []string{"C"},
	}
	s := res.Summary()
	if s == "no changes" {
		t.Error("expected non-empty summary")
	}
}
