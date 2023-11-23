package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Return a 200 http status
func healthProbe(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}
