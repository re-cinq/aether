package models

// Resource tracks the uilization and emission of a specific resource
type Resource struct {
	// The resource usage in percentage
	Usage Percentage

	// The total amount
	Total int

	// The unit representing this resource
	Unit ResourceUnit

	// The service name: Virtual Machine, RDS database etc..
	Service string

	// The provider used as source for this metric
	Provider Provider

	// Emissions at a specific point in time
	Emissions ResourceEmissions
}

type ResourceEmissions struct {

	// Current amount of emissions
	Current float64

	// The unit of the emission
	Emission EmissionUnit
}
