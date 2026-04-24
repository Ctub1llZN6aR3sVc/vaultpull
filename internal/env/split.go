package env

import (
	"fmt"
	"sort"
	"strings"
)

// SplitOptions controls how Split behaves.
type SplitOptions struct {
	// Delimiter is the character used to split values (default: ",").
	Delimiter string
	// Keys restricts splitting to only the specified keys.
	// If empty, all keys are considered.
	Keys []string
	// IndexedKeys emits KEY_0, KEY_1, ... instead of KEY_1, KEY_2, ...
	ZeroIndexed bool
	// DryRun returns the result without mutating anything.
	DryRun bool
}

// SplitResult holds the outcome of a Split operation.
type SplitResult struct {
	Expanded map[string][]string // original key -> split parts
	Errors   []string
}

// Summary returns a human-readable description of the split.
func (r SplitResult) Summary() string {
	if len(r.Expanded) == 0 {
		return "split: no keys expanded"
	}
	keys := make([]string, 0, len(r.Expanded))
	for k := range r.Expanded {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return fmt.Sprintf("split: expanded %d key(s): %s", len(keys), strings.Join(keys, ", "))
}

// Split expands delimited values into indexed keys.
// For example, FOO="a,b,c" becomes FOO_1="a", FOO_2="b", FOO_3="c".
// The original key is removed from the result.
func Split(secrets map[string]string, opts *SplitOptions) (map[string]string, SplitResult, error) {
	if opts == nil {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result, SplitResult{}, nil
	}

	delim := opts.Delimiter
	if delim == "" {
		delim = ","
	}

	targetKeys := toSet(opts.Keys)

	out := make(map[string]string, len(secrets))
	res := SplitResult{Expanded: make(map[string][]string)}

	for k, v := range secrets {
		if len(targetKeys) > 0 && !targetKeys[k] {
			out[k] = v
			continue
		}
		if !strings.Contains(v, delim) {
			out[k] = v
			continue
		}
		parts := strings.Split(v, delim)
		res.Expanded[k] = parts
		start := 1
		if opts.ZeroIndexed {
			start = 0
		}
		for i, part := range parts {
			newKey := fmt.Sprintf("%s_%d", k, i+start)
			out[newKey] = strings.TrimSpace(part)
		}
	}

	return out, res, nil
}
