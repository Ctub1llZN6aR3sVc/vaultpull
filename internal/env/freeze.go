package env

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// FreezeEntry records a frozen secret key with its value hash and timestamp.
type FreezeEntry struct {
	Key       string    `json:"key"`
	ValueHash string    `json:"value_hash"`
	FrozenAt  time.Time `json:"frozen_at"`
}

// FreezeResult holds the outcome of a Freeze operation.
type FreezeResult struct {
	Frozen  []string
	Skipped []string
}

// Freeze locks the given keys so their values are recorded and can be
// detected as changed. It writes a freeze file at path.
func Freeze(secrets map[string]string, keys []string, path string) (FreezeResult, error) {
	existing := map[string]FreezeEntry{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &existing)
	}

	var result FreezeResult
	for _, k := range keys {
		v, ok := secrets[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		existing[k] = FreezeEntry{
			Key:       k,
			ValueHash: hashString(v),
			FrozenAt:  time.Now().UTC(),
		}
		result.Frozen = append(result.Frozen, k)
	}

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return result, fmt.Errorf("freeze: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return result, fmt.Errorf("freeze: write: %w", err)
	}
	return result, nil
}

// CheckFrozen returns keys whose current values differ from the frozen hash.
func CheckFrozen(secrets map[string]string, path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("freeze: read: %w", err)
	}
	entries := map[string]FreezeEntry{}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("freeze: unmarshal: %w", err)
	}
	var changed []string
	for k, e := range entries {
		if v, ok := secrets[k]; ok {
			if hashString(v) != e.ValueHash {
				changed = append(changed, k)
			}
		}
	}
	sort.Strings(changed)
	return changed, nil
}

func hashString(s string) string {
	return fmt.Sprintf("%x", func() []byte {
		h := 0
		for _, c := range s {
			h = h*31 + int(c)
		}
		return []byte(fmt.Sprintf("%d", h))
	}())
}
