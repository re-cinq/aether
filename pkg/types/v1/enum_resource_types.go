package v1

import (
	"encoding/json"
	"errors"
)

// The resource type that we are collecting data for
type ResourceType string

// Error parsing the ResourceType
var ErrParsingResourceType = errors.New("unsupported ResourceType")

// ResourceTypes: Lookup map for listing all the supported resources
// as well as deserializing them
var ResourceTypes = map[string]ResourceType{
	cpuString:     Cpu,
	memoryString:  Memory,
	storageString: Storage,
	networkString: Network,
}

const (

	// CPU resource
	Cpu ResourceType = cpuString

	// Memory resource
	Memory ResourceType = memoryString

	// Storage resource
	Storage ResourceType = storageString

	// Network resource
	Network ResourceType = networkString

	// Constant string definitions
	cpuString     = "cpu"
	memoryString  = "memory"
	storageString = "storage"
	networkString = "network"
)

// Return the resource type as string
func (rt ResourceType) String() string {
	return string(rt)
}

// Custom deserialization for ResourceType
func (rt *ResourceType) UnmarshalJSON(data []byte) error {
	var value string

	// Unmarshall the bytes
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// Make sure the unmarshalled string value exists
	if resourceType, ok := ResourceTypes[value]; !ok {
		return ErrParsingResourceType
	} else {
		*rt = resourceType
	}

	return nil
}
