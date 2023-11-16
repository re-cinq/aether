package amazon

import v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"

type awsMetric struct {
	// metric name
	name string

	// metric kind
	kind v1.ResourceType

	// metric unit
	unit v1.ResourceUnit

	// value
	value float64

	// instance id
	instanceId string
}
