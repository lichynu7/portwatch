package ports

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAggregatorFlushOnMaxBatch(t *testing.T) {
	var mu sync.Mutex
	var got []Port

	cfg := AggregatorConfig{Window: 10 * time.Second, MaxBatch: 3}
	agg := NewAggregator(cfg, func(batch []Port) {
		mu.Lock()
		got = append(got, batch...)
		mu.Unlock()
	})

	for i := 0; i < 3; i++ {
		agg.Add(Port{Port: uint16(8000 + i)})
	}

	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 3 {
		t.Fatalf("expected 3 ports flushed, got %d", len(got))
	}
}

func TestAggregatorExplicitFlush(t *testing.T) {
	var mu sync.Mutex
	var got []Port

	cfg := AggregatorConfig{Window: 10 * time.Second, MaxBatch: 100}
	agg := NewAggregator(cfg, func(batch []Port) {
		mu.Lock()
		got = append(got, batch...)
		mu.Unlock()
	})

	agg.Add(Port{Port: 9000})
	agg.Add(Port{Port: 9001})
	agg.Flush()

	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 ports after explicit flush, got %d", len(got))
	}
}

func TestAggregatorRunCancels(t *testing.T) {
	var mu sync.Mutex
	var got []Port

	cfg := AggregatorConfig{Window: 50 * time.Millisecond, MaxBatch: 100}
	agg := NewAggregator(cfg, func(batch []Port) {
		mu.Lock()
		got = append(got, batch...)
		mu.Unlock()
	})

	agg.Add(Port{Port: 7777})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		agg.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not return after context cancel")
	}

	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected final flush on cancel, got none")
	}
}

func TestAggregatorEmptyFlushIsNoop(t *testing.T) {
	called := false
	cfg := AggregatorConfig{Window: time.Second, MaxBatch: 10}
	agg := NewAggregator(cfg, func(_ []Port) { called = true })
	agg.Flush()
	time.Sleep(10 * time.Millisecond)
	if called {
		t.Fatal("flush on empty buffer should not invoke flushFn")
	}
}
