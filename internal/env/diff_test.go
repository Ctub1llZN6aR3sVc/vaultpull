package env

import (
	"strings"
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(existing, incoming)
	if !result.IsEmpty() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar", "NEW_KEY": "new_val"}

	result := Diff(existing, incoming)
	if result.IsEmpty() {
		t.Fatal("expected changes")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "NEW_KEY" && c.Type == Added && c.NewValue == "new_val" {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW_KEY to be marked as added")
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "OLD_KEY": "old_val"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(existing, incoming)
	found := false
	for _, c := range result.Changes {
		if c.Key == "OLD_KEY" && c.Type == Removed && c.OldValue == "old_val" {
			found = true
		}
	}
	if !found {
		t.Error("expected OLD_KEY to be marked as removed")
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := Diff(existing, incoming)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != Changed || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_IsEmpty_False(t *testing.T) {
	result := Diff(map[string]string{}, map[string]string{"X": "1"})
	if result.IsEmpty() {
		t.Error("expected non-empty diff")
	}
}

func TestDiff_Summary_ContainsSymbols(t *testing.T) {
	existing := map[string]string{"A": "old", "B": "keep", "C": "gone"}
	incoming := map[string]string{"A": "new", "B": "keep", "D": "added"}

	result := Diff(existing, incoming)
	summary := result.Summary()

	if !strings.Contains(summary, "~ A") {
		t.Error("expected changed marker for A")
	}
	if !strings.Contains(summary, "+ D") {
		t.Error("expected added marker for D")
	}
	if !strings.Contains(summary, "- C") {
		t.Error("expected removed marker for C")
	}
}
