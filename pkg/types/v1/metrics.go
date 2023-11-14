package v1

import "k8s.io/klog/v2"

// Represents the metrics for a specific service
// The key is the unique name of the resource
type Metrics map[string]Resource

// Helper method for adding a label
func (m Metrics) Upsert(resource Resource) {

	// Make sure the map is initialised
	if m == nil {
		klog.Fatal("metrics map is nil")
	}

	// Assign the resource
	m[resource.Name] = resource

}
