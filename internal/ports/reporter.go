package ports

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// ReportFormat controls the output format of a Report.
type ReportFormat string

const (
	FormatTable ReportFormat = "table"
	FormatJSON  ReportFormat = "json"
)

// ReportEntry holds a single row in a port report.
type ReportEntry struct {
	Timestamp time.Time
	Port      uint16
	Protocol  string
	Process   string
	PID       int
	Service   string
	Severity  string
}

// Reporter writes formatted port alert reports to an io.Writer.
type Reporter struct {
	format ReportFormat
	out    io.Writer
}

// NewReporter creates a Reporter with the given format.
// If out is nil, os.Stdout is used.
func NewReporter(format ReportFormat, out io.Writer) (*Reporter, error) {
	if format != FormatTable && format != FormatJSON {
		return nil, fmt.Errorf("reporter: unsupported format %q", format)
	}
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{format: format, out: out}, nil
}

// Write renders entries to the configured output.
func (r *Reporter) Write(entries []ReportEntry) error {
	switch r.format {
	case FormatTable:
		return r.writeTable(entries)
	case FormatJSON:
		return r.writeJSON(entries)
	}
	return nil
}

func (r *Reporter) writeTable(entries []ReportEntry) error {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tPORT\tPROTO\tSERVICE\tPROCESS\tPID\tSEVERITY")
	for _, e := range entries {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%d\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Port, e.Protocol, e.Service, e.Process, e.PID, e.Severity)
	}
	return w.Flush()
}

func (r *Reporter) writeJSON(entries []ReportEntry) error {
	fmt.Fprintln(r.out, "[")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out,
			`  {"timestamp":%q,"port":%d,"protocol":%q,"service":%q,"process":%q,"pid":%d,"severity":%q}%s\n`,
			e.Timestamp.Format(time.RFC3339), e.Port, e.Protocol,
			e.Service, e.Process, e.PID, e.Severity, comma)
	}
	fmt.Fprintln(r.out, "]")
	return nil
}
