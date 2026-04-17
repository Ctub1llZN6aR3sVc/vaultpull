package env

import (
	"fmt"
	"sort"
	"strings"
)

// TagEntry associates a secret key with a set of string tags.
type TagEntry struct {
	Key  string
	Tags []string
}

// TagResult holds the tagged secrets and any keys that had no tags.
type TagResult struct {
	Tagged   map[string][]string // key -> tags
	Untagged []string
}

// Tag applies a tag map (key -> comma-separated tags) to a secrets map.
// Keys present in secrets but absent from tagMap are collected as Untagged.
func Tag(secrets map[string]string, tagMap map[string]string) TagResult {
	result := TagResult{
		Tagged: make(map[string][]string),
	}
	for k := range secrets {
		raw, ok := tagMap[k]
		if !ok || strings.TrimSpace(raw) == "" {
			result.Untagged = append(result.Untagged, k)
			continue
		}
		parts := strings.Split(raw, ",")
		var tags []string
		for _, p := range parts {
			t := strings.TrimSpace(p)
			if t != "" {
				tags = append(tags, t)
			}
		}
		result.Tagged[k] = tags
	}
	sort.Strings(result.Untagged)
	return result
}

// FilterByTag returns secrets whose keys have ALL of the requested tags.
func FilterByTag(secrets map[string]string, tagMap map[string]string, required ...string) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		raw := tagMap[k]
		tagSet := toTagSet(raw)
		if hasAllTags(tagSet, required) {
			out[k] = v
		}
	}
	return out
}

// SummaryByTag returns a human-readable summary line.
func SummaryByTag(result TagResult) string {
	return fmt.Sprintf("%d tagged, %d untagged", len(result.Tagged), len(result.Untagged))
}

func toTagSet(raw string) map[string]struct{} {
	s := make(map[string]struct{})
	for _, p := range strings.Split(raw, ",") {
		t := strings.TrimSpace(p)
		if t != "" {
			s[t] = struct{}{}
		}
	}
	return s
}

func hasAllTags(tagSet map[string]struct{}, required []string) bool {
	for _, r := range required {
		if _, ok := tagSet[r]; !ok {
			return false
		}
	}
	return true
}
