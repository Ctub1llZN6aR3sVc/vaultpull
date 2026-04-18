package env

import "sort"

// DedupeResult holds the outcome of a deduplication pass.
type DedupeResult struct {
	Removed []string // keys whose duplicate entries were collapsed
}

func (r DedupeResult) Summary() string {
	if len(r.Removed) == 0 {
		return "no duplicate keys found"
	}
	return fmt.Sprintf("%d duplicate key(s) removed: %s", len(r.Removed), joinKeys(r.Removed))
}

func joinKeys(keys []string) string {
	s := ""
	for i, k := range keys {
		if i > 0 {
			s += ", "
		}
		s += k
	}
	return s
}

// Dedupe returns a new map equal to secrets but records which keys appeared
// more than once in the ordered pairs slice. Since Go maps are already
// deduplicated by key, this variant operates on a []Pair representation.
//
// For the common map[string]string path it simply detects keys present in
// both `base` and `incoming` and, depending on keepFirst, retains the base
// value or the incoming value.
func Dedupe(base, incoming map[string]string, keepFirst bool) (map[string]string, DedupeResult) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	var dupes []string
	for k, v := range incoming {
		if _, exists := out[k]; exists {
			dupes = append(dupes, k)
			if !keepFirst {
				out[k] = v
			}
		} else {
			out[k] = v
		}
	}
	sort.Strings(dupes)
	return out, DedupeResult{Removed: dupes}
}
