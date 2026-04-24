package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds the configuration for the email alert handler.
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	To       []string
}

// emailHandler sends alert notifications via email.
type emailHandler struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailHandler creates a new Handler that sends alerts via SMTP email.
func NewEmailHandler(cfg EmailConfig) (Handler, error) {
	if cfg.SMTPHost == "" {
		return nil, fmt.Errorf("email handler: smtp host is required")
	}
	if len(cfg.To) == 0 {
		return nil, fmt.Errorf("email handler: at least one recipient is required")
	}

	var auth smtp.Auth
	if cfg.Username != "" {
		auth = smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	}

	return &emailHandler{cfg: cfg, auth: auth}, nil
}

// Send delivers the alert as an email message.
func (h *emailHandler) Send(a *Alert) error {
	subject := fmt.Sprintf("[portwatch] %s alert: port %d/%s",
		a.Severity, a.Port.Port, a.Port.Proto)

	body := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		strings.Join(h.cfg.To, ", "),
		h.cfg.From,
		subject,
		a.Message,
	)

	addr := fmt.Sprintf("%s:%d", h.cfg.SMTPHost, h.cfg.SMTPPort)
	err := smtp.SendMail(addr, h.auth, h.cfg.From, h.cfg.To, []byte(body))
	if err != nil {
		return fmt.Errorf("email handler: failed to send alert: %w", err)
	}
	return nil
}
