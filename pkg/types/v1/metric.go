package v1

import (
	"fmt"
	"time"

	"k8s.io/klog/v2"
)

// Metric tracks the uilization and emission of a specific resource
type Metric struct {

	// Unique name for the resource
	name string

	// The resource type
	resourceType ResourceType

	// The resource usage in percentage
	usage Percentage

	// The total amount of unit types
	unitAmount float64

	// The unit representing this resource
	unit ResourceUnit

	// Emissions at a specific point in time
	emissions ResourceEmissions

	// Time of update
	updatedAt time.Time

	// The resource specific labels
	labels Labels
}

// Creates a new metric
func NewMetric(name string) *Metric {
	// Make sure the service name is set
	if name == "" {
		klog.Error("failed to create metric, got an empty name")
		return nil
	}

	// Build the service
	return &Metric{
		name:      name,
		updatedAt: time.Now().UTC(),
		labels:    Labels{},
	}
}

// The resource unique name
// For instance a virtual machine could have multiple disks attached to it
// So the name could be the disk name
func (r *Metric) Name() string {
	return r.name
}

// The resource type
// - Cpu
// - Memory
// - Storage
// - Network
func (r *Metric) Type() ResourceType {
	return r.resourceType
}

// The resource usage in percentage
// It is a value between 0 and 100
func (r *Metric) Usage() Percentage {
	return r.usage
}

// The resource amount
// In case of a virtual machine for example is the total amount of cores
func (r *Metric) UnitAmount() float64 {
	return r.unitAmount
}

// The resource unit
// - In case of a cpu is the amount of cores
// - In case of ram is a multiple of bytes
func (r *Metric) Unit() ResourceUnit {
	return r.unit
}

// The calculated emissions for the resource
func (r *Metric) Emissions() *ResourceEmissions {
	return &r.emissions
}

// When the metric was updated last
func (r *Metric) LastUpdateTime() time.Time {
	return r.updatedAt
}

// Resource labels
func (r *Metric) Labels() Labels {
	return r.labels
}

// Creates a string representation for the resource
// Useful for logging or debugging
func (r *Metric) String() string {
	// Basic string
	out := fmt.Sprintf("type:%s name:%s | amount:%f %s | usage:%f%%", r.resourceType, r.name, r.unitAmount, r.unit, r.usage)

	// if we have emissions show them
	if r.emissions.value > 0 {
		out = fmt.Sprintf("%s => %f %s", out, r.emissions.value, r.emissions.unit)
	}

	return out
}

// -------------------------------------------------------------

// The resource type
// - Cpu
// - Memory
// - Storage
// - Network
func (r *Metric) SetType(resourceType ResourceType) *Metric {
	// Set the type
	r.resourceType = resourceType

	return r
}

// Adds the usage
// Examples:
// - 50.0 (indicates a 50% usage)
func (r *Metric) SetUsagePercentage(usage float64) *Metric {
	// Make sure we are not setting the usage to a negative value
	if usage < 0 {
		usage = 0
	}

	// It does not make much sense to have a value higher than 100%
	// So set it to 100 to make the CO2eq calculations easier
	if usage > 100 {
		usage = 100
	}

	// Assign the usage
	r.usage = Percentage(usage)

	// Allows to use it as a builder
	return r
}

// Adds the total amount of the resource:
// Examples:
// - total amount of core of a VM
// - disk size
func (r *Metric) SetUnitAmount(amount float64) *Metric {
	// No reason to have a negative amount
	if amount < 0 {
		amount = 0
	}

	// Assign the amount
	r.unitAmount = amount

	// Allows to use it as a builder
	return r
}

// Adds the resource unit
// Examples:
// - cores: in case of a CPU
// - Gb: in case of a Disk
// - Gb: in case of Ram
func (r *Metric) SetResourceUnit(unit ResourceUnit) *Metric {
	// Assign the resource unit
	r.unit = unit

	// Allows to use it as a builder
	return r
}

// Automatically update the last updated time to now
func (r *Metric) SetUpdatedAt() *Metric {
	// Assign the updated at
	r.updatedAt = time.Now().UTC()

	// Allows to use it as a builder
	return r
}

// Set the emissions for the resource
func (r *Metric) SetEmissions(emissions ResourceEmissions) *Metric {
	// Assign the amount
	r.emissions = emissions

	// Allows to use it as a builder
	return r
}

// Set the labels for the resource
func (r *Metric) SetLabels(labels Labels) *Metric {
	// Assign the amount
	r.labels = labels

	// Allows to use it as a builder
	return r
}
