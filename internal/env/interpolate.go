package env

import (
	"fmt"
	"regexp"
	"strings"
)

// InterpolateOptions controls how interpolation is applied.
type InterpolateOptions struct {
	// Strict returns an error if a referenced key is not found in the map.
	Strict bool
	// Fallback is returned for missing keys when Strict is false.
	Fallback string
	// DryRun returns the result without modifying the input.
	DryRun bool
}

// InterpolateResult holds the outcome of an interpolation pass.
type InterpolateResult struct {
	Resolved []string
	Missing  []string
}

func (r InterpolateResult) Summary() string {
	if len(r.Missing) == 0 {
		return fmt.Sprintf("interpolated %d key(s), no missing references", len(r.Resolved))
	}
	return fmt.Sprintf("interpolated %d key(s), %d missing reference(s): %s",
		len(r.Resolved), len(r.Missing), strings.Join(r.Missing, ", "))
}

var interpolateRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

// Interpolate replaces ${KEY} references within values using other entries in
// the same map. Self-referential keys are left unchanged to avoid infinite loops.
func Interpolate(secrets map[string]string, opts *InterpolateOptions) (map[string]string, InterpolateResult, error) {
	if opts == nil {
		opts = &InterpolateOptions{}
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	var res InterpolateResult
	missingSet := map[string]struct{}{}

	for key, val := range result {
		expanded := interpolateRe.ReplaceAllStringFunc(val, func(match string) string {
			ref := interpolateRe.FindStringSubmatch(match)[1]
			if ref == key {
				// avoid self-reference expansion
				return match
			}
			if replacement, ok := secrets[ref]; ok {
				return replacement
			}
			missingSet[ref] = struct{}{}
			return opts.Fallback
		})
		if expanded != val {
			res.Resolved = append(res.Resolved, key)
		}
		result[key] = expanded
	}

	for ref := range missingSet {
		res.Missing = append(res.Missing, ref)
	}

	if opts.Strict && len(res.Missing) > 0 {
		return nil, res, fmt.Errorf("interpolate: missing references: %s", strings.Join(res.Missing, ", "))
	}

	if opts.DryRun {
		return secrets, res, nil
	}
	return result, res, nil
}
