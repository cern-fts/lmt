package proxy

import (
	"io"
	"net/http"
)

// HealthCheckHandler returns HTTP 200.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Set additional headers defined in config.yml
	for header, val := range ResponseHeaders {
		w.Header().Set(header, val)
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
