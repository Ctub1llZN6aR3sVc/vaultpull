package env

import "sort"

// IndexEntry holds metadata about a secret key.
type IndexEntry struct {
	Key    string
	Source string
	Tags   []string
}

// IndexResult holds the indexed entries and a summary.
type IndexResult struct {
	Entries []IndexEntry
	Total   int
}

// IndexOptions controls how the index is built.
type IndexOptions struct {
	Source string
	Tags   map[string][]string // key -> tags
}

// Index builds an ordered index of secrets with optional metadata.
func Index(secrets map[string]string, opts *IndexOptions) IndexResult {
	if opts == nil {
		opts = &IndexOptions{}
	}
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]IndexEntry, 0, len(keys))
	for _, k := range keys {
		if _, ok := secrets[k]; !ok {
			continue
		}
		e := IndexEntry{
			Key:    k,
			Source: opts.Source,
		}
		if opts.Tags != nil {
			e.Tags = opts.Tags[k]
		}
		entries = append(entries, e)
	}
	return IndexResult{Entries: entries, Total: len(entries)}
}

// IndexSummary returns a human-readable summary line.
func IndexSummary(r IndexResult) string {
	if r.Total == 0 {
		return "index: no entries"
	}
	return "index: " + itoa(r.Total) + " entries"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := []byte{}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
