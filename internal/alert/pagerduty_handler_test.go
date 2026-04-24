package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPagerDutyHandlerMissingKey(t *testing.T) {
	_, err := NewPagerDutyHandler("")
	if err == nil {
		t.Fatal("expected error for empty integration key")
	}
}

func TestNewPagerDutyHandlerValid(t *testing.T) {
	h, err := NewPagerDutyHandler("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestPagerDutyHandlerSendsJSON(t *testing.T) {
	var received pagerDutyPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	h, _ := NewPagerDutyHandler("my-routing-key")
	h.client = ts.Client()

	// Redirect to test server by temporarily overriding via a custom transport.
	h.client = &http.Client{
		Transport: redirectTransport(ts.URL),
	}

	a := Alert{
		Message:   "unexpected port 9999 detected",
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	if err := h.Send(a); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.RoutingKey != "my-routing-key" {
		t.Errorf("routing key = %q, want %q", received.RoutingKey, "my-routing-key")
	}
	if received.EventAction != "trigger" {
		t.Errorf("event_action = %q, want trigger", received.EventAction)
	}
	if received.Payload.Summary != a.Message {
		t.Errorf("summary = %q, want %q", received.Payload.Summary, a.Message)
	}
	if received.Payload.Source != "portwatch" {
		t.Errorf("source = %q, want portwatch", received.Payload.Source)
	}
}

func TestPagerDutyHandlerNon2xxError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	h, _ := NewPagerDutyHandler("key")
	h.client = &http.Client{Transport: redirectTransport(ts.URL)}

	err := h.Send(Alert{Message: "test", Timestamp: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx status")
	}
}
