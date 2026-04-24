package env

import "fmt"

// InvertOptions controls the behaviour of Invert.
type InvertOptions struct {
	// FailOnConflict returns an error when two keys would map to the same
	// inverted key (i.e. duplicate values in the source map).
	FailOnConflict bool
	// DryRun returns the result without mutating any external state.
	DryRun bool
}

// InvertResult holds the outcome of an Invert call.
type InvertResult struct {
	Inverted   map[string]string
	Conflicts  []string
	Skipped    int
}

// Summary returns a human-readable description of the invert operation.
func (r InvertResult) Summary() string {
	if len(r.Conflicts) > 0 {
		return fmt.Sprintf("inverted %d keys, %d conflict(s): %v",
			len(r.Inverted), len(r.Conflicts), r.Conflicts)
	}
	return fmt.Sprintf("inverted %d keys, %d skipped", len(r.Inverted), r.Skipped)
}

// Invert swaps the keys and values of secrets so that each value becomes a key
// and each key becomes the corresponding value. When two source keys share the
// same value a conflict is recorded; with FailOnConflict set the function
// returns an error instead of silently keeping the last writer.
func Invert(secrets map[string]string, opts *InvertOptions) (InvertResult, error) {
	if opts == nil {
		opts = &InvertOptions{}
	}

	out := make(map[string]string, len(secrets))
	var conflicts []string
	skipped := 0

	for k, v := range secrets {
		if v == "" {
			skipped++
			continue
		}
		if existing, seen := out[v]; seen {
			conflicts = append(conflicts, fmt.Sprintf("%s<->%s", existing, k))
			if opts.FailOnConflict {
				return InvertResult{}, fmt.Errorf("invert: duplicate value %q maps to both %q and %q", v, existing, k)
			}
		}
		out[v] = k
	}

	return InvertResult{
		Inverted:  out,
		Conflicts: conflicts,
		Skipped:   skipped,
	}, nil
}
