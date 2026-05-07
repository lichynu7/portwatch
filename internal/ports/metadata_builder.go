package ports

import (
	"fmt"

	"github.com/user/portwatch/internal/config"
)

// BuildMetadataStore constructs a MetadataStore from the provided config.
// Returns nil when metadata collection is disabled.
func BuildMetadataStore(cfg config.MetadataConfig) *MetadataStore {
	if !cfg.Enabled {
		return nil
	}
	return NewMetadataStore(cfg.MaxKeysPerPort)
}

// AnnotatePort sets a metadata key-value pair on the store for a given Port.
// It is a no-op when store is nil.
func AnnotatePort(store *MetadataStore, p Port, key, value string) error {
	if store == nil {
		return nil
	}
	ok := store.Set(portKey(p), key, value)
	if !ok {
		return fmt.Errorf("metadata: key limit reached for port %s:%d", p.Addr, p.Port)
	}
	return nil
}

// GetAnnotation retrieves a metadata value for a Port.
// Returns empty string and false when store is nil or key is absent.
func GetAnnotation(store *MetadataStore, p Port, key string) (string, bool) {
	if store == nil {
		return "", false
	}
	e, ok := store.Get(portKey(p), key)
	if !ok {
		return "", false
	}
	return e.Value, true
}

// PurgePort removes all metadata for a Port from the store.
// It is a no-op when store is nil.
func PurgePort(store *MetadataStore, p Port) {
	if store == nil {
		return
	}
	store.Purge(portKey(p))
}
