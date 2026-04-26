// Package ports provides port scanning, filtering, and enrichment utilities.
package ports

import (
	"fmt"
	"net"
)

// PortInfo holds enriched metadata about an open port.
type PortInfo struct {
	Port     uint16
	Protocol string
	Address  string
	Process  *ProcessInfo
	Service  string
}

// Enricher adds metadata to raw port entries.
type Enricher struct {
	resolveService bool
	lookupProcess  bool
}

// EnricherOption configures an Enricher.
type EnricherOption func(*Enricher)

// WithServiceResolution enables well-known service name lookup.
func WithServiceResolution() EnricherOption {
	return func(e *Enricher) { e.resolveService = true }
}

// WithProcessLookup enables process info enrichment via /proc.
func WithProcessLookup() EnricherOption {
	return func(e *Enricher) { e.lookupProcess = true }
}

// NewEnricher constructs an Enricher with the provided options.
func NewEnricher(opts ...EnricherOption) *Enricher {
	e := &Enricher{}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Enrich takes a slice of raw Port values and returns enriched PortInfo entries.
func (e *Enricher) Enrich(ports []Port) []PortInfo {
	out := make([]PortInfo, 0, len(ports))
	for _, p := range ports {
		info := PortInfo{
			Port:     p.Port,
			Protocol: p.Protocol,
			Address:  p.Address,
		}
		if e.resolveService {
			info.Service = resolveServiceName(p.Port, p.Protocol)
		}
		if e.lookupProcess {
			proc, err := LookupProcess(p.Inode)
			if err == nil {
				info.Process = proc
			}
		}
		out = append(out, info)
	}
	return out
}

// resolveServiceName returns the well-known service name for a port/protocol pair.
func resolveServiceName(port uint16, proto string) string {
	name, err := net.LookupPort(proto, fmt.Sprintf("%d", port))
	if err != nil || name == 0 {
		return ""
	}
	// net.LookupPort returns the port number; use net.DefaultResolver for name
	svc, err := net.LookupCNAME(fmt.Sprintf("_%d._%s.local", port, proto))
	if err != nil || svc == "" {
		return knownServices[port]
	}
	return svc
}

// knownServices is a small built-in map for common ports.
var knownServices = map[uint16]string{
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
}
