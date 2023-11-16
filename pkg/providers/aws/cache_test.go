package amazon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwsCache(t *testing.T) {

	region := "my-region"
	service := "my-service"
	id := "my-id"

	resource := newAWSResource(region, service, id, "t2.micro", "spot", "test-instance")
	assert.Equal(t, region, resource.region)
	assert.Equal(t, service, resource.service)
	assert.Equal(t, id, resource.id)

	// Init the cache
	cache := newAWSCache()

	// Check the exists fails
	assert.False(t, cache.Exists(region, service, id))

	// Add a new resource
	cache.Add(resource)

	// Check the exists succeeds
	assert.True(t, cache.Exists(region, service, id))

	// Delete the resource
	cache.Delete(region, service, id)

	// Check the exists fails
	assert.False(t, cache.Exists(region, service, id))

	// Make sure we basically have an empty map
	_, exists := cache.entries[region]
	assert.False(t, exists)

}
