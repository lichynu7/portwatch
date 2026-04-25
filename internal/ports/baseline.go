package ports

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Baseline represents a set of ports considered "known good" for a host.
// It is persisted to disk and used to suppress alerts for expected listeners.
type Baseline struct {
	Entries map[string]BaselineEntry `json:"entries"`
	path    string
}

// BaselineEntry records why a port was added to the baseline.
type BaselineEntry struct {
	Port    uint16 `json:"port"`
	Proto   string `json:"proto"`
	Process string `json:"process"`
	Reason  string `json:"reason"`
}

// NewBaseline creates an empty baseline backed by the given file path.
func NewBaseline(path string) *Baseline {
	return &Baseline{
		Entries: make(map[string]BaselineEntry),
		path:    path,
	}
}

// LoadBaseline reads an existing baseline from disk, or returns an empty one
// if the file does not exist.
func LoadBaseline(path string) (*Baseline, error) {
	b := NewBaseline(path)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, b); err != nil {
		return nil, fmt.Errorf("baseline: parse %s: %w", path, err)
	}
	b.path = path
	return b, nil
}

// Add inserts a port into the baseline and persists the file.
func (b *Baseline) Add(entry BaselineEntry) error {
	key := baselineKey(entry.Proto, entry.Port)
	b.Entries[key] = entry
	return b.save()
}

// Contains reports whether the given proto/port pair is in the baseline.
func (b *Baseline) Contains(proto string, port uint16) bool {
	_, ok := b.Entries[baselineKey(proto, port)]
	return ok
}

// Remove deletes a port from the baseline and persists the file.
func (b *Baseline) Remove(proto string, port uint16) error {
	delete(b.Entries, baselineKey(proto, port))
	return b.save()
}

func (b *Baseline) save() error {
	if err := os.MkdirAll(filepath.Dir(b.path), 0o750); err != nil {
		return fmt.Errorf("baseline: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(b.path, data, 0o640)
}

func baselineKey(proto string, port uint16) string {
	return fmt.Sprintf("%s:%d", proto, port)
}
