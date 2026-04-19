package env

import (
	"fmt"
	"sort"
)

// KeyDiff represents the result of comparing key sets between two secret maps.
type KeyDiff struct {
	OnlyInA []string
	OnlyInB []string
	InBoth  []string
}

// DiffKeys compares the key sets of two secret maps and returns which keys
// appear only in a, only in b, or in both.
func DiffKeys(a, b map[string]string) KeyDiff {
	aSet := toStringSet(a)
	bSet := toStringSet(b)

	var onlyInA, onlyInB, inBoth []string

	for k := range aSet {
		if bSet[k] {
			inBoth = append(inBoth, k)
		} else {
			onlyInA = append(onlyInA, k)
		}
	}
	for k := range bSet {
		if !aSet[k] {
			onlyInB = append(onlyInB, k)
		}
	}

	sort.Strings(onlyInA)
	sort.Strings(onlyInB)
	sort.Strings(inBoth)

	return KeyDiff{
		OnlyInA: onlyInA,
		OnlyInB: onlyInB,
		InBoth:  inBoth,
	}
}

func toStringSet(m map[string]string) map[string]bool {
	s := make(map[string]bool, len(m))
	for k := range m {
		s[k] = true
	}
	return s
}

// KeyDiffSummary returns a human-readable summary of the key diff.
func KeyDiffSummary(d KeyDiff) string {
	if len(d.OnlyInA) == 0 && len(d.OnlyInB) == 0 {
		return "key sets are identical"
	}
	return fmt.Sprintf("%d only in A, %d only in B, %d shared",
		len(d.OnlyInA), len(d.OnlyInB), len(d.InBoth))
}
