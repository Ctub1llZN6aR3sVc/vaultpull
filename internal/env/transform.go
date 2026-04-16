package env

import (
	"strings"
)

// TransformOptions controls how secret keys are transformed before writing.
type TransformOptions struct {
	Uppercase  bool
	Prefix     string
	StripPrefix string
}

// Transform applies key transformations to a secrets map.
func Transform(secrets map[string]string, opts TransformOptions) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		key := k

		if opts.StripPrefix != "" && strings.HasPrefix(key, opts.StripPrefix) {
			key = strings.TrimPrefix(key, opts.StripPrefix)
		}

		if opts.Uppercase {
			key = strings.ToUpper(key)
		}

		if opts.Prefix != "" {
			key = opts.Prefix + key
		}

		if key != "" {
			out[key] = v
		}
	}
	return out
}
