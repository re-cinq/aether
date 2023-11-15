// Caches in memory the resources present in a specific AWS region
package amazon

import (
	"sync"
	"time"
)

// AWS Region
type awsRegion = string

// The name of the AWS service
type awsService = string

// Resource id
type awsResourceId = string

// The list of resources for a specific service and region
type awsResources = map[awsResourceId]awsResource

// The AWS Resource representation
type awsResource struct {
	// The AWS resource id
	id awsResourceId

	// The region where the resource is located
	region awsRegion

	// The service the resource belongs to
	service awsService

	// When was the last time it was updated
	lastUpdated time.Time
}

// / Helper which creates a new awsResource
func newAWSResource(region awsRegion, service awsService, id awsResourceId) *awsResource {
	return &awsResource{
		id:          id,
		service:     service,
		region:      region,
		lastUpdated: time.Now().UTC(),
	}
}

// The cache for the specific service
type serviceCache = map[awsService]awsResources

// We do not need to expose the cache outside of this modules
// since its main purpose is to optimise the AWS API queries
// and reduce them overall
type awsCache struct {
	// Sync for concurrent mutations of the cache
	lock sync.RWMutex

	// cache itseld
	entries map[awsRegion]serviceCache
}

// Initialise the empty cache
func newAWSCache() *awsCache {
	return &awsCache{
		lock:    sync.RWMutex{},
		entries: make(map[awsRegion]map[awsService]map[awsResourceId]awsResource),
	}
}

// Add an entry to the cache
func (cache *awsCache) Add(entry *awsResource) {

	// Lock the map cause we are writing
	cache.lock.Lock()

	// Add it
	serviceCache, exists := cache.entries[entry.region]
	// If the region does not exist, init the cache for that region
	if !exists {
		// Init the cache
		serviceCache = make(map[awsService]map[awsResourceId]awsResource)
		cache.entries[entry.region] = serviceCache
	}

	serviceEntries, exists := serviceCache[entry.service]
	// if the entries for the service are missing, initialise them
	if !exists {
		serviceEntries = make(map[awsResourceId]awsResource)
		serviceCache[entry.service] = serviceEntries
	}

	// Set the entry
	serviceEntries[entry.id] = *entry

	// Unlock it
	cache.lock.Unlock()

}

// Deletes an entry from the cache, in case it was removed from AWS
func (cache *awsCache) Delete(region awsRegion, service awsService, id awsResourceId) {

	// Lock the map cause we are writing
	cache.lock.Lock()

	// Add it
	serviceCache, exists := cache.entries[region]
	// If the region does not exist, init the cache for that region
	if !exists {
		// well...nothing to do here
		return
	}

	serviceEntries, exists := serviceCache[service]
	// if the entries for the service are missing, initialise them
	if !exists {
		// well...nothing to do here
		return
	}

	// Set the entry
	delete(serviceEntries, id)

	// Now check if the serviceEntries is empty so that we can remove the whole cache
	if len(serviceEntries) == 0 {
		delete(serviceCache, service)
	}

	// Now check if the serviceCache is empty so that we can remove the whole cache
	if len(serviceCache) == 0 {
		delete(cache.entries, region)
	}

	// Unlock it
	cache.lock.Unlock()

}

// Check if an entry exists
func (cache *awsCache) Exists(region awsRegion, service awsService, id awsResourceId) bool {

	// Check if the entry exists
	return cache.Get(region, service, id) != nil

}

// Get a specific resource
func (cache *awsCache) Get(region awsRegion, service awsService, id awsResourceId) *awsResource {

	// Read lock the cache
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	// Get the region specific cache
	regionCache, exists := cache.entries[region]
	if !exists {
		return nil
	}

	// Get the resources in the specific region
	resources, exists := regionCache[service]
	if !exists {
		return nil
	}

	// Check if the resource is present
	entry, exists := resources[id]
	if !exists {
		return nil
	}

	return &entry

}
