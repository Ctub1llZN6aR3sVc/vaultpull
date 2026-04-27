package env

import "fmt"

// WrapOptions controls how values are wrapped.
type WrapOptions struct {
	// Prefix is prepended to each value.
	Prefix string
	// Suffix is appended to each value.
	Suffix string
	// Keys restricts wrapping to specific keys. If empty, all keys are wrapped.
	Keys []string
	// DryRun returns a modified copy without mutating the input.
	DryRun bool
}

// WrapResult holds the result of a Wrap operation.
type WrapResult struct {
	Secrets  map[string]string
	Wrapped  []string
	Skipped  []string
}

// Summary returns a human-readable summary of the wrap operation.
func (r WrapResult) Summary() string {
	if len(r.Wrapped) == 0 {
		return "wrap: no keys wrapped"
	}
	return fmt.Sprintf("wrap: %d key(s) wrapped, %d skipped", len(r.Wrapped), len(r.Skipped))
}

// Wrap prepends and/or appends strings to secret values.
func Wrap(secrets map[string]string, opts *WrapOptions) WrapResult {
	if opts == nil {
		return WrapResult{Secrets: copyWrapMap(secrets)}
	}

	targetSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		targetSet[k] = true
	}

	out := copyWrapMap(secrets)
	var wrapped, skipped []string

	for k, v := range secrets {
		if len(targetSet) > 0 && !targetSet[k] {
			skipped = append(skipped, k)
			continue
		}
		newVal := opts.Prefix + v + opts.Suffix
		if !opts.DryRun {
			out[k] = newVal
		}
		wrapped = append(wrapped, k)
	}

	return WrapResult{
		Secrets: out,
		Wrapped: wrapped,
		Skipped: skipped,
	}
}

func copyWrapMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
