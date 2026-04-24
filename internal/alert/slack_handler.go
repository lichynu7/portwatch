package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackHandler sends alert notifications to a Slack incoming webhook URL.
type SlackHandler struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackHandler creates a SlackHandler that posts messages to the given
// Slack incoming webhook URL. Returns an error if the URL is empty.
func NewSlackHandler(webhookURL string) (*SlackHandler, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("slack handler: webhook URL must not be empty")
	}
	return &SlackHandler{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send formats the alert as a Slack message and POSTs it to the webhook URL.
func (s *SlackHandler) Send(a *Alert) error {
	message := fmt.Sprintf("[portwatch] %s — port %d (%s) on %s",
		a.Kind, a.Port.Port, a.Port.Proto, a.Port.Addr)

	payload := slackPayload{Text: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack handler: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack handler: post to slack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack handler: unexpected status %d", resp.StatusCode)
	}
	return nil
}
