// Caches in memory the resources present in a specific GCP accounts
package gcp

import (
	"sync"
	"time"
)

// GCP Region
type gcpRegion = string

// The name of the GCP service
type gcpService = string

// Resource id
type gcpResourceId = string

// The list of resources for a specific service and region
type gcpResources = map[gcpResourceId]gcpResource

// The GCP Resource representation
type gcpResource struct {
	// The GCP resource id
	id gcpResourceId

	// The region where the resource is located
	region gcpRegion

	// The service the resource belongs to
	service gcpService

	// For example spot, reserved
	lifecycle string

	// Amount of cores
	coreCount int

	// The instance kind for example n2-standard-8
	kind string

	// The name
	name string

	// When was the last time it was updated
	lastUpdated time.Time
}

// / Helper which creates a new awsResource
func newGCPResource(region gcpRegion, service gcpService, id gcpResourceId,
	kind, lifecycle, name string, coreCount int) *gcpResource {
	return &gcpResource{
		id:          id,
		service:     service,
		region:      region,
		lifecycle:   lifecycle,
		coreCount:   coreCount,
		kind:        kind,
		name:        name,
		lastUpdated: time.Now().UTC(),
	}
}

// The cache for the specific service
type serviceCache = map[gcpService]gcpResources

// We do not need to expose the cache outside of this modules
// since its main purpose is to optimise the GCP API queries
// and reduce them overall
type gcpCache struct {
	// Sync for concurrent mutations of the cache
	lock sync.RWMutex

	// cache itseld
	entries map[gcpRegion]serviceCache
}

// Initialise the empty cache
func newGCPCache() *gcpCache {
	return &gcpCache{
		lock:    sync.RWMutex{},
		entries: make(map[gcpRegion]map[gcpService]map[gcpResourceId]gcpResource),
	}
}

// Add an entry to the cache
func (cache *gcpCache) Add(entry *gcpResource) {

	// Lock the map cause we are writing
	cache.lock.Lock()

	// Add it
	serviceCache, exists := cache.entries[entry.region]
	// If the region does not exist, init the cache for that region
	if !exists {
		// Init the cache
		serviceCache = make(map[gcpService]map[gcpResourceId]gcpResource)
		cache.entries[entry.region] = serviceCache
	}

	serviceEntries, exists := serviceCache[entry.service]
	// if the entries for the service are missing, initialise them
	if !exists {
		serviceEntries = make(map[gcpResourceId]gcpResource)
		serviceCache[entry.service] = serviceEntries
	}

	// Set the entry
	serviceEntries[entry.name] = *entry

	// Unlock it
	cache.lock.Unlock()

}

// Deletes an entry from the cache, in case it was removed from AWS
func (cache *gcpCache) Delete(region gcpRegion, service gcpService, name string) {

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
	delete(serviceEntries, name)

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
func (cache *gcpCache) Exists(region gcpRegion, service gcpService, name string) bool {

	// Check if the entry exists
	return cache.Get(region, service, name) != nil

}

// Get a specific resource
func (cache *gcpCache) Get(region gcpRegion, service gcpService, name string) *gcpResource {

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
	entry, exists := resources[name]
	if !exists {
		return nil
	}

	return &entry

}
