package ports

import "github.com/user/portwatch/internal/config"

// SeverityFilter returns a FilterFunc that drops ports below the minimum severity.
func SeverityFilter(cfg config.SeverityFilterConfig) FilterFunc {
	return func(p Port) bool {
		return p.Severity >= cfg.MinSeverity
	}
}

// ApplySeverityFilter removes ports whose severity is below the configured minimum.
func ApplySeverityFilter(ports []Port, cfg config.SeverityFilterConfig) []Port {
	f := SeverityFilter(cfg)
	out := ports[:0:0]
	for _, p := range ports {
		if f(p) {
			out = append(out, p)
		}
	}
	return out
}
