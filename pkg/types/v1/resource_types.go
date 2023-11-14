package models

type ResourceType int

const (
	UnknownResourceType ResourceType = iota
	Cpu
	Memory
	Disk
	Bandwidth
)

func (rt ResourceType) String() string {
	switch rt {
	case Cpu:
		return "Cpu"
	case Memory:
		return "Memory"
	case Disk:
		return "Bandwidth"
	default:
		return "UnknownResourceType"
	}
}
