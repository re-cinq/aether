package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceParser(t *testing.T) {

	testData := `{
		"name": "cpu",
		"usage": 20.54,
		"total": 4.0,
		"unit": "core",
		"service": "virtual machine",
		"provider": "prometheus",
		"emissions": {
			"value": 1056.76,
			"unit": "gCO2eqkWh"
		}

	}`

	var testResource Resource
	err := json.Unmarshal([]byte(testData), &testResource)
	assert.Nil(t, err)

	assert.Equal(t, testResource.Name, "cpu")
	assert.Equal(t, testResource.Usage, Percentage(20.54))
	assert.Equal(t, testResource.Total, float64(4.0))
	assert.Equal(t, testResource.Unit, Core)
	assert.Equal(t, testResource.Service, "virtual machine")
	assert.Equal(t, testResource.Provider, Prometheus)
	assert.Equal(t, testResource.Emissions.Value, float64(1056.76))
	assert.Equal(t, testResource.Emissions.Emission, GCO2eqkWh)

}
