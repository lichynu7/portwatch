package alert

import (
	"fmt"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Alert represents a port monitoring alert event.
type Alert struct {
	Level     Level
	Port      uint16
	Protocol  string
	PID       int
	Message   string
	Timestamp time.Time
}

// Notifier is the interface implemented by alert output backends.
type Notifier interface {
	Notify(a Alert) error
}

// NewAlert constructs an Alert with the current timestamp.
func NewAlert(level Level, port uint16, protocol string, pid int, msg string) Alert {
	return Alert{
		Level:     level,
		Port:      port,
		Protocol:  protocol,
		PID:       pid,
		Message:   msg,
		Timestamp: time.Now(),
	}
}

// StdoutNotifier writes alerts to standard output.
type StdoutNotifier struct{}

// Notify prints the alert to stdout in a structured format.
func (s *StdoutNotifier) Notify(a Alert) error {
	fmt.Printf("[%s] %s | port=%d proto=%s pid=%d msg=%q\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Port,
		a.Protocol,
		a.PID,
		a.Message,
	)
	return nil
}
