package ports

import (
	"strings"

	"github.com/user/portwatch/internal/config"
)

// Label keys applied to Port entries.
const (
	LabelProtocol = "protocol"
	LabelPortRange = "port_range"
	LabelTrusted   = "trusted"
)

// Labeler attaches metadata labels to Port entries based on configurable rules.
type Labeler struct {
	cfg config.LabelerConfig
}

// NewLabeler creates a Labeler from the supplied config.
func NewLabeler(cfg config.LabelerConfig) (*Labeler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Labeler{cfg: cfg}, nil
}

// Label annotates each Port in the slice with computed labels and returns the
// augmented slice. The original slice is mutated in place.
func (l *Labeler) Label(ports []Port) []Port {
	for i := range ports {
		if ports[i].Labels == nil {
			ports[i].Labels = make(map[string]string)
		}
		l.applyProtocol(&ports[i])
		l.applyPortRange(&ports[i])
		l.applyTrusted(&ports[i])
	}
	return ports
}

func (l *Labeler) applyProtocol(p *Port) {
	proto := strings.ToLower(p.Protocol)
	if proto == "" {
		proto = "unknown"
	}
	p.Labels[LabelProtocol] = proto
}

func (l *Labeler) applyPortRange(p *Port) {
	switch {
	case p.Port < 1024:
		p.Labels[LabelPortRange] = "privileged"
	case p.Port < 49152:
		p.Labels[LabelPortRange] = "registered"
	default:
		p.Labels[LabelPortRange] = "ephemeral"
	}
}

func (l *Labeler) applyTrusted(p *Port) {
	for _, proc := range l.cfg.TrustedProcesses {
		if strings.EqualFold(p.Process, proc) {
			p.Labels[LabelTrusted] = "true"
			return
		}
	}
	p.Labels[LabelTrusted] = "false"
}
