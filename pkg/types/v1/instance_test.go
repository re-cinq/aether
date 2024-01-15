package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceOperations(t *testing.T) {
	// Make sure we get a nil back in case the instance name is empty
	failedInstance := NewInstance("", Prometheus)
	assert.Nil(t, failedInstance)

	// Mocking the instance id
	instanceID := "1234"

	// Mocking the metric
	r := NewMetric("cpu")
	r.SetUsage(170.4).SetUnitAmount(4.0).SetResourceUnit(Core)
	r.SetEmissions(NewResourceEmission(1024.57, GCO2eqkWh))
	r.SetUpdatedAt()

	// Create a new instance
	instance := NewInstance(instanceID, Prometheus).SetRegion("europe-west4").SetKind("n2-standard-8")

	// Make sure the region is assigned correctly
	assert.Equal(t, instanceID, instance.Name())
	assert.Equal(t, Prometheus, instance.Provider())
	assert.Equal(t, "europe-west4", instance.Region())
	assert.Equal(t, "n2-standard-8", instance.Kind())

	// Add the metrics
	instance.UpsertMetric(r)

	// Add a label
	instance.AddLabel("name", "test")

	// Make sure the label exists
	assert.True(t, instance.Labels().Exists("name"))

	// Check the resource was added
	existingResource, exists := instance.Metric(r.Name())

	// Make sure the resource exists
	assert.True(t, exists)

	// Make sure the resource is the same
	assert.Equal(t, *r, existingResource)

	// Test the Build functionality

	// Build the current view of the instance
	currentInstance := instance.Build()

	// Change the region
	instance.SetRegion("new-region")

	// Build the updated view of the
	updatedInstance := instance.Build()

	// Make sure the two are actually different
	assert.NotEqual(t, updatedInstance, currentInstance)
}
