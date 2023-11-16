package v1

import "k8s.io/klog/v2"

// Labels definition
type Labels map[string]string

// Helper method for adding a label
func (labels Labels) Add(key, value string) {
	// Make sure the map is initialized
	if labels == nil {
		klog.Fatal("labels map is nil")
	}

	// Assign the label
	labels[key] = value
}

// Helper method for getting a specific label
func (labels Labels) Get(key string) (string, bool) {
	// Make sure the map is initialized
	if labels == nil {
		return "", false
	}

	// Read the value
	value, exists := labels[key]

	// Return it
	return value, exists
}

// Helper method for getting a specific label
func (labels Labels) Exists(key string) bool {
	// Make sure the map is initialized
	if labels == nil {
		return false
	}

	// Read the value
	_, exists := labels[key]

	// Return it
	return exists
}

// Helper method for deleting a specific label
func (labels Labels) Delete(key string) {
	// Make sure the map is initialized
	if labels == nil {
		return
	}

	// Delete it
	delete(labels, key)
}
