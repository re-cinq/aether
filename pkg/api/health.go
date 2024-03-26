package api

import (
	"fmt"
	"net/http"
)

// Return a 200 http status
func healthProbe(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "up"}`)
}
