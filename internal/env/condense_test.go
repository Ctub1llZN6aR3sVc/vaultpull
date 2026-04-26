package env

import (
	"testing"
)

func TestCondense_NoOptions(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	out, res, err := Condense(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if res.TargetKey != "" {
		t.Errorf("expected empty TargetKey, got %q", res.TargetKey)
	}
}

func TestCondense_JoinsAllKeys(t *testing.T) {
	secrets := map[string]string{"B": "beta", "A": "alpha"}
	out, res, err := Condense(secrets, &CondenseOptions{TargetKey: "COMBINED", Separator: "|"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["COMBINED"] != "alpha|beta" && out["COMBINED"] != "beta|alpha" {
		t.Errorf("unexpected combined value: %q", out["COMBINED"])
	}
	if _, ok := out["A"]; ok {
		t.Error("expected key A to be removed")
	}
	_ = res
}

func TestCondense_SpecificKeys(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432", "NAME": "mydb"}
	out, res, err := Condense(secrets, &CondenseOptions{
		Keys:      []string{"HOST", "PORT"},
		TargetKey: "ADDR",
		Separator: ":",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ADDR"] != "HOST:PORT" && out["ADDR"] != "localhost:5432" {
		// sorted keys: HOST then PORT
		if out["ADDR"] != "localhost:5432" {
			t.Errorf("unexpected ADDR value: %q", out["ADDR"])
		}
	}
	if _, ok := out["HOST"]; ok {
		t.Error("expected HOST to be removed")
	}
	if _, ok := out["PORT"]; ok {
		t.Error("expected PORT to be removed")
	}
	if out["NAME"] != "mydb" {
		t.Error("expected NAME to be preserved")
	}
	if len(res.SourceKeys) != 2 {
		t.Errorf("expected 2 source keys, got %d", len(res.SourceKeys))
	}
}

func TestCondense_DefaultSeparator(t *testing.T) {
	secrets := map[string]string{"X": "foo", "Y": "bar"}
	out, _, err := Condense(secrets, &CondenseOptions{TargetKey: "Z"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v := out["Z"]
	if v != "foo,bar" && v != "bar,foo" {
		t.Errorf("expected comma-separated value, got %q", v)
	}
}

func TestCondense_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	out, res, err := Condense(secrets, &CondenseOptions{
		Keys:      []string{"A", "B"},
		TargetKey: "C",
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["C"]; ok {
		t.Error("dry run should not write target key")
	}
	if _, ok := out["A"]; !ok {
		t.Error("dry run should preserve source keys")
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
}

func TestCondense_Summary(t *testing.T) {
	res := CondenseResult{TargetKey: "OUT", SourceKeys: []string{"A", "B"}, Value: "1,2"}
	s := res.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}

	empty := CondenseResult{}
	if empty.Summary() == "" {
		t.Error("expected non-empty summary for empty result")
	}
}
