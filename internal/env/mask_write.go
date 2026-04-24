package env

import (
	"fmt"
	"sort"
	"strings"
)

// MaskWriteOptions controls how secrets are masked before writing.
type MaskWriteOptions struct {
	// AutoDetect masks keys that match sensitive patterns.
	AutoDetect bool
	// Keys is an explicit list of keys to mask.
	Keys []string
	// Placeholder replaces the masked value. Defaults to "***".
	Placeholder string
	// DryRun returns the masked map without modifying the original.
	DryRun bool
}

// MaskWriteResult holds the result of a MaskWrite operation.
type MaskWriteResult struct {
	Masked []string
}

func (r MaskWriteResult) Summary() string {
	if len(r.Masked) == 0 {
		return "mask_write: no keys masked"
	}
	sort.Strings(r.Masked)
	return fmt.Sprintf("mask_write: masked %d key(s): %s", len(r.Masked), strings.Join(r.Masked, ", "))
}

// MaskWrite returns a copy of secrets with sensitive values replaced by a placeholder.
// The original map is never mutated.
func MaskWrite(secrets map[string]string, opts MaskWriteOptions) (map[string]string, MaskWriteResult) {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	explicit := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = true
	}

	out := make(map[string]string, len(secrets))
	var result MaskWriteResult

	for k, v := range secrets {
		if explicit[k] || (opts.AutoDetect && IsSensitive(k)) {
			out[k] = placeholder
			result.Masked = append(result.Masked, k)
		} else {
			out[k] = v
		}
	}

	return out, result
}
