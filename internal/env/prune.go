package env

import "strings"

// PruneOptions controls how Prune behaves.
type PruneOptions struct {
	// RemoveEmpty removes keys whose values are empty strings.
	RemoveEmpty bool
	// RemoveKeys removes specific keys by exact match.
	RemoveKeys []string
	// RemovePrefix removes keys that match any of the given prefixes.
	RemovePrefix []string
	// DryRun returns the result without reporting a mutation intent.
	DryRun bool
}

// PruneResult holds the outcome of a Prune operation.
type PruneResult struct {
	Pruned []string
}

// Summary returns a human-readable description of what was pruned.
func (r PruneResult) Summary() string {
	if len(r.Pruned) == 0 {
		return "prune: nothing removed"
	}
	return "prune: removed keys: " + strings.Join(r.Pruned, ", ")
}

// Prune removes keys from secrets according to the provided options.
// It never mutates the input map.
func Prune(secrets map[string]string, opts PruneOptions) (map[string]string, PruneResult) {
	removeSet := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		removeSet[k] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	var pruned []string

	for k, v := range secrets {
		if _, ok := removeSet[k]; ok {
			pruned = append(pruned, k)
			continue
		}
		if opts.RemoveEmpty && v == "" {
			pruned = append(pruned, k)
			continue
		}
		if hasAnyPrefix(k, opts.RemovePrefix) {
			pruned = append(pruned, k)
			continue
		}
		out[k] = v
	}

	sortStrings(pruned)
	return out, PruneResult{Pruned: pruned}
}
