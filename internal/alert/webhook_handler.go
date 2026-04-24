package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookHandler delivers alerts as JSON POST requests to a configured URL.
type WebhookHandler struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Timestamp string `json:"timestamp"`
	Port      int    `json:"port"`
	Proto     string `json:"proto"`
	Message   string `json:"message"`
}

// NewWebhookHandler creates a WebhookHandler that posts to the given URL.
// Returns an error if url is empty.
func NewWebhookHandler(url string, timeout time.Duration) (*WebhookHandler, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook url must not be empty")
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &WebhookHandler{
		url: url,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// Send encodes the alert as JSON and POSTs it to the configured webhook URL.
func (w *WebhookHandler) Send(a Alert) error {
	payload := webhookPayload{
		Timestamp: a.Timestamp.UTC().Format(time.RFC3339),
		Port:      a.Port,
		Proto:     a.Proto,
		Message:   a.Message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.url)
	}
	return nil
}
