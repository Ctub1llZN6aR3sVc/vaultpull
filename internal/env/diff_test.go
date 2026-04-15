package env

import (
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	current := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(current, incoming)

	if !result.IsEmpty() {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v",
			result.Added, result.Removed, result.Changed)
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	current := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar", "NEW_KEY": "newval"}

	result := Diff(current, incoming)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(result.Added))
	}
	if result.Added["NEW_KEY"] != "newval" {
		t.Errorf("expected NEW_KEY=newval, got %q", result.Added["NEW_KEY"])
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	current := map[string]string{"FOO": "bar", "OLD_KEY": "oldval"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(current, incoming)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(result.Removed))
	}
	if result.Removed["OLD_KEY"] != "oldval" {
		t.Errorf("expected OLD_KEY=oldval in removed, got %q", result.Removed["OLD_KEY"])
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	current := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "baz"}

	result := Diff(current, incoming)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(result.Changed))
	}
	on, ok := result.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in changed")
	}
	if on.Old != "bar" || on.New != "baz" {
		t.Errorf("expected FOO old=bar new=baz, got old=%q new=%q", on.Old, on.New)
	}
}

func TestDiff_IsEmpty_False(t *testing.T) {
	current := map[string]string{}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(current, incoming)

	if result.IsEmpty() {
		t.Error("expected IsEmpty to return false")
	}
}
