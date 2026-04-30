package ports

import (
	"testing"
)

func TestGeoClassifyLoopback(t *testing.T) {
	gc := NewGeoClassifier()
	info := gc.Classify("127.0.0.1")
	if !info.IsLoopback {
		t.Fatalf("expected loopback, got %+v", info)
	}
	if info.Label != "loopback" {
		t.Fatalf("expected label loopback, got %s", info.Label)
	}
}

func TestGeoClassifyPrivate(t *testing.T) {
	gc := NewGeoClassifier()
	for _, ip := range []string{"10.0.0.1", "192.168.1.5", "172.16.0.3"} {
		info := gc.Classify(ip)
		if !info.IsPrivate {
			t.Fatalf("expected private for %s, got %+v", ip, info)
		}
		if info.Label != "private" {
			t.Fatalf("expected label private for %s, got %s", ip, info.Label)
		}
	}
}

func TestGeoClassifyPublic(t *testing.T) {
	gc := NewGeoClassifier()
	info := gc.Classify("8.8.8.8")
	if info.Label != "public" {
		t.Fatalf("expected public, got %s", info.Label)
	}
}

func TestGeoClassifyInvalid(t *testing.T) {
	gc := NewGeoClassifier()
	info := gc.Classify("not-an-ip")
	if info.Label != "invalid" {
		t.Fatalf("expected invalid, got %s", info.Label)
	}
}

func TestGeoClassifyLoopbackIPv6(t *testing.T) {
	gc := NewGeoClassifier()
	info := gc.Classify("::1")
	if !info.IsLoopback {
		t.Fatalf("expected loopback for ::1, got %+v", info)
	}
}

func TestEnrichWithGeo(t *testing.T) {
	ports := []Port{
		{IP: "127.0.0.1", Port: 8080},
		{IP: "10.0.0.1", Port: 9090},
		{IP: "8.8.8.8", Port: 53},
	}
	gc := NewGeoClassifier()
	enriched := EnrichWithGeo(ports, gc)

	if len(enriched) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(enriched))
	}

	expected := []string{"loopback", "private", "public"}
	for i, p := range enriched {
		label := p.Metadata["geo_label"]
		if label != expected[i] {
			t.Errorf("port %d: expected geo_label %s, got %s", i, expected[i], label)
		}
	}
}

func TestEnrichWithGeoPreservesExistingMetadata(t *testing.T) {
	ports := []Port{
		{IP: "127.0.0.1", Port: 80, Metadata: map[string]string{"service": "http"}},
	}
	gc := NewGeoClassifier()
	enriched := EnrichWithGeo(ports, gc)

	if enriched[0].Metadata["service"] != "http" {
		t.Fatalf("existing metadata was overwritten")
	}
	if enriched[0].Metadata["geo_label"] != "loopback" {
		t.Fatalf("geo_label not set")
	}
}
