package env

import (
	"testing"
)

func TestSplit_NoOptions(t *testing.T) {
	secrets := map[string]string{"FOO": "a,b,c", "BAR": "hello"}
	out, res, err := Split(secrets, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "a,b,c" {
		t.Errorf("expected FOO unchanged, got %q", out["FOO"])
	}
	if len(res.Expanded) != 0 {
		t.Errorf("expected no expansions, got %v", res.Expanded)
	}
}

func TestSplit_ExpandsAllKeys(t *testing.T) {
	secrets := map[string]string{"HOSTS": "a,b,c"}
	out, res, err := Split(secrets, &SplitOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["HOSTS"]; ok {
		t.Error("original key should be removed after split")
	}
	if out["HOSTS_1"] != "a" || out["HOSTS_2"] != "b" || out["HOSTS_3"] != "c" {
		t.Errorf("unexpected split output: %v", out)
	}
	if len(res.Expanded["HOSTS"]) != 3 {
		t.Errorf("expected 3 parts, got %v", res.Expanded["HOSTS"])
	}
}

func TestSplit_ZeroIndexed(t *testing.T) {
	secrets := map[string]string{"ITEMS": "x,y"}
	out, _, err := Split(secrets, &SplitOptions{ZeroIndexed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ITEMS_0"] != "x" || out["ITEMS_1"] != "y" {
		t.Errorf("unexpected zero-indexed output: %v", out)
	}
}

func TestSplit_CustomDelimiter(t *testing.T) {
	secrets := map[string]string{"PATHS": "/usr/bin:/usr/local/bin"}
	out, res, err := Split(secrets, &SplitOptions{Delimiter: ":"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PATHS_1"] != "/usr/bin" || out["PATHS_2"] != "/usr/local/bin" {
		t.Errorf("unexpected output: %v", out)
	}
	if len(res.Expanded["PATHS"]) != 2 {
		t.Errorf("expected 2 parts, got %v", res.Expanded)
	}
}

func TestSplit_RestrictedToKeys(t *testing.T) {
	secrets := map[string]string{"FOO": "a,b", "BAR": "c,d"}
	out, res, err := Split(secrets, &SplitOptions{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BAR"] != "c,d" {
		t.Errorf("BAR should be unchanged, got %q", out["BAR"])
	}
	if _, ok := out["FOO"]; ok {
		t.Error("FOO should have been replaced by indexed keys")
	}
	if len(res.Expanded) != 1 {
		t.Errorf("expected 1 expanded key, got %v", res.Expanded)
	}
}

func TestSplit_NoDelimiterInValue(t *testing.T) {
	secrets := map[string]string{"KEY": "singlevalue"}
	out, res, err := Split(secrets, &SplitOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "singlevalue" {
		t.Errorf("expected KEY unchanged, got %q", out["KEY"])
	}
	if len(res.Expanded) != 0 {
		t.Errorf("expected no expansions, got %v", res.Expanded)
	}
}

func TestSplit_SummaryNoExpansions(t *testing.T) {
	res := SplitResult{Expanded: map[string][]string{}}
	if res.Summary() != "split: no keys expanded" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}

func TestSplit_SummaryWithExpansions(t *testing.T) {
	res := SplitResult{Expanded: map[string][]string{"FOO": {"a", "b"}}}
	if res.Summary() != "split: expanded 1 key(s): FOO" {
		t.Errorf("unexpected summary: %s", res.Summary())
	}
}
