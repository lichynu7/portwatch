package ports

import (
	"math"
	"sync"
	"time"
)

// DecayConfig controls how scores decay over time for known ports.
type DecayConfig struct {
	Enabled  bool
	HalfLife time.Duration // time after which a score is halved
	MinScore float64       // floor; entries below this are evicted
}

// DefaultDecayConfig returns sensible defaults.
func DefaultDecayConfig() DecayConfig {
	return DecayConfig{
		Enabled:  true,
		HalfLife: 10 * time.Minute,
		MinScore: 1.0,
	}
}

type decayEntry struct {
	score    float64
	lastSeen time.Time
}

// ScoreDecayer applies exponential decay to per-port scores.
type ScoreDecayer struct {
	cfg     DecayConfig
	mu      sync.Mutex
	entries map[string]*decayEntry
	now     func() time.Time
}

// NewScoreDecayer creates a ScoreDecayer with the given config.
func NewScoreDecayer(cfg DecayConfig) *ScoreDecayer {
	return &ScoreDecayer{
		cfg:     cfg,
		entries: make(map[string]*decayEntry),
		now:     time.Now,
	}
}

// Update records a new raw score for key, applies decay to the previous
// accumulated value, adds the new score, and returns the resulting score.
func (d *ScoreDecayer) Update(key string, rawScore float64) float64 {
	if !d.cfg.Enabled {
		return rawScore
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	e, ok := d.entries[key]
	if !ok {
		e = &decayEntry{}
		d.entries[key] = e
	}

	if !e.lastSeen.IsZero() && d.cfg.HalfLife > 0 {
		elapsed := now.Sub(e.lastSeen)
		halves := float64(elapsed) / float64(d.cfg.HalfLife)
		e.score = e.score / math.Pow(2, halves)
	}
	e.score += rawScore
	e.lastSeen = now
	return e.score
}

// Get returns the current decayed score for key without updating it.
func (d *ScoreDecayer) Get(key string) float64 {
	if !d.cfg.Enabled {
		return 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	e, ok := d.entries[key]
	if !ok {
		return 0
	}
	now := d.now()
	elapsed := now.Sub(e.lastSeen)
	halves := float64(elapsed) / float64(d.cfg.HalfLife)
	return e.score / math.Pow(2, halves)
}

// Purge removes entries whose decayed score has fallen below MinScore.
func (d *ScoreDecayer) Purge() {
	if !d.cfg.Enabled {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	for k, e := range d.entries {
		elapsed := now.Sub(e.lastSeen)
		halves := float64(elapsed) / float64(d.cfg.HalfLife)
		decayed := e.score / math.Pow(2, halves)
		if decayed < d.cfg.MinScore {
			delete(d.entries, k)
		}
	}
}
