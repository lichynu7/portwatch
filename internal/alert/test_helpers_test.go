package alert

import (
	"net/http"
	"strings"
)

// redirectTransport rewrites all outbound request URLs to the given base URL,
// preserving the path. This allows tests to redirect HTTP clients to a local
// httptest.Server without modifying production code.
type redirectTransport string

func (base redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	target := strings.TrimRight(string(base), "/")
	cloned.URL.Scheme = "http"
	cloned.URL.Host = strings.TrimPrefix(target, "http://")
	return http.DefaultTransport.RoundTrip(cloned)
}
