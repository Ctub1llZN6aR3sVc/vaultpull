package env

import (
	"fmt"
	"strings"
)

// SubstituteOptions controls how variable substitution is applied.
type SubstituteOptions struct {
	// Strict causes an error when a referenced variable is not found.
	Strict bool
	// Fallback is the default value used for missing variables when not strict.
	Fallback string
	// DryRun returns a result without modifying the input.
	DryRun bool
}

// SubstituteResult holds the outcome of a substitution pass.
type SubstituteResult struct {
	Substituted []string
	Unresolved  []string
}

func (r SubstituteResult) Summary() string {
	if len(r.Unresolved) > 0 {
		return fmt.Sprintf("%d substituted, %d unresolved: %s",
			len(r.Substituted), len(r.Unresolved), strings.Join(r.Unresolved, ", "))
	}
	return fmt.Sprintf("%d substituted", len(r.Substituted))
}

// Substitute performs ${VAR} and $VAR style substitution within values,
// resolving references from the same secrets map.
func Substitute(secrets map[string]string, opts *SubstituteOptions) (map[string]string, SubstituteResult, error) {
	if opts == nil {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out, SubstituteResult{}, nil
	}

	out := make(map[string]string, len(secrets))
	var result SubstituteResult

	for k, v := range secrets {
		expanded, unresolved, err := expandValue(v, secrets, opts)
		if err != nil {
			return nil, result, fmt.Errorf("key %q: %w", k, err)
		}
		if expanded != v {
			result.Substituted = append(result.Substituted, k)
		}
		result.Unresolved = append(result.Unresolved, unresolved...)
		out[k] = expanded
	}

	if opts.DryRun {
		return secrets, result, nil
	}
	return out, result, nil
}

func expandValue(val string, lookup map[string]string, opts *SubstituteOptions) (string, []string, error) {
	var unresolved []string
	result := varPattern.ReplaceAllStringFunc(val, func(match string) string {
		key := extractKey(match)
		if v, ok := lookup[key]; ok {
			return v
		}
		unresolved = append(unresolved, key)
		if opts.Strict {
			return match // signal via unresolved
		}
		if opts.Fallback != "" {
			return opts.Fallback
		}
		return ""
	})
	if opts.Strict && len(unresolved) > 0 {
		return "", unresolved, fmt.Errorf("unresolved variables: %s", strings.Join(unresolved, ", "))
	}
	return result, unresolved, nil
}
