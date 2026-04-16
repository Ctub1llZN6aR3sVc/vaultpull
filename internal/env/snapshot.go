package env

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of secrets.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Profile   string            `json:"profile"`
	Secrets   map[string]string `json:"secrets"`
}

// SaveSnapshot writes a snapshot of secrets to a JSON file.
func SaveSnapshot(path, profile string, secrets map[string]string) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Profile:   profile,
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot unmarshal: %w", err)
	}
	return &snap, nil
}

// DiffSnapshot returns a Diff between a snapshot and current secrets.
func DiffSnapshot(snap *Snapshot, current map[string]string) DiffResult {
	if snap == nil {
		return Diff(map[string]string{}, current)
	}
	return Diff(snap.Secrets, current)
}
