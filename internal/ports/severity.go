package ports

import "strings"

// Severity represents the urgency level of a detected port alert.
type Severity int

const (
	SeverityInfo    Severity = iota // Known/expected port appeared
	SeverityWarning                 // Port outside common safe ranges
	SeverityCritical                // Port in a sensitive or privileged range
)

// String returns a human-readable label for the severity.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// SeverityConfig defines thresholds used when classifying port alerts.
type SeverityConfig struct {
	// PrivilegedMax is the upper bound (inclusive) of privileged ports.
	// Ports in [1, PrivilegedMax] are classified as Critical.
	PrivilegedMax uint16

	// EphemeralMin is the lower bound of ephemeral/dynamic ports.
	// Ports >= EphemeralMin are classified as Info.
	EphemeralMin uint16

	// CriticalProcesses is a list of process name substrings that always
	// trigger a Critical severity regardless of port number.
	CriticalProcesses []string
}

// DefaultSeverityConfig returns a SeverityConfig with sensible defaults.
func DefaultSeverityConfig() SeverityConfig {
	return SeverityConfig{
		PrivilegedMax:     1023,
		EphemeralMin:      32768,
		CriticalProcesses: []string{"nc", "ncat", "nmap", "socat"},
	}
}

// Classify returns the Severity for a given port entry using cfg.
func Classify(p Port, cfg SeverityConfig) Severity {
	// Process-name-based override takes priority.
	for _, name := range cfg.CriticalProcesses {
		if strings.Contains(strings.ToLower(p.Process), name) {
			return SeverityCritical
		}
	}

	switch {
	case p.Port <= cfg.PrivilegedMax && p.Port > 0:
		return SeverityCritical
	case p.Port >= cfg.EphemeralMin:
		return SeverityInfo
	default:
		return SeverityWarning
	}
}
