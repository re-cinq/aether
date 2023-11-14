package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testProviderStruct struct {
	TestProvider Provider `json:"provider"`
}

func TestProviderParser(t *testing.T) {

	testData := `{
		"provider": "prometheus"
	}`

	var testProvider testProviderStruct
	err := json.Unmarshal([]byte(testData), &testProvider)
	assert.Nil(t, err)

	assert.Equal(t, testProvider.TestProvider, Prometheus)

}
