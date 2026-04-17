package env

import (
	"testing"
)

func TestTag_TaggedAndUntagged(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "x", "API_KEY": "y", "PORT": "8080"}
	tagMap := map[string]string{"DB_PASS": "sensitive,db", "API_KEY": "sensitive"}
	res := Tag(secrets, tagMap)
	if len(res.Tagged) != 2 {
		t.Fatalf("expected 2 tagged, got %d", len(res.Tagged))
	}
	if len(res.Untagged) != 1 || res.Untagged[0] != "PORT" {
		t.Fatalf("expected [PORT] untagged, got %v", res.Untagged)
	}
}

func TestTag_AllUntaggedWhenMapEmpty(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	res := Tag(secrets, map[string]string{})
	if len(res.Tagged) != 0 {
		t.Fatal("expected no tagged keys")
	}
	if len(res.Untagged) != 2 {
		t.Fatalf("expected 2 untagged, got %d", len(res.Untagged))
	}
}

func TestTag_UntaggedIsSorted(t *testing.T) {
	secrets := map[string]string{"Z": "1", "A": "2", "M": "3"}
	res := Tag(secrets, map[string]string{})
	for i := 1; i < len(res.Untagged); i++ {
		if res.Untagged[i] < res.Untagged[i-1] {
			t.Fatal("untagged not sorted")
		}
	}
}

func TestFilterByTag_ReturnsMatchingKeys(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "x", "API_KEY": "y", "PORT": "8080"}
	tagMap := map[string]string{"DB_PASS": "sensitive,db", "API_KEY": "sensitive"}
	out := FilterByTag(secrets, tagMap, "sensitive")
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	if _, ok := out["PORT"]; ok {
		t.Fatal("PORT should not be included")
	}
}

func TestFilterByTag_MultipleRequiredTags(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "x", "API_KEY": "y"}
	tagMap := map[string]string{"DB_PASS": "sensitive,db", "API_KEY": "sensitive"}
	out := FilterByTag(secrets, tagMap, "sensitive", "db")
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; !ok {
		t.Fatal("expected DB_PASS")
	}
}

func TestFilterByTag_NoRequiredTagsReturnsAll(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	tagMap := map[string]string{}
	out := FilterByTag(secrets, tagMap)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestSummaryByTag(t *testing.T) {
	res := TagResult{Tagged: map[string][]string{"A": {"x"}}, Untagged: []string{"B", "C"}}
	s := SummaryByTag(res)
	if s != "1 tagged, 2 untagged" {
		t.Fatalf("unexpected summary: %s", s)
	}
}
