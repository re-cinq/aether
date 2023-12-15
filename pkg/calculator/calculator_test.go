package calculator

import (
	"fmt"
	"testing"

	v1 "github.com/re-cinq/cloud-carbon/pkg/types/v1"
)

func TestCalculateEmissions(t *testing.T) {
	c := calculate{
		cores:    4,
		usage:    v1.Percentage(24),
		minWatts: 1.3423402398570,
		maxWatts: 4.00498247528,
		chip:     35.23458732,
		pue:      1.0123,
		gridCO2e: 0.00023,
	}

	fmt.Println(c.operationalEmissions(5))
}

func TestCalculateEmbodiedEmissions(t *testing.T) {
	c := calculate{
		totalEmbodied: 1706.48,
	}
	fmt.Println(c.embodiedEmissions(5))
}
