package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsParser(t *testing.T) {

	// Create a valid resource
	cpuResource := NewMetric("cpu")
	assert.NotNil(t, cpuResource)

	cpuResource.SetUsagePercentage(20.54).SetTotal(4.0)
	cpuResource.SetType(Cpu).SetResourceUnit(Core)
	cpuResource.SetEmissions(NewResourceEmissions(1056.76, GCO2eqkWh))

	assert.Equal(t, cpuResource.Usage(), Percentage(20.54))
	assert.Equal(t, cpuResource.Total(), float64(4.0))
	assert.Equal(t, cpuResource.Unit(), Core)
	assert.Equal(t, cpuResource.Type(), Cpu)
	assert.Equal(t, cpuResource.emissions.Value(), float64(1056.76))
	assert.Equal(t, cpuResource.emissions.Unit(), GCO2eqkWh)

	metrics := Metrics{}
	metrics.Upsert(*cpuResource)

	existing := metrics[cpuResource.Name()]
	assert.Equal(t, *cpuResource, existing)

}
