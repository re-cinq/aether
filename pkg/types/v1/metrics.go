package v1

import "log/slog"

// Represents the metrics for a specific service
// The key is the unique name of the resource
type Metrics map[string]Metric

// Helper method for adding a specific metric
func (m Metrics) Upsert(metric *Metric) {
	// Make sure the map is initialized
	if m == nil {
		// TODO we should initilize this instead of showing an error
		// we should also change this to use pointers
		slog.Error("metrics map is nil")
	}

	// Assign the resource
	if metric != nil {
		m[metric.name] = *metric
	}
}
