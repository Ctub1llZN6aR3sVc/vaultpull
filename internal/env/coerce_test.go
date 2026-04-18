package env

import (
	"testing"
)

func TestCoerce_NoOptions(t *testing.T) {
	secrets := map[string]string{"KEY": "  value  ", "FLAG": "yes"}
	res := Coerce(secrets, CoerceOptions{})
	if res.Coerced["KEY"] != "  value  " {
		t.Errorf("expected untouched value, got %q", res.Coerced["KEY"])
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %v", res.Changes)
	}
}

func TestCoerce_TrimWhitespace(t *testing.T) {
	secrets := map[string]string{"KEY": "  hello  ", "CLEAN": "ok"}
	res := Coerce(secrets, CoerceOptions{TrimWhitespace: true})
	if res.Coerced["KEY"] != "hello" {
		t.Errorf("expected trimmed value, got %q", res.Coerced["KEY"])
	}
	if len(res.Changes) != 1 || res.Changes[0] != "KEY" {
		t.Errorf("expected KEY in changes, got %v", res.Changes)
	}
}

func TestCoerce_NormalizeBools(t *testing.T) {
	cases := map[string]string{
		"yes": "true", "1": "true", "on": "true",
		"no": "false", "0": "false", "off": "false",
		"true": "true", "false": "false",
	}
	for input, want := range cases {
		res := Coerce(map[string]string{"F": input}, CoerceOptions{NormalizeBools: true})
		if res.Coerced["F"] != want {
			t.Errorf("input %q: expected %q, got %q", input, want, res.Coerced["F"])
		}
	}
}

func TestCoerce_LowercaseValues(t *testing.T) {
	res := Coerce(map[string]string{"KEY": "HELLO"}, CoerceOptions{LowercaseValues: true})
	if res.Coerced["KEY"] != "hello" {
		t.Errorf("expected lowercase, got %q", res.Coerced["KEY"])
	}
}

func TestCoerce_UppercaseValues(t *testing.T) {
	res := Coerce(map[string]string{"KEY": "hello"}, CoerceOptions{UppercaseValues: true})
	if res.Coerced["KEY"] != "HELLO" {
		t.Errorf("expected uppercase, got %q", res.Coerced["KEY"])
	}
}

func TestCoerce_SummaryWithChanges(t *testing.T) {
	res := Coerce(map[string]string{"A": "  x  "}, CoerceOptions{TrimWhitespace: true})
	s := res.Summary()
	if s == "coerce: no changes" {
		t.Error("expected non-empty summary")
	}
}

func TestCoerce_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"K": "  v  "}
	Coerce(input, CoerceOptions{TrimWhitespace: true})
	if input["K"] != "  v  " {
		t.Error("input was mutated")
	}
}
