package env

import "strings"

// OmitOptions controls which keys are removed from secrets.
type OmitOptions struct {
	Keys     []string
	Prefixes []string
	Empty    bool // remove keys with empty values
}

// OmitResult summarises what was removed.
type OmitResult struct {
	Removed []string
}

func (r OmitResult) Summary() string {
	if len(r.Removed) == 0 {
		return "omit: no keys removed"
	}
	return "omit: removed " + strings.Join(r.Removed, ", ")
}

// Omit returns a new map with matching keys removed.
func Omit(secrets map[string]string, opts OmitOptions) (map[string]string, OmitResult) {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	var removed []string

	for k, v := range secrets {
		if _, drop := keySet[k]; drop {
			removed = append(removed, k)
			continue
		}
		if opts.Empty && v == "" {
			removed = append(removed, k)
			continue
		}
		dropped := false
		for _, p := range opts.Prefixes {
			if strings.HasPrefix(k, p) {
				removed = append(removed, k)
				dropped = true
				break
			}
		}
		if !dropped {
			out[k] = v
		}
	}

	return out, OmitResult{Removed: removed}
}
