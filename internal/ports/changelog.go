package ports

import (
	"sync"
	"time"
)

// ChangeType describes the kind of port change observed.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
)

// ChangeEvent records a single port appearing or disappearing.
type ChangeEvent struct {
	Port      Port
	Change    ChangeType
	Timestamp time.Time
}

// Changelog maintains a bounded, in-memory history of port change events.
type Changelog struct {
	mu      sync.RWMutex
	events  []ChangeEvent
	maxSize int
}

// NewChangelog creates a Changelog that retains at most maxSize events.
func NewChangelog(maxSize int) *Changelog {
	if maxSize <= 0 {
		maxSize = 256
	}
	return &Changelog{maxSize: maxSize}
}

// Record appends a ChangeEvent, evicting the oldest entry when the log is full.
func (c *Changelog) Record(ct ChangeType, p Port) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ev := ChangeEvent{Port: p, Change: ct, Timestamp: time.Now()}
	if len(c.events) >= c.maxSize {
		c.events = c.events[1:]
	}
	c.events = append(c.events, ev)
}

// Recent returns up to n most-recent events, newest first.
func (c *Changelog) Recent(n int) []ChangeEvent {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if n <= 0 || len(c.events) == 0 {
		return nil
	}
	start := len(c.events) - n
	if start < 0 {
		start = 0
	}
	slice := make([]ChangeEvent, len(c.events)-start)
	copy(slice, c.events[start:])
	// reverse so newest is first
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// Len returns the current number of stored events.
func (c *Changelog) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.events)
}

// Clear removes all stored events.
func (c *Changelog) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = c.events[:0]
}
