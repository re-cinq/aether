package calculator

import (
	"fmt"
	"testing"

	factors "github.com/re-cinq/cloud-carbon/pkg/types/v1/factors"
	"github.com/stretchr/testify/require"
)

func TestHourlyEmbodiedEmissions(t *testing.T) {
	assert := require.New(t)

	type testcase struct {
		description string
		specs       factors.Embodied
		expected    float64
	}
	tt := []testcase{
		{
			description: "e2-standard-2",
			specs: factors.Embodied{
				TotalVCPU:                 32.0,
				VCPU:                      2,
				TotalEmbodiedKiloWattCO2e: 12255.46,
			},
			expected: 0.014573178272450533,
		},
		{
			description: "n2-standard-2",
			specs: factors.Embodied{
				TotalVCPU:                 128.0,
				VCPU:                      2,
				TotalEmbodiedKiloWattCO2e: 1888.46,
			},
			expected: 0.0005614000665905632,
		},
		{
			description: "n2d-standard-2",
			specs: factors.Embodied{
				TotalVCPU:                 224.0,
				VCPU:                      2,
				TotalEmbodiedKiloWattCO2e: 2321.46,
			},
			expected: 0.00039435543052837574,
		},
		{
			description: "t2d-standard-2",
			specs: factors.Embodied{
				TotalVCPU:                 60.0,
				VCPU:                      2,
				TotalEmbodiedKiloWattCO2e: 1310.92,
			},
			expected: 0.0008313800101471335,
		},
		{
			description: "n1-standard-2",
			specs: factors.Embodied{
				TotalVCPU:                 96.0,
				VCPU:                      2,
				TotalEmbodiedKiloWattCO2e: 1677.48,
			},
			expected: 0.0006649067732115677,
		},
	}
	for _, test := range tt {
		t.Run(fmt.Sprintf("correct for %s", test.description), func(t *testing.T) {
			// #nosec G601
			res := hourlyEmbodiedEmissions(&test.specs)
			assert.Equal(test.expected, res)
		})
	}
}
