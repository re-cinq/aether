package models

type EmissionUnit int

const (
	// Unknown emission unit
	UnknownEmissionUnit EmissionUnit = iota
	// Grams of carbon per kilowatt hour
	gCO2eqkWh
)

func (e EmissionUnit) String() string {
	switch e {
	case gCO2eqkWh:
		return "gCO2eqkWh"

	default:
		return "UnknownEmissionUnit"
	}
}
