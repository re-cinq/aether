package v1

import "k8s.io/klog/v2"

// The service for which we are collecting the metrics
// The fields are not exported for several reasons:
//   - We will expose the information of this struct via prometheus endpoint
//     so we have no reasons so serialize this.
//   - Some of the fields in the struct need to be carefully validated, so we cannot allow
//     the structs fields to be passed manually
//   - If there is the need for serializing this then we can have a model dedicated for the View
//     functionality
type Service struct {

	// The provider used as source for this metric
	provider Provider

	// Unique name of the service
	// Can be the VM name
	name string

	// The region of the service
	region string

	// This is the kind of service
	// In case of a GCP virtual machine this could be n2-standard-8
	kind string

	// The metrics collection for the specific service
	metrics Metrics

	// Labels associated with the service
	labels Labels
}

// The provider of the service
// GCP, AWS, Prometheus etc...
func (s *Service) Provider() Provider {
	return s.provider
}

// Returns the service unique name
// Can be either a name or a unique id
func (s *Service) Name() string {
	return s.name
}

// Returns the service collected metrics
func (s *Service) Metrics() Metrics {
	return s.metrics
}

// Returns the region of the service
func (s *Service) Region() string {
	return s.region
}

// Returns the service kind
func (s *Service) Kind() string {
	return s.kind
}

// Returns the service collected labels
func (s *Service) Labels() Labels {
	return s.labels
}

// Get the specific metric
func (s *Service) Metric(metricName string) (Metric, bool) {
	// Get it
	resource, exists := s.metrics[metricName]

	// Return it
	return resource, exists
}

// Makes a copy of the current service
// This allows to reuse the service and just update its values and call build
// in case when a provider returns an array of resource metrics.
// See here as an example: pkg/providers/aws/cloudwatch.go (func getEc2Metrics())
func (s *Service) Build() Service {
	// Buils a copy of the service
	return Service{
		provider: s.provider,
		name:     s.name,
		region:   s.region,
		kind:     s.kind,
		metrics:  s.metrics,
		labels:   s.labels,
	}
}

// ---------------------------------------------------------------------

// Create a new service.
// We need both the name and the provider
func NewService(name string, provider Provider) *Service {
	// Make sure the service name is set
	if name == "" {
		klog.Error("failed to create service, got an empty name")
		return nil
	}

	// Build the service
	return &Service{
		name:     name,
		provider: provider,
		metrics:  Metrics{},
		labels:   Labels{},
	}
}

// Upsert the metric for the resource
func (s *Service) UpsertMetric(resource *Metric) *Service {
	// Upsert it
	s.metrics.Upsert(resource)

	return s
}

// Sets the region where the resource is located
// Examples:
// - europe-west4-a (GCP)
// - us-east-2 (AWS)
// - eu-east-rack-1 (Baremetal)
func (s *Service) SetRegion(region string) *Service {
	// Assign the region
	s.region = region

	// Allows to use it as a builder
	return s
}

// Sets the region where the resource is located
// Examples:
// - n2-standard-8 (GCP)
// - m6.2xlarge (AWS)
func (s *Service) SetKind(kind string) *Service {
	// Assign the region
	s.kind = kind

	// Allows to use it as a builder
	return s
}

// Insert a label to the service
func (s *Service) AddLabel(key, value string) *Service {
	// Insert the label
	s.labels.Add(key, value)

	return s
}
