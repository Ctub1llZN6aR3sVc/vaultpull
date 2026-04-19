package env

import "fmt"

// PickOptions controls how Pick behaves.
type PickOptions struct {
	// Keys is the explicit list of keys to retain.
	Keys []string
	// FailOnMissing causes Pick to return an error if any requested key is absent.
	FailOnMissing bool
	// DryRun returns the result without treating it as authoritative.
	DryRun bool
}

// PickResult holds the output of a Pick operation.
type PickResult struct {
	Picked  map[string]string
	Missing []string
}

// Summary returns a human-readable summary of the pick operation.
func (r PickResult) Summary() string {
	if len(r.Missing) == 0 {
		return fmt.Sprintf("picked %d key(s)", len(r.Picked))
	}
	return fmt.Sprintf("picked %d key(s), %d missing: %v", len(r.Picked), len(r.Missing), r.Missing)
}

// Pick returns a new map containing only the keys specified in opts.Keys.
// If FailOnMissing is set and a key is absent, an error is returned.
func Pick(secrets map[string]string, opts PickOptions) (PickResult, error) {
	out := make(map[string]string, len(opts.Keys))
	var missing []string

	for _, k := range opts.Keys {
		v, ok := secrets[k]
		if !ok {
			missing = append(missing, k)
			continue
		}
		out[k] = v
	}

	if opts.FailOnMissing && len(missing) > 0 {
		return PickResult{}, fmt.Errorf("pick: missing required keys: %v", missing)
	}

	return PickResult{Picked: out, Missing: missing}, nil
}
