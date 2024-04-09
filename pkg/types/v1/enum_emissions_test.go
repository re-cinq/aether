package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEmissionStruct struct {
	TestUnit EmissionUnit `json:"emissionUnit"`
}

func TestEmissionUnitParser(t *testing.T) {
	testData := `{
		"emissionUnit": "gCO2eq"
	}`

	var testEmissionUnit testEmissionStruct
	err := json.Unmarshal([]byte(testData), &testEmissionUnit)
	assert.Nil(t, err)

	assert.Equal(t, testEmissionUnit.TestUnit, GCO2eq)
}
