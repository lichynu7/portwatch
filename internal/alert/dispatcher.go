package alert

import (
	"fmt"
	"sync"
)

// Dispatcher fans out alerts to one or more registered Notifiers.
type Dispatcher struct {
	mu        sync.RWMutex
	notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the provided notifiers.
func NewDispatcher(notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{
		notifiers: notifiers,
	}
}

// Register adds a new Notifier to the dispatcher.
func (d *Dispatcher) Register(n Notifier) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.notifiers = append(d.notifiers, n)
}

// Dispatch sends the alert to all registered notifiers.
// It collects all errors encountered and returns a combined error,
// but still attempts delivery to every notifier.
func (d *Dispatcher) Dispatch(a Alert) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var errs []error
	for _, n := range d.notifiers {
		if err := n.Notify(a); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	return fmt.Errorf("%d notifiers failed dispatching alert; first error: %w", len(errs), errs[0])
}
