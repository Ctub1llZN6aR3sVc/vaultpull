package env

import (
	"testing"
)

func TestIndex_ReturnsAllKeys(t *testing.T) {
	secrets := map[string]string{"B": "2", "A": "1", "C": "3"}
	r := Index(secrets, nil)
	if r.Total != 3 {
		t.Fatalf("expected 3, got %d", r.Total)
	}
	if r.Entries[0].Key != "A" || r.Entries[1].Key != "B" || r.Entries[2].Key != "C" {
		t.Fatalf("expected sorted keys, got %v", r.Entries)
	}
}

func TestIndex_SetsSource(t *testing.T) {
	secrets := map[string]string{"KEY": "val"}
	r := Index(secrets, &IndexOptions{Source: "vault"})
	if r.Entries[0].Source != "vault" {
		t.Fatalf("expected source vault, got %s", r.Entries[0].Source)
	}
}

func TestIndex_SetsTags(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "x"}
	tags := map[string][]string{"DB_PASS": {"sensitive", "db"}}
	r := Index(secrets, &IndexOptions{Tags: tags})
	if len(r.Entries[0].Tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", r.Entries[0].Tags)
	}
}

func TestIndex_EmptySecrets(t *testing.T) {
	r := Index(map[string]string{}, nil)
	if r.Total != 0 || len(r.Entries) != 0 {
		t.Fatal("expected empty result")
	}
}

func TestIndexSummary_NoEntries(t *testing.T) {
	s := IndexSummary(IndexResult{})
	if s != "index: no entries" {
		t.Fatalf("unexpected: %s", s)
	}
}

func TestIndexSummary_WithEntries(t *testing.T) {
	r := Index(map[string]string{"A": "1", "B": "2"}, nil)
	s := IndexSummary(r)
	if s != "index: 2 entries" {
		t.Fatalf("unexpected: %s", s)
	}
}
