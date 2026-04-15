package env

import "strings"

// SensitiveKeyPatterns holds substrings that indicate a secret value.
var SensitiveKeyPatterns = []string{
	"PASSWORD",
	"SECRET",
	"TOKEN",
	"KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
	"API",
	"DSN",
	"CERT",
}

// IsSensitive returns true when the key name matches a known sensitive pattern.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range SensitiveKeyPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskValue replaces a sensitive value with a redacted placeholder.
// It reveals up to the first 4 characters when the value is long enough.
func MaskValue(value string) string {
	const placeholder = "****"
	if len(value) <= 4 {
		return placeholder
	}
	return value[:4] + placeholder
}

// MaskSecrets returns a copy of the provided map with sensitive values masked.
func MaskSecrets(secrets map[string]string) map[string]string {
	masked := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if IsSensitive(k) {
			masked[k] = MaskValue(v)
		} else {
			masked[k] = v
		}
	}
	return masked
}
