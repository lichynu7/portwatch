package ports

import "github.com/yourusername/portwatch/internal/config"

// BuildChangelog constructs a Changelog from the provided config.
// If the config is disabled, it returns nil so callers can skip recording.
func BuildChangelog(cfg config.ChangelogConfig) *Changelog {
	if !cfg.Enabled {
		return nil
	}
	return NewChangelog(cfg.MaxEvents)
}

// RecordDiff feeds the results of a port diff into a Changelog.
// added and removed are slices of Port values from snapshot.Diff.
// If cl is nil the call is a no-op, making it safe to call unconditionally.
func RecordDiff(cl *Changelog, added, removed []Port) {
	if cl == nil {
		return
	}
	for _, p := range added {
		cl.Record(ChangeAdded, p)
	}
	for _, p := range removed {
		cl.Record(ChangeRemoved, p)
	}
}
