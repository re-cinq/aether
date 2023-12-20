package v1

import bus "github.com/re-cinq/go-bus"

const (
	// Topic to be subscribed to when interested in metrics
	MetricsCollectedTopic bus.Topic = iota

	// Topic to be subscribed to when interested in emissions
	EmissionsCalculatedTopic
)
