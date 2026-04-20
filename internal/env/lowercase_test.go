package env

import (
	"testing"
)

func TestLowercase_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "BAR", "BAZ": "QUX"}
	res := Lowercase(secrets, LowercaseOptions{})
	if res.Output["FOO"] != "BAR" {
		t.Errorf("expected FOO=BAR, got %s", res.Output["FOO"])
	}
	if len(res.Changed) != 0 {
		t.Errorf("expected no changes, got %v", res.Changed)
	}
}

func TestLowercase_LowercasesKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := Lowercase(secrets, LowercaseOptions{LowercaseKeys: true})
	if _, ok := res.Output["foo"]; !ok {
		t.Error("expected key 'foo' in output")
	}
	if _, ok := res.Output["baz"]; !ok {
		t.Error("expected key 'baz' in output")
	}
	if len(res.Changed) != 2 {
		t.Errorf("expected 2 changes, got %d", len(res.Changed))
	}
}

func TestLowercase_LowercasesValues(t *testing.T) {
	secrets := map[string]string{"FOO": "HELLO", "BAR": "WORLD"}
	res := Lowercase(secrets, LowercaseOptions{LowercaseValues: true})
	if res.Output["FOO"] != "hello" {
		t.Errorf("expected hello, got %s", res.Output["FOO"])
	}
	if res.Output["BAR"] != "world" {
		t.Errorf("expected world, got %s", res.Output["BAR"])
	}
}

func TestLowercase_OnlyKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "HELLO", "BAR": "WORLD"}
	res := Lowercase(secrets, LowercaseOptions{
		LowercaseValues: true,
		OnlyKeys:        []string{"FOO"},
	})
	if res.Output["FOO"] != "hello" {
		t.Errorf("expected FOO=hello, got %s", res.Output["FOO"])
	}
	if res.Output["BAR"] != "WORLD" {
		t.Errorf("expected BAR=WORLD unchanged, got %s", res.Output["BAR"])
	}
}

func TestLowercase_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"FOO": "BAR"}
	Lowercase(secrets, LowercaseOptions{LowercaseKeys: true, LowercaseValues: true})
	if _, ok := secrets["FOO"]; !ok {
		t.Error("input map was mutated")
	}
}

func TestLowercase_SummaryNoChanges(t *testing.T) {
	res := LowercaseResult{Output: map[string]string{}, Changed: nil}
	if res.Summary() != "lowercase: no changes" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestLowercase_SummaryWithChanges(t *testing.T) {
	res := LowercaseResult{
		Output:  map[string]string{},
		Changed: []string{"FOO", "BAR"},
	}
	summary := res.Summary()
	if summary == "lowercase: no changes" {
		t.Error("expected non-empty summary")
	}
}
