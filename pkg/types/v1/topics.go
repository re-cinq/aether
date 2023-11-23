package v1

import "github.com/re-cinq/cloud-carbon/pkg/bus"

const (
	// Topic to be subscribed to when interested in metrics
	MetricsCollectedTopic bus.Topic = iota

	// Topic to be subscribed to when interested in emissions
	EmissionsCalculatedTopic
)
