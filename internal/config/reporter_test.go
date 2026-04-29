package config

import "testing"

func TestDefaultReporterConfig(t *testing.T) {
	cfg := DefaultReporterConfig()
	if cfg.Format != "table" {
		t.Errorf("expected format 'table', got %q", cfg.Format)
	}
	if cfg.OutputFile != "" {
		t.Errorf("expected empty OutputFile, got %q", cfg.OutputFile)
	}
	if len(cfg.IncludeFields) != 0 {
		t.Errorf("expected empty IncludeFields, got %v", cfg.IncludeFields)
	}
}

func TestReporterConfigValidateOK(t *testing.T) {
	cases := []struct {
		name   string
		format string
	}{
		{"table format", "table"},
		{"json format", "json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := ReporterConfig{Format: tc.format}
			if err := cfg.Validate(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestReporterConfigValidateEmpty(t *testing.T) {
	cfg := ReporterConfig{Format: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty format")
	}
}

func TestReporterConfigValidateInvalid(t *testing.T) {
	cfg := ReporterConfig{Format: "csv"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unsupported format 'csv'")
	}
}

func TestReporterConfigValidateWithOutputFile(t *testing.T) {
	cfg := ReporterConfig{
		Format:     "json",
		OutputFile: "/var/log/portwatch/report.json",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
