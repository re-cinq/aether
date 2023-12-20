package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAwsCache(t *testing.T) {
	region := "my-region"
	service := "my-service"
	id := "my-id"
	name := "test-instance"

	resource := newGCPResource(region, service, id, "n2-standard-8", "spot", name, 2)
	assert.Equal(t, region, resource.region)
	assert.Equal(t, service, resource.service)
	assert.Equal(t, id, resource.id)

	// Init the cache
	cache := newGCPCache()

	// Check the exists fails
	assert.False(t, cache.Exists(region, service, name))

	// Add a new resource
	cache.Add(resource)

	// Check the exists succeeds
	assert.True(t, cache.Exists(region, service, name))

	// Delete the resource
	cache.Delete(region, service, name)

	// Check the exists fails
	assert.False(t, cache.Exists(region, service, name))

	// Make sure we basically have an empty map
	_, exists := cache.entries[region]
	assert.False(t, exists)
}
