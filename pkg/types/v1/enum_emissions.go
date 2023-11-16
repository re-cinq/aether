package v1

import (
	"encoding/json"
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

// Custom deserialization for EmissionUnit
func (e *EmissionUnit) UnmarshalJSON(data []byte) error {
	var value string

	// Unmarshall the bytes
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// Make sure the unmarshalled string value exists
	if emissionUnit, ok := EmissionUnits[value]; !ok {
		return ErrParsingEmissionUnit
	} else {
		*e = emissionUnit
	}

	return nil
}
