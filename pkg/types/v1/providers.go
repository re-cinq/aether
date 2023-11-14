package models

type Provider int

const (
	UnknownProvider Provider = iota
	Aws
	Azure
	Baremetal
	GCP
)

func (p Provider) String() string {
	switch p {
	case Aws:
		return "Aws"
	case Azure:
		return "Azure"
	case Baremetal:
		return "Baremetal"
	case GCP:
		return "GCP"
	default:
		return "UnknownProvider"
	}
}
