package v1

import bus "github.com/re-cinq/go-bus"

type EmissionsCalculated struct {
	Instance Instance
}

// The topic this event is about
//
//nolint:all
func (e EmissionsCalculated) Topic() bus.Topic {
	return EmissionsCalculatedTopic
}

// Returns the unique name of the instance or service
//
//nolint:all
func (e EmissionsCalculated) Identifier() string {
	return e.Instance.Name
}
