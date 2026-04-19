package env

import (
	"fmt"
	"sort"
)

// ExtractOptions controls how secrets are extracted.
type ExtractOptions struct {
	Keys          []string
	Prefixes      []string
	StripPrefix   bool
	FailOnMissing bool
}

// ExtractResult holds the extracted secrets and metadata.
type ExtractResult struct {
	Secrets map[string]string
	Missing []string
}

// Summary returns a human-readable summary of the extraction.
func (r ExtractResult) Summary() string {
	if len(r.Missing) == 0 {
		return "extract: ok"
	}
	return "extract: missing keys: " + joinExtractKeys(r.Missing)
}

func joinExtractKeys(keys []string) string {
	out := ""
	for i, k := range keys {
		if i > 0 {
			out += ", "
		}
		out += k
	}
	return out
}

// Extract pulls a subset of secrets by explicit key list or prefix.
func Extract(secrets map[string]string, opts ExtractOptions) (ExtractResult, error) {
	out := make(map[string]string)
	var missing []string

	for _, k := range opts.Keys {
		v, ok := secrets[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		out[k] = v
	}

	for _, p := range opts.Prefixes {
		for k, v := range secrets {
			if len(k) >= len(p) && k[:len(p)] == p {
				newKey := k
				if opts.StripPrefix {
					newKey = k[len(p):]
				}
				if newKey != "" {
					out[newKey] = v
				}
			}
		}
	}

	if opts.FailOnMissing && len(missing) > 0 {
		sort.Strings(missing)
		return ExtractResult{}, fmt.Errorf("extract: missing required keys: %s", joinExtractKeys(missing))
	}

	sort.Strings(missing)
	return ExtractResult{Secrets: out, Missing: missing}, nil
}
