package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceOperations(t *testing.T) {
	// Make sure we get a nil back in case the service name is empty
	failedService := NewService("", Prometheus)
	assert.Nil(t, failedService)

	// Mocking the service id
	serviceID := "1234"

	// Mocking the metric
	r := NewMetric("cpu")
	r.SetUsagePercentage(170.4).SetUnitAmount(4.0).SetResourceUnit(Core)
	r.SetEmissions(NewResourceEmission(1024.57, GCO2eqkWh))
	r.SetUpdatedAt()

	// Make sure the usage validation succeeded
	assert.Equal(t, Percentage(100), r.Usage())

	// Create a new service
	service := NewService(serviceID, Prometheus).SetRegion("europe-west4-a").SetKind("n2-standard-8")

	// Make sure the region is assigned correctly
	assert.Equal(t, serviceID, service.Name())
	assert.Equal(t, Prometheus, service.Provider())
	assert.Equal(t, "europe-west4-a", service.Region())
	assert.Equal(t, "n2-standard-8", service.Kind())

	// Add the metrics
	service.UpsertMetric(r)

	// Add a label
	service.AddLabel("name", "test")

	// Make sure the label exists
	assert.True(t, service.Labels().Exists("name"))

	// Check the resource was added
	existingResource, exists := service.Metric(r.Name())

	// Make sure the resource exists
	assert.True(t, exists)

	// Make sure the resource is the same
	assert.Equal(t, *r, existingResource)

	// Test the Build functionality

	// Build the current view of the service
	currentService := service.Build()

	// Change the region
	service.SetRegion("new-region")

	// Build the updated view of the service
	updatedService := service.Build()

	// Make sure the two are actually different
	assert.NotEqual(t, updatedService, currentService)
}
