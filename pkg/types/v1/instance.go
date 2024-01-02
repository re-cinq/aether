package v1

import "k8s.io/klog/v2"

// The instance for which we are collecting the metrics
// The fields are not exported for several reasons:
//   - We will expose the information of this struct via prometheus endpoint
//     so we have no reasons so serialize this.
//   - Some of the fields in the struct need to be carefully validated, so we cannot allow
//     the structs fields to be passed manually
//   - If there is the need for serializing this then we can have a model dedicated for the View
//     functionality
type Instance struct {

	// The provider used as source for this metric
	provider Provider

	// The service type (Instance, Database etc..)
	service string

	// Unique name of the instance
	// Can be the VM name
	name string

	// The region of the instance
	region string

	// The instance zone
	zone string

	// This is the kind of service
	// In case of a GCP virtual machine this could be n2-standard-8
	kind string

	// The metrics collection for the specific service
	metrics Metrics

	// The emissions of the service during operation
	operationalEmissions ResourceEmissions

	// The embodied emissions for the service
	embodiedEmissions ResourceEmissions

	// Labels associated with the service
	labels Labels
}

// The provider of the service
// GCP, AWS, Prometheus etc...
func (i *Instance) Provider() Provider {
	return i.provider
}

// The service type
func (i *Instance) Service() string {
	return i.service
}

// Returns the instance unique name
// Can be either a name or a unique id
func (i *Instance) Name() string {
	return i.name
}

// Returns the embodied emissions for the instance
func (i *Instance) EmbodiedEmissions() ResourceEmissions {
	return i.embodiedEmissions
}

// Returns the operational emissions for the instance
func (i *Instance) OperationalEmissions() ResourceEmissions {
	return i.operationalEmissions
}

// Returns the instance collected metrics
func (i *Instance) Metrics() Metrics {
	return i.metrics
}

// Returns the region of the instance
func (i *Instance) Region() string {
	return i.region
}

// Returns the instance zone
func (i *Instance) Zone() string {
	return i.zone
}

// Returns the instance kind
func (i *Instance) Kind() string {
	return i.kind
}

// Returns the instance collected labels
func (i *Instance) Labels() Labels {
	return i.labels
}

// Get the specific metric
func (i *Instance) Metric(metricName string) (Metric, bool) {
	// Get it
	resource, exists := i.metrics[metricName]

	// Return it
	return resource, exists
}

// Makes a copy of the current instance
// This allows to reuse the instance and just update its values and call build
// in case when a provider returns an array of resource metrics.
// See here as an example: pkg/providers/aws/cloudwatch.go (func getEc2Metrics())
func (i *Instance) Build() Instance {
	// Builds a copy of the instance
	return Instance{
		provider: i.provider,
		name:     i.name,
		region:   i.region,
		zone:     i.zone,
		kind:     i.kind,
		metrics:  i.metrics,
		labels:   i.labels,
	}
}

// ---------------------------------------------------------------------

// Create a new instance.
// We need both the name and the provider
func NewInstance(name string, provider Provider) *Instance {
	// Make sure the instance name is set
	if name == "" {
		klog.Error("failed to create service, got an empty name")
		return nil
	}

	// Build the instance
	return &Instance{
		name:     name,
		provider: provider,
		metrics:  Metrics{},
		labels:   Labels{},
	}
}

// Upsert the metric for the resource
func (i *Instance) UpsertMetric(resource *Metric) *Instance {
	// Upsert it
	i.metrics.Upsert(resource)

	return i
}

// Set the embodied emissions for the instance
func (i *Instance) SetEmbodiedEmissions(embodied ResourceEmissions) *Instance {
	i.embodiedEmissions = embodied
	return i
}

// Set the operational emissions for the instance
func (i *Instance) SetOperationalEmissions(emissions ResourceEmissions) *Instance {
	i.operationalEmissions = emissions
	return i
}

// Set the service type
func (i *Instance) SetService(service string) *Instance {
	i.service = service
	return i
}

// Sets the region where the resource is located
// Examples:
// - europe-west4 (GCP)
// - us-east-2 (AWS)
// - eu-east-rack-1 (Baremetal)
func (i *Instance) SetRegion(region string) *Instance {
	// Assign the region
	i.region = region

	// Allows to use it as a builder
	return i
}

// SetZone sets the instance zone name
// Examples
// - europe-west4-a (GCP)
func (i *Instance) SetZone(zone string) *Instance {
	i.zone = zone

	return i
}

// Sets the instance kind
// Examples:
// - n2-standard-8 (GCP)
// - m6.2xlarge (AWS)
func (i *Instance) SetKind(kind string) *Instance {
	// Assign the region
	i.kind = kind

	// Allows to use it as a builder
	return i
}

// Insert a label to the service
func (i *Instance) AddLabel(key, value string) *Instance {
	// Insert the label
	i.labels.Add(key, value)

	return i
}
