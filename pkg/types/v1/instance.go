package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/re-cinq/aether/pkg/log"
)

// Used to speicy the state of an instance
type InstanceStatus string

const (
	// Instance is currently functioning
	InstancePending InstanceStatus = "pending"

	// Instance is currently functioning
	InstanceRunning InstanceStatus = "running"

	// Instance has been terminated and will be removed from the system
	InstanceTerminated InstanceStatus = "terminated"
)

// The instance for which we are collecting the metrics
type Instance struct {
	// unique identifier
	ID string

	// The provider used as source for this metric
	Provider Provider

	// The service type (Instance, Database etc..)
	Service string

	// Unique name of the instance
	// Can be the VM name
	Name string

	// The region of the instance
	// Examples:
	// - europe-west4 (GCP)
	// - us-east-2 (AWS)
	// - eu-east-rack-1 (Baremetal)
	Region string

	// The instance zone
	// - europe-west4-a (GCP)
	Zone string

	// This is the kind of service
	// Examples for VMs:
	// - n2-standard-8 (GCP)
	// - m6.2xlarge (AWS)
	Kind string

	// Status of the instance
	Status InstanceStatus

	// The metrics collection for the specific service
	// Operational emissions are stored here
	Metrics Metrics

	// The embodied emissions for the service
	EmbodiedEmissions ResourceEmissions

	// Labels associated with the service
	Labels Labels
}

// Create a new instance.
// We need both the name and the provider
func NewInstance(name string, provider Provider) *Instance {
	// Make sure the instance name is set
	if name == "" {
		slog.Error("failed to create service, got an empty name")
		return nil
	}

	// Build the instance
	return &Instance{
		Name:     name,
		Provider: provider,
		Metrics:  Metrics{},
		Labels:   Labels{},
	}
}

func (i *Instance) PrintPretty(ctx context.Context) {
	logger := log.FromContext(ctx)

	for _, m := range i.Metrics {
		logger.Debug(fmt.Sprintf(
			"Collected metric: %s %s %s %s | %s",
			i.Service,
			i.Region,
			i.Name,
			i.Kind,
			m.String(),
		))
	}
}
