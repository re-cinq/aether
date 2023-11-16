package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResourceUnitStruct struct {
	TestResourceUnit ResourceUnit `json:"unit"`
}

func TestResourceUnitParser(t *testing.T) {
	testData := `{
		"unit": "core"
	}`

	var testResourceUnit testResourceUnitStruct
	err := json.Unmarshal([]byte(testData), &testResourceUnit)
	assert.Nil(t, err)

	assert.Equal(t, testResourceUnit.TestResourceUnit, Core)
}
