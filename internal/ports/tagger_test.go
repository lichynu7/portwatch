package ports

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func baseTaggerCfg(rules ...config.TagRule) config.TaggerConfig {
	return config.TaggerConfig{Rules: rules}
}

func TestTaggerNoRules(t *testing.T) {
	tg := NewTagger(baseTaggerCfg())
	ports := []Port{{Port: 8080, Protocol: "tcp"}}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 0 {
		t.Fatalf("expected no tags, got %v", out[0].Tags)
	}
}

func TestTaggerMatchByPort(t *testing.T) {
	rule := config.TagRule{Port: 443, Tags: []string{"https", "tls"}}
	tg := NewTagger(baseTaggerCfg(rule))
	ports := []Port{{Port: 443, Protocol: "tcp"}}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", out[0].Tags)
	}
}

func TestTaggerNoMatchDifferentPort(t *testing.T) {
	rule := config.TagRule{Port: 443, Tags: []string{"https"}}
	tg := NewTagger(baseTaggerCfg(rule))
	ports := []Port{{Port: 80, Protocol: "tcp"}}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 0 {
		t.Fatalf("expected no tags, got %v", out[0].Tags)
	}
}

func TestTaggerMatchByProtocol(t *testing.T) {
	rule := config.TagRule{Protocol: "udp", Tags: []string{"udp-traffic"}}
	tg := NewTagger(baseTaggerCfg(rule))
	ports := []Port{
		{Port: 53, Protocol: "udp"},
		{Port: 53, Protocol: "tcp"},
	}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 1 || out[0].Tags[0] != "udp-traffic" {
		t.Fatalf("expected udp-traffic tag on udp port, got %v", out[0].Tags)
	}
	if len(out[1].Tags) != 0 {
		t.Fatalf("expected no tags on tcp port, got %v", out[1].Tags)
	}
}

func TestTaggerMatchByProcessName(t *testing.T) {
	rule := config.TagRule{ProcessName: "nginx", Tags: []string{"webserver"}}
	tg := NewTagger(baseTaggerCfg(rule))
	ports := []Port{{Port: 80, Protocol: "tcp", ProcessName: "nginx"}}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 1 || out[0].Tags[0] != "webserver" {
		t.Fatalf("expected webserver tag, got %v", out[0].Tags)
	}
}

func TestTaggerDeduplicatesTags(t *testing.T) {
	rule1 := config.TagRule{Port: 80, Tags: []string{"http"}}
	rule2 := config.TagRule{Protocol: "tcp", Tags: []string{"http"}}
	tg := NewTagger(baseTaggerCfg(rule1, rule2))
	ports := []Port{{Port: 80, Protocol: "tcp"}}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 1 {
		t.Fatalf("expected 1 deduplicated tag, got %v", out[0].Tags)
	}
}

func TestTaggerPreservesExistingTags(t *testing.T) {
	rule := config.TagRule{Port: 22, Tags: []string{"ssh"}}
	tg := NewTagger(baseTaggerCfg(rule))
	ports := []Port{{Port: 22, Protocol: "tcp", Tags: []string{"known"}}}
	out := tg.Tag(ports)
	if len(out[0].Tags) != 2 {
		t.Fatalf("expected 2 tags (known + ssh), got %v", out[0].Tags)
	}
}
