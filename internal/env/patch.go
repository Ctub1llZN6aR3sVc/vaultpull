package env

// PatchOptions controls how secrets are patched.
type PatchOptions struct {
	// Only patch keys that already exist in dst.
	ExistingOnly bool
	// DryRun returns the result without modifying dst.
	DryRun bool
}

// PatchResult summarises what changed.
type PatchResult struct {
	Patched []string
	Skipped []string
}

func (r PatchResult) Summary() string {
	if len(r.Patched) == 0 {
		return "patch: no keys updated"
	}
	return fmt.Sprintf("patch: %d key(s) updated, %d skipped", len(r.Patched), len(r.Skipped))
}

// Patch applies updates from patch into dst according to opts.
// It returns a new map and a PatchResult; dst is never mutated.
func Patch(dst, patch map[string]string, opts PatchOptions) (map[string]string, PatchResult) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var result PatchResult
	for k, v := range patch {
		if opts.ExistingOnly {
			if _, exists := dst[k]; !exists {
				result.Skipped = append(result.Skipped, k)
				continue
			}
		}
		if !opts.DryRun {
			out[k] = v
		}
		result.Patched = append(result.Patched, k)
	}
	sort.Strings(result.Patched)
	sort.Strings(result.Skipped)
	return out, result
}
