package env

import (
	"fmt"
	"io"
	"sort"
)

// PrintDiff writes a coloured/symbolic diff summary to w.
func PrintDiff(w io.Writer, result *DiffResult) {
	if result.IsEmpty() {
		fmt.Fprintln(w, "  (no changes)")
		return
	}

	// Sort for deterministic output.
	changes := make([]Change, len(result.Changes))
	copy(changes, result.Changes)
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	for _, c := range changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(w, "  + %s\n", c.Key)
		case Removed:
			fmt.Fprintf(w, "  - %s\n", c.Key)
		case Changed:
			fmt.Fprintf(w, "  ~ %s\n", c.Key)
		}
	}

	fmt.Fprintf(w, "\n  %s\n", result.Summary())
}
