package env

import "sort"

// CompareResult holds the result of comparing two secret maps.
type CompareResult struct {
	OnlyInA  []string
	OnlyInB  []string
	Differ   []string
	Identical []string
}

// IsClean returns true when both maps are identical.
func (r *CompareResult) IsClean() bool {
	return len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0 && len(r.Differ) == 0
}

// Compare compares two secret maps and categorises every key.
func Compare(a, b map[string]string) *CompareResult {
	r := &CompareResult{}

	keys := make(map[string]struct{})
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	for k := range keys {
		aVal, inA := a[k]
		bVal, inB := b[k]
		switch {
		case inA && !inB:
			r.OnlyInA = append(r.OnlyInA, k)
		case !inA && inB:
			r.OnlyInB = append(r.OnlyInB, k)
		case aVal == bVal:
			r.Identical = append(r.Identical, k)
		default:
			r.Differ = append(r.Differ, k)
		}
	}

	sort.Strings(r.OnlyInA)
	sort.Strings(r.OnlyInB)
	sort.Strings(r.Differ)
	sort.Strings(r.Identical)
	return r
}
