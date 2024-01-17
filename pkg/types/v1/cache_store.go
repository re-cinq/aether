package v1

import "time"

type Resource struct {
	// The resource ID
	ID string

	// The name
	Name string

	// The region where the resource is located
	Region string

	// The service the resource belongs to
	Service string

	// For example spot, reserved
	Lifecycle string

	// Amount of cores
	CoreCount int

	// The instance kind for example
	Kind string

	// When was the last time it was updated
	LastUpdated time.Time
}
