package env

import "fmt"

// AliasOptions controls how aliases are applied.
type AliasOptions struct {
	// Map of new_key -> existing_key. The value of existing_key is copied to new_key.
	Aliases map[string]string
	// If true, remove the original key after aliasing.
	RemoveSource bool
	// If true, do not mutate; return what would change.
	DryRun bool
}

// AliasResult holds the outcome of an Alias operation.
type AliasResult struct {
	Aliased []string
	Missing []string
}

func (r AliasResult) Summary() string {
	if len(r.Aliased) == 0 {
		return "alias: no keys aliased"
	}
	return fmt.Sprintf("alias: %d aliased, %d source keys missing", len(r.Aliased), len(r.Missing))
}

// Alias copies values from source keys to alias keys in secrets.
func Alias(secrets map[string]string, opts AliasOptions) (map[string]string, AliasResult) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var result AliasResult
	for newKey, srcKey := range opts.Aliases {
		v, ok := out[srcKey]
		if !ok {
			result.Missing = append(result.Missing, srcKey)
			continue
		}
		if !opts.DryRun {
			out[newKey] = v
			if opts.RemoveSource {
				delete(out, srcKey)
			}
		}
		result.Aliased = append(result.Aliased, newKey)
	}
	return out, result
}
