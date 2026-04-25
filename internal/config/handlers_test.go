package config

import (
	"testing"
)

func TestDefaultHandlers(t *testing.T) {
	h := DefaultHandlers()

	if !h.Log.Enabled {
		t.Error("expected log handler to be enabled by default")
	}
	if h.Log.Level != "info" {
		t.Errorf("expected log level 'info', got %q", h.Log.Level)
	}
	if h.Email.Port != 25 {
		t.Errorf("expected default SMTP port 25, got %d", h.Email.Port)
	}
	if h.Email.Enabled {
		t.Error("expected email handler to be disabled by default")
	}
	if h.Webhook.Enabled {
		t.Error("expected webhook handler to be disabled by default")
	}
	if h.Slack.Enabled {
		t.Error("expected slack handler to be disabled by default")
	}
	if h.PagerDuty.Enabled {
		t.Error("expected pagerduty handler to be disabled by default")
	}
}

func TestHandlerConfigFields(t *testing.T) {
	h := HandlerConfig{
		Email: EmailConfig{
			Enabled:    true,
			Host:       "smtp.example.com",
			Port:       587,
			From:       "alerts@example.com",
			Recipients: []string{"ops@example.com", "dev@example.com"},
		},
		Webhook: WebhookConfig{
			Enabled: true,
			URL:     "https://example.com/hook",
		},
		Slack: SlackConfig{
			Enabled:    true,
			WebhookURL: "https://hooks.slack.com/services/XXX",
		},
		PagerDuty: PagerDutyConfig{
			Enabled:        true,
			IntegrationKey: "abc123",
		},
	}

	if h.Email.Host != "smtp.example.com" {
		t.Errorf("unexpected email host: %q", h.Email.Host)
	}
	if len(h.Email.Recipients) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(h.Email.Recipients))
	}
	if h.Webhook.URL != "https://example.com/hook" {
		t.Errorf("unexpected webhook URL: %q", h.Webhook.URL)
	}
	if h.Slack.WebhookURL != "https://hooks.slack.com/services/XXX" {
		t.Errorf("unexpected slack webhook URL: %q", h.Slack.WebhookURL)
	}
	if h.PagerDuty.IntegrationKey != "abc123" {
		t.Errorf("unexpected PagerDuty key: %q", h.PagerDuty.IntegrationKey)
	}
}
