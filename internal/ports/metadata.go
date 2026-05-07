package ports

import (
	"sync"
	"time"
)

// MetadataEntry holds arbitrary key-value annotations attached to a port observation.
type MetadataEntry struct {
	Key       string
	Value     string
	UpdatedAt time.Time
}

// MetadataStore is a thread-safe store for per-port metadata annotations.
type MetadataStore struct {
	mu      sync.RWMutex
	entries map[string]map[string]MetadataEntry
	maxKeys int
}

// NewMetadataStore creates a MetadataStore with a maximum number of keys per port.
func NewMetadataStore(maxKeys int) *MetadataStore {
	if maxKeys <= 0 {
		maxKeys = 16
	}
	return &MetadataStore{
		entries: make(map[string]map[string]MetadataEntry),
		maxKeys: maxKeys,
	}
}

// Set stores a metadata key-value pair for the given port key.
// Returns false if the per-port key limit has been reached and the key is new.
func (m *MetadataStore) Set(portKey, key, value string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.entries[portKey]; !ok {
		m.entries[portKey] = make(map[string]MetadataEntry)
	}
	pk := m.entries[portKey]
	if _, exists := pk[key]; !exists && len(pk) >= m.maxKeys {
		return false
	}
	pk[key] = MetadataEntry{Key: key, Value: value, UpdatedAt: time.Now()}
	return true
}

// Get retrieves a metadata value for the given port key and metadata key.
func (m *MetadataStore) Get(portKey, key string) (MetadataEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pk, ok := m.entries[portKey]
	if !ok {
		return MetadataEntry{}, false
	}
	e, ok := pk[key]
	return e, ok
}

// All returns a copy of all metadata entries for the given port key.
func (m *MetadataStore) All(portKey string) []MetadataEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pk, ok := m.entries[portKey]
	if !ok {
		return nil
	}
	out := make([]MetadataEntry, 0, len(pk))
	for _, e := range pk {
		out = append(out, e)
	}
	return out
}

// Delete removes a single metadata key for the given port key.
func (m *MetadataStore) Delete(portKey, key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pk, ok := m.entries[portKey]; ok {
		delete(pk, key)
		if len(pk) == 0 {
			delete(m.entries, portKey)
		}
	}
}

// Purge removes all metadata for the given port key.
func (m *MetadataStore) Purge(portKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, portKey)
}
