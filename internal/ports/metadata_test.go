package ports

import (
	"testing"
)

func TestMetadataStoreSetAndGet(t *testing.T) {
	s := NewMetadataStore(8)
	ok := s.Set("127.0.0.1:8080", "owner", "alice")
	if !ok {
		t.Fatal("expected Set to return true")
	}
	e, found := s.Get("127.0.0.1:8080", "owner")
	if !found {
		t.Fatal("expected entry to be found")
	}
	if e.Value != "alice" {
		t.Fatalf("expected 'alice', got %q", e.Value)
	}
}

func TestMetadataStoreGetMissing(t *testing.T) {
	s := NewMetadataStore(8)
	_, found := s.Get("0.0.0.0:9000", "missing")
	if found {
		t.Fatal("expected not found for unknown key")
	}
}

func TestMetadataStoreMaxKeys(t *testing.T) {
	s := NewMetadataStore(2)
	s.Set("pk", "a", "1")
	s.Set("pk", "b", "2")
	ok := s.Set("pk", "c", "3") // should be rejected
	if ok {
		t.Fatal("expected Set to return false when limit reached")
	}
	_, found := s.Get("pk", "c")
	if found {
		t.Fatal("key 'c' should not have been stored")
	}
}

func TestMetadataStoreUpdateExistingKey(t *testing.T) {
	s := NewMetadataStore(2)
	s.Set("pk", "a", "1")
	s.Set("pk", "b", "2")
	// Updating existing key should succeed even at limit
	ok := s.Set("pk", "a", "updated")
	if !ok {
		t.Fatal("expected update of existing key to succeed")
	}
	e, _ := s.Get("pk", "a")
	if e.Value != "updated" {
		t.Fatalf("expected 'updated', got %q", e.Value)
	}
}

func TestMetadataStoreAll(t *testing.T) {
	s := NewMetadataStore(8)
	s.Set("pk", "x", "1")
	s.Set("pk", "y", "2")
	all := s.All("pk")
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestMetadataStoreAllEmpty(t *testing.T) {
	s := NewMetadataStore(8)
	if s.All("unknown") != nil {
		t.Fatal("expected nil for unknown port key")
	}
}

func TestMetadataStoreDelete(t *testing.T) {
	s := NewMetadataStore(8)
	s.Set("pk", "k", "v")
	s.Delete("pk", "k")
	_, found := s.Get("pk", "k")
	if found {
		t.Fatal("expected entry to be deleted")
	}
}

func TestMetadataStorePurge(t *testing.T) {
	s := NewMetadataStore(8)
	s.Set("pk", "a", "1")
	s.Set("pk", "b", "2")
	s.Purge("pk")
	if s.All("pk") != nil {
		t.Fatal("expected nil after purge")
	}
}

func TestMetadataStoreDefaultMaxKeys(t *testing.T) {
	s := NewMetadataStore(0) // should default to 16
	for i := 0; i < 16; i++ {
		s.Set("pk", string(rune('a'+i)), "v")
	}
	if len(s.All("pk")) != 16 {
		t.Fatalf("expected 16 entries")
	}
}
