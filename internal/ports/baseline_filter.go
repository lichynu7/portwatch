package ports

// ExcludeBaseline returns a FilterFunc that suppresses any port present in the
// provided Baseline. This allows operators to acknowledge known listeners and
// stop receiving repeated alerts for them.
func ExcludeBaseline(b *Baseline) FilterFunc {
	return func(p Port) bool {
		return b.Contains(p.Proto, p.Port)
	}
}

// ApplyBaseline filters ports slice, removing any entry found in the baseline.
// It returns only the ports that are NOT in the baseline (i.e., unexpected).
func ApplyBaseline(ports []Port, b *Baseline) []Port {
	if b == nil {
		return ports
	}
	out := make([]Port, 0, len(ports))
	for _, p := range ports {
		if !b.Contains(p.Proto, p.Port) {
			out = append(out, p)
		}
	}
	return out
}

// PartitionBaseline splits ports into two slices: known contains ports present
// in the baseline, and unknown contains ports that are not. This is useful when
// callers need to handle both groups separately rather than discarding known ports.
func PartitionBaseline(ports []Port, b *Baseline) (known, unknown []Port) {
	if b == nil {
		return nil, ports
	}
	for _, p := range ports {
		if b.Contains(p.Proto, p.Port) {
			known = append(known, p)
		} else {
			unknown = append(unknown, p)
		}
	}
	return known, unknown
}
