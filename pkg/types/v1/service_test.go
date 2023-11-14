package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceOperations(t *testing.T) {

	// Mocking the service id
	serviceId := "1234"

	// Mocking the resource
	r := NewResource("cpu")

	// Set additional values
	r.SetUsage(70.4).SetTotal(4.0).SetResourceUnit(Core)
	r.SetProvider(Prometheus).SetService("virtual machines")
	r.SetProvider(Prometheus).SetEmissions(NewResourceEmissions(1024.57, GCO2eqkWh))
	r.SetUpdatedAt()

	// Create a new service
	service := NewService(serviceId)

	// Add a metric
	service.UpsertMetric(*r)

	// Add a label
	service.AddLabel("name", "test")

	assert.True(t, service.Labels.Exists("name"))

}
