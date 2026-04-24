package env

import "sort"

// ReorderOptions controls how secrets are reordered.
type ReorderOptions struct {
	// Keys defines the explicit ordering of keys. Keys not listed appear at the end.
	Keys []string
	// Alphabetical sorts all keys alphabetically, ignoring Keys ordering.
	Alphabetical bool
	// DryRun returns the result without mutating anything.
	DryRun bool
}

// ReorderResult holds the output of a Reorder operation.
type ReorderResult struct {
	Ordered  []string
	Unlisted []string
}

// Summary returns a human-readable description of the reorder result.
func (r ReorderResult) Summary() string {
	if len(r.Unlisted) == 0 {
		return "all keys reordered as specified"
	}
	return "reordered; " + joinReorderKeys(r.Unlisted) + " appended at end"
}

func joinReorderKeys(keys []string) string {
	out := ""
	for i, k := range keys {
		if i > 0 {
			out += ", "
		}
		out += k
	}
	return out
}

// Reorder returns a new map with keys ordered according to opts.
// Because Go maps are unordered, the result also includes the ordered key slice.
func Reorder(secrets map[string]string, opts ReorderOptions) (map[string]string, ReorderResult) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var ordered []string
	var unlisted []string

	if opts.Alphabetical {
		for k := range secrets {
			ordered = append(ordered, k)
		}
		sort.Strings(ordered)
		return out, ReorderResult{Ordered: ordered, Unlisted: nil}
	}

	seen := make(map[string]bool)
	for _, k := range opts.Keys {
		if _, ok := secrets[k]; ok {
			ordered = append(ordered, k)
			seen[k] = true
		}
	}

	var rest []string
	for k := range secrets {
		if !seen[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	ordered = append(ordered, rest...)
	unlisted = rest

	return out, ReorderResult{Ordered: ordered, Unlisted: unlisted}
}
