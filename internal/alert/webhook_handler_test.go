package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewWebhookHandlerMissingURL(t *testing.T) {
	_, err := NewWebhookHandler("", 0)
	if err == nil {
		t.Fatal("expected error for empty url")
	}
}

func TestNewWebhookHandlerValid(t *testing.T) {
	wh, err := NewWebhookHandler("http://example.com/hook", 3*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wh.url != "http://example.com/hook" {
		t.Errorf("unexpected url: %s", wh.url)
	}
}

func TestWebhookHandlerSendsJSON(t *testing.T) {
	var received webhookPayload
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
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	wh, err := NewWebhookHandler(ts.URL, 3*time.Second)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	a, _ := NewAlert(8080, "tcp", "new listener detected")
	if err := wh.Send(a); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", received.Proto)
	}
	if received.Message != "new listener detected" {
		t.Errorf("unexpected message: %s", received.Message)
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestWebhookHandlerNon2xxError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	wh, _ := NewWebhookHandler(ts.URL, 3*time.Second)
	a, _ := NewAlert(9090, "udp", "test")
	if err := wh.Send(a); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
