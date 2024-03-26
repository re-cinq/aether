package v1

import (
	"github.com/re-cinq/cloud-carbon/pkg/bus"
)

const (
	// used to specify the event when metrics have been collected
	MetricsCollectedEvent bus.EventType = iota

	// used to speicfy the event when emissions for instances have been
	// calculated
	EmissionsCalculatedEvent
)
