package v1

type ResourceEmissions struct {

	// Current amount of emissions
	Value float64

	// The unit of the emission
	Unit EmissionUnit
}

// New instance of the resource emission
func NewResourceEmission(value float64, unit EmissionUnit) ResourceEmissions {
	return ResourceEmissions{
		Value: value,
		Unit:  unit,
	}
}
