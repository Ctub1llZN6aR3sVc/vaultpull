package env

import (
	"bufio"
	"os"
	"strings"
)

// Read parses an .env file and returns a map of key-value pairs.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
func Read(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
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
		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// unquoteValue strips surrounding single or double quotes from a value.
func unquoteValue(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
