package env

import (
	"fmt"
	"os"\n	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{filePath: filePath}
}

// Write writes the provided secrets map to the .env file.
// Existing file contents are overwritten.
func (w *Writer) Write(secrets map[string]string) error {
	var sb strings.Builder

	for key, value := range secrets {
		escaped := escapeValue(value)
		sb.WriteString(fmt.Sprintf("%s=%s\n", key, escaped))
	}

	return os.WriteFile(w.filePath, []byte(sb.String()), 0600)
}

// Merge writes secrets to the .env file, preserving existing keys
// that are not present in the provided secrets map.
func (w *Writer) Merge(secrets map[string]string) error {
	existing, err := Read(w.filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading existing env file: %w", err)
	}

	if existing == nil {
		existing = make(map[string]string)
	}

	for k, v := range secrets {
		existing[k] = v
	}

	return w.Write(existing)
}

// escapeValue wraps values containing spaces or special characters in quotes.
func escapeValue(value string) string {
	if strings.ContainsAny(value, " \t\n#") {
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(value, `"`, `\"`))
	}
	return value
}
