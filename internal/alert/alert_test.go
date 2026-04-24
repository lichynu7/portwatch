package alert

import (
	"errors"
	"testing"
	"time"
)

func TestNewAlert(t *testing.T) {
	before := time.Now()
	a := NewAlert(LevelWarning, 8080, "tcp", 1234, "unexpected listener")
	after := time.Now()

	if a.Level != LevelWarning {
		t.Errorf("expected level %s, got %s", LevelWarning, a.Level)
	}
	if a.Port != 8080 {
		t.Errorf("expected port 8080, got %d", a.Port)
	}
	if a.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", a.Protocol)
	}
	if a.PID != 1234 {
		t.Errorf("expected pid 1234, got %d", a.PID)
	}
	if a.Timestamp.Before(before) || a.Timestamp.After(after) {
		t.Error("timestamp out of expected range")
	}
}

// recordingNotifier captures alerts for inspection in tests.
type recordingNotifier struct {
	alerts []Alert
	errOn  int // return error on Nth call (1-based); 0 = never
	calls  int
}

func (r *recordingNotifier) Notify(a Alert) error {
	r.calls++
	r.alerts = append(r.alerts, a)
	if r.errOn > 0 && r.calls == r.errOn {
		return errors.New("notifier error")
	}
	return nil
}

func TestDispatcherDeliverToAll(t *testing.T) {
	n1 := &recordingNotifier{}
	n2 := &recordingNotifier{}
	d := NewDispatcher(n1, n2)

	a := NewAlert(LevelCritical, 22, "tcp", 999, "ssh exposed")
	if err := d.Dispatch(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(n1.alerts) != 1 || len(n2.alerts) != 1 {
		t.Errorf("expected each notifier to receive 1 alert")
	}
}

func TestDispatcherReturnsFirstError(t *testing.T) {
	failing := &recordingNotifier{errOn: 1}
	successful := &recordingNotifier{}
	d := NewDispatcher(failing, successful)

	a := NewAlert(LevelInfo, 80, "tcp", 0, "http listener")
	err := d.Dispatch(a)
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
	// Successful notifier must still have been called.
	if len(successful.alerts) != 1 {
		t.Error("successful notifier should still receive alert despite earlier error")
	}
}

func TestDispatcherRegister(t *testing.T) {
	d := NewDispatcher()
	n := &recordingNotifier{}
	d.Register(n)

	a := NewAlert(LevelWarning, 3000, "tcp", 42, "dev server")
	_ = d.Dispatch(a)

	if len(n.alerts) != 1 {
		t.Error("dynamically registered notifier did not receive alert")
	}
}
