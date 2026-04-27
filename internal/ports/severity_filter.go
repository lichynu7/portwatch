package ports

import "github.com/user/portwatch/internal/config"

// SeverityFilter returns a FilterFunc that excludes ports below the minimum
// severity level defined in cfg. Ports at or above the threshold are kept.
func SeverityFilter(cfg config.SeverityConfig, sc DefaultSeverityConfig) FilterFunc {
	return func(p Port) bool {
		level := Classify(p, sc)
		return level >= cfg.MinLevel
	}
}

// ApplySeverityFilter filters a slice of ports, retaining only those whose
// classified severity meets or exceeds the configured minimum level.
func ApplySeverityFilter(ports []Port, cfg config.SeverityConfig, sc DefaultSeverityConfig) []Port {
	f := SeverityFilter(cfg, sc)
	out := make([]Port, 0, len(ports))
	for _, p := range ports {
		if f(p) {
			out = append(out, p)
		}
	}
	return out
}
