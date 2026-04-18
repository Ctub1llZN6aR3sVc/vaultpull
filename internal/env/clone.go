package env

import "fmt"

// CloneOptions controls how secrets are cloned between maps.
type CloneOptions struct {
	Keys      []string // if set, only clone these keys
	Overwrite bool     // overwrite existing keys in dst
	DryRun    bool     // do not mutate dst
}

// CloneResult summarises a clone operation.
type CloneResult struct {
	Cloned  []string
	Skipped []string
}

func (r CloneResult) Summary() string {
	if len(r.Cloned) == 0 {
		return "clone: no keys copied"
	}
	return fmt.Sprintf("clone: %d copied, %d skipped", len(r.Cloned), len(r.Skipped))
}

// Clone copies secrets from src into dst according to opts.
func Clone(src, dst map[string]string, opts CloneOptions) (map[string]string, CloneResult) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range src {
			keys = append(keys, k)
		}
	}

	var result CloneResult
	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if _, exists := out[k]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if !opts.DryRun {
			out[k] = v
		}
		result.Cloned = append(result.Cloned, k)
	}
	return out, result
}
