package env

import (
	"fmt"
	"strings"
)

// CondenseOptions controls how Condense behaves.
type CondenseOptions struct {
	// Separator is placed between joined values. Defaults to ",".
	Separator string
	// Keys lists the specific keys whose values should be condensed into one.
	// If empty, all keys are joined.
	Keys []string
	// TargetKey is the key written with the combined value.
	// If empty, the first key (sorted) is used.
	TargetKey string
	// DryRun reports what would change without mutating the map.
	DryRun bool
}

// CondenseResult describes what Condense did.
type CondenseResult struct {
	TargetKey  string
	SourceKeys []string
	Value      string
	DryRun     bool
}

func (r CondenseResult) Summary() string {
	if len(r.SourceKeys) == 0 {
		return "condense: no keys matched"
	}
	action := "condensed"
	if r.DryRun {
		action = "would condense"
	}
	return fmt.Sprintf("condense: %s %d keys into %q", action, len(r.SourceKeys), r.TargetKey)
}

// Condense joins values from selected keys into a single key.
// The source keys (excluding the target) are removed from the result.
func Condense(secrets map[string]string, opts *CondenseOptions) (map[string]string, CondenseResult, error) {
	if opts == nil {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out, CondenseResult{}, nil
	}

	sep := opts.Separator
	if sep == "" {
		sep = ","
	}

	candidate := opts.Keys
	if len(candidate) == 0 {
		for k := range secrets {
			candidate = append(candidate, k)
		}
	}
	sortStrings(candidate)

	var parts []string
	var used []string
	for _, k := range candidate {
		v, ok := secrets[k]
		if !ok {
			continue
		}
		parts = append(parts, v)
		used = append(used, k)
	}

	target := opts.TargetKey
	if target == "" && len(used) > 0 {
		target = used[0]
	}

	joined := strings.Join(parts, sep)
	res := CondenseResult{
		TargetKey:  target,
		SourceKeys: used,
		Value:      joined,
		DryRun:     opts.DryRun,
	}

	if opts.DryRun {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out, res, nil
	}

	out := make(map[string]string)
	for k, v := range secrets {
		out[k] = v
	}
	for _, k := range used {
		delete(out, k)
	}
	if target != "" {
		out[target] = joined
	}
	return out, res, nil
}
