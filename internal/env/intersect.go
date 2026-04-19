package env

// IntersectOptions controls how Intersect behaves.
type IntersectOptions struct {
	// KeepValue selects which map's value to use when a key exists in both.
	// "a" (default) keeps the value from a, "b" keeps the value from b.
	KeepValue string
}

// IntersectResult holds the result of an Intersect operation.
type IntersectResult struct {
	Secrets map[string]string
	Kept    []string
}

// Summary returns a human-readable summary.
func (r IntersectResult) Summary() string {
	if len(r.Kept) == 0 {
		return "intersect: no common keys found"
	}
	return "intersect: " + itoa(len(r.Kept)) + " common key(s) retained"
}

// Intersect returns only the keys present in both a and b.
// The value is taken from a unless opts.KeepValue == "b".
func Intersect(a, b map[string]string, opts IntersectOptions) IntersectResult {
	out := make(map[string]string)
	var kept []string

	for k := range a {
		if _, ok := b[k]; ok {
			if opts.KeepValue == "b" {
				out[k] = b[k]
			} else {
				out[k] = a[k]
			}
			kept = append(kept, k)
		}
	}

	sortStrings(kept)
	return IntersectResult{Secrets: out, Kept: kept}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
