package env

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for env keys.
type SortOrder string

const (
	SortAlpha      SortOrder = "alpha"       // alphabetical ascending
	SortAlphaDesc  SortOrder = "alpha_desc"  // alphabetical descending
	SortByPrefix   SortOrder = "prefix"      // group by common prefix, then alpha
	SortNatural    SortOrder = "natural"     // natural / insertion order (no-op)
)

// SortOptions controls how Sort behaves.
type SortOptions struct {
	Order     SortOrder
	// PrefixSep is the separator used to detect prefix groups (default "_").
	PrefixSep string
}

// SortResult holds the ordered key list and a human-readable summary.
type SortResult struct {
	Keys    []string
	Summary string
}

// Sort returns a SortResult whose Keys slice contains all keys from secrets
// ordered according to opts. The secrets map itself is not mutated.
func Sort(secrets map[string]string, opts SortOptions) SortResult {
	if opts.PrefixSep == "" {
		opts.PrefixSep = "_"
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	switch opts.Order {
	case SortAlphaDesc:
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})
	case SortByPrefix:
		keys = sortByPrefix(keys, opts.PrefixSep)
	case SortNatural:
		// preserve map iteration — already random, but callers can pre-seed order
		// nothing to do
	default: // SortAlpha and anything unrecognised
		sort.Strings(keys)
	}

	summary := buildSortSummary(keys, opts.Order)
	return SortResult{Keys: keys, Summary: summary}
}

// sortByPrefix groups keys sharing the same prefix (text before the first
// occurrence of sep) and sorts groups alphabetically, with keys inside each
// group also sorted alphabetically.
func sortByPrefix(keys []string, sep string) []string {
	groups := map[string][]string{}
	order := []string{}
	seen := map[string]bool{}

	for _, k := range keys {
		prefix := k
		if idx := strings.Index(k, sep); idx != -1 {
			prefix = k[:idx]
		}
		if !seen[prefix] {
			seen[prefix] = true
			order = append(order, prefix)
		}
		groups[prefix] = append(groups[prefix], k)
	}

	sort.Strings(order)

	result := make([]string, 0, len(keys))
	for _, prefix := range order {
		g := groups[prefix]
		sort.Strings(g)
		result = append(result, g...)
	}
	return result
}

func buildSortSummary(keys []string, order SortOrder) string {
	if len(keys) == 0 {
		return "no keys to sort"
	}
	return strings.Join([]string{
		"sorted",
		string(order),
		"—",
		strconv(len(keys)),
		"keys",
	}, " ")
}

// strconv is a tiny helper to avoid importing fmt just for Itoa.
func strconv(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
