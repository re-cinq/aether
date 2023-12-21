package v1

type ResourceEmissions struct {

	// Current amount of emissions
	value float64

	// The unit of the emission
	unit EmissionUnit
}

// New instance of the resource emission
func NewResourceEmission(value float64, unit EmissionUnit) ResourceEmissions {
	return ResourceEmissions{
		value: value,
		unit:  unit,
	}
}

// Sets the emission unit.
// Currently we support only Grams of carbon per kilowatt hour, however in the future
// we might allow different scales, like KiloGrams of carbon per kilowatt hour
func (re *ResourceEmissions) SetUnit(unit EmissionUnit) *ResourceEmissions {
	// Set the unit
	re.unit = unit

	return re
}

// Set the emission value based on the Unit
func (re *ResourceEmissions) SetValue(value float64) *ResourceEmissions {
	// Negative emissions are not allowed
	if value < 0 {
		value = 0
	}

	// Set the value now
	re.value = value

	return re
}

// Return the emission value
func (re *ResourceEmissions) Value() float64 {
	return re.value
}

// Return the emission unit
func (re *ResourceEmissions) Unit() EmissionUnit {
	return re.unit
}
