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
	r.SetUsagePercentage(20.54).SetTotal(4.0).SetResourceUnit(Core).SetType(CPU)
	r.SetEmissions(NewResourceEmissions(1056.76, GCO2eqkWh))
	r.SetUpdatedAt()

	// Validate the data
	assert.Equal(t, r.Name(), "cpu")
	assert.Equal(t, r.Usage(), Percentage(20.54))
	assert.Equal(t, r.Total(), float64(4.0))
	assert.Equal(t, r.Type(), CPU)
	assert.Equal(t, r.Unit(), Core)
	assert.Equal(t, r.Emissions().Value(), float64(1056.76))
	assert.Equal(t, r.Emissions().Unit(), GCO2eqkWh)

	// Now create a metric with an empty name which should return nil
	r = NewMetric("")
	assert.Nil(t, r)
}
