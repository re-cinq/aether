package source

import v1 "github.com/re-cinq/aether/pkg/types/v1"

// Source is an interface for fetching metrics that can be calculated
type Source interface {
	// Setup is used for the initial loading of a Source
	// and initialization functionality should happen at this point
	Setup() error

	// Stop is the "teardown" that will be used for graceful shutdown
	Stop() error

	// Fetch is the business logic that should return a list of instances
	// that have metrics attached to them mainly cpu, memory, storage and network
	Fetch() ([]v1.Instance, error)
}
