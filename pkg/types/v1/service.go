package v1

import "k8s.io/klog/v2"

// The service for which we are collecting the metrics
type Service struct {

	// Unique ID of the server
	Id string `json:"id"`

	// The metrics collection for the specific service
	Metrics Metrics `json:"metrics"`

	// Labels associated with the service
	Labels Labels `json:"labels"`
}

// Create a new service
func NewService(id string) *Service {

	// Make sure the service ID exists
	if id == "" {
		klog.Error("failed to create service, got an empty id")
		return nil
	}

	// Build the service
	return &Service{
		Id:      id,
		Metrics: Metrics{},
		Labels:  Labels{},
	}
}

// Upsert the metric for the resource
func (s *Service) UpsertMetric(resource Resource) {

	// Upsert it
	s.Metrics.Upsert(resource)

}

// Insert a label to the service
func (s *Service) AddLabel(key, value string) {

	// Insert the label
	s.Labels.Add(key, value)

}
