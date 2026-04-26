package ports

// DeduplicateAlerts returns a FilterFunc that suppresses repeated alerts
// for the same port within the deduplicator's configured window.
//
// Usage:
//
//	filter := DeduplicateAlerts(dedup)
//	filtered := filter(ports)
func DeduplicateAlerts(d *Deduplicator) FilterFunc {
	return func(ports []Port) []Port {
		if d == nil {
			return ports
		}
		result := ports[:0]
		for _, p := range ports {
			key := portKey(p)
			if !d.IsDuplicate(key) {
				result = append(result, p)
			}
		}
		return result
	}
}

// portKey builds a stable string key for a Port used by deduplication.
// It reuses the same format as snapshot.portKey for consistency.
func portKey(p Port) string {
	return p.Protocol + ":" + p.Address
}
