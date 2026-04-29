package ports

import (
	"sync"
	"time"
)

// AggregatorConfig controls how alerts are batched before dispatch.
type AggregatorConfig struct {
	// Window is how long to collect alerts before flushing.
	Window time.Duration
	// MaxBatch is the maximum number of alerts per flush.
	MaxBatch int
}

// DefaultAggregatorConfig returns sensible defaults.
func DefaultAggregatorConfig() AggregatorConfig {
	return AggregatorConfig{
		Window:   5 * time.Second,
		MaxBatch: 20,
	}
}

// Aggregator batches Port alerts within a time window and flushes them as
// a slice to the registered flush function.
type Aggregator struct {
	cfg     AggregatorConfig
	mu      sync.Mutex
	buf     []Port
	flushFn func([]Port)
}

// NewAggregator creates an Aggregator with the given config and flush callback.
func NewAggregator(cfg AggregatorConfig, flushFn func([]Port)) *Aggregator {
	return &Aggregator{
		cfg:     cfg,
		flushFn: flushFn,
	}
}

// Add appends a port to the current batch. If the batch reaches MaxBatch the
// buffer is flushed immediately.
func (a *Aggregator) Add(p Port) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.buf = append(a.buf, p)
	if len(a.buf) >= a.cfg.MaxBatch {
		a.locked_flush()
	}
}

// Flush drains the buffer and calls flushFn with any accumulated ports.
func (a *Aggregator) Flush() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.locked_flush()
}

// locked_flush must be called with a.mu held.
func (a *Aggregator) locked_flush() {
	if len(a.buf) == 0 {
		return
	}
	batch := make([]Port, len(a.buf))
	copy(batch, a.buf)
	a.buf = a.buf[:0]
	go a.flushFn(batch)
}

// Run starts the periodic flush loop. It blocks until ctx is cancelled.
func (a *Aggregator) Run(ctx interface{ Done() <-chan struct{} }) {
	ticker := time.NewTicker(a.cfg.Window)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.Flush()
		case <-ctx.Done():
			a.Flush()
			return
		}
	}
}
