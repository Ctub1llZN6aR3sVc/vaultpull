package env

import (
	"testing"
)

func TestMerge_OverwriteReplacesExisting(t *testing.T) {
	base := map[string]string{"FOO": "old", "BAR": "keep"}
	incoming := map[string]string{"FOO": "new", "BAZ": "added"}

	result := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOverwrite})

	if result["FOO"] != "new" {
		t.Errorf("expected FOO=new, got %s", result["FOO"])
	}
	if result["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %s", result["BAR"])
	}
	if result["BAZ"] != "added" {
		t.Errorf("expected BAZ=added, got %s", result["BAZ"])
	}
}

func TestMerge_KeepExistingPreservesValues(t *testing.T) {
	base := map[string]string{"FOO": "old", "BAR": "keep"}
	incoming := map[string]string{"FOO": "new", "BAZ": "added"}

	result := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyKeepExisting})

	if result["FOO"] != "old" {
		t.Errorf("expected FOO=old (preserved), got %s", result["FOO"])
	}
	if result["BAZ"] != "added" {
		t.Errorf("expected BAZ=added (new key), got %s", result["BAZ"])
	}
}

func TestMerge_DoesNotMutateInputs(t *testing.T) {
	base := map[string]string{"FOO": "original"}
	incoming := map[string]string{"FOO": "changed"}

	Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOverwrite})

	if base["FOO"] != "original" {
		t.Error("base map was mutated")
	}
}

func TestMerge_EmptyBase(t *testing.T) {
	base := map[string]string{}
	incoming := map[string]string{"KEY": "val"}

	result := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOverwrite})

	if result["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %s", result["KEY"])
	}
}

func TestMerge_EmptyIncoming(t *testing.T) {
	base := map[string]string{"KEY": "val"}
	incoming := map[string]string{}

	result := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOverwrite})

	if result["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %s", result["KEY"])
	}
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
}
