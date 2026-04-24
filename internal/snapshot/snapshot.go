package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Snapshot holds a point-in-time record of open ports.
type Snapshot struct {
	Timestamp time.Time    `json:"timestamp"`
	Ports     []ports.Port `json:"ports"`
}

// New creates a new Snapshot from the provided port list.
func New(ps []ports.Port) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ps,
	}
}

// Save writes the snapshot to a JSON file at the given path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from a JSON file at the given path.
// If the file does not exist it returns an empty snapshot and no error.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return &Snapshot{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// Diff compares a previous snapshot against a current one and returns
// ports that are new (added) and ports that have disappeared (removed).
func Diff(prev, curr *Snapshot) (added []ports.Port, removed []ports.Port) {
	prevSet := make(map[string]struct{}, len(prev.Ports))
	for _, p := range prev.Ports {
		prevSet[portKey(p)] = struct{}{}
	}

	currSet := make(map[string]struct{}, len(curr.Ports))
	for _, p := range curr.Ports {
		currSet[portKey(p)] = struct{}{}
	}

	for _, p := range curr.Ports {
		if _, ok := prevSet[portKey(p)]; !ok {
			added = append(added, p)
		}
	}

	for _, p := range prev.Ports {
		if _, ok := currSet[portKey(p)]; !ok {
			removed = append(removed, p)
		}
	}
	return
}

func portKey(p ports.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.LocalPort)
}
