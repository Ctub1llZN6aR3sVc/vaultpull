package env

import "sort"

// UnionOptions controls how Union merges two secret maps.
type UnionOptions struct {
	// PreferA means keys present in both maps take their value from A.
	// When false, B's value wins.
	PreferA bool
}

// UnionResult holds the result of a Union operation.
type UnionResult struct {
	Secrets  map[string]string
	OnlyInA  []string
	OnlyInB  []string
	InBoth   []string
}

// Summary returns a human-readable description of the union.
func (r UnionResult) Summary() string {
	if len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 {
		return "union: all keys shared between both sources"
	}
	s := "union:"
	if len(r.OnlyInA) > 0 {
		s += " " + itoa(len(r.OnlyInA)) + " only in A"
	}
	if len(r.OnlyInB) > 0 {
		s += " " + itoa(len(r.OnlyInB)) + " only in B"
	}
	if len(r.InBoth) > 0 {
		s += " " + itoa(len(r.InBoth)) + " shared"
	}
	return s
}

// Union merges maps a and b, returning all keys from both.
// When a key exists in both maps, opts.PreferA controls which value is kept.
func Union(a, b map[string]string, opts UnionOptions) UnionResult {
	out := make(map[string]string, len(a)+len(b))
	var onlyInA, onlyInB, inBoth []string

	for k, v := range a {
		if bv, ok := b[k]; ok {
			inBoth = append(inBoth, k)
			if opts.PreferA {
				out[k] = v
			} else {
				out[k] = bv
			}
		} else {
			onlyInA = append(onlyInA, k)
			out[k] = v
		}
	}

	for k, v := range b {
		if _, ok := a[k]; !ok {
			onlyInB = append(onlyInB, k)
			out[k] = v
		}
	}

	sort.Strings(onlyInA)
	sort.Strings(onlyInB)
	sort.Strings(inBoth)

	return UnionResult{
		Secrets: out,
		OnlyInA: onlyInA,
		OnlyInB: onlyInB,
		InBoth:  inBoth,
	}
}
