package config

import (
	"testing"
)

func TestDefaultTaggerConfig(t *testing.T) {
	cfg := DefaultTaggerConfig()
	if !cfg.Enabled {
		t.Fatal("expected tagger to be enabled by default")
	}
	if len(cfg.Rules) == 0 {
		t.Fatal("expected at least one default rule")
	}
}

func TestTaggerConfigValidateOK(t *testing.T) {
	cfg := DefaultTaggerConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestTaggerConfigValidateEmptyTags(t *testing.T) {
	cfg := TaggerConfig{
		Enabled: true,
		Rules: []TagRule{
			{Port: 8080, Tags: []string{}},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for rule with no tags")
	}
}

func TestTaggerConfigValidateNoRulesOK(t *testing.T) {
	cfg := TaggerConfig{Enabled: false, Rules: nil}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error for empty rules: %v", err)
	}
}

func TestTagRuleFields(t *testing.T) {
	rule := TagRule{
		Port:        9200,
		Protocol:    "tcp",
		ProcessName: "elasticsearch",
		Tags:        []string{"search", "elastic"},
	}
	if rule.Port != 9200 {
		t.Errorf("expected port 9200, got %d", rule.Port)
	}
	if rule.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", rule.Protocol)
	}
	if rule.ProcessName != "elasticsearch" {
		t.Errorf("expected process elasticsearch, got %s", rule.ProcessName)
	}
	if len(rule.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(rule.Tags))
	}
}
