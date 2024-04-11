package calculator

import (
	"context"
	"testing"
	"time"

	v1 "github.com/re-cinq/aether/pkg/types/v1"
	data "github.com/re-cinq/emissions-data/pkg/types/v2"
	"github.com/stretchr/testify/assert"
)

func params() *parameters {
	m := &v1.Metric{
		Name:  "basic",
		Usage: 27,
	}
	return &parameters{
		grid:   7,
		pue:    1.2,
		metric: m,
		// using t3.micro AWS instance as default
		powerCPU: []data.Wattage{
			{
				Percentage: 0,
				Wattage:    1.21,
			},
			{
				Percentage: 10,
				Wattage:    3.05,
			},
			{
				Percentage: 50,
				Wattage:    7.16,
			},
			{
				Percentage: 100,
				Wattage:    9.96,
			},
		},
		powerRAM: []data.Wattage{
			{
				Percentage: 0,
				Wattage:    0.15,
			},
			{
				Percentage: 10,
				Wattage:    0.24,
			},
			{
				Percentage: 50,
				Wattage:    0.62,
			},
			{
				Percentage: 100,
				Wattage:    1.00,
			},
		},
		vCPU:           2,
		embodiedFactor: 1000,
	}
}

func TestCalculateCPU(t *testing.T) {
	type testcase struct {
		name     string
		interval time.Duration // this is nanoseconds
		params   *parameters
		expRes   float64
		hasErr   bool
		expErr   string
	}

	for _, test := range []*testcase{
		func() *testcase {
			// Default test case
			return &testcase{
				name:     "default t3.micro at 27% usage over 5m",
				interval: 5 * time.Minute,
				params:   params(),
				expRes:   0.007453764094512195,
			}
		}(),
		func() *testcase {
			// vCPUs not set in params, but set in metric
			p := params()
			p.vCPU = 0
			p.metric.UnitAmount = 2
			p.metric.Unit = v1.VCPU
			return &testcase{
				name:     "vCPU set in metric, not params",
				interval: 5 * time.Minute,
				params:   p,
				expRes:   0.007453764094512195,
			}
		}(),
		func() *testcase {
			// vCPUs not set
			p := params()
			p.vCPU = 0
			return &testcase{
				name:     "vCPU not set",
				interval: 5 * time.Minute,
				params:   p,
				hasErr:   true,
				expErr:   "error vCPU set to 0",
				expRes:   0.00,
			}
		}(),

		func() *testcase {
			// Calculate the default values over
			// a 30 second interval, instead of
			// 5 minutes
			return &testcase{
				name:     "default 30 second interval",
				interval: 30 * time.Second,
				params:   params(),
				expRes:   0.0007453764094512195,
			}
		}(),

		func() *testcase {
			// Calculate the default values over
			// one hour
			return &testcase{
				name:     "1 hour interval",
				interval: 1 * time.Hour,
				params:   params(),
				expRes:   0.08944516913414635,
			}
		}(),

		func() *testcase {
			// calculate with 4 vCPUs
			p := params()
			p.vCPU = 4
			return &testcase{
				name:     "4 vCPU",
				interval: 5 * time.Minute,
				params:   p,
				expRes:   0.01490752818902439,
			}
		}(),

		func() *testcase {
			p := params()
			p.pue = 1.0
			// test if PUE is exactly 1
			return &testcase{
				name:     "PUE is exactly 1.0",
				interval: 5 * time.Minute,
				params:   p,
				expRes:   0.006211470078760163,
			}
		}(),

		func() *testcase {
			p := params()
			p.grid = 402
			// test an extremely high grid CO2e
			// This value was collected from azures
			// Germany West Central region
			return &testcase{
				name:     "High grid CO2e",
				interval: 5 * time.Minute,
				params:   p,
				expRes:   0.42805902371341464,
			}
		}(),

		func() *testcase {
			// create a relatively large server with higher
			// than typical min and max watts, 32 vCPUs, and
			// a utilization of 90%
			p := params()
			p.metric.UnitAmount = 32
			p.metric.Usage = 90
			return &testcase{
				name:     "large server and large workload",
				interval: 5 * time.Minute,
				params:   p,
				expRes:   0.013218390243902438,
			}
		}(),
	} {
		t.Run(test.name, func(t *testing.T) {
			res, err := cpu(context.TODO(), test.interval, test.params)
			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCubicSplineInterpolation(t *testing.T) {
	type testcase struct {
		name     string
		powerCPU []data.Wattage
		usage    float64
		expRes   float64
		hasErr   bool
		expErr   string
	}
	testcases := []testcase{
		{
			name: "t3.micro instance at 27%",
			powerCPU: []data.Wattage{
				{
					Percentage: 0,
					Wattage:    1.21,
				},
				{
					Percentage: 10,
					Wattage:    3.05,
				},
				{
					Percentage: 50,
					Wattage:    7.16,
				},
				{
					Percentage: 100,
					Wattage:    9.96,
				},
			},
			usage:  27.00,
			expRes: 0.005324117210365854,
		},
		{
			name:     "empty wattage",
			powerCPU: []data.Wattage{},
			usage:    27.01,
			expRes:   0,
			hasErr:   true,
			expErr:   "error: cannot calculate CPU energy, no wattage found",
		},
		{
			name: "at exactly 10% utilization",
			powerCPU: []data.Wattage{
				{
					Percentage: 0,
					Wattage:    1.21,
				},
				{
					Percentage: 10,
					Wattage:    3.05,
				},
				{
					Percentage: 50,
					Wattage:    7.16,
				},
				{
					Percentage: 100,
					Wattage:    9.96,
				},
			},
			usage:  10,
			expRes: 0.0030499999999999993,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			res, err := cubicSplineInterpolation(test.powerCPU, test.usage)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
		})
	}
}

func TestCalculateMemory(t *testing.T) {
	type testcase struct {
		name   string
		params *parameters
		expRes float64
		hasErr bool
		expErr string
	}
	for _, test := range []*testcase{
		func() *testcase {
			// pass: default test case
			return &testcase{
				name:   "default t3.micro at 27%",
				params: params(),
				expRes: 0.0033801701414634144,
			}
		}(),
		func() *testcase {
			// fail: powerRAM wattage not set
			p := params()
			p.powerRAM = []data.Wattage{}
			return &testcase{
				name:   "default t3.micro at 27%",
				params: p,
				expRes: 0,
				hasErr: true,
				expErr: "RAM wattage data not found for memory calculation",
			}
		}(),
	} {
		t.Run(test.name, func(t *testing.T) {
			res, err := memory(context.TODO(), test.params)
			assert.Equalf(t, test.expRes, res, "Result should be: %v, got: %v", test.expRes, res)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
