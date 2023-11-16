package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResourceTypeStruct struct {
	TestResource ResourceType `json:"type"`
}

func TestResourceTypeParser(t *testing.T) {

	testData := `{
		"type": "cpu"
	}`

	var testResource testResourceTypeStruct
	err := json.Unmarshal([]byte(testData), &testResource)
	assert.Nil(t, err)

	assert.Equal(t, testResource.TestResource, CPU)

}
