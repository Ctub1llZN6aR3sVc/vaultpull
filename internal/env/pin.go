package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PinEntry records a pinned version of a secret key.
type PinEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// PinFile holds all pinned entries.
type PinFile struct {
	Entries []PinEntry `json:"entries"`
}

// SavePins writes pinned entries to a JSON file.
func SavePins(path string, pins []PinEntry) error {
	pf := PinFile{Entries: pins}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadPins reads pinned entries from a JSON file.
func LoadPins(path string) ([]PinEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("pin: read: %w", err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("pin: unmarshal: %w", err)
	}
	return pf.Entries, nil
}

// ApplyPins overlays pinned values onto secrets, returning a new map.
func ApplyPins(secrets map[string]string, pins []PinEntry) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	now := time.Now()
	for _, p := range pins {
		if p.ExpiresAt.IsZero() || now.Before(p.ExpiresAt) {
			out[p.Key] = p.Value
		}
	}
	return out
}
