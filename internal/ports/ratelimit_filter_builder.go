package ports

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/config"
)

// BuildRateLimitFilter constructs an AlertRateLimitFilter from the application
// config, returning nil when the stage is disabled.
func BuildRateLimitFilter(cfg config.RateLimitFilterConfig) (*AlertRateLimitFilter, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	d, err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("ratelimit_filter: %w", err)
	}

	rlCfg := RateLimitFilterConfig{
		Window:  d,
		MaxHits: cfg.MaxHits,
	}
	return NewAlertRateLimitFilter(rlCfg), nil
}

// RateLimitStage returns a pipeline-compatible transform function that applies
// the rate-limit filter to each batch of ports. If filter is nil the batch is
// passed through unchanged.
func RateLimitStage(filter *AlertRateLimitFilter) func([]Port) []Port {
	if filter == nil {
		return func(ports []Port) []Port { return ports }
	}
	return func(ports []Port) []Port {
		return ApplyRateLimitFilter(ports, filter, time.Now())
	}
}
