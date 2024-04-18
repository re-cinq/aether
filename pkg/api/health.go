package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type health struct {
	Status        string   `json:"status"`
	Exporters     string   `json:"exporters"`
	Sources       string   `json:"sources"`
	FailedPlugins []string `json:"failedPlugins,omitempty"`
}

// Return a 200 http status
func (a *API) healthProbe(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h := &health{
		Status:    "up",
		Exporters: "up",
		Sources:   "up",
	}

	// we need to ping every registered exporter
	for _, p := range a.exporters.Plugins {
		err := p.Client.Ping()
		if err != nil {
			h.Exporters = "down"
			h.FailedPlugins = append(h.FailedPlugins, fmt.Sprintf("exporter:%s", p.Name))
		}
	}

	// we need to ping every registered source
	for _, p := range a.sources.Plugins {
		err := p.Client.Ping()
		if err != nil {
			h.Sources = "down"
			h.FailedPlugins = append(h.FailedPlugins, fmt.Sprintf("source:%s", p.Name))
		}
	}

	body, err := json.Marshal(h)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": %q}`, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(body))
}
