package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Return a 200 http status
func (a *API) healthProbe(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	plugins := "up"
	var failedPlugins []string

	// we need to ping every registered plugin
	for _, p := range a.plugins.Plugins {
		err := p.Client.Ping()
		if err != nil {
			plugins = "down"
			failedPlugins = append(failedPlugins, p.Name)
		}
	}

	failedPluginsJSON, err := json.Marshal(failedPlugins)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": %q}`, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "up", "plugins": %q, failedPlugins: %s}`, plugins, failedPluginsJSON)
}
