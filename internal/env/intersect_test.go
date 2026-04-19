package env

import (
	"testing"
)

func TestIntersect_CommonKeys(t *testing.T) {
	a := map[string]string{"FOO": "foo_a", "BAR": "bar_a", "ONLY_A": "x"}
	b := map[string]string{"FOO": "foo_b", "BAR": "bar_b", "ONLY_B": "y"}

	res := Intersect(a, b, IntersectOptions{})

	if _, ok := res.Secrets["FOO"]; !ok {
		t.Error("expected FOO in result")
	}
	if _, ok := res.Secrets["BAR"]; !ok {
		t.Error("expected BAR in result")
	}
	if _, ok := res.Secrets["ONLY_A"]; ok {
		t.Error("ONLY_A should not be in result")
	}
	if _, ok := res.Secrets["ONLY_B"]; ok {
		t.Error("ONLY_B should not be in result")
	}
}

func TestIntersect_KeepsValueFromA(t *testing.T) {
	a := map[string]string{"KEY": "from_a"}
	b := map[string]string{"KEY": "from_b"}

	res := Intersect(a, b, IntersectOptions{KeepValue: "a"})
	if res.Secrets["KEY"] != "from_a" {
		t.Errorf("expected from_a, got %s", res.Secrets["KEY"])
	}
}

func TestIntersect_KeepsValueFromB(t *testing.T) {
	a := map[string]string{"KEY": "from_a"}
	b := map[string]string{"KEY": "from_b"}

	res := Intersect(a, b, IntersectOptions{KeepValue: "b"})
	if res.Secrets["KEY"] != "from_b" {
		t.Errorf("expected from_b, got %s", res.Secrets["KEY"])
	}
}

func TestIntersect_NoCommonKeys(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}

	res := Intersect(a, b, IntersectOptions{})
	if len(res.Secrets) != 0 {
		t.Errorf("expected empty result, got %v", res.Secrets)
	}
	if res.Summary() != "intersect: no common keys found" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestIntersect_DoesNotMutateInputs(t *testing.T) {
	a := map[string]string{"K": "v"}
	b := map[string]string{"K": "v2"}

	Intersect(a, b, IntersectOptions{})
	if a["K"] != "v" {
		t.Error("a was mutated")
	}
}

func TestIntersect_SummaryWithKeys(t *testing.T) {
	a := map[string]string{"X": "1", "Y": "2"}
	b := map[string]string{"X": "3", "Y": "4"}

	res := Intersect(a, b, IntersectOptions{})
	if res.Summary() != "intersect: 2 common key(s) retained" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
