package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastResult holds the outcome of a cast operation.
type CastResult struct {
	Key      string
	Original string
	Casted   string
	Type     string
}

// CastOptions controls how values are cast.
type CastOptions struct {
	// Types maps key names to desired types: "int", "bool", "float", "string"
	Types map[string]string
	Strict bool // return error on uncastable value
}

// Cast coerces secret values to the specified types.
func Cast(secrets map[string]string, opts CastOptions) (map[string]string, []CastResult, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var results []CastResult
	for key, typ := range opts.Types {
		v, ok := out[key]
		if !ok {
			continue
		}
		casted, err := castValue(v, typ)
		if err != nil {
			if opts.Strict {
				return nil, nil, fmt.Errorf("cast: key %q value %q cannot be cast to %s: %w", key, v, typ, err)
			}
			continue
		}
		results = append(results, CastResult{Key: key, Original: v, Casted: casted, Type: typ})
		out[key] = casted
	}
	return out, results, nil
}

func castValue(v, typ string) (string, error) {
	switch strings.ToLower(typ) {
	case "int":
		n, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(n)), nil
	case "bool":
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return "", err
		}
		return strconv.FormatBool(b), nil
	case "float":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case "string":
		return v, nil
	default:
		return "", fmt.Errorf("unknown type %q", typ)
	}
}
