package env

import (
	"fmt"
	"sort"
	"strings"
)

// CollectOptions controls how secrets are collected into a structured map.
type CollectOptions struct {
	// GroupBy is a key whose value is used to namespace collected entries.
	// If empty, all keys are placed at the top level.
	GroupBy string
	// Separator is placed between the group name and the key. Defaults to "_".
	Separator string
	// DryRun returns the result without mutating the input.
	DryRun bool
}

// CollectResult holds the output of a Collect operation.
type CollectResult struct {
	Out     map[string]string
	Groups  []string
	Summary string
}

// Collect groups secrets under a namespace derived from a key's value.
// All keys are prefixed with "<groupValue><sep><key>" when GroupBy is set.
// Keys that do not exist in secrets are left unchanged.
func Collect(secrets map[string]string, opts *CollectOptions) (CollectResult, error) {
	if opts == nil {
		return CollectResult{Out: copyCollectMap(secrets)}, nil
	}

	sep := opts.Separator
	if sep == "" {
		sep = "_"
	}

	out := make(map[string]string, len(secrets))
	groupSet := map[string]struct{}{}

	if opts.GroupBy == "" {
		for k, v := range secrets {
			out[k] = v
		}
		return CollectResult{Out: out, Summary: "no grouping applied"}, nil
	}

	groupVal, ok := secrets[opts.GroupBy]
	if !ok {
		return CollectResult{}, fmt.Errorf("collect: group-by key %q not found in secrets", opts.GroupBy)
	}

	groupVal = strings.ToUpper(strings.TrimSpace(groupVal))

	for k, v := range secrets {
		if k == opts.GroupBy {
			out[k] = v
			continue
		}
		newKey := groupVal + sep + k
		out[newKey] = v
		groupSet[groupVal] = struct{}{}
	}

	groups := make([]string, 0, len(groupSet))
	for g := range groupSet {
		groups = append(groups, g)
	}
	sort.Strings(groups)

	summary := fmt.Sprintf("collected %d keys under group %q", len(out)-1, groupVal)

	return CollectResult{Out: out, Groups: groups, Summary: summary}, nil
}

func copyCollectMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
