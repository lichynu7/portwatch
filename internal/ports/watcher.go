package ports

import (
	"context"
	"fmt"
	"time"
)

// WatchConfig holds configuration for the port watcher.
type WatchConfig struct {
	Interval  time.Duration
	ProcFS    string
	Baseline  *Baseline
	Throttle  *Throttle
	Filter    *FilterConfig
}

// PortEvent represents a detected change in open ports.
type PortEvent struct {
	Port    Port
	EventType string // "new" or "closed"
}

// Watcher continuously scans for port changes and emits events.
type Watcher struct {
	cfg     WatchConfig
	scanner *Scanner
	events  chan PortEvent
}

// NewWatcher creates a Watcher using the provided configuration.
func NewWatcher(cfg WatchConfig) (*Watcher, error) {
	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("watcher: interval must be positive")
	}
	procFS := cfg.ProcFS
	if procFS == "" {
		procFS = "/proc"
	}
	return &Watcher{
		cfg:     cfg,
		scanner: NewScanner(procFS),
		events:  make(chan PortEvent, 64),
	}, nil
}

// Events returns the read-only channel of port events.
func (w *Watcher) Events() <-chan PortEvent {
	return w.events
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	defer close(w.events)

	prev, err := w.scan()
	if err != nil {
		return fmt.Errorf("watcher: initial scan: %w", err)
	}

	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			curr, err := w.scan()
			if err != nil {
				continue
			}
			w.diff(prev, curr)
			prev = curr
		}
	}
}

func (w *Watcher) scan() ([]Port, error) {
	ports, err := w.scanner.Scan()
	if err != nil {
		return nil, err
	}
	if w.cfg.Filter != nil {
		ports = ExcludeSafe(ports, w.cfg.Filter)
	}
	if w.cfg.Baseline != nil {
		ports = ApplyBaseline(ports, w.cfg.Baseline)
	}
	return ports, nil
}

func (w *Watcher) diff(prev, curr []Port) {
	prevMap := indexPorts(prev)
	currMap := indexPorts(curr)

	for key, p := range currMap {
		if _, ok := prevMap[key]; !ok {
			w.emit(PortEvent{Port: p, EventType: "new"})
		}
	}
	for key, p := range prevMap {
		if _, ok := currMap[key]; !ok {
			w.emit(PortEvent{Port: p, EventType: "closed"})
		}
	}
}

func (w *Watcher) emit(e PortEvent) {
	if w.cfg.Throttle != nil {
		key := fmt.Sprintf("%s:%d", e.EventType, e.Port.Port)
		if !w.cfg.Throttle.Allow(key) {
			return
		}
	}
	select {
	case w.events <- e:
	default:
	}
}

func indexPorts(ports []Port) map[string]Port {
	m := make(map[string]Port, len(ports))
	for _, p := range ports {
		key := fmt.Sprintf("%s:%d", p.Proto, p.Port)
		m[key] = p
	}
	return m
}
