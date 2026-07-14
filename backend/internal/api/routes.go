package api

import "net/http"

// RegisterRoutes attaches all HTTP routes to the given mux. Method-qualified
// patterns (Go 1.22+) make the mux return 405 automatically when the path
// matches but the method does not.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/calculate", HandleCalculate)
	mux.HandleFunc("GET /health", HandleHealth)
}
