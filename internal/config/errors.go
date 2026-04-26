package config

import "fmt"

// errorf is a package-local helper that formats a config validation error.
func errorf(format string, args ...any) error {
	return fmt.Errorf("config: "+format, args...)
}
