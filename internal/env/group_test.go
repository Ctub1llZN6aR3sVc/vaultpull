package env

import (
	"testing"
)

func TestGroup_PartitionsByPrefix(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "vaultpull",
		"OTHER":    "value",
	}
	res := Group(secrets, []string{"DB", "APP"}, "_")
	if res.Groups["DB"]["HOST"] != "localhost" {
		t.Errorf("expected DB HOST=localhost")
	}
	if res.Groups["DB"]["PORT"] != "5432" {
		t.Errorf("expected DB PORT=5432")
	}
	if res.Groups["APP"]["NAME"] != "vaultpull" {
		t.Errorf("expected APP NAME=vaultpull")
	}
	if res.Ungrouped["OTHER"] != "value" {
		t.Errorf("expected OTHER in ungrouped")
	}
}

func TestGroup_UngroupedWhenNoPrefixes(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := Group(secrets, nil, "_")
	if len(res.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %d", len(res.Ungrouped))
	}
	if len(res.Groups) != 0 {
		t.Errorf("expected no groups")
	}
}

func TestGroup_DefaultSeparator(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "h"}
	res := Group(secrets, []string{"DB"}, "")
	if res.Groups["DB"]["HOST"] != "h" {
		t.Errorf("expected default separator to work")
	}
}

func TestGroup_Summary(t *testing.T) {
	secrets := map[string]string{"A_X": "1", "B_Y": "2", "B_Z": "3"}
	res := Group(secrets, []string{"A", "B"}, "_")
	s := res.Summary()
	if s != "A(1), B(2)" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestGroup_SummaryNoGroups(t *testing.T) {
	res := &GroupResult{Groups: map[string]map[string]string{}, Ungrouped: map[string]string{}}
	if res.Summary() != "no groups found" {
		t.Errorf("expected empty summary message")
	}
}

func TestGroup_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "h"}
	origLen := len(secrets)
	Group(secrets, []string{"DB"}, "_")
	if len(secrets) != origLen {
		t.Errorf("input map was mutated")
	}
}
