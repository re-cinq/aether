package models

type ResourceUnit int

const (
	// UnknownResourceUnit
	UnknownResourceUnit ResourceUnit = iota

	// CPU Cores
	Cores

	// -------------------------------------------
	// Used for both Ram and Disk

	// Kb: Kilobytes
	Kb

	// Mb: Megabytes
	Mb

	// Gb: Gigabytes
	Gb

	// // -------------------------------------------
	// Used for bandwidth

	// Kbs: Kilobytes per second
	Kbs

	// Mbs: Megabytes per second
	Mbs

	// Gbs: Gigabytes per second
	Gbs
)

func (ru ResourceUnit) String() string {
	switch ru {
	case Cores:
		return "Cores"
	case Kb:
		return "Kb"
	case Mb:
		return "Mb"
	case Gb:
		return "Gb"
	case Kbs:
		return "Kbs"
	case Mbs:
		return "Mbs"
	case Gbs:
		return "Gbs"
	default:
		return "UnknownResourceUnit"
	}
}
