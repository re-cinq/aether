package v1

import (
	"encoding/json"
	"errors"
)

// ResourceUnit The unit of the resource type that we are collecting data for
type ResourceUnit string

// ErrParsingResourceUnits Error parsing the ResourceUnits
var ErrParsingResourceUnits = errors.New("unsupported ResourceUnits")

// ResourceUnits: Lookup map for listing all the supported resource units
// as well as deserializing them
var ResourceUnits = map[string]ResourceUnit{
	coreString: Core,
	kbString:   Kb,
	mbString:   Mb,
	gbString:   Gb,
	tbString:   Tb,
	kbsString:  Kbs,
	mbsString:  Mbs,
	gbsString:  Gbs,
	tbsString:  Tbs,
}

const (

	// CPU Core
	Core ResourceUnit = coreString

	// -------------------------------------------
	// Used for both Ram and Disk

	// Kb: Kilobytes
	Kb ResourceUnit = kbString

	// Mb: Megabytes
	Mb ResourceUnit = mbString

	// Gb: Gigabytes
	Gb ResourceUnit = gbString

	// Tb: Terabytes
	Tb ResourceUnit = tbString

	// Pb: Petabytes
	Pb ResourceUnit = pbString

	// // -------------------------------------------
	// Used for bandwidth

	// Kbs: Kilobits per second
	Kbs ResourceUnit = kbsString

	// Mbs: Megabits per second
	Mbs ResourceUnit = mbsString

	// Gbs: Gigabits per second
	Gbs ResourceUnit = gbsString

	// Tbs: Terabits per second
	Tbs ResourceUnit = tbsString

	// Static strings
	coreString = "core"
	kbString   = "kb"
	mbString   = "mb"
	gbString   = "gb"
	tbString   = "tb"
	pbString   = "pb"

	kbsString = "kbit/s"
	mbsString = "mbit/s"
	gbsString = "gbit/s"
	tbsString = "tbit/s"
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
