package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func TestNewSlackHandlerMissingURL(t *testing.T) {
	_, err := NewSlackHandler("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}
}

func TestNewSlackHandlerValid(t *testing.T) {
	h, err := NewSlackHandler("https://hooks.slack.com/services/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestSlackHandlerSendsJSON(t *testing.T) {
	var received slackPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h, err := NewSlackHandler(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := NewAlert(KindNew, ports.Port{Port: 8080, Proto: "tcp", Addr: "0.0.0.0"})
	if err := h.Send(a); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if received.Text == "" {
		t.Error("expected non-empty slack message text")
	}
}

func TestSlackHandlerNon2xxError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	h, err := NewSlackHandler(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := NewAlert(KindNew, ports.Port{Port: 9090, Proto: "tcp", Addr: "127.0.0.1"})
	if err := h.Send(a); err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSlackHandlerClosedServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	ts.Close() // close immediately so the request will fail

	h, err := NewSlackHandler(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := NewAlert(KindNew, ports.Port{Port: 443, Proto: "tcp", Addr: "0.0.0.0"})
	if err := h.Send(a); err == nil {
		t.Fatal("expected error when sending to closed server, got nil")
	}
}
