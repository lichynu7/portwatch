package ports

import (
	"net"
	"strings"
)

// GeoInfo holds a minimal geographic/network classification for an IP address.
type GeoInfo struct {
	IP        string
	IsPrivate bool
	IsLoopback bool
	IsLinkLocal bool
	Label     string
}

// GeoClassifier classifies IP addresses into network locality buckets.
type GeoClassifier struct{}

// NewGeoClassifier returns a new GeoClassifier.
func NewGeoClassifier() *GeoClassifier {
	return &GeoClassifier{}
}

// Classify returns a GeoInfo for the given IP string.
// It performs purely in-process classification with no external calls.
func (g *GeoClassifier) Classify(ipStr string) GeoInfo {
	info := GeoInfo{IP: ipStr}

	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		info.Label = "invalid"
		return info
	}

	switch {
	case ip.IsLoopback():
		info.IsLoopback = true
		info.Label = "loopback"
	case ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast():
		info.IsLinkLocal = true
		info.Label = "link-local"
	case ip.IsPrivate():
		info.IsPrivate = true
		info.Label = "private"
	default:
		info.Label = "public"
	}

	return info
}

// EnrichWithGeo annotates each Port in the slice with a geo label stored in
// the Metadata map under the key "geo_label".
func EnrichWithGeo(ports []Port, gc *GeoClassifier) []Port {
	out := make([]Port, 0, len(ports))
	for _, p := range ports {
		info := gc.Classify(p.IP)
		if p.Metadata == nil {
			p.Metadata = make(map[string]string)
		}
		p.Metadata["geo_label"] = info.Label
		out = append(out, p)
	}
	return out
}
