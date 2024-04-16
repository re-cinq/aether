package v1

import "context"

// Source is an interface for fetching metrics that can be calculated
type Source interface {
	// Stop is the "teardown" that will be used for graceful shutdown
	Stop(context.Context) error

	// Fetch is the business logic that should return a list of instances
	// that have metrics attached to them mainly cpu, memory, storage and network
	Fetch(context.Context) ([]*Instance, error)
}
