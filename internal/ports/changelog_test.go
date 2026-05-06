package ports

import (
	"testing"
	"time"
)

func makeChangePort(port uint16) Port {
	return Port{LocalPort: port, Protocol: "tcp"}
}

func TestChangelogRecordAndLen(t *testing.T) {
	cl := NewChangelog(10)
	if cl.Len() != 0 {
		t.Fatalf("expected 0 events, got %d", cl.Len())
	}
	cl.Record(ChangeAdded, makeChangePort(8080))
	if cl.Len() != 1 {
		t.Fatalf("expected 1 event, got %d", cl.Len())
	}
}

func TestChangelogEvictsOldest(t *testing.T) {
	cl := NewChangelog(3)
	for i := uint16(1); i <= 5; i++ {
		cl.Record(ChangeAdded, makeChangePort(i))
	}
	if cl.Len() != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", cl.Len())
	}
	recent := cl.Recent(3)
	// newest first: ports 5, 4, 3
	expected := []uint16{5, 4, 3}
	for i, ev := range recent {
		if ev.Port.LocalPort != expected[i] {
			t.Errorf("event[%d]: want port %d, got %d", i, expected[i], ev.Port.LocalPort)
		}
	}
}

func TestChangelogRecentNewestFirst(t *testing.T) {
	cl := NewChangelog(10)
	cl.Record(ChangeAdded, makeChangePort(100))
	time.Sleep(time.Millisecond)
	cl.Record(ChangeRemoved, makeChangePort(200))

	recent := cl.Recent(2)
	if len(recent) != 2 {
		t.Fatalf("expected 2 events, got %d", len(recent))
	}
	if recent[0].Port.LocalPort != 200 {
		t.Errorf("expected newest port 200 first, got %d", recent[0].Port.LocalPort)
	}
	if recent[0].Change != ChangeRemoved {
		t.Errorf("expected ChangeRemoved, got %s", recent[0].Change)
	}
}

func TestChangelogRecentLimitsCap(t *testing.T) {
	cl := NewChangelog(10)
	for i := uint16(1); i <= 8; i++ {
		cl.Record(ChangeAdded, makeChangePort(i))
	}
	recent := cl.Recent(3)
	if len(recent) != 3 {
		t.Fatalf("expected 3 events, got %d", len(recent))
	}
}

func TestChangelogRecentZeroReturnsNil(t *testing.T) {
	cl := NewChangelog(10)
	cl.Record(ChangeAdded, makeChangePort(9000))
	if got := cl.Recent(0); got != nil {
		t.Errorf("expected nil for n=0, got %v", got)
	}
}

func TestChangelogClear(t *testing.T) {
	cl := NewChangelog(10)
	cl.Record(ChangeAdded, makeChangePort(443))
	cl.Record(ChangeRemoved, makeChangePort(80))
	cl.Clear()
	if cl.Len() != 0 {
		t.Errorf("expected 0 after Clear, got %d", cl.Len())
	}
}

func TestChangelogDefaultMaxSize(t *testing.T) {
	cl := NewChangelog(0) // should default to 256
	for i := 0; i < 300; i++ {
		cl.Record(ChangeAdded, makeChangePort(uint16(i%65535+1)))
	}
	if cl.Len() != 256 {
		t.Errorf("expected 256 with default max, got %d", cl.Len())
	}
}
