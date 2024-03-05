package v1

import (
	"encoding/json"
	"errors"
)

// ResourceUnit The unit of the resource type that we are collecting data for
type ResourceUnit string

// ErrParsingResourceUnits Error parsing the ResourceUnits
var ErrParsingResourceUnits = errors.New("unsupported ResourceUnits")

// ResourceUnits Lookup map for listing all the supported resource units
// as well as deserializing them
var ResourceUnits = map[string]ResourceUnit{
	vCPUString: VCPU,
	kbString:   KB,
	mbString:   MB,
	gbString:   GB,
	tbString:   TB,
	kbsString:  KBs,
	mbsString:  MBs,
	gbsString:  GBs,
	tbsString:  TBs,
}

const (

	// vCPU
	VCPU ResourceUnit = vCPUString

	// -------------------------------------------
	// Used for both Ram and Disk

	// KB: Kilobytes
	KB ResourceUnit = kbString

	// MB: Megabytes
	MB ResourceUnit = mbString

	// GB: Gigabytes
	GB ResourceUnit = gbString

	// TB: Terabytes
	TB ResourceUnit = tbString

	// PB: Petabytes
	PB ResourceUnit = pbString

	// // -------------------------------------------
	// Used for bandwidth

	// KBs: Kilobits per second
	KBs ResourceUnit = kbsString

	// MBs: Megabits per second
	MBs ResourceUnit = mbsString

	// GBs: Gigabits per second
	GBs ResourceUnit = gbsString

	// TBs: Terabits per second
	TBs ResourceUnit = tbsString

	// Static strings
	vCPUString = "vCPU"
	kbString   = "KB"
	mbString   = "MB"
	gbString   = "GB"
	tbString   = "TB"
	pbString   = "PB"

	kbsString = "KB/s"
	mbsString = "MB/s"
	gbsString = "GB/s"
	tbsString = "TB/s"
)

// Returns a string representation of the Resource unit
func (ru ResourceUnit) String() string {
	return string(ru)
}

// Custom deserialization for ResourceUnit
func (ru *ResourceUnit) UnmarshalJSON(data []byte) error {
	var value string

	// Unmarshall the bytes
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// Make sure the unmarshalled string value exists
	if resourceUnit, ok := ResourceUnits[value]; !ok {
		return ErrParsingResourceUnits
	} else {
		*ru = resourceUnit
	}

	return nil
}
