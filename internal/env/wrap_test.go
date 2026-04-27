package env

import (
	"testing"
)

func TestWrap_NoOptions(t *testing.T) {
	secrets := map[string]string{"KEY": "value"}
	res := Wrap(secrets, nil)
	if res.Secrets["KEY"] != "value" {
		t.Fatalf("expected value unchanged, got %q", res.Secrets["KEY"])
	}
	if len(res.Wrapped) != 0 {
		t.Fatalf("expected no wrapped keys, got %d", len(res.Wrapped))
	}
}

func TestWrap_AddsPrefixAndSuffix(t *testing.T) {
	secrets := map[string]string{"A": "hello", "B": "world"}
	res := Wrap(secrets, &WrapOptions{Prefix: "[", Suffix: "]"})
	if res.Secrets["A"] != "[hello]" {
		t.Errorf("expected [hello], got %q", res.Secrets["A"])
	}
	if res.Secrets["B"] != "[world]" {
		t.Errorf("expected [world], got %q", res.Secrets["B"])
	}
	if len(res.Wrapped) != 2 {
		t.Errorf("expected 2 wrapped, got %d", len(res.Wrapped))
	}
}

func TestWrap_RestrictedToKeys(t *testing.T) {
	secrets := map[string]string{"A": "foo", "B": "bar"}
	res := Wrap(secrets, &WrapOptions{Prefix: ">>>", Keys: []string{"A"}})
	if res.Secrets["A"] != ">>>foo" {
		t.Errorf("expected >>>foo, got %q", res.Secrets["A"])
	}
	if res.Secrets["B"] != "bar" {
		t.Errorf("expected bar unchanged, got %q", res.Secrets["B"])
	}
	if len(res.Wrapped) != 1 {
		t.Errorf("expected 1 wrapped, got %d", len(res.Wrapped))
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestWrap_DryRunDoesNotMutate(t *testing.T) {
	secrets := map[string]string{"X": "original"}
	res := Wrap(secrets, &WrapOptions{Prefix: "pre_", Suffix: "_suf", DryRun: true})
	if res.Secrets["X"] != "original" {
		t.Errorf("dry run should not mutate, got %q", res.Secrets["X"])
	}
	if len(res.Wrapped) != 1 {
		t.Errorf("expected wrapped count 1, got %d", len(res.Wrapped))
	}
}

func TestWrap_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"K": "v"}
	original := map[string]string{"K": "v"}
	Wrap(secrets, &WrapOptions{Prefix: "X"})
	if secrets["K"] != original["K"] {
		t.Errorf("input mutated: got %q", secrets["K"])
	}
}

func TestWrap_SummaryNoWrapped(t *testing.T) {
	res := WrapResult{}
	if res.Summary() != "wrap: no keys wrapped" {
		t.Errorf("unexpected summary: %q", res.Summary())
	}
}

func TestWrap_SummaryWithWrapped(t *testing.T) {
	res := WrapResult{Wrapped: []string{"A", "B"}, Skipped: []string{"C"}}
	expected := "wrap: 2 key(s) wrapped, 1 skipped"
	if res.Summary() != expected {
		t.Errorf("expected %q, got %q", expected, res.Summary())
	}
}
