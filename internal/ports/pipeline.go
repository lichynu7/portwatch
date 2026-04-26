// Package ports provides port scanning, enrichment, filtering, and alerting
// primitives for the portwatch daemon.
//
// Pipeline wires together the scanning, enrichment, deduplication, throttling,
// and baseline-filtering stages into a single reusable processing chain.
package ports

import (
	"context"
	"fmt"
	"log"

	"github.com/example/portwatch/internal/config"
)

// Pipeline holds the ordered set of processing stages applied to every batch
// of raw ports returned by the Watcher. Each stage may transform, enrich, or
// drop individual Port entries before they are forwarded to alert handlers.
type Pipeline struct {
	watcher  *Watcher
	enricher *Enricher
	baseline *Baseline
	dedup    *Deduplicator
	throttle *Throttle
	filter   *FilterConfig
	out      chan []Port
}

// PipelineConfig groups all sub-configs required to construct a Pipeline.
type PipelineConfig struct {
	Watcher  config.WatcherConfig
	Enricher config.EnricherConfig
	Baseline config.BaselineConfig
	Dedup    config.DedupConfig
	Throttle config.ThrottleConfig
	Filter   *FilterConfig // optional; nil disables safe-port filtering
}

// NewPipeline constructs a Pipeline from the provided configuration.
// It initialises each stage and returns an error if any stage fails to
// initialise (e.g. invalid intervals, missing baseline file).
func NewPipeline(cfg PipelineConfig) (*Pipeline, error) {
	w, err := NewWatcher(cfg.Watcher)
	if err != nil {
		return nil, fmt.Errorf("pipeline: watcher: %w", err)
	}

	var opts []EnricherOption
	if cfg.Enricher.ResolveServices {
		opts = append(opts, WithServiceResolution())
	}
	if cfg.Enricher.LookupProcess {
		opts = append(opts, WithProcessLookup())
	}
	e := NewEnricher(opts...)

	bl, err := LoadBaseline(cfg.Baseline.Path)
	if err != nil {
		// A missing baseline is not fatal; start with an empty one.
		log.Printf("pipeline: baseline not loaded (%v); starting empty", err)
		bl = NewBaseline()
	}

	dd, err := NewDeduplicator(cfg.Dedup)
	if err != nil {
		return nil, fmt.Errorf("pipeline: deduplicator: %w", err)
	}

	th, err := NewThrottle(cfg.Throttle)
	if err != nil {
		return nil, fmt.Errorf("pipeline: throttle: %w", err)
	}

	return &Pipeline{
		watcher:  w,
		enricher: e,
		baseline: bl,
		dedup:    dd,
		throttle: th,
		filter:   cfg.Filter,
		out:      make(chan []Port, cfg.Watcher.BufferSize),
	}, nil
}

// Out returns the read-only channel on which processed port batches are
// delivered. Consumers should range over this channel until it is closed.
func (p *Pipeline) Out() <-chan []Port {
	return p.out
}

// Run starts the watcher and begins processing port batches through the
// pipeline stages. It blocks until ctx is cancelled, then closes Out().
func (p *Pipeline) Run(ctx context.Context) error {
	defer close(p.out)

	watchCh := p.watcher.Ports()

	go func() {
		if err := p.watcher.Run(ctx); err != nil {
			log.Printf("pipeline: watcher exited: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case batch, ok := <-watchCh:
			if !ok {
				return nil
			}
			processed := p.process(batch)
			if len(processed) == 0 {
				continue
			}
			select {
			case p.out <- processed:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// process applies each pipeline stage in order to a raw batch of ports and
// returns the surviving entries ready for alerting.
func (p *Pipeline) process(batch []Port) []Port {
	// Stage 1: enrich with service names and process info.
	batch = p.enricher.Enrich(batch)

	// Stage 2: remove ports that are in the known-good baseline.
	batch = ApplyBaseline(p.baseline, batch)

	// Stage 3: remove ports that are considered safe by the filter rules.
	if p.filter != nil {
		batch = ExcludeSafe(p.filter, batch)
	}

	// Stage 4: deduplicate — suppress ports already seen within the window.
	batch = DeduplicateAlerts(p.dedup, batch)

	// Stage 5: throttle — enforce per-port burst limits.
	batch = p.throttleFilter(batch)

	return batch
}

// throttleFilter applies the throttle stage, returning only ports that pass.
func (p *Pipeline) throttleFilter(batch []Port) []Port {
	out := batch[:0]
	for _, port := range batch {
		key := portKey(port)
		if p.throttle.Allow(key) {
			out = append(out, port)
		}
	}
	return out
}
