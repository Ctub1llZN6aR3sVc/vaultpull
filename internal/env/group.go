package env

import (
	"fmt"
	"sort"
	"strings"
)

// GroupResult holds secrets partitioned by group prefix.
type GroupResult struct {
	Groups map[string]map[string]string
	Ungrouped map[string]string
}

// Summary returns a human-readable summary of groups.
func (g *GroupResult) Summary() string {
	if len(g.Groups) == 0 {
		return "no groups found"
	}
	names := make([]string, 0, len(g.Groups))
	for k := range g.Groups {
		names = append(names, k)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, n := range names {
		parts = append(parts, fmt.Sprintf("%s(%d)", n, len(g.Groups[n])))
	}
	return strings.Join(parts, ", ")
}

// Group partitions secrets by key prefix using sep as the delimiter.
// Keys without a matching prefix land in Ungrouped.
func Group(secrets map[string]string, prefixes []string, sep string) *GroupResult {
	if sep == "" {
		sep = "_"
	}
	result := &GroupResult{
		Groups:    make(map[string]map[string]string),
		Ungrouped: make(map[string]string),
	}
	for _, p := range prefixes {
		result.Groups[p] = make(map[string]string)
	}
	for k, v := range secrets {
		matched := false
		for _, p := range prefixes {
			token := p + sep
			if strings.HasPrefix(k, token) {
				stripped := strings.TrimPrefix(k, token)
				result.Groups[p][stripped] = v
				matched = true
				break
			}
		}
		if !matched {
			result.Ungrouped[k] = v
		}
	}
	return result
}
