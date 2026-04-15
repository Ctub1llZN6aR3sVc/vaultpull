package env

import (
	"fmt"
	"time"
)

// ExpiryEntry records when a secret key was last synced and when it expires.
type ExpiryEntry struct {
	Key       string
	SyncedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired returns true if the entry's expiry time is before now.
func (e ExpiryEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// CheckExpiry inspects a map of secrets against a list of expiry entries and
// returns a slice of keys that have expired.
func CheckExpiry(secrets map[string]string, entries []ExpiryEntry) []string {
	expired := []string{}
	index := make(map[string]ExpiryEntry, len(entries))
	for _, e := range entries {
		index[e.Key] = e
	}
	for key := range secrets {
		if e, ok := index[key]; ok && e.IsExpired() {
			expired = append(expired, key)
		}
	}
	return expired
}

// BuildExpiryEntries creates ExpiryEntry records for each key in secrets,
// setting the expiry to syncedAt plus ttl.
func BuildExpiryEntries(secrets map[string]string, syncedAt time.Time, ttl time.Duration) []ExpiryEntry {
	if ttl <= 0 {
		return nil
	}
	entries := make([]ExpiryEntry, 0, len(secrets))
	for key := range secrets {
		entries = append(entries, ExpiryEntry{
			Key:       key,
			SyncedAt:  syncedAt,
			ExpiresAt: syncedAt.Add(ttl),
		})
	}
	return entries
}

// ExpiryWarnings returns human-readable warning strings for expired keys.
func ExpiryWarnings(expired []string) []string {
	warnings := make([]string, 0, len(expired))
	for _, key := range expired {
		warnings = append(warnings, fmt.Sprintf("secret %q has expired and should be re-synced", key))
	}
	return warnings
}
