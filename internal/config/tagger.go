package config

import "fmt"

// TagRule defines a matching rule that attaches one or more tags to a port
// when all non-zero fields match.
type TagRule struct {
	// Port matches an exact port number. Zero means any port.
	Port uint16 `toml:"port" yaml:"port"`
	// Protocol matches "tcp" or "udp" (case-insensitive). Empty means any.
	Protocol string `toml:"protocol" yaml:"protocol"`
	// ProcessName is a substring match against the owning process name.
	ProcessName string `toml:"process_name" yaml:"process_name"`
	// Tags is the list of tags to attach when this rule matches.
	Tags []string `toml:"tags" yaml:"tags"`
}

// TaggerConfig holds all tagging rules.
type TaggerConfig struct {
	// Enabled controls whether the tagger stage runs in the pipeline.
	Enabled bool `toml:"enabled" yaml:"enabled"`
	// Rules is the ordered list of tag rules to evaluate.
	Rules []TagRule `toml:"rules" yaml:"rules"`
}

// DefaultTaggerConfig returns a TaggerConfig with sensible defaults.
func DefaultTaggerConfig() TaggerConfig {
	return TaggerConfig{
		Enabled: true,
		Rules: []TagRule{
			{Port: 22, Protocol: "tcp", Tags: []string{"ssh"}},
			{Port: 80, Protocol: "tcp", Tags: []string{"http"}},
			{Port: 443, Protocol: "tcp", Tags: []string{"https", "tls"}},
			{Port: 3306, Tags: []string{"database", "mysql"}},
			{Port: 5432, Tags: []string{"database", "postgres"}},
		},
	}
}

// Validate checks the TaggerConfig for consistency.
func (c TaggerConfig) Validate() error {
	for i, rule := range c.Rules {
		if len(rule.Tags) == 0 {
			return fmt.Errorf("tagger rule[%d] has no tags defined", i)
		}
	}
	return nil
}
