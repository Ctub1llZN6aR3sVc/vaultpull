package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Imported int
	Skipped  int
	Keys     []string
}

// ImportOptions controls how an import is performed.
type ImportOptions struct {
	Overwrite bool
	DryRun    bool
}

// ImportFromFile reads key=value pairs from src and merges them into dst map.
// Lines starting with '#' and blank lines are ignored.
func ImportFromFile(path string, dst map[string]string, opts ImportOptions) (ImportResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return ImportResult{}, fmt.Errorf("import: open %s: %w", path, err)
	}
	defer f.Close()

	incoming := map[string]string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := unquoteValue(strings.TrimSpace(parts[1]))
		incoming[key] = val
	}
	if err := scanner.Err(); err != nil {
		return ImportResult{}, fmt.Errorf("import: scan %s: %w", path, err)
	}

	result := ImportResult{}
	for k, v := range incoming {
		if _, exists := dst[k]; exists && !opts.Overwrite {
			result.Skipped++
			continue
		}
		if !opts.DryRun {
			dst[k] = v
		}
		result.Imported++
		result.Keys = append(result.Keys, k)
	}
	return result, nil
}
