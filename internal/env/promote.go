package env

import "fmt"

// PromoteOptions controls how secrets are promoted between profiles.
type PromoteOptions struct {
	Overwrite bool
	DryRun    bool
}

// PromoteResult summarises what changed during a promotion.
type PromoteResult struct {
	Added     []string
	Skipped   []string
	Overwrote []string
}

func (r PromoteResult) String() string {
	return fmt.Sprintf("added=%d skipped=%d overwrote=%d",
		len(r.Added), len(r.Skipped), len(r.Overwrote))
}

// Promote copies keys from src into dst according to opts.
// It returns the merged map and a result summary.
// When DryRun is true the returned map is a copy and dst is never mutated.
func Promote(dst, src map[string]string, opts PromoteOptions) (map[string]string, PromoteResult) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var res PromoteResult
	for k, v := range src {
		if existing, exists := out[k]; exists {
			if !opts.Overwrite {
				res.Skipped = append(res.Skipped, k)
				continue
			}
			_ = existing
			res.Overwrote = append(res.Overwrote, k)
		} else {
			res.Added = append(res.Added, k)
		}
		if !opts.DryRun {
			out[k] = v
		}
	}
	return out, res
}
