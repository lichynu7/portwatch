package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseHexPort(t *testing.T) {
	tests := []struct {
		input    string
		wantPort int
		wantErr  bool
	}{
		{"00000000:0050", 80, false},
		{"00000000:01BB", 443, false},
		{"00000000:270F", 9999, false},
		{"invalid", 0, true},
		{"00000000:ZZZZ", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseHexPort(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseHexPort(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.wantPort {
				t.Errorf("parseHexPort(%q) = %d, want %d", tt.input, got, tt.wantPort)
			}
		})
	}
}

func TestScannerParseNetFile(t *testing.T) {
	// Write a minimal fake /proc/net/tcp file
	dir := t.TempDir()
	tcpFile := filepath.Join(dir, "tcp")

	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 00000000:01BB 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12346 1 0000000000000000 100 0 0 10 0
   2: 0100007F:0035 00000000:0000 01 00000000:00000000 00:00000000 00000000   101        0 12347 1 0000000000000000 100 0 0 10 0
`
	if err := os.WriteFile(tcpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	s := &Scanner{procNetPath: dir}
	listeners, err := s.parseNetFile(tcpFile, "tcp")
	if err != nil {
		t.Fatalf("parseNetFile returned error: %v", err)
	}

	if len(listeners) != 2 {
		t.Fatalf("expected 2 listeners (state 0A), got %d", len(listeners))
	}

	expectedPorts := []int{80, 443}
	for i, l := range listeners {
		if l.Port != expectedPorts[i] {
			t.Errorf("listener[%d].Port = %d, want %d", i, l.Port, expectedPorts[i])
		}
		if l.Protocol != "tcp" {
			t.Errorf("listener[%d].Protocol = %q, want \"tcp\"", i, l.Protocol)
		}
	}
}

func TestScannerMissingFile(t *testing.T) {
	s := &Scanner{procNetPath: "/nonexistent/path"}
	// Scan should not return an error even if proto files are missing
	listeners, err := s.Scan()
	if err != nil {
		t.Errorf("Scan() expected no error on missing files, got: %v", err)
	}
	if len(listeners) != 0 {
		t.Errorf("expected 0 listeners, got %d", len(listeners))
	}
	_ = fmt.Sprintf("ok") // suppress import warning
}
