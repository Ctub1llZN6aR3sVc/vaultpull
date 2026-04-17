package env

import (
	"fmt"
	"os"
	"strings"
)

// ResolveOptions controls how secret values are resolved.
type ResolveOptions struct {
	// FallbackToEnv allows falling back to OS environment variables when a key
	// is missing from the secrets map.
	FallbackToEnv bool
	// Required lists keys that must be present after resolution.
	Required []string
	// Defaults provides fallback values for missing keys (applied before env lookup).
	Defaults map[string]string
}

// ResolveResult holds the resolved secrets and any warnings.
type ResolveResult struct {
	Secrets  map[string]string
	Warnings []string
}

// Resolve merges defaults, vault secrets, and optionally OS env vars,
// then validates required keys are present.
func Resolve(secrets map[string]string, opts ResolveOptions) (ResolveResult, error) {
	out := make(map[string]string, len(secrets))
	var warnings []string

	// Apply defaults first.
	for k, v := range opts.Defaults {
		out[k] = v
	}

	// Overlay vault secrets.
	for k, v := range secrets {
		out[k] = v
	}

	// Fall back to OS env for missing keys listed in Required.
	if opts.FallbackToEnv {
		for _, key := range opts.Required {
			if _, ok := out[key]; !ok {
				if val, found := os.LookupEnv(key); found {
					out[key] = val
					warnings = append(warnings, fmt.Sprintf("key %q resolved from OS environment", key))
				}
			}
		}
	}

	// Validate required keys.
	var missing []string
	for _, key := range opts.Required {
		if v, ok := out[key]; !ok || strings.TrimSpace(v) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return ResolveResult{}, fmt.Errorf("resolve: missing required keys: %s", strings.Join(missing, ", "))
	}

	return ResolveResult{Secrets: out, Warnings: warnings}, nil
}
