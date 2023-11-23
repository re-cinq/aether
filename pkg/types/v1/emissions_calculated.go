package v1

import "github.com/re-cinq/cloud-carbon/pkg/bus"

type EmissionsCalculated struct {
	Instance Service
}

// The topic this event is about
func (e *EmissionsCalculated) Topic() bus.Topic {
	return EmissionsCalculatedTopic
}

// Returns the unique name of the instance or service
func (e *EmissionsCalculated) Id() string {
	return e.Instance.name
}
