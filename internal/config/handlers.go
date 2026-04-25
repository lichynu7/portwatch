package config

// HandlerConfig holds the configuration for all supported alert handlers.
type HandlerConfig struct {
	Log       LogConfig       `toml:"log"`
	Email     EmailConfig     `toml:"email"`
	Webhook   WebhookConfig   `toml:"webhook"`
	Slack     SlackConfig     `toml:"slack"`
	PagerDuty PagerDutyConfig `toml:"pagerduty"`
}

// LogConfig configures the built-in log handler.
type LogConfig struct {
	Enabled bool   `toml:"enabled"`
	Level   string `toml:"level"`
}

// EmailConfig configures the SMTP email handler.
type EmailConfig struct {
	Enabled    bool     `toml:"enabled"`
	Host       string   `toml:"host"`
	Port       int      `toml:"port"`
	From       string   `toml:"from"`
	Recipients []string `toml:"recipients"`
}

// WebhookConfig configures the generic webhook handler.
type WebhookConfig struct {
	Enabled bool   `toml:"enabled"`
	URL     string `toml:"url"`
}

// SlackConfig configures the Slack webhook handler.
type SlackConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
}

// PagerDutyConfig configures the PagerDuty events handler.
type PagerDutyConfig struct {
	Enabled        bool   `toml:"enabled"`
	IntegrationKey string `toml:"integration_key"`
}

// DefaultHandlers returns a HandlerConfig with sensible defaults.
func DefaultHandlers() HandlerConfig {
	return HandlerConfig{
		Log: LogConfig{
			Enabled: true,
			Level:   "info",
		},
		Email: EmailConfig{
			Port: 25,
		},
	}
}

// EnabledHandlers returns the names of all handlers that are currently enabled.
func (h HandlerConfig) EnabledHandlers() []string {
	var enabled []string
	if h.Log.Enabled {
		enabled = append(enabled, "log")
	}
	if h.Email.Enabled {
		enabled = append(enabled, "email")
	}
	if h.Webhook.Enabled {
		enabled = append(enabled, "webhook")
	}
	if h.Slack.Enabled {
		enabled = append(enabled, "slack")
	}
	if h.PagerDuty.Enabled {
		enabled = append(enabled, "pagerduty")
	}
	return enabled
}
