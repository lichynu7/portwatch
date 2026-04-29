package ports

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var testEntries = []ReportEntry{
	{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Port:      8080,
		Protocol:  "tcp",
		Process:   "nginx",
		PID:       1234,
		Service:   "http-alt",
		Severity:  "warning",
	},
	{
		Timestamp: time.Date(2024, 1, 15, 10, 1, 0, 0, time.UTC),
		Port:      443,
		Protocol:  "tcp",
		Process:   "nginx",
		PID:       1234,
		Service:   "https",
		Severity:  "info",
	},
}

func TestNewReporterInvalidFormat(t *testing.T) {
	_, err := NewReporter("xml", nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNewReporterDefaultsToStdout(t *testing.T) {
	r, err := NewReporter(FormatTable, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestReporterTableOutput(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewReporter(FormatTable, &buf)
	if err := r.Write(testEntries); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TIMESTAMP") {
		t.Error("expected header in table output")
	}
	if !strings.Contains(out, "8080") {
		t.Error("expected port 8080 in output")
	}
	if !strings.Contains(out, "nginx") {
		t.Error("expected process name in output")
	}
}

func TestReporterJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewReporter(FormatJSON, &buf)
	if err := r.Write(testEntries); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"port":8080`) {
		t.Error("expected port field in JSON output")
	}
	if !strings.Contains(out, `"severity":"warning"`) {
		t.Error("expected severity field in JSON output")
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Error("expected JSON array start")
	}
}

func TestReporterEmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewReporter(FormatTable, &buf)
	if err := r.Write(nil); err != nil {
		t.Fatalf("unexpected error on empty entries: %v", err)
	}
}
