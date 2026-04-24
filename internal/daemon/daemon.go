// Package daemon implements the main watch loop for portwatch.
package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon watches open ports at a configured interval and dispatches alerts
// when unexpected listeners appear or disappear.
type Daemon struct {
	cfg        *config.Config
	scanner    *ports.Scanner
	dispatcher *alert.Dispatcher
}

// New creates a Daemon with the provided configuration and dispatcher.
func New(cfg *config.Config, dispatcher *alert.Dispatcher) *Daemon {
	return &Daemon{
		cfg:        cfg,
		scanner:    ports.NewScanner(),
		dispatcher: dispatcher,
	}
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	log.Printf("portwatch daemon starting (interval=%s, snapshot=%s)",
		d.cfg.Interval, d.cfg.SnapshotPath)

	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch daemon stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := d.tick(); err != nil {
				log.Printf("tick error: %v", err)
			}
		}
	}
}

// tick performs a single scan-diff-alert cycle.
func (d *Daemon) tick() error {
	current, err := d.scanner.Scan()
	if err != nil {
		return err
	}

	filtered := ports.FilterByPort(current, d.cfg.IgnorePorts)

	prev, err := snapshot.Load(d.cfg.SnapshotPath)
	if err != nil {
		// No previous snapshot — save baseline and return.
		snap := snapshot.New(filtered)
		return snap.Save(d.cfg.SnapshotPath)
	}

	added, removed := snapshot.Diff(prev, snapshot.New(filtered))

	for _, p := range added {
		a := alert.NewAlert(alert.KindAppeared, p)
		if err := d.dispatcher.Deliver(ctx, a); err != nil {
			log.Printf("alert delivery error: %v", err)
		}
	}
	for _, p := range removed {
		a := alert.NewAlert(alert.KindVanished, p)
		if err := d.dispatcher.Deliver(ctx, a); err != nil {
			log.Printf("alert delivery error: %v", err)
		}
	}

	snap := snapshot.New(filtered)
	return snap.Save(d.cfg.SnapshotPath)
}
