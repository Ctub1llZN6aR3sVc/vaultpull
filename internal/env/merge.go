package env

// MergeStrategy defines how conflicts are resolved when merging secret maps.
type MergeStrategy int

const (
	// MergeStrategyOverwrite replaces existing values with incoming values.
	MergeStrategyOverwrite MergeStrategy = iota
	// MergeStrategyKeepExisting preserves existing values and only adds new keys.
	MergeStrategyKeepExisting
)

// MergeOptions configures the behaviour of Merge.
type MergeOptions struct {
	Strategy MergeStrategy
}

// Merge combines base and incoming secret maps according to opts.
// The returned map is always a new map; neither input is mutated.
func Merge(base, incoming map[string]string, opts MergeOptions) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}

	switch opts.Strategy {
	case MergeStrategyKeepExisting:
		for k, v := range incoming {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	default: // MergeStrategyOverwrite
		for k, v := range incoming {
			result[k] = v
		}
	}

	return result
}
