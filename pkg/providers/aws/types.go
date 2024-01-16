package amazon

import (
	"time"
)

// The AWS Resource representation
type resource struct {
	// The AWS resource id
	id string

	// The region where the resource is located
	region string

	// The service the resource belongs to
	service string

	// For example spot, reserved
	lifecycle string

	// Amount of cores
	coreCount int

	// The instance kind for example
	kind string

	// The name
	name string

	// When was the last time it was updated
	lastUpdated time.Time
}
