package env

// Chain resolves secrets by trying multiple source maps in order,
// returning the first map that contains a non-empty value for each key.
// Keys found in earlier sources take priority over later ones.
func Chain(keys []string, sources ...map[string]string) map[string]string {
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		for _, src := range sources {
			if val, ok := src[key]; ok && val != "" {
				result[key] = val
				break
			}
		}
	}
	return result
}

// ChainAll resolves all keys present in any source map.
// Earlier sources take priority.
func ChainAll(sources ...map[string]string) map[string]string {
	// Collect all keys
	seen := make(map[string]struct{})
	for _, src := range sources {
		for k := range src {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return Chain(keys, sources...)
}
