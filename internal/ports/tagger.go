package ports

import (
	"strings"

	"github.com/user/portwatch/internal/config"
)

// Tag represents a string label attached to a port event.
type Tag = string

// Tagger applies a set of user-defined and automatic tags to a Port.
type Tagger struct {
	cfg config.TaggerConfig
}

// NewTagger constructs a Tagger from the given configuration.
func NewTagger(cfg config.TaggerConfig) *Tagger {
	return &Tagger{cfg: cfg}
}

// Tag enriches each Port in the slice with matching tags and returns the
// updated slice. The original slice is modified in place.
func (t *Tagger) Tag(ports []Port) []Port {
	for i := range ports {
		ports[i].Tags = t.tagsFor(ports[i])
	}
	return ports
}

func (t *Tagger) tagsFor(p Port) []string {
	seen := make(map[string]struct{})
	var tags []string

	add := func(tag string) {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return
		}
		if _, ok := seen[tag]; !ok {
			seen[tag] = struct{}{}
			tags = append(tags, tag)
		}
	}

	// Carry forward any tags already on the port.
	for _, existing := range p.Tags {
		add(existing)
	}

	// Apply rules from config.
	for _, rule := range t.cfg.Rules {
		if ruleMatches(rule, p) {
			for _, tag := range rule.Tags {
				add(tag)
			}
		}
	}

	return tags
}

func ruleMatches(rule config.TagRule, p Port) bool {
	if rule.Port != 0 && rule.Port != p.Port {
		return false
	}
	if rule.Protocol != "" && !strings.EqualFold(rule.Protocol, p.Protocol) {
		return false
	}
	if rule.ProcessName != "" && !strings.Contains(strings.ToLower(p.ProcessName), strings.ToLower(rule.ProcessName)) {
		return false
	}
	return true
}
