package env

import (
	"fmt"
	"strings"
)

// ComputeOptions controls how computed keys are derived from existing secrets.
type ComputeOptions struct {
	// Rules maps a new key name to a Go text/template-style expression.
	// Expressions may reference existing keys using {{.KEY_NAME}} syntax.
	Rules map[string]string

	// Overwrite allows computed values to replace existing keys.
	// When false, existing keys are preserved.
	Overwrite bool

	// DryRun returns the result without modifying the input map.
	DryRun bool

	// FailOnError causes Compute to return an error if any expression
	// references a missing key or produces an empty value.
	FailOnError bool
}

// ComputeResult holds the outcome of a Compute operation.
type ComputeResult struct {
	// Added lists keys that were newly created.
	Added []string

	// Skipped lists keys that were not written because Overwrite was false.
	Skipped []string

	// Errors lists keys whose expressions could not be evaluated.
	Errors []string
}

// Summary returns a human-readable description of the compute result.
func (r ComputeResult) Summary() string {
	parts := []string{}
	if len(r.Added) > 0 {
		parts = append(parts, fmt.Sprintf("%d added", len(r.Added)))
	}
	if len(r.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", len(r.Skipped)))
	}
	if len(r.Errors) > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", len(r.Errors)))
	}
	if len(parts) == 0 {
		return "compute: no changes"
	}
	return "compute: " + strings.Join(parts, ", ")
}

// Compute evaluates each rule in opts.Rules against the provided secrets map
// and injects the resulting key/value pairs. Expressions are resolved using
// the same {{.KEY}} / $KEY interpolation supported by env.RenderTemplate.
//
// Example:
//
//	secrets := map[string]string{"HOST": "db.example.com", "PORT": "5432"}
//	opts := &ComputeOptions{
//		Rules: map[string]string{
//			"DATABASE_URL": "postgres://{{.HOST}}:{{.PORT}}/mydb",
//		},
//	}
//	result, _ := Compute(secrets, opts)
//	// secrets["DATABASE_URL"] == "postgres://db.example.com:5432/mydb"
func Compute(secrets map[string]string, opts *ComputeOptions) (ComputeResult, error) {
	var result ComputeResult

	if opts == nil || len(opts.Rules) == 0 {
		return result, nil
	}

	// Work on a copy so we don't mutate the caller's map during evaluation.
	working := make(map[string]string, len(secrets))
	for k, v := range secrets {
		working[k] = v
	}

	// Sort rule keys for deterministic evaluation order.
	keys := make([]string, 0, len(opts.Rules))
	for k := range opts.Rules {
		keys = append(keys, k)
	}
	sortStrings(keys)

	for _, key := range keys {
		expr := opts.Rules[key]

		value, err := RenderTemplate(expr, working)
		if err != nil {
			result.Errors = append(result.Errors, key)
			if opts.FailOnError {
				return result, fmt.Errorf("compute: rule %q failed: %w", key, err)
			}
			continue
		}

		if value == "" && opts.FailOnError {
			result.Errors = append(result.Errors, key)
			return result, fmt.Errorf("compute: rule %q produced an empty value", key)
		}

		if _, exists := working[key]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, key)
			continue
		}

		working[key] = value
		result.Added = append(result.Added, key)
	}

	if !opts.DryRun {
		for k, v := range working {
			secrets[k] = v
		}
	}

	return result, nil
}
