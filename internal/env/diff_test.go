package env

import (
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(old, incoming)
	if !result.IsEmpty() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	old := map[string]string{}
	incoming := map[string]string{"NEW_KEY": "value"}

	result := Diff(old, incoming)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != ChangeAdded {
		t.Errorf("expected ChangeAdded, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].Key != "NEW_KEY" {
		t.Errorf("expected key NEW_KEY, got %s", result.Changes[0].Key)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	old := map[string]string{"GONE": "bye"}
	incoming := map[string]string{}

	result := Diff(old, incoming)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != ChangeRemoved {
		t.Errorf("expected ChangeRemoved, got %s", result.Changes[0].Type)
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	old := map[string]string{"KEY": "old"}
	incoming := map[string]string{"KEY": "new"}

	result := Diff(old, incoming)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != ChangeUpdated {
		t.Errorf("expected ChangeUpdated, got %s", c.Type)
	}
	if c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
	}
}

func TestDiff_IsEmpty_False(t *testing.T) {
	old := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "2"}

	result := Diff(old, incoming)
	if result.IsEmpty() {
		t.Error("expected IsEmpty to be false")
	}
}

func TestChange_String(t *testing.T) {
	cases := []struct {
		change Change
		want   string
	}{
		{Change{Key: "FOO", Type: ChangeAdded}, "+ FOO"},
		{Change{Key: "BAR", Type: ChangeRemoved}, "- BAR"},
		{Change{Key: "BAZ", Type: ChangeUpdated}, "~ BAZ"},
	}
	for _, tc := range cases {
		got := tc.change.String()
		if got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}
