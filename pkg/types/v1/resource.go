package v1

import (
	"time"

	"k8s.io/klog/v2"
)

// Resource tracks the uilization and emission of a specific resource
type Resource struct {

	// Unique name for the resource
	// For instance a virtual machine could have multiple disks attached to it
	Name string `json:"name"`

	// The resource usage in percentage
	Usage Percentage `json:"usage"`

	// The total amount
	Total float64 `json:"total"`

	// The unit representing this resource
	Unit ResourceUnit `json:"unit"`

	// The service name: Virtual Machine, RDS database etc..
	Service string `json:"service"`

	// The provider used as source for this metric
	Provider Provider `json:"provider"`

	// Emissions at a specific point in time
	Emissions ResourceEmissions `json:"emissions"`

	// Time of update
	UpdatedAt time.Time `json:"updatedAt"`

	// The resource specific labels
	Labels Labels `json:"labels"`
}

type ResourceEmissions struct {

	// Current amount of emissions
	Value float64 `json:"value"`

	// The unit of the emission
	Emission EmissionUnit `json:"unit"`
}

// New instance of the resource emission
func NewResourceEmissions(value float64, unit EmissionUnit) ResourceEmissions {
	return ResourceEmissions{
		Value:    value,
		Emission: unit,
	}
}

// Creates a new resource
func NewResource(name string) *Resource {
	// Make sure the service name is set
	if name == "" {
		klog.Error("failed to create resource, got an empty name")
		return nil
	}

	// Build the service
	return &Resource{
		Name: name,
	}
}

// Adds the usage
func (r *Resource) SetUsage(usage float64) *Resource {

	// Assign the usage
	r.Usage = Percentage(usage)

	// Allows to use it as a builder
	return r
}

// Adds the total amount of the resource:
// Examples:
// - total amount of core of a VM
// - disk size
func (r *Resource) SetTotal(total float64) *Resource {

	// Assign the total
	r.Total = total

	// Allows to use it as a builder
	return r
}

// Adds the resource unit
// Examples:
// - cores: in case of a CPU
// - Gb: in case of a Disk
// - Gb: in case of Ram
func (r *Resource) SetResourceUnit(unit ResourceUnit) *Resource {

	// Assign the total
	r.Unit = unit

	// Allows to use it as a builder
	return r
}

// Adds the type of service
// Examples:
// - virtual machine
// - RDS
// - SQS
func (r *Resource) SetService(service string) *Resource {

	// Assign the total
	r.Service = service

	// Allows to use it as a builder
	return r
}

// Adds the provider for this specific resource
// Examples:
// - Prometheus
// - Aws
func (r *Resource) SetProvider(provider Provider) *Resource {

	// Assign the total
	r.Provider = provider

	// Allows to use it as a builder
	return r
}

// Automatically update the last updated time to now
func (r *Resource) SetUpdatedAt() *Resource {

	// Assign the updated at
	r.UpdatedAt = time.Now().UTC()

	// Allows to use it as a builder
	return r
}

// Set the emissions for the resource
func (r *Resource) SetEmissions(emissions ResourceEmissions) *Resource {

	// Assign the total
	r.Emissions = emissions

	// Allows to use it as a builder
	return r
}
