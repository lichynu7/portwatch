package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const pagerDutyEventsURL = "https://events.pagerduty.com/v2/enqueue"

// PagerDutyHandler sends alerts to PagerDuty via the Events API v2.
type PagerDutyHandler struct {
	integrationKey string
	client         *http.Client
}

type pagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     pagerDutyDetails  `json:"payload"`
}

type pagerDutyDetails struct {
	Summary   string `json:"summary"`
	Source    string `json:"source"`
	Severity  string `json:"severity"`
	Timestamp string `json:"timestamp"`
}

// NewPagerDutyHandler creates a PagerDutyHandler.
// integrationKey is the PagerDuty Events API v2 integration key.
func NewPagerDutyHandler(integrationKey string) (*PagerDutyHandler, error) {
	if integrationKey == "" {
		return nil, fmt.Errorf("pagerduty: integration key must not be empty")
	}
	return &PagerDutyHandler{
		integrationKey: integrationKey,
		client:         &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers an alert to PagerDuty.
func (h *PagerDutyHandler) Send(a Alert) error {
	body := pagerDutyPayload{
		RoutingKey:  h.integrationKey,
		EventAction: "trigger",
		Payload: pagerDutyDetails{
			Summary:   a.Message,
			Source:    "portwatch",
			Severity:  "warning",
			Timestamp: a.Timestamp.UTC().Format(time.RFC3339),
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("pagerduty: marshal payload: %w", err)
	}

	resp, err := h.client.Post(pagerDutyEventsURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("pagerduty: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("pagerduty: unexpected status %d", resp.StatusCode)
	}
	return nil
}
