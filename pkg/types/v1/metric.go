package v1

import (
	"fmt"
	"log/slog"
	"time"
)

// Metric tracks the uilization and emission of a specific resource
type Metric struct {

	// Unique name for the resource
	// For instance a virtual machine could have multiple disks attached to it
	// So the name could be the disk name
	Name string

	// The resource type
	// - Cpu
	// - Memory
	// - Storage
	// - Network
	ResourceType ResourceType

	// The resource usage in percentage
	// It is a value between 0 and 100
	Usage float64

	// The total amount of unit types
	// - total amount of vCPUs of a VM
	// - disk size
	UnitAmount float64

	// The unit type representing this resource
	// Examples:
	// - vCPUs: in case of a CPU
	// - Gb: in case of a Disk
	// - Gb: in case of Ram
	Unit ResourceUnit

	// The energy consumption calculated
	// for the metric. This is then multiplied
	// by the pue and grid coefficient to get
	// the Emissions data
	Energy float64

	// Emissions at a specific point in time
	Emissions ResourceEmissions

	// Time of update
	UpdatedAt time.Time

	// The resource specific labels
	Labels Labels
}

// Creates a new metric
func NewMetric(name string) *Metric {
	// Make sure the service name is set
	if name == "" {
		slog.Error("failed to create metric, got an empty name")
		return nil
	}

	// Build the service
	return &Metric{
		Name:      name,
		UpdatedAt: time.Now().UTC(),
		Labels:    Labels{},
	}
}

// Creates a string representation for the resource
// Useful for logging or debugging
func (r *Metric) String() string {
	// Basic string
	out := fmt.Sprintf("type:%s name:%s | amount:%f %s | usage:%f%%", r.ResourceType, r.Name, r.UnitAmount, r.Unit, r.Usage)

	// if we have emissions show them
	if r.Emissions.Value > 0 {
		out = fmt.Sprintf("%s => %f %s", out, r.Emissions.Value, r.Emissions.Unit)
	}

	return out
}

// Automatically update the last updated time to now
func (r *Metric) SetUpdatedAt() {
	// Assign the updated at
	r.UpdatedAt = time.Now().UTC()
}
