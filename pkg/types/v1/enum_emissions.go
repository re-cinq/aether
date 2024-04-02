package v1

import (
	"errors"
)

// EmissionUnits Defines the unit of emission of CO2 equivalent
type EmissionUnit string

// ErrParsingEmissionUnit Error parsing the EmissionUnit
var ErrParsingEmissionUnit = errors.New("unsupported EmissionUnit")

// EmissionUnits Lookup map for listing all the supported emissions
// as well as deserializing them
var EmissionUnits = map[string]EmissionUnit{
	GCO2eqkWhString: GCO2eqkWhString,
}

const (
	// Grams of carbon per kilowatt hour
	GCO2eqkWh EmissionUnit = GCO2eqkWhString

	// Constant string definitions
	GCO2eqkWhString = "gCO2eqkWh"
)

// Return the emission unit as string
func (e EmissionUnit) String() string {
	return string(e)
}
