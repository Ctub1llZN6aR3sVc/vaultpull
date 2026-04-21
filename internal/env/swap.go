package env

import "fmt"

// SwapOptions controls the behaviour of Swap.
type SwapOptions struct {
	// Pairs maps old key names to new key names.
	// The value from the old key is moved to the new key.
	Pairs map[string]string
	// FailOnMissing returns an error when a source key is absent.
	FailOnMissing bool
	// DryRun returns the result without modifying anything.
	DryRun bool
}

// SwapResult describes the outcome of a Swap call.
type SwapResult struct {
	Swapped []string
	Missing []string
}

// Summary returns a human-readable summary of the swap result.
func (r SwapResult) Summary() string {
	if len(r.Swapped) == 0 {
		return "swap: no keys swapped"
	}
	return fmt.Sprintf("swap: %d key(s) swapped, %d missing", len(r.Swapped), len(r.Missing))
}

// Swap renames keys in secrets according to opts.Pairs.
// The old key is removed and its value is written under the new key.
func Swap(secrets map[string]string, opts SwapOptions) (map[string]string, SwapResult, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var result SwapResult

	for oldKey, newKey := range opts.Pairs {
		val, ok := out[oldKey]
		if !ok {
			result.Missing = append(result.Missing, oldKey)
			if opts.FailOnMissing {
				return nil, result, fmt.Errorf("swap: key %q not found", oldKey)
			}
			continue
		}
		if !opts.DryRun {
			delete(out, oldKey)
			out[newKey] = val
		}
		result.Swapped = append(result.Swapped, oldKey)
	}

	return out, result, nil
}
