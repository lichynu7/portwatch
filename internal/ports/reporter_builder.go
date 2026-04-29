package ports

import (
	"fmt"
	"os"

	"github.com/example/portwatch/internal/config"
)

// BuildReporter constructs a Reporter from a ReporterConfig.
// If cfg.OutputFile is set, the file is opened (or created) for appending.
// The caller is responsible for closing the returned *os.File if non-nil.
func BuildReporter(cfg config.ReporterConfig) (*Reporter, *os.File, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("reporter: invalid config: %w", err)
	}

	format := ReportFormat(cfg.Format)

	if cfg.OutputFile == "" {
		r, err := NewReporter(format, os.Stdout)
		return r, nil, err
	}

	f, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, fmt.Errorf("reporter: cannot open output file %q: %w", cfg.OutputFile, err)
	}

	r, err := NewReporter(format, f)
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return r, f, nil
}

// EntryFromPort converts a Port into a ReportEntry ready for rendering.
func EntryFromPort(p Port) ReportEntry {
	process := ""
	pid := 0
	if p.Process != nil {
		process = p.Process.Name
		pid = p.Process.PID
	}
	return ReportEntry{
		Timestamp: p.SeenAt,
		Port:      p.Port,
		Protocol:  p.Protocol,
		Process:   process,
		PID:       pid,
		Service:   p.Service,
		Severity:  p.Severity.String(),
	}
}
