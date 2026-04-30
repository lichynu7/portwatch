package ports

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Fingerprint represents a stable hash of a set of open ports,
// used to detect changes between scan cycles.
type Fingerprint struct {
	Hash  string
	Count int
}

// NewFingerprint computes a deterministic SHA-256 fingerprint
// over the provided ports slice. Order of ports does not matter.
func NewFingerprint(ports []Port) Fingerprint {
	if len(ports) == 0 {
		return Fingerprint{Hash: emptySHA256(), Count: 0}
	}

	keys := make([]string, 0, len(ports))
	for _, p := range ports {
		keys = append(keys, fingerprintKey(p))
	}
	sort.Strings(keys)

	h := sha256.New()
	h.Write([]byte(strings.Join(keys, "\n")))
	return Fingerprint{
		Hash:  fmt.Sprintf("%x", h.Sum(nil)),
		Count: len(ports),
	}
}

// Equal reports whether two fingerprints represent the same port set.
func (f Fingerprint) Equal(other Fingerprint) bool {
	return f.Hash == other.Hash
}

// String returns a short human-readable representation.
func (f Fingerprint) String() string {
	if len(f.Hash) > 12 {
		return fmt.Sprintf("%s... (%d ports)", f.Hash[:12], f.Count)
	}
	return fmt.Sprintf("%s (%d ports)", f.Hash, f.Count)
}

func fingerprintKey(p Port) string {
	return fmt.Sprintf("%s:%d:%s", p.Protocol, p.Port, p.State)
}

func emptySHA256() string {
	h := sha256.New()
	return fmt.Sprintf("%x", h.Sum(nil))
}
