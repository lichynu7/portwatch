package ports

import "strings"

// FilterConfig holds criteria for filtering scanned ports.
type FilterConfig struct {
	// AllowedPorts is a set of ports that should never trigger alerts.
	AllowedPorts map[uint16]bool
	// AllowedProcesses is a set of process name substrings considered safe.
	AllowedProcesses []string
}

// NewFilterConfig creates a FilterConfig from slices of allowed ports and process names.
func NewFilterConfig(ports []uint16, processes []string) *FilterConfig {
	allowed := make(map[uint16]bool, len(ports))
	for _, p := range ports {
		allowed[p] = true
	}
	return &FilterConfig{
		AllowedPorts:     allowed,
		AllowedProcesses: processes,
	}
}

// IsSafe returns true if the given port entry is considered safe according
// to the filter configuration.
func (fc *FilterConfig) IsSafe(port Port) bool {
	if fc == nil {
		return false
	}
	if fc.AllowedPorts[port.Port] {
		return true
	}
	for _, proc := range fc.AllowedProcesses {
		if proc != "" && strings.Contains(port.Process, proc) {
			return true
		}
	}
	return false
}

// ExcludeSafe filters out ports deemed safe, returning only unexpected listeners.
func ExcludeSafe(ports []Port, fc *FilterConfig) []Port {
	if fc == nil {
		return ports
	}
	out := make([]Port, 0, len(ports))
	for _, p := range ports {
		if !fc.IsSafe(p) {
			out = append(out, p)
		}
	}
	return out
}
