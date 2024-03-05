package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceFunctionalities(t *testing.T) {
	// Create a valid resource
	r := NewMetric("cpu")
	assert.NotNil(t, r)

	// Assign the various values to the resource
	r.Usage = 20.54
	r.UnitAmount = 4.0
	r.Unit = VCPU
	r.ResourceType = CPU
	r.Emissions = NewResourceEmission(1056.76, GCO2eqkWh)

	// Validate the data
	assert.Equal(t, r.Name, "cpu")
	assert.Equal(t, r.Usage, 20.54)
	assert.Equal(t, r.UnitAmount, float64(4.0))
	assert.Equal(t, r.ResourceType, CPU)
	assert.Equal(t, r.Unit, VCPU)
	assert.Equal(t, r.Emissions.Value, float64(1056.76))
	assert.Equal(t, r.Emissions.Unit, GCO2eqkWh)

	// Now create a metric with an empty name which should return nil
	r = NewMetric("")
	assert.Nil(t, r)
}
