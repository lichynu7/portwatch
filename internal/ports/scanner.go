package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Listener represents an open port with its associated process info.
type Listener struct {
	Protocol string
	LocalAddr string
	Port     int
	PID      int
	Process  string
}

// Scanner reads open ports from the /proc/net filesystem.
type Scanner struct {
	procNetPath string
}

// NewScanner creates a Scanner with the default /proc/net path.
func NewScanner() *Scanner {
	return &Scanner{procNetPath: "/proc/net"}
}

// Scan returns all currently open TCP and UDP listeners.
func (s *Scanner) Scan() ([]Listener, error) {
	var listeners []Listener

	for _, proto := range []string{"tcp", "tcp6", "udp", "udp6"} {
		path := fmt.Sprintf("%s/%s", s.procNetPath, proto)
		entries, err := s.parseNetFile(path, proto)
		if err != nil {
			// Non-fatal: some protocols may not be available
			continue
		}
		listeners = append(listeners, entries...)
	}

	return listeners, nil
}

// parseNetFile parses a single /proc/net/{proto} file.
func (s *Scanner) parseNetFile(path, proto string) ([]Listener, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var listeners []Listener
	scanner := bufio.NewScanner(f)

	// Skip header line
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		// State 0A = LISTEN for TCP, UDP is always stateless
		state := fields[3]
		if strings.HasPrefix(proto, "tcp") && state != "0A" {
			continue
		}

		local := fields[1]
		port, err := parseHexPort(local)
		if err != nil {
			continue
		}

		listeners = append(listeners, Listener{
			Protocol:  proto,
			LocalAddr: local,
			Port:      port,
		})
	}

	return listeners, scanner.Err()
}

// parseHexPort extracts the port from a hex-encoded "addr:port" string.
func parseHexPort(addrPort string) (int, error) {
	parts := strings.Split(addrPort, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid addr:port %q", addrPort)
	}
	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return 0, err
	}
	return int(port), nil
}

// FilterByPort returns only the listeners matching the given port number.
func FilterByPort(listeners []Listener, port int) []Listener {
	var result []Listener
	for _, l := range listeners {
		if l.Port == port {
			result = append(result, l)
		}
	}
	return result
}
