package env

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchState holds the last known hash of an env file.
type WatchState struct {
	Path    string
	Hash    string
	CheckAt time.Time
}

// HashFile returns a SHA-256 hex digest of the file at path.
// Returns an empty string if the file does not exist.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("watch: open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("watch: hash %s: %w returns true when path differs from the recordedw *WatchState) Changed() (bool, error) {
	current, err := HashFile(w.Path)
	if err != nil {
		return false, err
	}
	return current != w.Hash, nil
}

// Refresh updates the stored hash to the current file state.
func (w *WatchState) Refresh() error {
	h, err := HashFile(w.Path)
	if err != nil {
		return err
	}
	w.Hash = h
	w.CheckAt = time.Now()
	return nil
}

// NewWatchState creates a WatchState initialised with the current hash.
func NewWatchState(path string) (*WatchState, error) {
	w := &WatchState{Path: path}
	if err := w.Refresh(); err != nil {
		return nil, err
	}
	return w, nil
}
