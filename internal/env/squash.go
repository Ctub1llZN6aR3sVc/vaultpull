package env

import (
	"fmt"
	"sort"
	"strings"
)

// SquashOptions controls how Squash behaves.
type SquashOptions struct {
	// Separator is placed between values when joining duplicates.
	// Defaults to ",".
	Separator string
	// KeepFirst retains only the first occurrence instead of joining.
	KeepFirst bool
	// KeepLast retains only the last occurrence instead of joining.
	KeepLast bool
}

// SquashResult describes what Squash did.
type SquashResult struct {
	Squashed []string // keys whose values were merged or deduplicated
}

func (r SquashResult) Summary() string {
	if len(r.Squashed) == 0 {
		return "squash: no duplicate keys found"
	}
	sort.Strings(r.Squashed)
	return fmt.Sprintf("squash: merged keys: %s", strings.Join(r.Squashed, ", "))
}

// Squash merges duplicate keys from multiple secret maps into a single map.
// When the same key appears in more than one source the values are joined
// using Separator (default ","), or the first/last value is kept when the
// corresponding flag is set.
func Squash(sources []map[string]string, opts *SquashOptions) (map[string]string, SquashResult) {
	if opts == nil {
		opts = &SquashOptions{}
	}
	sep := opts.Separator
	if sep == "" {
		sep = ","
	}

	// Track insertion order for deterministic output.
	order := []string{}
	seen := map[string]bool{}
	accum := map[string][]string{}

	for _, src := range sources {
		// Sort keys within each source for deterministic processing.
		keys := make([]string, 0, len(src))
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if !seen[k] {
				seen[k] = true
				order = append(order, k)
			}
			accum[k] = append(accum[k], src[k])
		}
	}

	result := map[string]string{}
	var squashed []string

	for _, k := range order {
		vals := accum[k]
		if len(vals) > 1 {
			squashed = append(squashed, k)
		}
		switch {
		case opts.KeepFirst:
			result[k] = vals[0]
		case opts.KeepLast:
			result[k] = vals[len(vals)-1]
		default:
			result[k] = strings.Join(vals, sep)
		}
	}

	return result, SquashResult{Squashed: squashed}
}
