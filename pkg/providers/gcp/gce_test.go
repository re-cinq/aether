package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCEInstanceTypeExtraction(t *testing.T) {

	instanceTypeUrl := "https://www.googleapis.com/compute/v1/projects/cloud-carbon-project/zones/europe-north1-a/machineTypes/e2-micro"
	expected := "e2-micro"

	parsed := getValueFromURL(instanceTypeUrl)
	assert.Equal(t, expected, parsed)

}

func TestGCEInstanceZoneExtraction(t *testing.T) {

	zoneUrl := "https://www.googleapis.com/compute/v1/projects/cloud-carbon-project/zones/europe-north1-a"
	expected := "europe-north1-a"

	parsed := getValueFromURL(zoneUrl)
	assert.Equal(t, expected, parsed)

}
